package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("z5lghx2r3tm45n1")
		if err != nil {
			return err
		}

		collection.ListRule = types.Pointer("@request.auth.id = @collection.permissions.user.id && id ?= @collection.permissions.read.id")

		collection.ViewRule = types.Pointer("@request.auth.id = @collection.permissions.user.id && id ?= @collection.permissions.read.id")

		collection.CreateRule = types.Pointer("@request.auth.id = @collection.permissions.user.id && @collection.permissions.create = true")

		collection.UpdateRule = types.Pointer("@request.auth.id = @collection.permissions.user.id && id ?= @collection.permissions.update.id")

		collection.DeleteRule = types.Pointer("@request.auth.id = @collection.permissions.user.id && id ?= @collection.permissions.delete.id")

		// add
		new_created_by := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "naciykac",
			"name": "created_by",
			"type": "relation",
			"required": false,
			"unique": false,
			"options": {
				"collectionId": "27do0wbcuyfmbmx",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": [
					"username"
				]
			}
		}`), new_created_by)
		collection.Schema.AddField(new_created_by)

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

		// remove
		collection.Schema.RemoveField("naciykac")

		return dao.SaveCollection(collection)
	})
}
