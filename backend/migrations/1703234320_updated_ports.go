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

		collection, err := dao.FindCollectionByNameOrId("cti4l8f4mz8df3r")
		if err != nil {
			return err
		}

		// add
		new_link := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "i3uvw7ij",
			"name": "link",
			"type": "url",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"exceptDomains": null,
				"onlyDomains": null
			}
		}`), new_link)
		collection.Schema.AddField(new_link)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("cti4l8f4mz8df3r")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("i3uvw7ij")

		return dao.SaveCollection(collection)
	})
}
