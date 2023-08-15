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
		new_sol_enabled := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "ueieydw5",
			"name": "sol_enabled",
			"type": "bool",
			"required": false,
			"unique": false,
			"options": {}
		}`), new_sol_enabled)
		collection.Schema.AddField(new_sol_enabled)

		// add
		new_sol_auth := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "yfszwhpw",
			"name": "sol_auth",
			"type": "bool",
			"required": false,
			"unique": false,
			"options": {}
		}`), new_sol_auth)
		collection.Schema.AddField(new_sol_auth)

		// add
		new_sol_user := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "dbcmnrmp",
			"name": "sol_user",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_sol_user)
		collection.Schema.AddField(new_sol_user)

		// add
		new_sol_password := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "hvz8stfy",
			"name": "sol_password",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_sol_password)
		collection.Schema.AddField(new_sol_password)

		// add
		new_sol_port := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "6kfqheid",
			"name": "sol_port",
			"type": "number",
			"required": true,
			"unique": false,
			"options": {
				"min": 1,
				"max": 65535
			}
		}`), new_sol_port)
		collection.Schema.AddField(new_sol_port)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("z5lghx2r3tm45n1")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("ueieydw5")

		// remove
		collection.Schema.RemoveField("yfszwhpw")

		// remove
		collection.Schema.RemoveField("dbcmnrmp")

		// remove
		collection.Schema.RemoveField("hvz8stfy")

		// remove
		collection.Schema.RemoveField("6kfqheid")

		return dao.SaveCollection(collection)
	})
}
