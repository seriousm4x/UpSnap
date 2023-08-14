package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("djqp3uxrac2ores")
		if err != nil {
			return err
		}

		collection.CreateRule = types.Pointer("@request.auth.id ?= @collection.permissions.user.id && @collection.permissions.create = true")

		collection.UpdateRule = types.Pointer("@request.auth.id ?= @collection.permissions.user.id && @collection.permissions.create = true")

		collection.DeleteRule = types.Pointer("@request.auth.id ?= @collection.permissions.user.id && @collection.permissions.create = true")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("djqp3uxrac2ores")
		if err != nil {
			return err
		}

		collection.CreateRule = nil

		collection.UpdateRule = nil

		collection.DeleteRule = nil

		return dao.SaveCollection(collection)
	})
}
