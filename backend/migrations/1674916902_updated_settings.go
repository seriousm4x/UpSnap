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

		collection, err := dao.FindCollectionByNameOrId("nmj3ko20gzkg8n3")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("7ehvglvv")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("nmj3ko20gzkg8n3")
		if err != nil {
			return err
		}

		// add
		del_notifications := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "7ehvglvv",
			"name": "notifications",
			"type": "bool",
			"required": false,
			"unique": false,
			"options": {}
		}`), del_notifications)
		collection.Schema.AddField(del_notifications)

		return dao.SaveCollection(collection)
	})
}
