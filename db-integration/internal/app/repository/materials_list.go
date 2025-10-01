package repository

import (
	"db-integration/internal/app/ds"
	"errors"
	"strings"

	"gorm.io/gorm"
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

// Получаем черновой заказ пользователя
func (r *Repository) GetDraftOrder(userID int) (*ds.MaterialOrder, error) {
	var order ds.MaterialOrder
	err := r.db.Where("creator_id = ? AND request_status = ?", userID, "черновик").First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // черновик отсутствует
		}
		return nil, err
	}
	return &order, nil
}

// Создаём новый черновой заказ
func (r *Repository) CreateDraftOrder(userID int) (*ds.MaterialOrder, error) {
	order := ds.MaterialOrder{
		CreatorID:     userID,
		ModeratorID:   userID, // можно назначить себя модератором, либо 0/NULL
		RequestStatus: "черновик",
	}
	if err := r.db.Create(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// Добавляем материал в заказ, если его там нет
func (r *Repository) AddMaterialToOrder(orderID int, materialID int) error {
	var count int64
	err := r.db.Model(&ds.MaterialMaterialOrder{}).Where("material_order_id = ? AND material_id = ?", orderID, materialID).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // материал уже добавлен
	}

	item := ds.MaterialMaterialOrder{
		MaterialID:      materialID,
		MaterialOrderID: orderID,
	}
	return r.db.Create(&item).Error
}

// Получаем количество материалов в заказе
func (r *Repository) GetOrderMaterialsCount(orderID int) (int64, error) {
	var count int64
	err := r.db.Model(&ds.MaterialMaterialOrder{}).Where("material_order_id = ?", orderID).Count(&count).Error
	return count, err
}

// SetOrderStatus обновляет статус заказа по его ID
func (r *Repository) SetOrderStatus(orderID int, status string) error {
	return r.db.Model(&ds.MaterialOrder{}).Where("id = ?", orderID).Update("request_status", status).Error
}
