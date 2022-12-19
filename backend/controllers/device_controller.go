package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/seriousm4x/UpSnap/models"
	"github.com/seriousm4x/UpSnap/queries"
)

func GetDevices(c *fiber.Ctx) error {
	var devices []models.Device

	if err := queries.GetAllDevices(&devices); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"devices": devices,
	})
}

func CreateDevice(c *fiber.Ctx) error {
	var device models.Device

	if err := c.BodyParser(&device); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if err := queries.CreateDevice(&device); err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":  false,
		"msg":    "created",
		"device": device,
	})
}

func PatchDevice(c *fiber.Ctx) error {
	var params map[string]interface{}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "strconv: " + err.Error(),
		})
	}

	if err := c.BodyParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "paramsparser: " + err.Error(),
		})
	}

	if err := queries.PatchDevice(&params, id); err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   "patching: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "patched",
	})
}

func DeleteDevice(c *fiber.Ctx) error {
	var device models.Device

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if err := queries.GetOneDevice(&device, id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if err := queries.DeleteDevice(&device, id); err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "deleted",
	})
}
