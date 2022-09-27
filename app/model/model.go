package model

type KeyplayMetadata struct {
	Metadata interface{} `json:"metadata"`
	Id       string      `json:"id"`
}

type KeyplayAttribute struct {
	KeyData []string `json:"keyData"`
	Id      string   `json:"id"`
}
