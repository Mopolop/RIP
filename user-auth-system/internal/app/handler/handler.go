package handler

import (
	"html/template"
	"user-auth-system/internal/app/config"
	redisclient "user-auth-system/internal/app/redis"
	"user-auth-system/internal/app/repository"
	"user-auth-system/internal/app/role"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
	Config     *config.Config
	Redis      *redisclient.Client
}

func NewHandler(r *repository.Repository, cfg *config.Config, redis *redisclient.Client) *Handler {
	return &Handler{
		Repository: r,
		Config:     cfg,
		Redis:      redis,
	}
}

// RegisterHandler регистрируем маршруты
func (h *Handler) RegisterHandler(router *gin.Engine) {
	// ---------------------------
	// Публичные маршруты
	// ---------------------------
	router.GET("/api/materials", h.GetMaterialsAPI)
	router.GET("/api/materials/:id", h.GetMaterialAPI)
	router.POST("/sign_up", h.Register)
	router.POST("/api/users/login", h.LoginUserAPI)
	router.POST("/api/users/logout", h.LogoutUserAPI)

	// ---------------------------
	// Защищённые маршруты для всех авторизованных (User + Admin)
	// ---------------------------
	auth := router.Group("/api")
	auth.Use(h.WithAuthCheck(role.User, role.Admin))
	{
		// Пользователи
		auth.GET("/users/:id", h.GetUserAPI)
		auth.PUT("/users/:id", h.UpdateUserAPI)

		// Материалы и заказы (все методы кроме админских)
		auth.GET("/materials_order/:id", h.GetMaterialsOrder)
		auth.GET("/api/material/:id", h.GetMaterialAPI)
		auth.POST("/orders/draft/add/:id", h.AddMaterialToDraftOrderAPI)
		auth.GET("/orders", h.GetOrdersAPI)
		auth.GET("/orders/:id", h.GetOrderWithMaterialsAPI)
		auth.GET("/orders/draft/cart", h.GetDraftCartAPI)
		auth.PUT("/orders/:id", h.UpdateMaterialOrderAPI)
		auth.PUT("/orders/:id/form", h.FormMaterialOrderAPI)
		auth.PUT("/orders/materials/:order_id/:material_id/wall_length", h.UpdateWallLengthAPI)
		auth.POST("/orders/delete/:id", h.DeleteMaterialsOrderAPI)
		auth.DELETE("/orders/:order_id/material/:material_id", h.DeleteMaterialFromOrderAPI)
	}

	// ---------------------------
	// Только Admin
	// ---------------------------
	admin := router.Group("/api")
	admin.Use(h.WithAuthCheck(role.Admin))
	{
		admin.POST("/material", h.CreateMaterialAPI)
		admin.PUT("/material/:id", h.UpdateMaterialAPI)
		admin.POST("/material/:id/image", h.UploadMaterialImage)
		admin.POST("/material/:id/delete", h.DeleteMaterialLogicalAPI)
		admin.PUT("/orders/:id/complete", h.CompleteOrRejectOrderAPI)
	}
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
