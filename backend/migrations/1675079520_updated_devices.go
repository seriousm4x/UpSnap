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
		edit_ports := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "ilrwvlev",
			"name": "ports",
			"type": "relation",
			"required": false,
			"unique": false,
			"options": {
				"collectionId": "cti4l8f4mz8df3r",
				"cascadeDelete": true,
				"maxSelect": null,
				"displayFields": null
			}
		}`), edit_ports)
		collection.Schema.AddField(edit_ports)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("z5lghx2r3tm45n1")
		if err != nil {
			return err
		}

		// update
		edit_ports := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "ilrwvlev",
			"name": "ports",
			"type": "relation",
			"required": false,
			"unique": false,
			"options": {
				"collectionId": "cti4l8f4mz8df3r",
				"cascadeDelete": false,
				"maxSelect": null,
				"displayFields": null
			}
		}`), edit_ports)
		collection.Schema.AddField(edit_ports)

		return dao.SaveCollection(collection)
	})
}
