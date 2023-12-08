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

		// add
		new_wake_cmd := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "murx1eum",
			"name": "wake_cmd",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_wake_cmd)
		collection.Schema.AddField(new_wake_cmd)

		// add
		new_ping_cmd := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "kefxqjbj",
			"name": "ping_cmd",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_ping_cmd)
		collection.Schema.AddField(new_ping_cmd)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("z5lghx2r3tm45n1")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("murx1eum")

		// remove
		collection.Schema.RemoveField("kefxqjbj")

		return dao.SaveCollection(collection)
	})
}
