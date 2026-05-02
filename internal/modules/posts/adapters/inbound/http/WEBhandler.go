package http

import (
	"context"
	conventorLib "hexagonalapp/internal/platform/helpers/stnccollection"
	pagination "hexagonalapp/internal/platform/helpers/stnchelper"
	"math"
	"net/http"
	"strconv"

	"hexagonalapp/internal/modules/posts/app"
	"hexagonalapp/internal/modules/posts/domain"
	"hexagonalapp/internal/platform/adapters/inbound/http/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/session"
)

type WebHandler struct {
	service *app.Service
	store   *session.Store
}


func NewWEB(service *app.Service) *WebHandler {
	return &WebHandler{service: service}
}

func (h *WebHandler) Register(r fiber.Router) {
	r.Get("/post/create", h.Create)
	r.Post("/post/store", h.Store)
	r.Get("/post/:id", h.GetPost)
	r.Get("/post/:id/edit", h.Edit)
	r.Post("/post/:id/update", h.Update)
	r.Post("/post/:id/delete", h.Delete)
	r.Delete("/post/:id/delete", h.Delete)
	r.Get("/posts/list", h.ListAllPosts)
	r.Get("/post/user/:user_id", h.ListByUser)

	// Classic pagination route for posts (mirrors users module)
	r.Get("posts/list/classic_pagination", h.ListPostsWithPagination)
}

func (h *WebHandler) Create(c fiber.Ctx) error {
	return c.Render("posts/create", h.baseData(c, fiber.Map{
		"PageTitle":   "New Post Create",
		"FormAction":  "/web/post/store",
		"SubmitLabel": "Create post",
		"Post":        domain.CreatePostInput{},
	}))
}

func (h *WebHandler) Store(c fiber.Ctx) error {
	input, validationData, err := domain.BindInput(c)
	if err != nil {
		middleware.SetFlash(h.store, c, "", "validation errors")
		return domain.RenderCreateWithErrors(c, input, validationData, err)
	}

	post, validationData, err := h.service.CreatePost(c.Context(), input)
	if err != nil {
		return domain.RenderCreateWithErrors(c, input, validationData, err)
	}
	ID := conventorLib.UintToString(post.ID)
	middleware.SetFlash(h.store, c, "created successfully", "")
	return c.Redirect().To("/web/post/" + ID + "/edit")
}

func (h *WebHandler) Edit(c fiber.Ctx) error {
	postID := c.Params("id")

	post, err := h.service.GetPost(context.Background(), postID)
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("post not found")
	}
	ID := conventorLib.UintToString(post.ID)
	return c.Render("posts/edit", h.baseData(c, fiber.Map{
		"PageTitle":   "Edit post",
		"FormAction":  "/web/post/%s/update" + ID,
		"SubmitLabel": "Edit post",
		"FormMode":    "edit",
		"PostID":      post.ID,
		"Post": domain.CreatePostInput{
			UserID:  post.UserID,
			Title:   post.Title,
			Content: post.Content,
		},
	}))
}

func (h *WebHandler) Update(c fiber.Ctx) error {
	postID := c.Params("id")
	input, validationData, err := domain.BindInput(c)
	if err != nil {
		return domain.RenderEditWithErrors(c, postID, input, validationData, err)
	}

	post, validationData, err := h.service.UpdatePost(c.Context(), postID, input)
	if err != nil {
		return domain.RenderEditWithErrors(c, postID, input, validationData, err)
	}
	ID := conventorLib.UintToString(post.ID)
	middleware.SetFlash(h.store, c, "updated successfully", "")
	return c.Redirect().To("/web/post/" + ID + "/edit")
}

func (h *WebHandler) Delete(c fiber.Ctx) error {
	postID := c.Params("id")
	if err := h.service.DeletePost(c.Context(), postID); err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	middleware.SetFlash(h.store, c, "post deleted successfully", "")
	return c.Redirect().To("/web/posts/list")
}

func (h *WebHandler) GetPost(c fiber.Ctx) error {
	postID := c.Params("id")
	post, err := h.service.GetPost(context.Background(), postID)
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("post not found")
	}

	ID := conventorLib.UintToString(post.ID)
	return c.Render("posts/show", h.baseData(c, fiber.Map{
		"PageTitle":   "Show post",
		"FormAction":  "/web/post/%s/update" + ID,
		"SubmitLabel": "Edit post",
		"FormMode":    "edit",
		"PostID":      post.ID,
		"UserID":      post.UserID,
		"Title":       post.Title,
		"Content":     post.Content,
	}))
}

func (h *WebHandler) ListAllPosts(c fiber.Ctx) error {
	posts, err := h.service.ListPosts(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	flash := middleware.ConsumeFlash(h.store, c)
	return c.Render("posts/listAllPosts", h.baseData(c, fiber.Map{
		"PageTitle":    "Posts",
		"Posts":        posts,
		"FlashSuccess": flash.Success,
		"FlashError":   flash.Error,
		"CsrfToken":    csrf.TokenFromContext(c),
	}))
}

func (h *WebHandler) ListByUser(c fiber.Ctx) error {
	posts, err := h.service.ListPostsByUser(c.Context(), c.Params("user_id"))
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	return c.Render("posts/listAllPosts", h.baseData(c, fiber.Map{
		"PageTitle": "Posts",
		"Posts":     posts,
	}))
}

// parsePage same behaviour as users module: invalid or <1 -> 1
func parsePage(v string) int {
	page, err := strconv.Atoi(v)
	if err != nil || page < 1 {
		return 1
	}
	return page
}

func (h *WebHandler) ListPostsWithPagination(c fiber.Ctx) error {
	page := parsePage(c.Query("page", "1"))

	// Use the same page size as users handler for consistency
	limit := 3
	offset := (page - 1) * limit

	// Simple approach: get all posts from service and slice.
	// This mirrors the users implementation style; if dataset grows,
	// consider adding a paginated repository method instead.
	allPosts, err := h.service.ListPosts(c.Context())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	totalItems := int64(len(allPosts))

	// create page slice safely
	var postsPage []domain.Post
	if offset >= len(allPosts) {
		postsPage = []domain.Post{}
	} else {
		end := offset + limit
		if end > len(allPosts) {
			end = len(allPosts)
		}
		postsPage = allPosts[offset:end]
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	if totalPages < 1 {
		totalPages = 1
	}

	var pages []int
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}

	paging := pagination.Pagination{
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: page,
		Pages:       pages,
		HasNext:     page < totalPages,
		HasPrev:     page > 1,
		NextPage:    page + 1,
		PrevPage:    page - 1,
	}

	flash := middleware.ConsumeFlash(h.store, c)
	return c.Render("posts/listAllPosts", h.baseData(c, fiber.Map{
		"Title":        "Posts",
		"Posts":        postsPage,
		"Pagination":   paging,
		"TotalRecord":  totalItems,
		"FlashSuccess": flash.Success,
		"FlashError":   flash.Error,
		"CsrfToken":    csrf.TokenFromContext(c),
	}))
}

func (h *WebHandler) baseData(c fiber.Ctx, data fiber.Map) fiber.Map {
	flashPop := middleware.PopFlash(h.store, c)
	flash := middleware.ConsumeFlash(h.store, c)
	csrfToken := csrf.TokenFromContext(c)
	base := fiber.Map{
		"CsrfToken":    csrfToken,
		"FlashType":    flashPop.Type,
		"FlashMessage": flashPop.Message,
		"FlashSuccess": flash.Success,
		"FlashError":   flash.Error,
	}
	for k, v := range data {
		base[k] = v
	}
	return base
}
