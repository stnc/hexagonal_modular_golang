package http

import (
	"context"
	"fmt"

	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/session"
	// validator "gopkg.in/go-playground/validator.v9"
	"hexagonalapp/internal/modules/user/adapters/inbound/http/middleware"
	"hexagonalapp/internal/modules/user/adapters/inbound/http/viewmodels"
	"hexagonalapp/internal/modules/user/app"
	"math"
"hexagonalapp/internal/modules/user/domain"
)

/*
getProduct

getProductByID

getAllProducts
*/

type Pagination struct {
	TotalItems  int64
	TotalPages  int
	CurrentPage int
	NextPage    int
	PrevPage    int
	Pages       []int
	HasNext     bool
	HasPrev     bool
	PageList    []int
}

type WebHandler struct {
	service *app.Service
	store   *session.Store
}

func NewWEB(service *app.Service) *WebHandler {
	return &WebHandler{service: service}
}

func (h *WebHandler) Register(r fiber.Router) {
	r.Get("/user/create", h.Create) // new  // #TODO  REadme ekle ve sidebar ekle
	r.Post("/user/store", h.Store)  //create
	r.Get("/users/:id", h.GetUser)

	r.Get("/list/list_users_with_pagination", h.ListUsersWithPagination)
	r.Get("/list/normal_users", h.ListAllUsers)
	r.Get("/list/datatable", h.ListUsersDatatable)
	r.Get("/ajaxdatatable", h.ListUsersDatatableAjax)
}

func (h *WebHandler) Create(c fiber.Ctx) error {
	flash := middleware.ConsumeFlash(h.store, c)
	return c.Render("users/create", fiber.Map{
		"Title":        "New User Create",
		"User":         viewmodels.UserForm{},
		"FormAction":   "/web/user/store",
		"SubmitLabel":  "Create user",
		"FlashSuccess": flash.Success,
		"FlashError":   flash.Error,
		"CsrfToken":    csrf.TokenFromContext(c),
	})
}

func (h *WebHandler) Store(c fiber.Ctx) error {

	input, validationData, err := domain.BindInput(c)
	if err != nil {
		return domain.RenderCreateWithErrors(c, input, validationData, err)
	}

	user, validationData, err := h.service.CreateUser(c.Context(), input)
	if err != nil {
		return domain.RenderCreateWithErrors(c, input, validationData, err)
	}

	middleware.SetFlash(h.store, c, "success", fmt.Sprintf("%s created successfully", user.Name))
	return c.Redirect().To("/web/user/edit/"+strconv.FormatUint(uint64( user.ID), 10))

	// input := app.CreateUserInput{
	// 	Name:  c.FormValue("name"),
	// 	Email: c.FormValue("email"),
	// }
	// 	fmt.Println(input)
	// 	user, err := h.service.CreateUser(c.Context(), input)
	// fmt.Println(user)
	// fmt.Println(err)

	// if err != nil {
	// 	return h.handleFormError(c, "users/create", "/users/create", viewmodels.UserForm{Name: input.Name, Email: input.Email}, err)
	// }
	// if err := middleware.SetFlash(h.store, c, fmt.Sprintf("%s added succesfull", user.Name), ""); err != nil {
	// 	return err
	// }
	// return c.Redirect().To("/web/user/create")
}




func (h *WebHandler) GetUser(c fiber.Ctx) error {
	user, err := h.service.GetUser(context.Background(), c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

func parsePage(v string) int {
	page, err := strconv.Atoi(v)
	if err != nil || page < 1 {
		return 1
	}
	return page
}

func (h *WebHandler) ListUsersWithPagination(c fiber.Ctx) error {

	page := parsePage(c.Query("page", "1"))

	var totalItems int64
	h.service.Count(&totalItems)

	if page < 1 {
		page = 1
	}

	limit := 3 // every page  10 item
	offset := (page - 1) * limit

	users, err := h.service.ListUsersPagination(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	// HTML'de butonları oluşturmak için sayfa numaralarını içeren bir slice
	var pages []int
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}

	paging := Pagination{
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: page,
		Pages:       pages, // 1, 2, 3... şeklindeki dizi
		HasNext:     page < totalPages,
		HasPrev:     page > 1,
		NextPage:    page + 1,
		PrevPage:    page - 1,
	}

	flash := middleware.ConsumeFlash(h.store, c)
	return c.Render("users/ListUsersWithPagination", fiber.Map{
		"Title":        "Users",
		"Users":        users,
		"Pagination":   paging,
		"TotalRecord":  totalItems,
		"FlashSuccess": flash.Success,
		"FlashError":   flash.Error,
		"CsrfToken":    csrf.TokenFromContext(c),
	})
}

func (h *WebHandler) ListAllUsers(c fiber.Ctx) error {

	users, err := h.service.ListUsers(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	flash := middleware.ConsumeFlash(h.store, c)
	return c.Render("users/listAllUsers", fiber.Map{
		"Title":        "Users",
		"Users":        users,
		"FlashSuccess": flash.Success,
		"FlashError":   flash.Error,
		"CsrfToken":    csrf.TokenFromContext(c),
	})
}
func (h *WebHandler) ListUsersDatatable(c fiber.Ctx) error {
	flash := middleware.ConsumeFlash(h.store, c)
	return c.Render("users/datatable", fiber.Map{
		"Title": "Users",
		// "Users":        users,
		"FlashSuccess": flash.Success,
		"FlashError":   flash.Error,
		"CsrfToken":    csrf.TokenFromContext(c),
	})
}

func (h *WebHandler) ListUsersDatatableAjax(c fiber.Ctx) error {

	users, totalRecords, filteredRecords, err := h.service.ListDataTable(c, context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"draw":            c.Query("draw"),
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            users,
	})
}
