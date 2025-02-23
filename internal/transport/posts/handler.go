package posts

import (
	"errors"
	"fmt"

	"github.com/Lucky112/social/internal/models"
	"github.com/Lucky112/social/internal/transport/jwt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Обработчик HTTP-запросов на создание и просмотр постов
type PostsHandler struct {
	service  PostsService
	validate *validator.Validate
}

func NewPostsHandler(service PostsService) PostsHandler {
	return PostsHandler{
		service:  service,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

// Обработчик HTTP-запросов на создание поста
func (h *PostsHandler) CreatePost(c *fiber.Ctx) error {
	var payload post

	err := c.BodyParser(&payload)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(
			postError{fmt.Sprintf("failed to parse body: %v", err)},
		)
		return nil
	}

	err = h.validate.Struct(payload)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(
			postError{fmt.Sprintf("invalid body: %v", err)},
		)
		return nil
	}

	userId, err := jwt.ExtractUserId(c)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(
			postError{fmt.Sprintf("failed to extract user id: %v", err)},
		)
		return nil
	}
	payload.userId = userId

	p, err := payload.toModel()
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(
			postError{fmt.Sprintf("failed to parse post: %v", err)},
		)
		return nil
	}

	id, err := h.service.Add(c.Context(), p)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(
			postError{fmt.Sprintf("failed to save post: %v", err)},
		)
		return nil
	}

	err = c.Status(fiber.StatusCreated).JSON(
		postResponse{id},
	)
	if err != nil {
		return fmt.Errorf("sending response: %v", err)
	}

	return nil
}

// Обработчик HTTP-запросов на конкретный пост
func (h *PostsHandler) GetPostById(c *fiber.Ctx) error {
	postID := c.Params("post_id")
	userID := c.Params("user_id")

	p, err := h.service.Get(c.Context(), userID, postID)
	if err != nil {
		if errors.Is(err, models.PostNotFound) {
			c.Status(fiber.StatusNotFound).JSON(
				postError{err.Error()},
			)
			return nil
		}

		c.Status(fiber.StatusInternalServerError).JSON(
			postError{fmt.Sprintf("failed to find post: %v", err)},
		)
		return nil
	}

	payload := fromModel(p)

	err = c.JSON(payload)
	if err != nil {
		return fmt.Errorf("sending response: %v", err)
	}

	return nil
}

// Обработчик HTTP-запросов на список постов
func (h *PostsHandler) GetPosts(c *fiber.Ctx) error {
	userID := c.Params("user_id")

	posts, err := h.service.GetAll(c.Context(), userID)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(
			postError{fmt.Sprintf("failed to get all posts: %v", err)},
		)
		return nil
	}

	payload := make([]*post, len(posts))

	for i, p := range posts {
		payload[i] = fromModel(p)
	}

	err = c.JSON(payload)
	if err != nil {
		return fmt.Errorf("sending response: %v", err)
	}
	return nil
}
