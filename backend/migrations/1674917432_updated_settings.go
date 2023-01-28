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
		new_website_title := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "slrjkipw",
			"name": "website_title",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_website_title)
		collection.Schema.AddField(new_website_title)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("nmj3ko20gzkg8n3")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("slrjkipw")

		return dao.SaveCollection(collection)
	})
}
