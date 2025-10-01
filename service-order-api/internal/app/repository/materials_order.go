package repository

import (
	"db-integration/internal/app/ds"
	"fmt"
)

func (r *Repository) GetOrderByID(id int) (ds.MaterialOrder, []ds.MaterialMaterialOrder, error) {
	var order ds.MaterialOrder
	if err := r.db.First(&order, id).Error; err != nil {
		return ds.MaterialOrder{}, nil, fmt.Errorf("заказ с ID=%d не найден", id)
	}

	var mmos []ds.MaterialMaterialOrder
	if err := r.db.Preload("Material").Where("material_order_id = ?", id).Find(&mmos).Error; err != nil {
		return ds.MaterialOrder{}, nil, err
	}

	return order, mmos, nil
}
