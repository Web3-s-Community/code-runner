package models

// https://goswagger.io/use/spec/model.html

// swagger:model
type Submission struct {
	Code   string `json:"code"`
	CodeID string `json:"code_id"`
	Uuid   string `json:"uuid"`
}

// swagger:model
type SubmissionPlayground struct {
	Code     string `json:"code"`
	CodeTest string `json:"code_test"`
	FileName string `json:"file_name"`
}

// swagger:model
type ResponseRunResult struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Stdout  string `json:"stdout"`
	Stderr  string `json:"stderr"`
}
