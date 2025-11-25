package servicepostgre

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
)

type RoleService interface {
    GetAll(c *fiber.Ctx) error
    GetByID(c *fiber.Ctx) error
    Create(c *fiber.Ctx) error
    Update(c *fiber.Ctx) error
    Delete(c *fiber.Ctx) error
}

type roleService struct {
    repo repo.RoleRepository
}

func NewRoleService(r repo.RoleRepository) RoleService {
    return &roleService{r}
}

func (s *roleService) GetAll(c *fiber.Ctx) error {
    data, err := s.repo.GetAll()
    if err != nil {
        return c.Status(500).JSON(err.Error())
    }
    return c.JSON(data)
}

func (s *roleService) GetByID(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid uuid"})
    }

    data, err := s.repo.GetByID(id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "not found"})
    }

    return c.JSON(data)
}

func (s *roleService) Create(c *fiber.Ctx) error {
    var req m.CreateRoleRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON("invalid body")
    }

    role := m.Role{
        ID:          uuid.New(),
        Name:        req.Name,
        Description: req.Description,
    }

    newRole, err := s.repo.Create(role)
    if err != nil {
        return c.Status(500).JSON(err.Error())
    }

    return c.Status(201).JSON(newRole)
}

func (s *roleService) Update(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        return c.Status(400).JSON("invalid uuid")
    }

    var req m.UpdateRoleRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON("invalid body")
    }

    update := m.Role{
        Name:        req.Name,
        Description: req.Description,
    }

    updated, err := s.repo.Update(id, update)
    if err != nil {
        return c.Status(500).JSON(err.Error())
    }

    return c.JSON(updated)
}

func (s *roleService) Delete(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        return c.Status(400).JSON("invalid uuid")
    }

    if err := s.repo.Delete(id); err != nil {
        return c.Status(404).JSON("not found")
    }

    return c.JSON(fiber.Map{"message": "deleted"})
}