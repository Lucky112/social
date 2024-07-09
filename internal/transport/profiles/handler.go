package profiles

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Обработчик HTTP-запросов на регистрацию и аутентификацию пользователей
type ProfilesHandler struct {
	storage ProfilesStorage
}

func NewProfilesHandler(storage ProfilesStorage) ProfilesHandler {
	return ProfilesHandler{
		storage: storage,
	}
}

// Обработчик HTTP-запросов на создание анкеты
func (h *ProfilesHandler) CreateProfile(c *fiber.Ctx) error {
	var payload profile
	fmt.Println(string(c.Body()))
	err := c.BodyParser(&payload)
	if err != nil {
		return fmt.Errorf("parsing body: %v", err)
	}

	p, err := payload.toModel()
	if err != nil {
		return fmt.Errorf("parsing profile: %v", err)
	}
	id := generateId()

	err = h.storage.Add(id, p)
	if err != nil {
		return fmt.Errorf("storing profile: %v", err)
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

	p, err := h.storage.Get(id)
	if err != nil {
		return fmt.Errorf("getting profile: %v", err)
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
	profiles, err := h.storage.GetAll()
	if err != nil {
		return fmt.Errorf("getting profile: %v", err)
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

func generateId() string {
	return uuid.NewString()
}
