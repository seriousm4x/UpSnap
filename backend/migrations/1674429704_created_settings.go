package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "nmj3ko20gzkg8n3",
			"created": "2023-01-22 23:21:44.281Z",
			"updated": "2023-01-22 23:21:44.281Z",
			"name": "settings",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "ysutxavs",
					"name": "interval",
					"type": "number",
					"required": true,
					"unique": false,
					"options": {
						"min": 1,
						"max": null
					}
				}
			],
			"listRule": null,
			"viewRule": null,
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("nmj3ko20gzkg8n3")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
