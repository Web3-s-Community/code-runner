package ginalgorand

import (
	"fmt"
	"net/http"

	"code-runner/module/algorand/biz"
	"code-runner/module/code/models"

	"github.com/gin-gonic/gin"
)

func ExecuteCodePlaygroundHandler(c *gin.Context) {
	var _submission models.SubmissionPlayground

	// Call BindJSON to bind the received JSON to
	// _submission.
	if err := c.BindJSON(&_submission); err != nil {
		return
	}

	output, err := biz.ExecuteCodeTest(_submission)
	if err != nil {
		// return json with model.ResponseRunResult
		c.JSON(http.StatusOK, models.ResponseRunResult{
			Status:  0,
			Message: "Execution failed",
			Stdout:  output,
			Stderr:  fmt.Sprintf("Execution failed: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, models.ResponseRunResult{
		Status:  1,
		Message: "Execution successful",
		Stdout:  output,
		Stderr:  "",
	})
}
