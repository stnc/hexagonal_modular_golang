package stnchelper



import (
	"fmt"
	"gorm.io/gorm"
)

// Search fonksiyonu: Arama kriterlerini ekler
func ApplyOrder(db *gorm.DB, orderColumn string, orderDir string, defaultOrder string) *gorm.DB {
	if orderColumn == "" {
		return db.Order(defaultOrder)
	}
	return db.Order(fmt.Sprintf("%s %s", orderColumn, orderDir))
}

// ApplySearch: Arama terimini verilen kolonlarda (OR mantığıyla) arar.
func SearchScope(search string, columns []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if search == "" {
			return db
		}

		s := "%" + search + "%"
		// İlk koşulu Where ile başlat
		db = db.Where(fmt.Sprintf("%s LIKE ?", columns[0]), s)

		// Kalanları Or ile ekle
		for i := 1; i < len(columns); i++ {
			db = db.Or(fmt.Sprintf("%s LIKE ?", columns[i]), s)
		}
		return db
	}
}

// Pagination fonksiyonu: Offset ve Limit ekler
func ApplyPagination(db *gorm.DB, start int, length int) *gorm.DB {
	return db.Offset(start).Limit(length)
}