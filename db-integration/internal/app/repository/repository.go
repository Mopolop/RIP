package repository

import (
	"db-integration/internal/app/ds"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Авто-миграция моделей
	err = db.AutoMigrate(
		&ds.Material{},
		&ds.MaterialOrder{},
		&ds.MaterialMaterialOrder{},
	)
	if err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

func (r *Repository) GetOrderMaterialsCount(orderID int) (int, error) {
	var count int64
	if err := r.db.Model(&ds.MaterialMaterialOrder{}).
		Where("material_order_id = ?", orderID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
