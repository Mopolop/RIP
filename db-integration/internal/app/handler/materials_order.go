package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// Получение заказа по ID
func (h *Handler) GetOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	order, mmos, err := h.Repository.GetOrderByID(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	ctx.HTML(http.StatusOK, "materials_order.html", gin.H{
		"order":     order,
		"materials": mmos, // вот сюда прокидываем список материалов
	})
}
