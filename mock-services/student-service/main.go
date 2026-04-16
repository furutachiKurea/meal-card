// mock-services/student-service 学籍验证 Mock 服务
// 端口 :9090，硬编码学生/教职工名单，修改名单直接改 students map。
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type person struct {
	Name string
	Type string // "student" | "staff"
}

// students 硬编码的学生/教职工名单，key 为 12 位学号/工号
var students = map[string]person{
	"202305133513": {Name: "张三", Type: "student"},
	"202305133514": {Name: "李四", Type: "student"},
	"202305133515": {Name: "王五", Type: "student"},
	"202201000001": {Name: "陈老师", Type: "staff"},
	"202201000002": {Name: "刘老师", Type: "staff"},
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	idNumber := r.URL.Query().Get("idNumber")
	w.Header().Set("Content-Type", "application/json")

	p, ok := students[idNumber]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{"valid": false})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"valid":    true,
		"idNumber": idNumber,
		"name":     p.Name,
		"type":     p.Type,
	})
}

func main() {
	http.HandleFunc("/validate", validateHandler)
	log.Println("学籍验证 Mock 服务启动，监听 :9090")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
