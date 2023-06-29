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
		new_groups := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "rk2vrn57",
			"name": "groups",
			"type": "relation",
			"required": false,
			"unique": false,
			"options": {
				"collectionId": "djqp3uxrac2ores",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": null,
				"displayFields": []
			}
		}`), new_groups)
		collection.Schema.AddField(new_groups)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("z5lghx2r3tm45n1")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("rk2vrn57")

		return dao.SaveCollection(collection)
	})
}
