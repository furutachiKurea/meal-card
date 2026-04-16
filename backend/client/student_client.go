// Package client 提供外部服务调用客户端
package client

import (
	"backend/service"
)

// MockStudentValidator 硬编码学籍数据，实现 service.StudentValidator 接口，
// 用于开发和演示，无需启动任何外部服务
type MockStudentValidator struct{}

// NewMockStudentValidator 创建 MockStudentValidator
func NewMockStudentValidator() *MockStudentValidator {
	return &MockStudentValidator{}
}

// 硬编码的学籍数据，key 为证件号
var studentDB = map[string]*service.StudentInfo{
	// 学生
	"202100010001": {IDNumber: "202100010001", Name: "张三", Type: "student"},
	"202100010002": {IDNumber: "202100010002", Name: "李四", Type: "student"},
	"202100010003": {IDNumber: "202100010003", Name: "王五", Type: "student"},
	"202100010004": {IDNumber: "202100010004", Name: "赵六", Type: "student"},
	"202100010005": {IDNumber: "202100010005", Name: "陈七", Type: "student"},
	// 教职工
	"T20010001": {IDNumber: "T20010001", Name: "刘老师", Type: "staff"},
	"T20010002": {IDNumber: "T20010002", Name: "孙老师", Type: "staff"},
	"T20010003": {IDNumber: "T20010003", Name: "周主任", Type: "staff"},
}

// Validate 查询硬编码学籍数据，证件号不在名单中返回 STUDENT_NOT_FOUND
func (v *MockStudentValidator) Validate(idNumber string) (*service.StudentInfo, error) {
	if info, ok := studentDB[idNumber]; ok {
		return info, nil
	}
	return nil, &service.BizError{Code: service.ErrCodeStudentNotFound, Message: "证件号不存在，非本校学生/教职工"}
}
