package handler

import (
	"db-integration/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{Repository: r}
}

func (h *Handler) GetMaterials(ctx *gin.Context) {
	searchQuery := ctx.Query("material_search")
	var materials []repository.Material
	var err error

	if searchQuery == "" {
		materials, err = h.Repository.GetMaterials()
	} else {
		materials, err = h.Repository.GetMaterialsByTitle(searchQuery)
	}

	if err != nil {
		logrus.Error(err)
	}

	// получаем заказ, чтобы посчитать количество товаров
	order, _ := h.Repository.GetOrder()
	orderCount := len(order.Materials)

	ctx.HTML(http.StatusOK, "materials_list.html", gin.H{
		"materials":  materials,
		"query":      searchQuery,
		"orderCount": orderCount,
	})
}

// Получение конкретного материала / услуги по ID
func (h *Handler) GetMaterial(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusBadRequest, "Invalid ID")
		return
	}

	material, err := h.Repository.GetMaterial(id)
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusNotFound, "Material not found")
		return
	}

	ctx.HTML(http.StatusOK, "detailed_material.html", gin.H{
		"material": material,
	})
}

func (h *Handler) GetOrder(ctx *gin.Context) {
	idStr := ctx.Param("id") // теперь URL: /materials_order/1
	id, err := strconv.Atoi(idStr)
	if err != nil || id != 1 {
		logrus.Error("Invalid order ID")
		ctx.String(http.StatusBadRequest, "Invalid order ID")
		return
	}

	order, err := h.Repository.GetOrder()
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusInternalServerError, "Unable to get order")
		return
	}

	ctx.HTML(http.StatusOK, "materials_order.html", gin.H{
		"order": order,
	})
}
