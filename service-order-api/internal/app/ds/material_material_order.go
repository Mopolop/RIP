package ds

import "database/sql"

// MaterialMaterialOrder представляет связь "материал ↔ заказ" (многие ко многим)
type MaterialMaterialOrder struct {
	MaterialID          int             `gorm:"primaryKey" json:"material_id"`
	MaterialOrderID     int             `gorm:"primaryKey" json:"material_order_id"`
	WallLength          sql.NullFloat64 `json:"wall_length,omitempty"`
	MaterialConsumption int             `json:"material_consumption"`
	MortarConsumption   sql.NullFloat64 `json:"mortar_consumption"`

	Material      *Material      `gorm:"foreignKey:MaterialID;references:ID" json:"material,omitempty"`
	MaterialOrder *MaterialOrder `gorm:"foreignKey:MaterialOrderID;references:ID" json:"material_order,omitempty"`
}
