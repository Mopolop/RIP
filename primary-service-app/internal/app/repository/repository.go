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
	ID           int
	Title        string
	Description  string
	Image        string
	Consumption  string
	Count        string
	MainMaterial string
	CountPerM2   string
	CountPerM3   string
	NetWeight    string
	LengthMM     string
	HeightMM     string
	WidthMM      string
	Country      string
}

type Order struct {
	ID        int
	Materials []Material
}

func (r *Repository) GetMaterials() ([]Material, error) {
	materials := []Material{
		{
			ID:          1,
			Title:       "Кирпич строительный рядовой полнотелый красный",
			Description: "250×120×65 мм",
			Consumption: "0,03 м³ на 1 м³ кладки",
			Image:       "http://localhost:9000/materials/brick-red.png",
			Count:       "512 шт./м³",

			MainMaterial: "Керамика",
			CountPerM2:   "62",
			CountPerM3:   "512",
			NetWeight:    "3.67",
			LengthMM:     "250",
			HeightMM:     "65",
			WidthMM:      "120",
			Country:      "Россия",
		},
		{
			ID:          2,
			Title:       "Блок газобетонный",
			Description: "600×250×100 мм",
			Consumption: "0,05 м³ на 1 м³ кладки",
			Image:       "http://localhost:9000/materials/block-gaz.png",
			Count:       "67 шт./м³",

			MainMaterial: "Газобетон",
			CountPerM2:   "12",
			CountPerM3:   "67",
			NetWeight:    "2.5",
			LengthMM:     "600",
			HeightMM:     "100",
			WidthMM:      "250",
			Country:      "Россия",
		},
		{
			ID:          3,
			Title:       "Блок керамический",
			Description: "250×219×510 мм",
			Consumption: "0,04 м³ на 1 м³ кладки",
			Image:       "http://localhost:9000/materials/block-keram.png",
			Count:       "36 шт./м³",

			MainMaterial: "Керамика",
			CountPerM2:   "15",
			CountPerM3:   "36",
			NetWeight:    "4.2",
			LengthMM:     "250",
			HeightMM:     "510",
			WidthMM:      "219",
			Country:      "Россия",
		},
		{
			ID:          4,
			Title:       "Кирпич строительный силикатный полнотелый белый",
			Description: "250×120×65 мм",
			Consumption: "0,5 м³ на 1 м³ кладки",
			Image:       "http://localhost:9000/materials/brick-white.png",
			Count:       "512 шт./м³",

			MainMaterial: "Силикат",
			CountPerM2:   "60",
			CountPerM3:   "512",
			NetWeight:    "3.9",
			LengthMM:     "250",
			HeightMM:     "65",
			WidthMM:      "120",
			Country:      "Россия",
		},
		{
			ID:          5,
			Title:       "Кирпич шамотный огнеупорный",
			Description: "230×114×65 мм",
			Consumption: "0,035 м³ на 1 м³ кладки",
			Image:       "http://localhost:9000/materials/fireclay_brick.png",
			Count:       "587 шт./м³",

			MainMaterial: "Шамот",
			CountPerM2:   "55",
			CountPerM3:   "587",
			NetWeight:    "3.5",
			LengthMM:     "230",
			HeightMM:     "65",
			WidthMM:      "114",
			Country:      "Россия",
		},
		{
			ID:          6,
			Title:       "Кирпич фасадный клинкерный пустотелый красный гладкий",
			Description: "250×120×65 мм",
			Consumption: "0,03 м³ на 1 м³ кладки",
			Image:       "http://localhost:9000/materials/facade_brick.png",
			Count:       "513 шт./м³",

			MainMaterial: "Клинкер",
			CountPerM2:   "63",
			CountPerM3:   "513",
			NetWeight:    "3.8",
			LengthMM:     "250",
			HeightMM:     "65",
			WidthMM:      "120",
			Country:      "Германия",
		},
		{
			ID:          7,
			Title:       "Блок керамзитобетонный пустотелый Стандарт Керамзит",
			Description: "390×190×188 мм",
			Consumption: "0,02 м³ на 1 м³ кладки",
			Image:       "http://localhost:9000/materials/block-keramb.png",
			Count:       "72 шт./м³",

			MainMaterial: "Керамзитобетон",
			CountPerM2:   "14",
			CountPerM3:   "72",
			NetWeight:    "5.1",
			LengthMM:     "390",
			HeightMM:     "188",
			WidthMM:      "190",
			Country:      "Россия",
		},
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

func (r *Repository) GetOrder() (Order, error) {
	materials, _ := r.GetMaterials()

	order := Order{
		ID:        1,
		Materials: []Material{materials[0], materials[1]}, // единственная заявка
	}

	return order, nil
}
