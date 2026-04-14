package postgres

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
	"hexagonalapp/internal/modules/user/domain"
	"hexagonalapp/internal/platform/helpers/stnchelper"
	"strconv"
	"time"
)

type UserModel struct {
	ID        uint      `gorm:"primaryKey;auto_increment;column:id"`
	Name      string    `gorm:"column:name;not null"`
	Email     string    `gorm:"column:email;not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (UserModel) TableName() string { return "user" }

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(&UserModel{})
}

func (r *Repository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}


func (r *Repository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.User{}, id).Error
}

func (r *Repository) ExistsByEmail(ctx context.Context, email string, excludeID uint) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&domain.User{}).Where("email = ?", email)
	if excludeID > 0 {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}



func (r *Repository) GetByID(ctx context.Context, id string) (domain.User, error) {
	var m UserModel
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, err
		}
		return domain.User{}, err
	}
	return toDomain(m), nil
}



func (r *Repository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}


func (r *Repository) List(ctx context.Context) ([]domain.User, error) {
	var rows []domain.User
	if err := r.db.WithContext(ctx).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
} // Order fonksiyonu: Sıralama kriterlerini ekler

/*
func (r *Repository) ListDataTable(c fiber.Ctx, ctx context.Context) ([]domain.User, int64, int64, error) {
    var rows []domain.User

    var totalRecords int64
    // Toplam kayıt sayısı (filtresiz)
    r.db.Model(&domain.User{}).Count(&totalRecords)

    startStr := c.Query("start", "0")

    start, _ := strconv.Atoi(startStr)

    lengthStr := c.Query("length", "10")

    length, _ := strconv.Atoi(lengthStr)

    orderIdxStr := c.Query("order[0][column]", "0")

    orderIdx, _ := strconv.Atoi(orderIdxStr)

    search := c.Query("search[value]")

    orderDir := c.Query("order[0][dir]", "desc")
    // for sorter and search columns
    columns := []string{"id", "name", "email"}

    query := r.db.WithContext(ctx).Model(&domain.User{})

    query = query.Scopes(stnchelper.SearchScope(search, columns))

    var filteredRecords int64

    query.Count(&filteredRecords)

    query = stnchelper.ApplyOrder(query, columns[orderIdx], orderDir, "id desc")

    query = stnchelper.ApplyPagination(query, start, length)

    if err := query.Find(&rows).Error; err != nil {
        return nil, 0, 0, err
    }

    return rows, totalRecords, filteredRecords, nil
}

*/

func (r *Repository) ListDataTable(c fiber.Ctx, ctx context.Context) ([]domain.User, int64, int64, error) {
	var rows []domain.User

	// 1. Girdi ayrıştırma ve hata yönetimi
	start, _ := strconv.Atoi(c.Query("start", "0"))
	length, _ := strconv.Atoi(c.Query("length", "10"))
	orderIdx, _ := strconv.Atoi(c.Query("order[0][column]", "0"))
	search := c.Query("search[value]")
	orderDir := c.Query("order[0][dir]", "desc")

	// Sütun tanımlarını güvenli yönetin
	columns := []string{"id", "name", "email"}
	if orderIdx < 0 || orderIdx >= len(columns) {
		orderIdx = 0
	}

	// 2. Toplam kayıt sayısı (Bunu veritabanı yükünü azaltmak için
	// sadece ilk çağrıda cache'leyebilir veya ayrı bir metot yapabilirsiniz)
	var totalRecords int64
	if err := r.db.Model(&domain.User{}).Count(&totalRecords).Error; err != nil {
		return nil, 0, 0, err
	}

	// 3. Sorgu oluşturma (Scoped)
	query := r.db.WithContext(ctx).Model(&domain.User{})

	// Arama filtreleme
	if search != "" {
		query = query.Scopes(stnchelper.SearchScope(search, columns))
	}

	// Filtrelenmiş kayıt sayısı
	var filteredRecords int64
	if err := query.Count(&filteredRecords).Error; err != nil {
		return nil, 0, 0, err
	}

	// 4. Sıralama ve Sayfalama
	query = stnchelper.ApplyOrder(query, columns[orderIdx], orderDir, "id desc")
	query = stnchelper.ApplyPagination(query, start, length)

	// 5. Veriyi çekme
	if err := query.Find(&rows).Error; err != nil {
		return nil, 0, 0, err
	}

	return rows, totalRecords, filteredRecords, nil
}

func (r *Repository) ListPagination(ctx context.Context, limit, offset int) ([]domain.User /*int64,*/, error) {
	var users []domain.User

	if err := r.db.WithContext(ctx).Order("id desc").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) Count(total *int64) {
	var users domain.User
	var count int64
	r.db.Model(users).Count(&count)
	*total = count
}




func toModel(u domain.User) UserModel {
	return UserModel{ID: u.ID, Name: u.Name, Email: u.Email, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt}
}

func toDomain(m UserModel) domain.User {
	return domain.User{ID: m.ID, Name: m.Name, Email: m.Email, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}
