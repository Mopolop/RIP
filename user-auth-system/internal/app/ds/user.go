package ds

import (
	"user-auth-system/internal/app/role"
)

type User struct {
	ID       int       `gorm:"primaryKey;autoIncrement"` // обязательный PK для GORM
	Login    string    `gorm:"unique;not null"`
	Role     role.Role `gorm:"type:int"` // хранить enum как int
	Password string
}

type RegisterReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Структура ответа
type RegisterResp struct {
	Ok bool `json:"ok"`
}
