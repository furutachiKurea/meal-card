// Package client 提供外部服务调用客户端
package client

import (
	"backend/service"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// HttpStudentValidator 通过 HTTP 调用学籍验证 Mock 服务，实现 service.StudentValidator 接口
type HttpStudentValidator struct {
	baseURL    string
	httpClient *http.Client
}

// NewHttpStudentValidator 创建 HttpStudentValidator，baseURL 从环境变量 STUDENT_SERVICE_URL 读取，
// 默认 http://localhost:9090
func NewHttpStudentValidator() *HttpStudentValidator {
	baseURL := os.Getenv("STUDENT_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:9090"
	}
	return &HttpStudentValidator{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

type validateResponse struct {
	Valid    bool   `json:"valid"`
	IDNumber string `json:"idNumber"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}

// Validate 调用学籍验证服务校验证件号
func (v *HttpStudentValidator) Validate(idNumber string) (*service.StudentInfo, error) {
	url := fmt.Sprintf("%s/validate?idNumber=%s", v.baseURL, idNumber)
	resp, err := v.httpClient.Get(url)
	if err != nil {
		return nil, &service.BizError{Code: service.ErrCodeStudentServiceError, Message: "学籍验证服务暂时不可用"}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, &service.BizError{Code: service.ErrCodeStudentNotFound, Message: "证件号不存在，非本校学生/教职工"}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, &service.BizError{Code: service.ErrCodeStudentServiceError, Message: "学籍验证服务暂时不可用"}
	}

	var body validateResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, &service.BizError{Code: service.ErrCodeStudentServiceError, Message: "学籍验证服务返回异常"}
	}

	return &service.StudentInfo{
		IDNumber: body.IDNumber,
		Name:     body.Name,
		Type:     body.Type,
	}, nil
}
