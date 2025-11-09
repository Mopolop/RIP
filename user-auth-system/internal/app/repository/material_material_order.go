package repository

import (
	"context"
	"user-auth-system/internal/app/ds"
)

func (r *Repository) DeleteMaterialFromOrder(ctx context.Context, materialID, orderID int) error {
	return r.db.WithContext(ctx).
		Where("material_id = ? AND material_order_id = ?", materialID, orderID).
		Delete(&ds.MaterialMaterialOrder{}).Error
}

func (r *Repository) UpdateWallLength(ctx context.Context, materialID, orderID int, wallLength float64) error {
	return r.db.WithContext(ctx).
		Model(&ds.MaterialMaterialOrder{}).
		Where("material_id = ? AND material_order_id = ?", materialID, orderID).
		Update("wall_length", wallLength).Error
}
