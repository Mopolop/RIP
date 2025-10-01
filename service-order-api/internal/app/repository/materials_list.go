package repository

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"unicode"

	"db-integration/internal/app/ds"
	"errors"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func (r *Repository) GetMaterials() ([]ds.Material, error) {
	var materials []ds.Material
	result := r.db.Order("id ASC").Find(&materials)
	if result.Error != nil {
		return nil, result.Error
	}
	return materials, nil
}

func (r *Repository) GetMaterialsByTitle(title string) ([]ds.Material, error) {
	var materials []ds.Material
	result := r.db.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(title)+"%").Find(&materials)
	if result.Error != nil {
		return nil, result.Error
	}
	return materials, nil
}

// Получаем черновой заказ пользователя
func (r *Repository) GetDraftOrder(userID int) (*ds.MaterialOrder, error) {
	var order ds.MaterialOrder
	err := r.db.Where("creator_id = ? AND request_status = ?", userID, "черновик").First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // черновик отсутствует
		}
		return nil, err
	}
	return &order, nil
}

// Создаём новый черновой заказ
func (r *Repository) CreateDraftOrder(userID int) (*ds.MaterialOrder, error) {
	order := ds.MaterialOrder{
		CreatorID:     userID,
		ModeratorID:   userID, // можно назначить себя модератором, либо 0/NULL
		RequestStatus: "черновик",
	}
	if err := r.db.Create(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// Добавляем материал в заказ, если его там нет
func (r *Repository) AddMaterialToOrder(orderID int, materialID int) error {
	var count int64
	err := r.db.Model(&ds.MaterialMaterialOrder{}).Where("material_order_id = ? AND material_id = ?", orderID, materialID).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // материал уже добавлен
	}

	item := ds.MaterialMaterialOrder{
		MaterialID:      materialID,
		MaterialOrderID: orderID,
	}
	return r.db.Create(&item).Error
}

// Получаем количество материалов в заказе
func (r *Repository) GetOrderMaterialsCount(orderID int) (int64, error) {
	var count int64
	err := r.db.Model(&ds.MaterialMaterialOrder{}).Where("material_order_id = ?", orderID).Count(&count).Error
	return count, err
}

// SetOrderStatus обновляет статус заказа по его ID
func (r *Repository) SetOrderStatus(orderID int, status string) error {
	return r.db.Model(&ds.MaterialOrder{}).Where("id = ?", orderID).Update("request_status", status).Error
}

// Получаем один материал по ID
func (r *Repository) GetMaterialByID(id int) (*ds.Material, error) {
	var material ds.Material
	err := r.db.First(&material, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &material, nil
}

// Получаем список материалов с опциональной фильтрацией по названию
func (r *Repository) GetMaterialsFiltered(title string) ([]ds.Material, error) {
	var materials []ds.Material
	query := r.db.Model(&ds.Material{})
	if title != "" {
		query = query.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(title)+"%")
	}
	err := query.Find(&materials).Error
	if err != nil {
		return nil, err
	}
	return materials, nil
}

// Создаём новый материал
func (r *Repository) CreateMaterial(material *ds.Material) error {
	return r.db.Create(material).Error
}

// Обновляем материал по ID
func (r *Repository) UpdateMaterial(id int, updated *ds.Material) error {
	var material ds.Material
	if err := r.db.First(&material, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // материал не найден
		}
		return err
	}

	// Обновляем все поля
	return r.db.Model(&material).Updates(updated).Error
}

// Удаляем материал по ID
func (r *Repository) DeleteMaterial(id int) error {
	var material ds.Material
	err := r.db.First(&material, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // материал не найден, можно считать успешным
		}
		return err
	}
	return r.db.Delete(&material).Error
}

// UploadMaterialImage загружает новое изображение материала в MinIO, удаляет старое, обновляет image_url
func (r *Repository) UploadMaterialImage(id int, fileHeader *multipart.FileHeader) error {
	// Получаем материал из БД
	var material ds.Material
	if err := r.db.First(&material, id).Error; err != nil {
		return err
	}

	// Удаляем старое изображение из MinIO, если есть
	if material.Image != "" {
		parts := strings.Split(material.Image, "/")
		objectName := parts[len(parts)-1]
		_ = r.minioClient.RemoveObject(context.Background(), r.bucketName, objectName, minio.RemoveObjectOptions{})
	}
	// Открываем новый файл
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	// Расширение файла
	ext := filepath.Ext(fileHeader.Filename)
	base := strings.TrimSuffix(fileHeader.Filename, ext)

	// Переводим в латиницу
	latinBase := toLatin(base)

	// Генерация имени файла
	objectName := fmt.Sprintf("material-%s%s", latinBase, ext)

	// Загружаем в MinIO
	_, err = r.minioClient.PutObject(
		context.Background(),
		r.bucketName,
		objectName,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{ContentType: fileHeader.Header.Get("Content-Type")},
	)
	if err != nil {
		return err
	}

	// Формируем URL
	imageURL := fmt.Sprintf("http://%s/%s/%s", r.minioClient.EndpointURL().Host, r.bucketName, objectName)

	// Обновляем поле ImageURL в БД
	return r.db.Model(&ds.Material{}).Where("id = ?", id).Update("image", imageURL).Error
}

// toLatin переводит строку в латиницу, оставляет только ASCII буквы и цифры
func toLatin(s string) string {
	var out strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) && r <= unicode.MaxASCII {
			out.WriteRune(unicode.ToLower(r))
		} else if unicode.IsDigit(r) {
			out.WriteRune(r)
		}
	}
	return out.String()
}
