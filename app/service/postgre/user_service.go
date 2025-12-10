package servicepostgre

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	utils "github.com/Ahlul-Mufi/data-prestasi/utils/postgre"
)

type UserService interface {
    Login(c *fiber.Ctx) error
    Profile(c *fiber.Ctx) error
}

type userService struct {
    repo repo.UserRepository
}

func NewUserService(r repo.UserRepository) UserService {
    return &userService{r}
}


func (s *userService) Login(c *fiber.Ctx) error {
    var req m.LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON("invalid body")
    }

    user, err := s.repo.FindByIdentity(req.Identity)
    if err != nil {
        return c.Status(401).JSON("invalid identity")
    }
    if !utils.CheckPassword(user.PasswordHash, req.Password) {
        return c.Status(401).JSON("invalid credentials")
    }
    token, err := utils.GenerateToken(user) 

    if err != nil {
        return c.Status(500).JSON("token generation failed")
    }

    return c.JSON(m.LoginResponse{
        Token: token,
        User:  user,
    })
}

func (s *userService) Profile(c *fiber.Ctx) error {
    userIDStr := c.Locals("user_id")

    if userIDStr == nil {
        return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
    }

    userID, err := uuid.Parse(userIDStr.(string))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid user ID format"})
    }

    user, err := s.repo.FindByID(userID) 
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "user not found"})
    }

    return c.JSON(user)
}