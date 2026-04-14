package app

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"gopkg.in/go-playground/validator.v9"

	"hexagonalapp/internal/modules/user/domain"
	"hexagonalapp/internal/modules/user/ports"

	"time"
	// "hexagonalapp/internal/platform/id"
)

type Service struct {
	repo     ports.Repository
	cache    ports.Cache
	audit    ports.Audit
	validate *validator.Validate
}

func New(repo ports.Repository, cache ports.Cache, audit ports.Audit) *Service {
	return &Service{repo: repo, cache: cache, audit: audit}
}

func (s *Service) CreateUser(ctx context.Context, input domain.CreateUserInput) (*domain.User, fiber.Map, error) {
	if err := domain.ValidateInput(input); err != nil {
		return &domain.User{}, fiber.Map{}, err
	}

	exists, err := s.repo.ExistsByEmail(ctx, input.Email, 0)
	if err != nil {
		return nil, fiber.Map{}, err
	}
	if exists {
		return nil, fiber.Map{"ErrorEmail": "email already used"}, domain.ErrEmailAlreadyUsed
	}

	now := time.Now().UTC()
	user := &domain.User{
		ID:        input.ID,
		Name:      input.Name,
		Email:     input.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	user2 := domain.User{Name: input.Name, Email: input.Email}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fiber.Map{}, err
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, user2)
	}
	if s.audit != nil {
		_ = s.audit.Record(ctx, "user.created", user2)
	}

	return user, fiber.Map{}, nil
}

/* alt taraf orginal dir */

func (s *Service) GetUser(ctx context.Context, userID string) (domain.User, error) {
	if user, ok, err := s.cache.Get(ctx, userID); err == nil && ok {
		return user, nil
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}
	_ = s.cache.Set(ctx, user)
	return user, nil
}

func (s *Service) ListUsers(ctx context.Context) ([]domain.User, error) {
	return s.repo.List(ctx)
}

func (s *Service) ListDataTable(c fiber.Ctx, ctx context.Context) ([]domain.User, int64, int64, error) {
	return s.repo.ListDataTable(c, ctx)
}

func (s *Service) ListUsersPagination(ctx context.Context, limit, offset int) ([]domain.User, error) {
	return s.repo.ListPagination(ctx, limit, offset)
}

func (f *Service) Count(postTotalCount *int64) {
	f.repo.Count(postTotalCount)
}

/*
//buradaki alttaki kisim ismini json daki verisinden okumak icin yazildi belki ilerde eklenebilir
//validate.RegisterTagNameFunc(func(fld reflect.StructField) string {

func (v *Service) Validate() map[string]string {
	var (
		validate *validator.Validate
		uni      *ut.UniversalTranslator
	)

	tr := en.New()
	uni = ut.New(tr, tr)
	trans, _ := uni.GetTranslator("en")

	validate = validator.New()

	// tr_translations.RegisterDefaultTranslations(validate, trans)
	en_translations.RegisterDefaultTranslations(validate, trans)

	errorLog := make(map[string]string)

	// JSON tag'ini field name olarak kullan
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			return fld.Name
		}
		return name
	})

	err := validate.Struct(v)

	if err != nil {
		errs := err.(validator.ValidationErrors)

		for _, e := range errs {
			// can translate each error one at a time.
			lng := strings.Replace(e.Translate(trans), e.Field(), "This", 1)
			errorLog[e.Field()+"_error"] = e.Translate(trans)
			// errorLog[e.Field()] = e.Translate(trans)
			errorLog[e.Field()] = lng
			errorLog[e.Field()+"_valid"] = "is-invalid"
		}
	}
	// fmt.Println(errorLog)
	return errorLog
}
*/
