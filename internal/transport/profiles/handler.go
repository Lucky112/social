package profiles

import (
	"errors"
	"fmt"

	"github.com/Lucky112/social/internal/models"
	"github.com/gofiber/fiber/v2"
)

// Обработчик HTTP-запросов на создание и просмотр анкет
type ProfilesHandler struct {
	service ProfilesService
}

func NewProfilesHandler(service ProfilesService) ProfilesHandler {
	return ProfilesHandler{
		service: service,
	}
}

// Обработчик HTTP-запросов на создание анкеты
func (h *ProfilesHandler) CreateProfile(c *fiber.Ctx) error {
	var payload profile
	fmt.Println(string(c.Body()))
	err := c.BodyParser(&payload)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(
			profileError{fmt.Sprintf("failed to parse body: %v", err)},
		)
		return nil
	}

	userId, ok := c.Context().Value("user").(string)
	if !ok {
		c.Status(fiber.StatusBadRequest).JSON(
			profileError{fmt.Sprintf("failed to extract user id: %v", err)},
		)
		return nil
	}
	payload.userId = userId

	p, err := payload.toModel()
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(
			profileError{fmt.Sprintf("failed to parse profile: %v", err)},
		)
		return nil
	}

	id, err := h.service.Add(c.Context(), p)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(
			profileError{fmt.Sprintf("failed to save profile: %v", err)},
		)
		return nil
	}

	err = c.JSON(profileResponse{id})
	if err != nil {
		return fmt.Errorf("sending response: %v", err)
	}

	return nil
}

// Обработчик HTTP-запросов на конкретную анкет
func (h *ProfilesHandler) GetProfile(c *fiber.Ctx) error {
	id := c.Params("id")

	p, err := h.service.Get(c.Context(), id)
	if err != nil {
		if errors.Is(err, models.ProfileNotFound) {
			c.Status(fiber.StatusNotFound).JSON(
				profileError{err.Error()},
			)
			return nil
		}

		c.Status(fiber.StatusInternalServerError).JSON(
			profileError{fmt.Sprintf("failed to find profile: %v", err)},
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

// Обработчик HTTP-запросов на список анкет
func (h *ProfilesHandler) GetProfiles(c *fiber.Ctx) error {
	profiles, err := h.service.GetAll(c.Context())
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(
			profileError{fmt.Sprintf("failed to get all profiles: %v", err)},
		)
		return nil
	}

	payload := make([]*profile, len(profiles))

	for i, p := range profiles {
		payload[i] = fromModel(p)
	}

	err = c.JSON(payload)
	if err != nil {
		return fmt.Errorf("sending response: %v", err)
	}
	return nil
}
