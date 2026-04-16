// Package service 定义学籍验证接口，供 CardService 依赖注入
package service

// StudentInfo 学籍验证结果
type StudentInfo struct {
	IDNumber string
	Name     string
	Type     string // "student" | "staff"
}

// StudentValidator 学籍验证接口，CardService 依赖此接口，不感知底层实现
type StudentValidator interface {
	Validate(idNumber string) (*StudentInfo, error)
}
