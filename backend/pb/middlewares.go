package pb

import (
	"fmt"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
)

func RequireUpSnapPermission() *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Func: func(e *core.RequestEvent) error {
			admin := e.HasSuperuserAuth()
			if admin {
				return e.Next()
			}

			user := e.Auth
			if user == nil {
				return apis.NewUnauthorizedError("The request requires admin or record authorization token to be set.", nil)
			}

			deviceId := e.Request.PathValue("id")

			// find record where user has device with power permission
			res, err := App.FindFirstRecordByFilter("permissions",
				fmt.Sprintf("user.id = '%s' && power.id ?= '%s'", user.Id, deviceId))
			if res == nil || err != nil {
				return apis.NewForbiddenError("You are not allowed to perform this request.", nil)
			}

			return e.Next()
		},
	}
}
