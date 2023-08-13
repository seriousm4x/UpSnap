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

		collection, err := dao.FindCollectionByNameOrId("27do0wbcuyfmbmx")
		if err != nil {
			return err
		}

		collection.ListRule = types.Pointer("@request.auth.id = id")

		collection.ViewRule = types.Pointer("@request.auth.id = id")

		collection.UpdateRule = types.Pointer("@request.auth.id = id")

		collection.DeleteRule = types.Pointer("@request.auth.id = id")

		// add
		new_avatar := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "i1qnezoa",
			"name": "avatar",
			"type": "number",
			"required": false,
			"unique": false,
			"options": {
				"min": 0,
				"max": 9
			}
		}`), new_avatar)
		collection.Schema.AddField(new_avatar)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("27do0wbcuyfmbmx")
		if err != nil {
			return err
		}

		collection.ListRule = nil

		collection.ViewRule = nil

		collection.UpdateRule = nil

		collection.DeleteRule = nil

		// remove
		collection.Schema.RemoveField("i1qnezoa")

		return dao.SaveCollection(collection)
	})
}
