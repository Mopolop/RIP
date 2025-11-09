package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"user-auth-system/internal/app/ds"

	"github.com/gin-gonic/gin"
)

// POST /orders/delete/:id - пометить заказ как удалённый
func (h *Handler) DeleteMaterialsOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.SetOrderStatus(ctx.Request.Context(), id, "удален"); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// После удаления можно редиректить на главную
	ctx.Redirect(http.StatusSeeOther, "/")
}

// GET /api/orders/draft/cart
// GET /api/orders/draft/cart
// GET /api/orders/draft/cart
func (h *Handler) GetDraftCartAPI(ctx *gin.Context) {
	userID, _ := h.getUserFromContext(ctx)

	// Если пользователь не найден — как в GetCartJSON
	if userID == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"request_id": -1,
			"count":      0,
		})
		return
	}

	// Ищем черновик
	order, err := h.Repository.GetDraftOrder(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"request_id": -1,
			"count":      0,
		})
		return
	}

	if order == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"request_id": -1,
			"count":      0,
		})
		return
	}

	count, err := h.Repository.GetOrderMaterialsCount(ctx.Request.Context(), order.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"request_id": -1,
			"count":      0,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"request_id": order.ID,
		"count":      count,
	})
}

// GetOrdersAPI godoc
// @Summary      Получить список заказов
// @Description  Возвращает список заказов с возможностью фильтрации по статусу и диапазону дат. Разрешённые статусы: "сформирован", "завершен", "отклонен".
// @Tags		 Заявки с материалами
// @Accept       json
// @Produce      json
// @Param        status  query     string  false  "Статус заказа (сформирован, завершен, отклонен), можно указать несколько через запятую"
// @Param        start   query     string  false  "Дата начала фильтрации (формат YYYY-MM-DD)"
// @Param        end     query     string  false  "Дата окончания фильтрации (формат YYYY-MM-DD)"
// @Success      200     {object}  ds.OrdersListResponse  "Успешный ответ со списком заказов"
// @Failure      500     {object}  ds.ErrorResponse       "Ошибка на сервере"
// @Security BearerAuth
// @Router       /api/orders [get]
func (h *Handler) GetOrdersAPI(ctx *gin.Context) {
	status := ctx.Query("status")
	start := ctx.Query("start") // формат YYYY-MM-DD
	end := ctx.Query("end")

	// Получаем ID авторизованного пользователя
	userID, _ := h.getUserFromContext(ctx)

	orders, err := h.Repository.GetOrdersFilteredForUser(ctx.Request.Context(), status, start, end, userID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
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

	order, materials, err := h.Repository.GetOrderByID(ctx.Request.Context(), orderID)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	// Проверяем, что запрашиваемый заказ принадлежит текущему пользователю или пользователь — Admin
	userID, roleID := h.getUserFromContext(ctx)
	isAdmin := false
	if roleID == 2 { // role.Admin == 2
		isAdmin = true
	}
	if !isAdmin && order.CreatorID != userID {
		// не владелец и не админ — запрещено
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	// Формируем DTO материалов
	var materialsDTO []ds.MaterialInOrder
	for _, mmo := range materials {
		var wl *float64
		if mmo.WallLength.Valid {
			wl = &mmo.WallLength.Float64
		}
		materialsDTO = append(materialsDTO, ds.MaterialInOrder{
			ID:          mmo.MaterialID,
			Title:       mmo.Material.Title,
			Consumption: mmo.Material.Consumption,
			Count:       mmo.Material.Count,
			Image:       mmo.Material.Image,
			WallLength:  wl,
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
	if order.DateFinish != nil {
		dateFinish = order.DateFinish
	}

	resp := ds.OrderWithMaterials{
		ID:            order.ID,
		CreatorID:     order.CreatorID,
		Status:        order.RequestStatus,
		CeilingHeight: ceilingHeight,
		WallThickness: wallThickness,
		DateCreate:    order.DateCreate,
		DateForm:      order.DateForm,
		DateFinish:    dateFinish,
		Materials:     materialsDTO,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"order": resp,
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

	if err := h.Repository.UpdateMaterialOrder(ctx.Request.Context(), id, req); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	order, _, err := h.Repository.GetOrderByID(ctx.Request.Context(), id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"order": order,
	})
}

func (h *Handler) FormMaterialOrderAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.FormMaterialOrder(ctx.Request.Context(), orderID); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Возвращаем обновлённый заказ
	order, _, err := h.Repository.GetOrderByID(ctx.Request.Context(), orderID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"order": order,
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

	if err := h.Repository.CompleteOrRejectOrder(ctx.Request.Context(), orderID, req); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	order, materials, err := h.Repository.GetOrderByID(ctx.Request.Context(), orderID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
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
	if err := h.Repository.SoftDeleteOrder(ctx.Request.Context(), orderID); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"orderID": orderID,
	})
}
