package util

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	Environment          string        `mapstructure:"ENVIRONMENT"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	MigrationURL         string        `mapstructure:"MIGRATION_URL"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	ServiceAddress       string        `mapstructure:"SERVICE_ADDRESS"`
	EmailSenderName      string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress   string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword  string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
	FolderAlgoPath       string        `mapstructure:"FOLDER_ALGO_PATH"`
	CertFilePath         string        `mapstructure:"CERT_FILE_PATH"`
	KeyFilePath          string        `mapstructure:"KEY_FILE_PATH"`
	NPMPath              string        `mapstructure:"NPM_PATH"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func LoggerOutput(config Config) zerolog.LevelWriter {
	fAll, _ := os.OpenFile(
		"web3school.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	// https: //github.com/rs/zerolog/issues/150
	rootLogger := zerolog.New(fAll).With().Str("some_key", "some_val").Timestamp().Logger()
	// consoleLogger := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}
	consoleLogger := zerolog.ConsoleWriter{Out: os.Stderr}
	// consoleLogger := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}

	// https://zerolog.io/#multiple-log-output/
	multi := zerolog.MultiLevelWriter(os.Stdout, rootLogger)
	switch config.Environment {
	case "development":
		// https://github.com/rs/zerolog?tab=readme-ov-file#pretty-logging
		multi = zerolog.MultiLevelWriter(rootLogger, consoleLogger)
	}

	return multi
}

func CORSConfig() gin.HandlerFunc {
	return cors.New(cors.Config{
		// AllowOrigins:     []string{"https://foo.com"},
		AllowAllOrigins:  true,
		AllowMethods:     []string{"POST", "GET", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
		MaxAge: 12 * time.Hour,
	})
}
