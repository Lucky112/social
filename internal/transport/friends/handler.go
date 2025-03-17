package friends

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Lucky112/social/internal/models"
	"github.com/Lucky112/social/internal/transport/jwt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Обработчик HTTP-запросов на добавление и удаление друзей
type FriendsHandler struct {
	service  FriendsService
	profiles ProfileService
	validate *validator.Validate
}

func NewFriendsHandler(service FriendsService, profiles ProfileService) FriendsHandler {
	return FriendsHandler{
		service:  service,
		profiles: profiles,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (h *FriendsHandler) GetFriends(c *fiber.Ctx) error {
	userId, err := jwt.ExtractUserId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			friendError{fmt.Sprintf("failed to extract user id: %v", err)},
		)
	}

	friends, err := h.service.Get(c.Context(), userId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			friendError{fmt.Sprintf("failed to get friends: %v", err)},
		)
	}

	err = c.Status(fiber.StatusOK).JSON(
		fromModel(friends),
	)
	if err != nil {
		return fmt.Errorf("sending response: %v", err)
	}

	return nil
}

func (h *FriendsHandler) AddFriend(c *fiber.Ctx) error {
	userId, err := jwt.ExtractUserId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			friendError{fmt.Sprintf("failed to extract user id: %v", err)},
		)
	}

	friendID := c.Params("friend_id")
	if friendID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			friendError{"friend_id path parameter is required"},
		)
	}

	_, err = h.profiles.Get(c.Context(), friendID)
	if err != nil {
		if errors.Is(err, models.ProfileNotFound) {
			c.Status(fiber.StatusBadRequest).JSON(
				friendError{fmt.Sprintf("looking for friend's profile: %v", err)},
			)
			return nil
		}

		c.Status(fiber.StatusInternalServerError).JSON(
			friendError{fmt.Sprintf("failed to find friend's profile: %v", err)},
		)
		return nil
	}

	err = h.service.Add(c.Context(), userId, friendID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			friendError{fmt.Sprintf("failed to add friend '%s': %v", friendID, err)},
		)
	}

	err = c.SendStatus(fiber.StatusCreated)
	if err != nil {
		return fmt.Errorf("sending response: %v", err)
	}

	return nil
}

func (h *FriendsHandler) DeleteFriend(c *fiber.Ctx) error {
	userId, err := jwt.ExtractUserId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			friendError{fmt.Sprintf("failed to extract user id: %v", err)},
		)
	}

	friendID := c.Params("friend_id")
	if friendID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			friendError{"friend_id path parameter is required"},
		)
	}

	err = h.service.Delete(c.Context(), userId, friendID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			friendError{fmt.Sprintf("failed to delete friend '%s': %v", friendID, err)},
		)
	}

	err = c.SendStatus(fiber.StatusNoContent)
	if err != nil {
		return fmt.Errorf("sending response: %v", err)
	}

	return nil
}
