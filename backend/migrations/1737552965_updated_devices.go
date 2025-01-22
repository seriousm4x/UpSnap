package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("z5lghx2r3tm45n1")
		if err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(8, []byte(`{
			"hidden": false,
			"id": "select355930381",
			"maxSelect": 1,
			"name": "link_open",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "select",
			"values": [
				"same_tab",
				"new_tab"
			]
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("z5lghx2r3tm45n1")
		if err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("select355930381")

		return app.Save(collection)
	})
}
