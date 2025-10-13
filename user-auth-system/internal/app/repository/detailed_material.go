package repository

import "db-integration/internal/app/ds"

func (r *Repository) GetMaterial(id int) (ds.Material, error) {
	var material ds.Material
	result := r.db.First(&material, id)
	if result.Error != nil {
		return ds.Material{}, result.Error
	}
	return material, nil
}
