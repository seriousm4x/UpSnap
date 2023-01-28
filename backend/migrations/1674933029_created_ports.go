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
			"id": "cti4l8f4mz8df3r",
			"created": "2023-01-28 19:10:29.223Z",
			"updated": "2023-01-28 19:10:29.223Z",
			"name": "ports",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "8nwuncgg",
					"name": "number",
					"type": "number",
					"required": true,
					"unique": false,
					"options": {
						"min": null,
						"max": 65535
					}
				},
				{
					"system": false,
					"id": "o0he3pu6",
					"name": "name",
					"type": "text",
					"required": true,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "w0bh39gv",
					"name": "status",
					"type": "bool",
					"required": false,
					"unique": false,
					"options": {}
				}
			],
			"listRule": "",
			"viewRule": "",
			"createRule": "",
			"updateRule": "",
			"deleteRule": "",
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("cti4l8f4mz8df3r")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
