package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Material struct {
	ID          int
	Title       string // наименование материала
	Description string // описание или размеры
	Consumption string // расход раствора на 1 м³ кладки
	Image       string // название файла изображения (ключ в Minio)
	Count       string // количество шт./м³
}

func (r *Repository) GetMaterials() ([]Material, error) {
	materials := []Material{
		{
			ID:          1,
			Title:       "Кирпич строительный рядовой полнотелый красный",
			Description: "250×120×65 мм",
			Consumption: "0,03 м³ на 1 м³ кладки",
			Image:       "brick-red.png",
			Count:       "512 шт./м³",
		},
		{
			ID:          2,
			Title:       "Блок газобетонный",
			Description: "600×250×100 мм",
			Consumption: "0,05 м³ на 1 м³ кладки",
			Image:       "block-gaz.png",
			Count:       "67 шт./м³",
		},
		{
			ID:          3,
			Title:       "Блок керамический",
			Description: "250×219×510 мм",
			Consumption: "0,04 м³ на 1 м³ кладки",
			Image:       "block-keram.png",
			Count:       "36 шт./м³",
		},
		{
			ID:          4,
			Title:       "Кирпич строительный силикатный полнотелый белый",
			Description: "250×120×65 мм",
			Consumption: "0,5 м³ на 1 м³ кладки",
			Image:       "brick-white.png",
			Count:       "512 шт./м³",
		},
		{
			ID:          5,
			Title:       "Кирпич шамотный огнеупорный",
			Description: "230×114×65 мм",
			Consumption: "0,035 м³ на 1 м³ кладки",
			Image:       "fireclay_brick.png",
			Count:       "587 шт./м³",
		},
		{
			ID:          6,
			Title:       "Кирпич фасадный клинкерный пустотелый красный гладкий",
			Description: "250×120×65 мм",
			Consumption: "0,03 м³ на 1 м³ кладки",
			Image:       "facade_brick.png",
			Count:       "513 шт./м³",
		},
		{
			ID:          7,
			Title:       "Блок керамзитобетонный пустотелый Стандарт Керамзит",
			Description: "390×190×188 мм",
			Consumption: "0,02 м³ на 1 м³ кладки",
			Image:       "block-keram	b.png",
			Count:       "72 шт./м³",
		},
	}

	if len(materials) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}
	return materials, nil
}

func (r *Repository) GetMaterial(id int) (Material, error) {
	materials, err := r.GetMaterials()
	if err != nil {
		return Material{}, err
	}

	for _, m := range materials {
		if m.ID == id {
			return m, nil
		}
	}
	return Material{}, fmt.Errorf("материал не найден")
}

func (r *Repository) GetMaterialsByTitle(title string) ([]Material, error) {
	materials, err := r.GetMaterials()
	if err != nil {
		return []Material{}, err
	}

	var result []Material
	for _, m := range materials {
		if strings.Contains(strings.ToLower(m.Title), strings.ToLower(title)) {
			result = append(result, m)
		}
	}

	return result, nil
}
