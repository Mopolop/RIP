package repository

import (
	"db-integration/internal/app/ds"
	"strings"
)

func (r *Repository) GetMaterials() ([]ds.Material, error) {
	var materials []ds.Material
	result := r.db.Find(&materials)
	if result.Error != nil {
		return nil, result.Error
	}
	return materials, nil
}

func (r *Repository) GetMaterialsByTitle(title string) ([]ds.Material, error) {
	var materials []ds.Material
	result := r.db.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(title)+"%").Find(&materials)
	if result.Error != nil {
		return nil, result.Error
	}
	return materials, nil
}
