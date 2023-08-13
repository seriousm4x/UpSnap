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

		// remove
		collection.Schema.RemoveField("qqvyfrex")

		// add
		new_status := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "qmiatdwa",
			"name": "status",
			"type": "select",
			"required": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"values": [
					"pending",
					"online",
					"offline"
				]
			}
		}`), new_status)
		collection.Schema.AddField(new_status)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("z5lghx2r3tm45n1")
		if err != nil {
			return err
		}

		// add
		del_status := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "qqvyfrex",
			"name": "status",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), del_status)
		collection.Schema.AddField(del_status)

		// remove
		collection.Schema.RemoveField("qmiatdwa")

		return dao.SaveCollection(collection)
	})
}
