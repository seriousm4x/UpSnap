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

		collection, err := dao.FindCollectionByNameOrId("z5lghx2r3tm45n1")
		if err != nil {
			return err
		}

		// update
		edit_shutdown := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "1a7yrwo9",
			"name": "shutdown",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_shutdown)
		collection.Schema.AddField(edit_shutdown)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("z5lghx2r3tm45n1")
		if err != nil {
			return err
		}

		// update
		edit_shutdown := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "1a7yrwo9",
			"name": "shutdown_cmd",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_shutdown)
		collection.Schema.AddField(edit_shutdown)

		return dao.SaveCollection(collection)
	})
}
