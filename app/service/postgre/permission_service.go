package servicepostgre

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
)

type PermissionService interface {
    GetAll(c *fiber.Ctx) error
    GetByID(c *fiber.Ctx) error
    Create(c *fiber.Ctx) error
    Update(c *fiber.Ctx) error
    Delete(c *fiber.Ctx) error
}

type permissionService struct {
    repo repo.PermissionRepository
}

func NewPermissionService(r repo.PermissionRepository) PermissionService {
    return &permissionService{r}
}

func (s *permissionService) GetAll(c *fiber.Ctx) error {
    data, err := s.repo.GetAll()
    if err != nil {
        return c.Status(500).JSON(err.Error())
    }
    return c.JSON(data)
}

func (s *permissionService) GetByID(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        return c.Status(400).JSON("invalid uuid")
    }

    data, err := s.repo.GetByID(id)
    if err != nil {
        return c.Status(404).JSON("not found")
    }

    return c.JSON(data)
}

func (s *permissionService) Create(c *fiber.Ctx) error {
    var req m.CreatePermissionRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON("invalid body")
    }

    p := m.Permission{
        Name:        req.Name,
        Resource:    req.Resource,
        Action:      req.Action,
        Description: req.Description,
    }

    newP, err := s.repo.Create(p)
    if err != nil {
        return c.Status(500).JSON(err.Error())
    }

    return c.Status(201).JSON(newP)
}

func (s *permissionService) Update(c *fiber.Ctx) error {
    id, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON("invalid uuid")
    }

    var req m.UpdatePermissionRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON("invalid body")
    }

    p := m.Permission{
        Name:        req.Name,
        Resource:    req.Resource,
        Action:      req.Action,
        Description: req.Description,
    }

    updated, err := s.repo.Update(id, p)
    if err != nil {
        return c.Status(500).JSON(err.Error())
    }

    return c.JSON(updated)
}

func (s *permissionService) Delete(c *fiber.Ctx) error {
    id, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON("invalid uuid")
    }

    if err := s.repo.Delete(id); err != nil {
        return c.Status(404).JSON("not found")
    }

    return c.JSON(fiber.Map{"message": "deleted"})
}
