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

		collection, err := dao.FindCollectionByNameOrId("z5lghx2r3tm45n1")
		if err != nil {
			return err
		}

		collection.ListRule = types.Pointer("@request.auth.id ?= @collection.permissions.user.id && id ?= @collection.permissions.read.id")

		collection.ViewRule = types.Pointer("@request.auth.id ?= @collection.permissions.user.id && id ?= @collection.permissions.read.id")

		collection.CreateRule = types.Pointer("@request.auth.id ?= @collection.permissions.user.id && @collection.permissions.create = true")

		collection.UpdateRule = types.Pointer("@request.auth.id ?= @collection.permissions.user.id && id ?= @collection.permissions.update.id")

		collection.DeleteRule = types.Pointer("@request.auth.id ?= @collection.permissions.user.id && id ?= @collection.permissions.delete.id")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("z5lghx2r3tm45n1")
		if err != nil {
			return err
		}

		collection.ListRule = nil

		collection.ViewRule = nil

		collection.CreateRule = nil

		collection.UpdateRule = nil

		collection.DeleteRule = nil

		return dao.SaveCollection(collection)
	})
}
