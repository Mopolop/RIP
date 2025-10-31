package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"service-order-api/internal/app/ds"
	"strconv"
	"time"
)

// POST /orders/delete/:id - пометить заказ как удалённый
func (h *Handler) DeleteMaterialsOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.SetOrderStatus(id, "удален"); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// После удаления можно редиректить на главную
	ctx.Redirect(http.StatusSeeOther, "/")
}

// GET /api/orders/draft/cart
func (h *Handler) GetDraftCartAPI(ctx *gin.Context) {
	// Пока без авторизации — используем userID = 1
	userID := 1

	// Ищем черновик
	order, err := h.Repository.GetDraftOrder(userID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	if order == nil {
		// Если черновика нет — возвращаем пустую корзину
		ctx.JSON(http.StatusOK, gin.H{
			"status":    "success",
			"orderID":   0,
			"itemCount": 0,
		})
		return
	}

	// Считаем количество услуг
	count, err := h.Repository.GetOrderMaterialsCount(order.ID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"orderID":   order.ID,
		"itemCount": count,
	})
}

func (h *Handler) GetOrdersAPI(ctx *gin.Context) {
	start := ctx.Query("start")
	end := ctx.Query("end")

	fixedStatuses := ctx.Query("status")
	if fixedStatuses == "" {
		fixedStatuses = "завершен,отклонен,отменен"
	}

	orders, err := h.Repository.GetOrdersFiltered(fixedStatuses, start, end)
	if err != nil {
		if err.Error() == "not_found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "заказы со статусом 'черновик' или 'удален' не доступны",
			})
			return
		}
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"orders": orders,
	})
}

func (h *Handler) GetOrderWithMaterialsAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	order, materials, err := h.Repository.GetOrderByID(orderID)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	// Формируем DTO материалов
	var materialsDTO []ds.MaterialInOrder
	for _, mmo := range materials {
		var wl *float64
		if mmo.WallLength.Valid {
			wl = &mmo.WallLength.Float64
		}

		var mortarConsumption *float64
		if mmo.MortarConsumption.Valid {
			mortarConsumption = &mmo.MortarConsumption.Float64
		}

		materialsDTO = append(materialsDTO, ds.MaterialInOrder{
			ID:                  mmo.MaterialID,
			Title:               mmo.Material.Title,
			Consumption:         mmo.Material.Consumption,
			Count:               mmo.Material.Count,
			Image:               mmo.Material.Image,
			WallLength:          wl,
			MaterialConsumption: mmo.MaterialConsumption,
			MortarConsumption:   mortarConsumption,
		})
	}

	// DTO заказа
	var ceilingHeight, wallThickness *float64
	if order.CeilingHeight.Valid {
		ceilingHeight = &order.CeilingHeight.Float64
	}
	if order.WallThickness.Valid {
		wallThickness = &order.WallThickness.Float64
	}

	var dateFinish *time.Time
	if order.DateFinish.Valid {
		dateFinish = &order.DateFinish.Time
	}

	resp := ds.OrderWithMaterials{
		ID:            order.ID,
		CreatorID:     order.CreatorID,
		Status:        order.RequestStatus,
		CeilingHeight: ceilingHeight,
		WallThickness: wallThickness,
		DateCreate:    order.DateCreate,
		DateForm:      &order.DateForm,
		DateFinish:    dateFinish,
		Materials:     materialsDTO,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"order":  resp,
	})
}

func (h *Handler) UpdateMaterialOrderAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var raw map[string]interface{}
	if err := ctx.BindJSON(&raw); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	allowed := map[string]bool{
		"ceiling_height": true,
		"wall_thickness": true,
	}

	for k := range raw {
		if !allowed[k] {
			h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("недопустимое поле: %s", k))
			return
		}
	}

	var req ds.UpdateOrderRequest
	if v, ok := raw["ceiling_height"]; ok {
		if f, ok := v.(float64); ok {
			req.CeilingHeight = &f
		}
	}
	if v, ok := raw["wall_thickness"]; ok {
		if f, ok := v.(float64); ok {
			req.WallThickness = &f
		}
	}

	if err := h.Repository.UpdateMaterialOrder(id, req); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// 🔹 Формируем ответ только с обновлёнными полями
	updated := gin.H{
		"id": id,
	}
	if req.CeilingHeight != nil {
		updated["ceiling_height"] = *req.CeilingHeight
	}
	if req.WallThickness != nil {
		updated["wall_thickness"] = *req.WallThickness
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"updated": updated,
	})
}

func (h *Handler) FormMaterialOrderAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.FormMaterialOrder(orderID); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Возвращаем обновлённый заказ
	order, _, err := h.Repository.GetOrderByID(orderID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"order":  order,
	})
}

func (h *Handler) CompleteOrRejectOrderAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var req ds.CompleteOrderRequest
	if err := ctx.BindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.CompleteOrRejectOrder(orderID, req); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	order, materials, err := h.Repository.GetOrderByID(orderID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"order":     order,
		"materials": materials,
	})
}

func (h *Handler) DeleteMaterialsOrderAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("некорректный ID заказа"))
		return
	}

	// Проставляем статус "удален" и дату завершения
	if err := h.Repository.SoftDeleteOrder(orderID); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"orderID": orderID,
	})
}
