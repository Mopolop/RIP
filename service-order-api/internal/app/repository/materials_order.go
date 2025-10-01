package repository

import (
	"db-integration/internal/app/ds"
	"fmt"
	"time"
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

func (r *Repository) GetOrdersFiltered(status, start, end string) ([]ds.OrderResponse, error) {
	var orders []ds.OrderResponse

	query := r.db.
		Table("material_orders mo").
		Select(`mo.id, 
		        mo.request_status as status, 
		        mo.date_create, 
		        mo.date_form, 
		        mo.date_finish, 
		        u1.login as moderator, 
		        u2.login as creator`).
		Joins("LEFT JOIN users u1 ON u1.id = mo.moderator_id").
		Joins("LEFT JOIN users u2 ON u2.id = mo.creator_id").
		Where("mo.request_status NOT IN ?", []string{"удален", "черновик"})

	// фильтр по статусу
	if status != "" {
		query = query.Where("mo.request_status = ?", status)
	}

	// фильтр по диапазону дат
	if start != "" && end != "" {
		query = query.Where("mo.date_create BETWEEN ? AND ?", start, end)
	}

	if err := query.Scan(&orders).Error; err != nil {
		return nil, fmt.Errorf("ошибка при получении заказов: %w", err)
	}

	return orders, nil
}

func (r *Repository) UpdateMaterialOrder(orderID int, req ds.UpdateOrderRequest) error {
	updates := make(map[string]interface{})

	if req.CeilingHeight != nil {
		updates["ceiling_height"] = *req.CeilingHeight
	}
	if req.WallThickness != nil {
		updates["wall_thickness"] = *req.WallThickness
	}

	if len(updates) == 0 {
		return nil // ничего менять не нужно
	}

	return r.db.Model(&ds.MaterialOrder{}).Where("id = ?", orderID).Updates(updates).Error
}

func (r *Repository) FormMaterialOrder(orderID int) error {
	// Проверяем, что все wall_length заполнены
	var count int64
	if err := r.db.Model(&ds.MaterialMaterialOrder{}).
		Where("material_order_id = ? AND wall_length IS NULL", orderID).
		Count(&count).Error; err != nil {
		return fmt.Errorf("ошибка проверки wall_length: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("нельзя сформировать заказ: не все wall_length заполнены")
	}

	// Обновляем заказ: статус и date_form
	updates := map[string]interface{}{
		"request_status": "сформирован",
		"date_form":      time.Now(),
	}

	// Обновляем только если текущий статус черновик
	if err := r.db.Model(&ds.MaterialOrder{}).
		Where("id = ? AND request_status = ?", orderID, "черновик").
		Updates(updates).Error; err != nil {
		return fmt.Errorf("ошибка обновления заказа: %w", err)
	}

	return nil
}
