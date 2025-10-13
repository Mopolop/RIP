package handler

import (
	"html/template"
	"service-order-api/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// RegisterHandler регистрируем маршруты
func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/", h.GetMaterials)
	router.GET("/detailed_material/:id", h.GetMaterial)
	router.GET("/materials_order/:id", h.GetMaterialsOrder)
	router.GET("/api/material/:id", h.GetMaterialAPI)
	router.GET("/api/materials", h.GetMaterialsAPI)
	router.GET("/api/orders/draft/cart", h.GetDraftCartAPI)
	router.GET("/api/orders", h.GetOrdersAPI)
	router.GET("/api/orders/:id", h.GetOrderWithMaterialsAPI)
	router.GET("/api/users/:id", h.GetUserAPI)

	router.POST("/orders/draft/add/:id", h.AddMaterialToDraftOrder)
	router.POST("/orders/delete/:id", h.DeleteMaterialsOrder)
	router.POST("/api/material", h.CreateMaterialAPI)
	router.POST("/api/orders/draft/add/:id", h.AddMaterialToDraftOrderAPI)
	router.POST("/api/material/:id/image", h.UploadMaterialImage)
	router.POST("/api/material/:id/delete", h.DeleteMaterialLogicalAPI)
	router.POST("/api/orders/delete/:id", h.DeleteMaterialsOrderAPI)
	router.POST("/api/users/register", h.RegisterUserAPI)
	router.POST("/api/users/login", h.LoginUserAPI)
	router.POST("/api/users/logout", h.LogoutUserAPI)

	router.PUT("/api/material/:id", h.UpdateMaterialAPI)
	router.PUT("/api/orders/:id", h.UpdateMaterialOrderAPI)
	router.PUT("/api/orders/:id/form", h.FormMaterialOrderAPI)
	router.PUT("/api/orders/:id/complete", h.CompleteOrRejectOrderAPI)
	router.PUT("/api/orders/materials/:order_id/:material_id/wall_length", h.UpdateWallLengthAPI)
	router.PUT("/api/users/:id", h.UpdateUserAPI)

	router.DELETE("/api/orders/:order_id/material/:material_id", h.DeleteMaterialFromOrderAPI)

}

// RegisterStatic регистрирует статические файлы и шаблоны
func (h *Handler) RegisterStatic(router *gin.Engine) {
	// Функция для склонения русских слов (1 — одна, 2-4 — несколько, 5+ — many)
	plural := func(n int, one, few, many string) string {
		nn := n % 100
		if nn >= 11 && nn <= 19 {
			return many
		}
		i := nn % 10
		if i == 1 {
			return one
		}
		if i >= 2 && i <= 4 {
			return few
		}
		return many
	}

	tmpl := template.Must(template.New("templates").Funcs(template.FuncMap{
		"plural": plural,
	}).ParseGlob("templates/*"))

	router.SetHTMLTemplate(tmpl)
	router.Static("/static", "./resources")
}

// errorHandler для удобного вывода ошибок
func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
