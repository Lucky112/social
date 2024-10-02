package profiles

import (
	"errors"
	"fmt"

	"github.com/Lucky112/social/internal/models"
	"github.com/Lucky112/social/internal/transport/jwt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Обработчик HTTP-запросов на создание и просмотр анкет
type ProfilesHandler struct {
	service  ProfilesService
	validate *validator.Validate
}

func NewProfilesHandler(service ProfilesService) ProfilesHandler {
	return ProfilesHandler{
		service:  service,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

// Обработчик HTTP-запросов на создание анкеты
func (h *ProfilesHandler) CreateProfile(c *fiber.Ctx) error {
	var payload profile

	err := c.BodyParser(&payload)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(
			profileError{fmt.Sprintf("failed to parse body: %v", err)},
		)
		return nil
	}

	err = h.validate.Struct(payload)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(
			profileError{fmt.Sprintf("invalid body: %v", err)},
		)
		return nil
	}

	userId, err := jwt.ExtractUserId(c)
	if err != nil {
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

	err = c.Status(fiber.StatusCreated).JSON(
		profileResponse{id},
	)
	if err != nil {
		return fmt.Errorf("sending response: %v", err)
	}

	return nil
}

// Обработчик HTTP-запросов на конкретную анкет
func (h *ProfilesHandler) GetProfileById(c *fiber.Ctx) error {
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

// Обработчик HTTP-запросов на поиск анкеты по параметрам
func (h *ProfilesHandler) SearchProfile(c *fiber.Ctx) error {
	name := c.Query("name")
	surname := c.Query("surname")

	params := &models.SearchParams{
		NamePrefix:    name,
		SurnamePrefix: surname,
	}

	profiles, err := h.service.Search(c.Context(), params)
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
