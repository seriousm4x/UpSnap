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

		// add
		new_lazy_ping := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "yjyq5pvg",
			"name": "lazy_ping",
			"type": "bool",
			"required": false,
			"unique": false,
			"options": {}
		}`), new_lazy_ping)
		collection.Schema.AddField(new_lazy_ping)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("nmj3ko20gzkg8n3")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("yjyq5pvg")

		return dao.SaveCollection(collection)
	})
}
