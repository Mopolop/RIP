package ds

import "database/sql"

// MaterialMaterialOrder представляет связь "материал ↔ заказ" (многие ко многим)
type MaterialMaterialOrder struct {
	MaterialID          int             `gorm:"not null;uniqueIndex:idx_material_order" json:"material_id"`
	MaterialOrderID     int             `gorm:"not null;uniqueIndex:idx_material_order" json:"material_order_id"`
	WallLength          sql.NullFloat64 `json:"wall_length,omitempty"`
	MaterialConsumption int             `json:"material_consumption"` // шт.
	MortarConsumption   sql.NullFloat64 `json:"mortar_consumption"`   // м³

	Material      *Material      `gorm:"foreignKey:MaterialID;references:ID" json:"material,omitempty"`
	MaterialOrder *MaterialOrder `gorm:"foreignKey:MaterialOrderID;references:ID" json:"material_order,omitempty"`
}
