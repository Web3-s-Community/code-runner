package biz

import (
	"code-runner/module/code/models"
	"code-runner/util"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	srcPath string = "algo/projects/algo-playground"
)

type AlgoStorage interface {
	AlgoConfig(
		ctx context.Context,
		config *util.Config,
		moreKeys ...string,
	) error
}

type AlgoBiz struct {
	store AlgoStorage
}

func NewAlgoBiz(store AlgoStorage) *AlgoBiz {
	return &AlgoBiz{store: store}
}

func ExecuteCodeTest(submission models.SubmissionPlayground) (string, error) {
	output := ""
	code := submission.Code
	codeTest := submission.CodeTest

	fileName := submission.FileName
	err := SaveCodeToFile(srcPath+"/contracts", fileName+".algo.ts", code)
	if err != nil {
		return "", fmt.Errorf("Error saving code to file: %s", err.Error())
	}
	err = SaveCodeToFile(srcPath+"/__test__", fileName+".test.ts", codeTest)
	if err != nil {
		return "", fmt.Errorf("Error saving code test to file: %s", err.Error())
	}

	config, err := util.LoadConfig(".")
	if err != nil {
		return "", fmt.Errorf("Cannot load config: %s", err.Error())
	}

	FolderAlgoPath := config.FolderAlgoPath

	// cmd := exec.Command("npm", "run", "test")

	cmd := exec.Command(config.NPMPath, "run", "test", "--", "__test__/"+fileName+".test.ts")
	// cmd := exec.Command("npm", "run", "test-uuid", fileTestUserPath)
	// https://stackoverflow.com/questions/43135919/how-to-run-a-shell-command-in-a-specific-folder
	cmd.Dir = FolderAlgoPath + srcPath

	// // Start with the current environment
	// cmd.Env = os.Environ()

	// // Set the UUID environment variable
	// cmd.Env = append(cmd.Env, "UUID="+uuid)

	// // Set the CLASSNAME environment variable
	// className := toCamelCase(codeID)
	// cmd.Env = append(cmd.Env, "CLASSNAME="+className)
	// // Set the CODEID environment variable
	// cmd.Env = append(cmd.Env, "CODEID="+codeID)

	// // Set up standard input and output buffers to interact with the command
	// var stdout, stderr bytes.Buffer
	// // cmd.Stdin = bytes.NewReader(code)
	// cmd.Stdout = &stdout
	// cmd.Stderr = &stderr

	// // Run the command
	// err = cmd.Run()
	// if err != nil {
	// 	return "", fmt.Errorf("execution failed: %s\nstderr: %s", err.Error(), stderr.String())
	// }

	// output = stdout.String()

	// // https://mzampetakis.com/posts/Exec-Command/#capturing-the-commands-output
	// dt, err := cmd.Output()
	// if err != nil {
	// 	return "", fmt.Errorf("execution failed: %s\nstderr: %s", err.Error(), "")
	// }
	// output = string(dt)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Run Test Fail: %s\n\n %s", err.Error(), string(out))
	}

	output = string(out)

	// https://stackoverflow.com/questions/1877045/how-do-you-get-the-output-of-a-system-command-in-go
	// if cmd, e := exec.Run("/bin/ls", nil, nil, exec.DevNull, exec.Pipe, exec.MergeWithStdout); e == nil {
	// 	b, _ := ioutil.ReadAll(cmd.Stdout)
	// 	println("output: " + string(b))
	// }

	// // https://zetcode.com/golang/exec-command/
	// stdout, err := cmd.StdoutPipe()

	// if err != nil {
	// 	return "", fmt.Errorf("execution StdoutPipe failed: %s\nstderr: %s", err.Error(), "")
	// }

	// if err := cmd.Start(); err != nil {
	// 	return "", fmt.Errorf("execution Start failed: %s\nstderr: %s", err.Error(), "")
	// }

	// data, err := io.ReadAll(stdout)

	// if err != nil {
	// 	return "", fmt.Errorf("execution ReadAll failed: %s\nstderr: %s", err.Error(), "")
	// }

	// if err := cmd.Wait(); err != nil {
	// 	return "", fmt.Errorf("execution Wait failed: %s\nstderr: %s", err.Error(), "")
	// }

	// output = string(data)

	// https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html

	return output, nil
}

func SaveCodeToFile(folderPath string, fileName string, code string) error {
	// Create the Test if it doesn't exist
	err := os.MkdirAll(folderPath, 0755)
	if err != nil {
		return fmt.Errorf("Error creating Folder User Test: %s\n", err.Error())
	}

	filePath := folderPath + "/" + fileName
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Error creating file: %s\n", err.Error())
	}
	defer file.Close()

	// Copy the request body to the file
	_, err = io.Copy(file, strings.NewReader(code))
	if err != nil {
		return fmt.Errorf("Failed to write request to file: %s", err.Error())
	}

	return nil
}

func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = cases.Title(language.English).String(parts[i])
	}
	return strings.Join(parts, "")
}
