package handler

import (
	"db-integration/internal/app/ds"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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
	status := ctx.Query("status")
	start := ctx.Query("start") // формат YYYY-MM-DD
	end := ctx.Query("end")

	orders, err := h.Repository.GetOrdersFiltered(status, start, end)
	if err != nil {
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
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Получаем заказ и материалы
	order, materials, err := h.Repository.GetOrderByID(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	// Конвертация sql.NullTime в *time.Time для DateForm и DateFinish
	var dateForm *time.Time
	if !order.DateForm.IsZero() {
		dateForm = &order.DateForm
	}

	var dateFinish *time.Time
	if order.DateFinish.Valid {
		dateFinish = &order.DateFinish.Time
	}

	// Формируем ответ
	orderResp := ds.OrderResponse{
		ID:         order.ID,
		Status:     order.RequestStatus,
		DateCreate: order.DateCreate,
		DateForm:   dateForm,
		DateFinish: dateFinish,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"order":     orderResp,
		"materials": materials,
	})
}

func (h *Handler) UpdateMaterialOrderAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Читаем тело запроса как map
	var raw map[string]interface{}
	if err := ctx.BindJSON(&raw); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Разрешённые поля
	allowed := map[string]bool{
		"ceiling_height": true,
		"wall_thickness": true,
	}

	// Проверяем лишние поля
	for k := range raw {
		if !allowed[k] {
			h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("недопустимое поле: %s", k))
			return
		}
	}

	// Преобразуем map в DTO
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

	order, _, err := h.Repository.GetOrderByID(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"order":  order,
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
