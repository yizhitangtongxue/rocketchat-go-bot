package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

var ollamaHost string = "http://localhost:11434/api/generate"

// 定义要发送的数据结构
type GenerateRequest struct {
	Model   string         `json:"model"`
	Prompt  string         `json:"prompt"`
	Stream  bool           `json:"stream"`
	Options RequestOptions `json:"options"`
}

type RequestOptions struct {
	Temperature float64 `json:"temperature"`
}

type ResponseStruct struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Response           string    `json:"response"`
	Done               bool      `json:"done"`
	Context            []int     `json:"context"`
	TotalDuration      int64     `json:"total_duration"`
	LoadDuration       int64     `json:"load_duration"`
	PromptEvalCount    int       `json:"prompt_eval_count"`
	PromptEvalDuration int64     `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"`
	EvalDuration       int64     `json:"eval_duration"`
}

func SendMessage(model string, prompt string) (ResponseStruct, error) {
	// 准备请求数据
	requestData := GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
		Options: RequestOptions{
			Temperature: 0.8,
		},
	}

	// 将数据序列化为JSON格式
	jsonPayload, err := json.Marshal(requestData)
	if err != nil {
		return ResponseStruct{}, err
	}

	// 创建HTTP请求
	req, err := http.NewRequest(http.MethodPost, ollamaHost, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return ResponseStruct{}, err
	}

	// 设置请求头，通常需要设置Content-Type为application/json
	req.Header.Set("Content-Type", "application/json")

	// 创建HTTP客户端并执行请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ResponseStruct{}, err
	}
	defer resp.Body.Close()

	// 读取并打印响应体
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ResponseStruct{}, err
	}

	var response ResponseStruct
	err = json.Unmarshal([]byte(respBody), &response)
	if err != nil {
		return ResponseStruct{}, err
	}
	return response, nil
}
