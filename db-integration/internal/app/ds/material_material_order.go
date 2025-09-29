package ds

import "database/sql"

// MaterialMaterialOrder представляет связь "материал ↔ заказ" (многие ко многим)
type MaterialMaterialOrder struct {
	ID              int `gorm:"primaryKey"`
	MaterialID      int `gorm:"not null"`
	MaterialOrderID int `gorm:"not null"`
	WallLength      sql.NullFloat64

	Material      Material      `gorm:"foreignKey:MaterialID;references:ID"`
	MaterialOrder MaterialOrder `gorm:"foreignKey:MaterialOrderID;references:ID"`
}
