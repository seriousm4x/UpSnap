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
		new_wake_confirm := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "chtoblvl",
			"name": "wake_confirm",
			"type": "bool",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {}
		}`), new_wake_confirm)
		collection.Schema.AddField(new_wake_confirm)

		// add
		new_shutdown_confirm := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "muqy3axq",
			"name": "shutdown_confirm",
			"type": "bool",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {}
		}`), new_shutdown_confirm)
		collection.Schema.AddField(new_shutdown_confirm)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("z5lghx2r3tm45n1")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("chtoblvl")

		// remove
		collection.Schema.RemoveField("muqy3axq")

		return dao.SaveCollection(collection)
	})
}
