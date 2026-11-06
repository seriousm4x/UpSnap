package pb

import (
	"ganigeorgiev.com/pocketbase/dbx"
	"ganigeorgiev.com/pocketbase/pocketbase/apis"
	"ganigeorgiev.com/pocketbase/pocketbase/core"
	"ganigeorgiev.com/pocketbase/pocketbase/tools/hook"
)

// SkipSuccessfulHealthCheckLog prevents routine Docker health checks from being
// persisted while retaining failed health checks for diagnostics.
func SkipSuccessfulHealthCheckLog(e *core.RequestEvent) error {
	if e.Request.URL.Path == "/api/health" {
		return apis.SkipSuccessActivityLog().Func(e)
	}

	return e.Next()
}

func RequireScanDevicesPermission() *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Func: func(e *core.RequestEvent) error {
			if e.HasSuperuserAuth() {
				return e.Next()
			}

			user := e.Auth
			if user == nil {
				return apis.NewUnauthorizedError("The request requires superuser or record authorization token to be set.", nil)
			}

			res, err := e.App.FindFirstRecordByFilter(
				"permissions",
				"user.id = {:userId} && create = true",
				dbx.Params{"userId": user.Id},
			)
			if res == nil || err != nil {
				return apis.NewForbiddenError("You are not allowed to perform this request.", nil)
			}

			return e.Next()
		},
	}
}

func RequireUpSnapPermission() *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Func: func(e *core.RequestEvent) error {
			if e.HasSuperuserAuth() {
				return e.Next()
			}

			user := e.Auth
			if user == nil {
				return apis.NewUnauthorizedError("The request requires superuser or record authorization token to be set.", nil)
			}

			deviceId := e.Request.PathValue("id")

			// find record where user has device with power permission
			res, err := e.App.FindFirstRecordByFilter("permissions", "user.id = {:userId} && power.id ?= {:deviceId}", dbx.Params{
				"userId":   user.Id,
				"deviceId": deviceId,
			})
			if res == nil || err != nil {
				return apis.NewForbiddenError("You are not allowed to perform this request.", nil)
			}

			return e.Next()
		},
	}
}
