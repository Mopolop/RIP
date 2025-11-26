package repository

import (
	"context"
	"fmt"
	"strings"
	"time"
	"user-auth-system/internal/app/ds"
)

func (r *Repository) GetOrderByID(ctx context.Context, id int) (ds.MaterialOrder, []ds.MaterialMaterialOrder, error) {
	var order ds.MaterialOrder
	if err := r.db.WithContext(ctx).First(&order, id).Error; err != nil {
		return ds.MaterialOrder{}, nil, fmt.Errorf("заказ с ID=%d не найден", id)
	}

	var mmos []ds.MaterialMaterialOrder
	if err := r.db.WithContext(ctx).Preload("Material").Where("material_order_id = ?", id).Find(&mmos).Error; err != nil {
		return ds.MaterialOrder{}, nil, err
	}

	return order, mmos, nil
}

func (r *Repository) GetOrdersFiltered(status, start, end string) ([]ds.OrderResponse, error) {
	var orders []ds.OrderResponse

	// Разрешённые статусы для выдачи
	allowedStatuses := map[string]bool{
		"сформирован": true,
		"завершен":    true,
		"отклонен":    true,
	}

	query := r.db.
		Table("material_orders mo").
		Select(`mo.id, 
				mo.request_status as status, 
				mo.date_create, 
				mo.date_form, 
				mo.date_finish, 
				u1.login as moderator, 
				u2.login as creator, 
				(SELECT COUNT(*) FROM material_material_orders mmo WHERE mmo.material_order_id = mo.id AND (mmo.mortar_consumption IS NOT NULL OR mmo.material_consumption <> 0)) as results_count`).
		Joins("LEFT JOIN users u1 ON u1.id = mo.moderator_id").
		Joins("LEFT JOIN users u2 ON u2.id = mo.creator_id")

	// фильтр по статусу
	if status != "" {
		statuses := []string{}
		for _, s := range strings.Split(status, ",") {
			s = strings.TrimSpace(s)
			if allowedStatuses[s] { // оставляем только разрешённые
				statuses = append(statuses, s)
			}
		}
		if len(statuses) == 0 {
			// Если после фильтра ничего не осталось — не выдаём ни один заказ
			return []ds.OrderResponse{}, nil
		}
		query = query.Where("mo.request_status IN ?", statuses)
	} else {
		// Если статус не указан — выдаём все разрешённые статусы
		query = query.Where("mo.request_status IN ?", []string{"сформирован", "завершен", "отклонен"})
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

// GetOrdersFilteredForUser возвращает заказы указанного пользователя (creator_id = userID)
func (r *Repository) GetOrdersFilteredForUser(ctx context.Context, status, start, end string, userID int) ([]ds.OrderResponse, error) {
	var orders []ds.OrderResponse

	// Разрешённые статусы для выдачи
	allowedStatuses := map[string]bool{
		"сформирован": true,
		"завершен":    true,
		"отклонен":    true,
	}

	query := r.db.WithContext(ctx).
		Table("material_orders mo").
		Select(`mo.id, 
				mo.request_status as status, 
				mo.date_create, 
				mo.date_form, 
				mo.date_finish, 
				u1.login as moderator, 
				u2.login as creator, 
				(SELECT COUNT(*) FROM material_material_orders mmo WHERE mmo.material_order_id = mo.id AND (mmo.mortar_consumption IS NOT NULL OR mmo.material_consumption <> 0)) as results_count`).
		Joins("LEFT JOIN users u1 ON u1.id = mo.moderator_id").
		Joins("LEFT JOIN users u2 ON u2.id = mo.creator_id").
		Where("mo.creator_id = ?", userID)

	// фильтр по статусу
	if status != "" {
		statuses := []string{}
		for _, s := range strings.Split(status, ",") {
			s = strings.TrimSpace(s)
			if allowedStatuses[s] { // оставляем только разрешённые
				statuses = append(statuses, s)
			}
		}
		if len(statuses) == 0 {
			return []ds.OrderResponse{}, nil
		}
		query = query.Where("mo.request_status IN ?", statuses)
	} else {
		// Если статус не указан — выдаём все разрешённые статусы
		query = query.Where("mo.request_status IN ?", []string{"сформирован", "завершен", "отклонен"})
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

func (r *Repository) UpdateMaterialOrder(ctx context.Context, orderID int, req ds.UpdateOrderRequest) error {
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

	return r.db.WithContext(ctx).Model(&ds.MaterialOrder{}).Where("id = ?", orderID).Updates(updates).Error
}

func (r *Repository) FormMaterialOrder(ctx context.Context, orderID int) error {
	// Проверяем, что все wall_length заполнены
	var count int64
	if err := r.db.WithContext(ctx).Model(&ds.MaterialMaterialOrder{}).
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
	if err := r.db.WithContext(ctx).Model(&ds.MaterialOrder{}).
		Where("id = ? AND request_status = ?", orderID, "черновик").
		Updates(updates).Error; err != nil {
		return fmt.Errorf("ошибка обновления заказа: %w", err)
	}

	return nil
}

func (r *Repository) CompleteOrRejectOrder(ctx context.Context, orderID int, req ds.CompleteOrderRequest) error {
	var order ds.MaterialOrder
	if err := r.db.WithContext(ctx).First(&order, orderID).Error; err != nil {
		return fmt.Errorf("заказ с ID=%d не найден", orderID)
	}

	// Обновляем статус, модератора и дату завершения
	updates := map[string]interface{}{
		"request_status": req.Status,
		"moderator_id":   req.ModeratorID,
		"date_finish":    time.Now(),
	}

	if err := r.db.WithContext(ctx).Model(&ds.MaterialOrder{}).Where("id = ?", orderID).Updates(updates).Error; err != nil {
		return fmt.Errorf("не удалось обновить заказ: %w", err)
	}

	// Если заказ отклонён — прекращаем выполнение, не считаем расход
	if req.Status == "отклонен" {
		return nil
	}

	return nil
}

func (r *Repository) SoftDeleteOrder(ctx context.Context, orderID int) error {
	updates := map[string]interface{}{
		"request_status": "удален",
		"date_form":      time.Now(), // дата завершения
	}

	return r.db.WithContext(ctx).Model(&ds.MaterialOrder{}).Where("id = ?", orderID).Updates(updates).Error
}

// UpdateMaterialMaterialOrderResults применяет результаты расчёта для записей material_material_orders
func (r *Repository) UpdateMaterialMaterialOrderResults(ctx context.Context, orderID int, results []ds.MMResult) error {
	for _, res := range results {
		updates := map[string]interface{}{
			"material_consumption": res.MaterialConsumption,
			"mortar_consumption":   res.MortarConsumption,
		}
		if err := r.db.WithContext(ctx).Model(&ds.MaterialMaterialOrder{}).
			Where("material_order_id = ? AND material_id = ?", orderID, res.MaterialID).
			Updates(updates).Error; err != nil {
			return err
		}
	}
	return nil
}
