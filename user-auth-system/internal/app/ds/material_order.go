package ds

import (
	"database/sql"
	"time"
)

type MaterialOrder struct {
	ID            int             `gorm:"primaryKey;autoIncrement"`       // первичный ключ
	CeilingHeight sql.NullFloat64 `gorm:"type:NUMERIC(5,2);default:null"` // высота потолка
	WallThickness sql.NullFloat64 `gorm:"type:NUMERIC(5,2);default:null"` // толщина стены
	CreatorID     int             `gorm:"not null"`                       // ID создателя заказа
	ModeratorID   *int            `gorm:""`                               // ID модератора заказа (nullable)
	RequestStatus string          `gorm:"type:varchar(50);not null"`      // статус заказа
	DateCreate    time.Time       `gorm:"not null;autoCreateTime"`        // дата создания
	DateForm      *time.Time      `gorm:"default:null"`                   // дата формирования заявки
	DateFinish    *time.Time      `gorm:"default:null"`                   // дата завершения (может быть null)

	// Связи с пользователями
	Creator   *User `gorm:"foreignKey:CreatorID;references:ID"`   // автор заказа (обязательно)
	Moderator *User `gorm:"foreignKey:ModeratorID;references:ID"` // модератор (nullable)
}

type OrderResponse struct {
	ID         int        `json:"id"`
	Status     string     `json:"status"`
	DateCreate time.Time  `json:"date_create"`
	DateForm   *time.Time `json:"date_form,omitempty"`
	DateFinish *time.Time `json:"date_finish,omitempty"`
}

type OrdersListResponse struct {
	Status string          `json:"status" example:"success"`
	Orders []OrderResponse `json:"orders"`
}

type UpdateOrderRequest struct {
	CeilingHeight *float64 `json:"ceiling_height,omitempty"`
	WallThickness *float64 `json:"wall_thickness,omitempty"`
}

type OrderWithMaterials struct {
	ID            int               `json:"id"`
	CreatorID     int               `json:"creator_id"`
	Status        string            `json:"status"`
	CeilingHeight *float64          `json:"ceiling_height,omitempty"`
	WallThickness *float64          `json:"wall_thickness,omitempty"`
	DateCreate    time.Time         `json:"date_create"`
	DateForm      *time.Time        `json:"date_form,omitempty"`
	DateFinish    *time.Time        `json:"date_finish,omitempty"`
	Materials     []MaterialInOrder `json:"materials"`
}

type CompleteOrderRequest struct {
	Status      string `json:"status" binding:"required,oneof=завершен отклонен"`
	ModeratorID int    `json:"moderator_id" binding:"required"`
}
