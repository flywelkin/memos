package gml

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"net/http"
	"time"
)

const name = ""

type ClaimsMessage struct {
	ApiKey    string `json:"api_key"`
	exp       int64  `json:"exp"`
	timestamp int64  `json:"timestamp"`
	jwt.RegisteredClaims
}

type ChatRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Id        string    `json:"id"`
	Created   int       `json:"created"`
	Model     string    `json:"model"`
	RequestId string    `json:"request_id"`
	Usage     Usage     `json:"usage"`
	Choices   []Choices `json:"choices"`
	Error     Error     `json:"error"`
}

type Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Choices struct {
	Index        int     `json:"index"`
	FinishReason string  `json:"finish_reason"`
	Message      Message `json:"message"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// 生成认证token
func generateToken(key string, secret string) (string, error) {
	claims := ClaimsMessage{
		ApiKey:    key,
		exp:       int64(time.Now().Unix() + 1000),
		timestamp: int64(time.Now().Unix()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	token.Header["alg"] = "HS256"
	token.Header["sign_type"] = "SIGN"

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, err
}

// 发送请求
func Chat(message ChatRequest) (string, error) {
	uri := "https://open.bigmodel.cn/api/paas/v4/chat/completions"

	// 生成认证信息
	token, err := generateToken("e9c74c03ecdd287a16c84762031e4721", "EGHN08zlnVWtsFiC")
	if err != nil {
		return "", err
	}
	messages := []ChatRequest{message}
	// 封装请求
	values := map[string]any{
		"model":    "glm-3-turbo",
		"messages": messages,
	}

	jsonValue, err := json.Marshal(values)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}

	// 设置头信息
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	// 发起请求
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// 读取数据
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	fmt.Println(string(responseBody))

	// 序列化
	chatCompletionResponse := ChatResponse{}
	err = json.Unmarshal(responseBody, &chatCompletionResponse)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if chatCompletionResponse.Error.Code != "" {
		return chatCompletionResponse.Error.Message, nil
	}
	return chatCompletionResponse.Choices[0].Message.Content, nil
}
