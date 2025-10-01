package ds

// Material представляет строительный материал
type Material struct {
	ID           int    `gorm:"primaryKey"`        // первичный ключ
	Title        string `gorm:"type:varchar(255)"` // название
	Description  string `gorm:"type:varchar(255)"` // описание
	Image        string `gorm:"type:varchar(255)"` // ссылка на изображение
	Consumption  string `gorm:"type:varchar(255)"` // расход
	Count        string `gorm:"type:varchar(255)"` // количество
	MainMaterial string `gorm:"type:varchar(255)"` // основной материал
	CountPerM2   string `gorm:"type:varchar(50)"`  // шт. на м²
	CountPerM3   string `gorm:"type:varchar(50)"`  // шт. на м³
	NetWeight    string `gorm:"type:varchar(50)"`  // масса
	LengthMM     string `gorm:"type:varchar(50)"`  // длина
	HeightMM     string `gorm:"type:varchar(50)"`  // высота
	WidthMM      string `gorm:"type:varchar(50)"`  // ширина
	Country      string `gorm:"type:varchar(100)"` // страна
}
