package handler

import (
	"net/http"
	"strconv"

	"db-integration/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

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

	// Для примера используем userID = 1
	userID := 1

	// Ищем черновой заказ пользователя
	order, err := h.Repository.GetDraftOrder(userID)
	var orderCount int64 = 0
	if err != nil {
		logrus.Warn("Не удалось получить черновой заказ: ", err)
	} else if order != nil {
		// Получаем количество материалов в черновике
		orderCount, err = h.Repository.GetOrderMaterialsCount(order.ID)
		if err != nil {
			logrus.Warn("Не удалось получить количество материалов в черновике: ", err)
			orderCount = 0
		}
	}

	ctx.HTML(http.StatusOK, "materials_list.html", gin.H{
		"materials":  materials,
		"query":      searchQuery,
		"orderCount": orderCount,
		// передаём ID чернового заказа в шаблон (0 если заказа нет)
		"orderID": func() int {
			if order != nil {
				return order.ID
			}
			return 0
		}(),
	})
}

// POST /orders/draft/add/:id
func (h *Handler) AddMaterialToDraftOrder(ctx *gin.Context) {
	materialIDStr := ctx.Param("id")
	materialID, err := strconv.Atoi(materialIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Для примера используем userID = 1
	userID := 1

	// Получаем черновой заказ
	order, err := h.Repository.GetDraftOrder(userID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	if order == nil {
		order, err = h.Repository.CreateDraftOrder(userID)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	// Добавляем материал
	if err := h.Repository.AddMaterialToOrder(order.ID, materialID); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Сохраняем count в сессии или просто редиректим на текущую страницу
	ctx.Redirect(http.StatusSeeOther, ctx.Request.Referer())
}

// Получение заказа по ID
func (h *Handler) GetMaterialsOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Если id == 0 — отвечаем унифицированным сообщением, не перенаправляя/не показывая внутреннюю ошибку
	if id == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":      "error",
			"description": "заказ не найден или удален",
		})
		return
	}

	order, mmos, err := h.Repository.GetOrderByID(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	// Если статус заказа не черновик, считаем, что заказа нет/он удалён
	if order.RequestStatus != "черновик" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":      "error",
			"description": "заказ не найден или удален",
		})
		return
	}

	ctx.HTML(http.StatusOK, "materials_order.html", gin.H{
		"order":     order,
		"materials": mmos, // вот сюда прокидываем список материалов
	})
}

// GET /api/materials/:id
func (h *Handler) GetMaterialAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	material, err := h.Repository.GetMaterialByID(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if material == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":      "error",
			"description": "материал не найден",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"material": material,
	})
}

// GET /api/materials?title=<название>
func (h *Handler) GetMaterialsAPI(ctx *gin.Context) {
	title := ctx.Query("title")

	materials, err := h.Repository.GetMaterialsFiltered(title)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"materials": materials,
	})
}

// POST /api/material
func (h *Handler) CreateMaterialAPI(ctx *gin.Context) {
	var input ds.Material

	// Привязываем JSON из запроса
	if err := ctx.ShouldBindJSON(&input); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Создаём материал через репозиторий
	if err := h.Repository.CreateMaterial(&input); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":   "success",
		"material": input,
	})
}

// PUT /api/material/:id
func (h *Handler) UpdateMaterialAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var input ds.Material
	if err := ctx.ShouldBindJSON(&input); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.UpdateMaterial(id, &input)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"material": input,
	})
}

// DELETE /api/material/:id
func (h *Handler) DeleteMaterialAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteMaterial(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "материал успешно удалён",
	})
}

// POST /api/orders/draft/add/:id
func (h *Handler) AddMaterialToDraftOrderAPI(ctx *gin.Context) {
	// Получаем ID материала из URL
	materialIDStr := ctx.Param("id")
	materialID, err := strconv.Atoi(materialIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Для примера используем userID = 1
	userID := 1

	// Получаем черновой заказ пользователя
	order, err := h.Repository.GetDraftOrder(userID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Если чернового заказа нет — создаём новый
	if order == nil {
		order, err = h.Repository.CreateDraftOrder(userID)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	// Добавляем материал в заказ
	if err := h.Repository.AddMaterialToOrder(order.ID, materialID); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Получаем новое количество материалов в заказе
	count, _ := h.Repository.GetOrderMaterialsCount(order.ID)

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "материал добавлен в черновой заказ",
		"orderID":   order.ID,
		"itemCount": count,
	})
}

// Загрузить/заменить изображение материала
func (h *Handler) UploadMaterialImage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid material id"})
		return
	}

	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no image file"})
		return
	}

	if err := h.Repository.UploadMaterialImage(id, fileHeader); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "image uploaded"})
}
