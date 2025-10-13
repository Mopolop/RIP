package handler

import (
	"net/http"
	"strconv"

	"db-integration/internal/app/ds"
	"github.com/gin-gonic/gin"
)

// Получение конкретного материала по ID
func (h *Handler) GetMaterial(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var material ds.Material
	material, err = h.Repository.GetMaterial(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	ctx.HTML(http.StatusOK, "detailed_material.html", gin.H{
		"material": material,
	})
}
