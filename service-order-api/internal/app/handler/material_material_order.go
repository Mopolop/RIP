package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// DELETE /api/orders/:order_id/material/:material_id
func (h *Handler) DeleteMaterialFromOrderAPI(ctx *gin.Context) {
	orderIDStr := ctx.Param("order_id")
	materialIDStr := ctx.Param("material_id")

	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("некорректный ID заказа"))
		return
	}

	materialID, err := strconv.Atoi(materialIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("некорректный ID материала"))
		return
	}

	if err := h.Repository.DeleteMaterialFromOrder(materialID, orderID); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":      "success",
		"material_id": materialID,
		"order_id":    orderID,
		"message":     "услуга удалена из заявки",
	})
}

// PUT /api/orders/:order_id/material/:material_id/wall_length
func (h *Handler) UpdateWallLengthAPI(ctx *gin.Context) {
	orderIDStr := ctx.Param("order_id")
	materialIDStr := ctx.Param("material_id")

	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("некорректный ID заказа"))
		return
	}

	materialID, err := strconv.Atoi(materialIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("некорректный ID материала"))
		return
	}

	var req struct {
		WallLength float64 `json:"wall_length"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("некорректное тело запроса"))
		return
	}

	if err := h.Repository.UpdateWallLength(materialID, orderID, req.WallLength); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":      "success",
		"material_id": materialID,
		"order_id":    orderID,
		"wall_length": req.WallLength,
	})
}
