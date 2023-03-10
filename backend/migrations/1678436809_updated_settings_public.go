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

		collection, err := dao.FindCollectionByNameOrId("t1rajoyzl6691g3")
		if err != nil {
			return err
		}

		// add
		new_setup_completed := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "vrxqcsnv",
			"name": "setup_completed",
			"type": "bool",
			"required": false,
			"unique": false,
			"options": {}
		}`), new_setup_completed)
		collection.Schema.AddField(new_setup_completed)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("t1rajoyzl6691g3")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("vrxqcsnv")

		return dao.SaveCollection(collection)
	})
}
