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
		edit_sol_port := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "6kfqheid",
			"name": "sol_port",
			"type": "number",
			"required": false,
			"unique": false,
			"options": {
				"min": 1,
				"max": 65535
			}
		}`), edit_sol_port)
		collection.Schema.AddField(edit_sol_port)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("z5lghx2r3tm45n1")
		if err != nil {
			return err
		}

		// update
		edit_sol_port := &schema.SchemaField{}
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
		}`), edit_sol_port)
		collection.Schema.AddField(edit_sol_port)

		return dao.SaveCollection(collection)
	})
}
