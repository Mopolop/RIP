package handler

import (
	"net/http"

	"db-integration/internal/app/ds"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Получение списка материалов
func (h *Handler) GetMaterials(ctx *gin.Context) {
	searchQuery := ctx.Query("material_search")
	var materials []ds.Material
	var err error

	if searchQuery == "" {
		materials, err = h.Repository.GetMaterials()
	} else {
		materials, err = h.Repository.GetMaterialsByTitle(searchQuery)
	}

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// получаем количество материалов в заказе (например, заказ с ID=1)
	orderCount, err := h.Repository.GetOrderMaterialsCount(1)
	if err != nil {
		logrus.Warn("Не удалось получить количество материалов в заказе: ", err)
		orderCount = 0
	}

	ctx.HTML(http.StatusOK, "materials_list.html", gin.H{
		"materials":  materials,
		"query":      searchQuery,
		"orderCount": orderCount,
	})

}
