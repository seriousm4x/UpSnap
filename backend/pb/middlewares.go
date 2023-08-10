package pb

import (
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
)

func RequireUpSnapPermission() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			admin, _ := c.Get(apis.ContextAdminKey).(*models.Admin)
			if admin != nil {
				return next(c)
			}

			user, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
			if user == nil {
				return apis.NewUnauthorizedError("The request requires admin or record authorization token to be set.", nil)
			}

			deviceId := c.PathParam("id")

			// find record where user has device with power permission
			res, err := App.Dao().FindFirstRecordByFilter("permissions",
				fmt.Sprintf("user.id = '%s' && power.id ?= '%s'", user.Id, deviceId))
			if res == nil || err != nil {
				return apis.NewForbiddenError("You are not allowed to perform this request.", nil)
			}

			return next(c)
		}
	}
}
