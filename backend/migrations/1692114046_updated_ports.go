package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("cti4l8f4mz8df3r")
		if err != nil {
			return err
		}

		// update
		edit_number := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "8nwuncgg",
			"name": "number",
			"type": "number",
			"required": true,
			"unique": false,
			"options": {
				"min": 1,
				"max": 65535
			}
		}`), edit_number)
		collection.Schema.AddField(edit_number)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("cti4l8f4mz8df3r")
		if err != nil {
			return err
		}

		// update
		edit_number := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "8nwuncgg",
			"name": "number",
			"type": "number",
			"required": true,
			"unique": false,
			"options": {
				"min": null,
				"max": 65535
			}
		}`), edit_number)
		collection.Schema.AddField(edit_number)

		return dao.SaveCollection(collection)
	})
}
