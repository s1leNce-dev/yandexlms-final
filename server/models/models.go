package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Login    string `json:"login" gorm:"unique"`
	Password string `json:"password"`
}

type Expression struct {
	gorm.Model
	Status string  `json:"status"`
	Result float64 `json:"result"`
	Tasks  []Task  `json:"-"`
}

type Task struct {
	gorm.Model
	ExpressionID  uint   `json:"-"`
	Expression    string `json:"expression"`
	OperationTime int    `json:"operation_time"`
	Done          bool   `json:"-"`
}
