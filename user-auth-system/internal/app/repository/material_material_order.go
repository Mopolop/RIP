package repository

import (
	"user-auth-system/internal/app/ds"
)

func (r *Repository) DeleteMaterialFromOrder(materialID, orderID int) error {
	return r.db.
		Where("material_id = ? AND material_order_id = ?", materialID, orderID).
		Delete(&ds.MaterialMaterialOrder{}).Error
}

func (r *Repository) UpdateWallLength(materialID, orderID int, wallLength float64) error {
	return r.db.
		Model(&ds.MaterialMaterialOrder{}).
		Where("material_id = ? AND material_order_id = ?", materialID, orderID).
		Update("wall_length", wallLength).Error
}
