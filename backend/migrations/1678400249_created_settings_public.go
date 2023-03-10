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
			"id": "t1rajoyzl6691g3",
			"created": "2023-03-09 22:17:29.463Z",
			"updated": "2023-03-09 22:17:29.463Z",
			"name": "settings_public",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "s1grje98",
					"name": "website_title",
					"type": "text",
					"required": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "tlzz87pj",
					"name": "favicon",
					"type": "file",
					"required": false,
					"unique": false,
					"options": {
						"maxSelect": 1,
						"maxSize": 5242880,
						"mimeTypes": [
							"image/x-icon",
							"image/png",
							"image/svg+xml",
							"image/gif",
							"image/jpeg"
						],
						"thumbs": []
					}
				}
			],
			"listRule": "",
			"viewRule": "",
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

		collection, err := dao.FindCollectionByNameOrId("t1rajoyzl6691g3")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
