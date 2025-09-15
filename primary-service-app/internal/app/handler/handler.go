package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"primary-service-app/internal/app/repository"
	"strconv"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{Repository: r}
}

func (h *Handler) GetOrders(ctx *gin.Context) {
	searchQuery := ctx.Query("query") // получение строки поиска
	var materials []repository.Material
	var err error

	if searchQuery == "" {
		materials, err = h.Repository.GetMaterials() // все материалы
	} else {
		materials, err = h.Repository.GetMaterialsByTitle(searchQuery) // фильтр по названию
	}

	if err != nil {
		logrus.Error(err)
	}

	// Передаём в шаблон index.html: список материалов и строку поиска
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"materials": materials,
		"query":     searchQuery,
	})
}

func (h *Handler) GetOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	material, err := h.Repository.GetMaterial(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "order.html", gin.H{
		"material": material,
	})
}
