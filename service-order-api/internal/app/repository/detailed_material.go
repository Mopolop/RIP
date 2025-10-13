package repository

import "service-order-api/internal/app/ds"

func (r *Repository) GetMaterial(id int) (ds.Material, error) {
	var material ds.Material
	result := r.db.First(&material, id)
	if result.Error != nil {
		return ds.Material{}, result.Error
	}
	return material, nil
}
