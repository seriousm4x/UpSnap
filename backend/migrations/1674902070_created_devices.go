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
			"id": "z5lghx2r3tm45n1",
			"created": "2023-01-28 10:34:30.781Z",
			"updated": "2023-01-28 10:34:30.781Z",
			"name": "devices",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "tiqcmnjo",
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
					"id": "1si6ajha",
					"name": "ip",
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
					"id": "fyqmpon6",
					"name": "mac",
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
					"id": "s8c5z7n0",
					"name": "netmask",
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
					"id": "qqvyfrex",
					"name": "status",
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
					"id": "ilrwvlev",
					"name": "ports",
					"type": "relation",
					"required": false,
					"unique": false,
					"options": {
						"maxSelect": null,
						"collectionId": "cti4l8f4mz8df3r",
						"cascadeDelete": false
					}
				},
				{
					"system": false,
					"id": "gdctb8hj",
					"name": "link",
					"type": "url",
					"required": false,
					"unique": false,
					"options": {
						"exceptDomains": null,
						"onlyDomains": null
					}
				},
				{
					"system": false,
					"id": "kzvdbbws",
					"name": "wake_cron",
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
					"id": "dcd3vuc3",
					"name": "wake_cron_enabled",
					"type": "bool",
					"required": false,
					"unique": false,
					"options": {}
				},
				{
					"system": false,
					"id": "91bs6clw",
					"name": "shutdown_cron",
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
					"id": "vnlymcuw",
					"name": "shutdown_cron_enabled",
					"type": "bool",
					"required": false,
					"unique": false,
					"options": {}
				},
				{
					"system": false,
					"id": "1a7yrwo9",
					"name": "shutdown_cmd",
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
					"id": "iarksm3l",
					"name": "password",
					"type": "text",
					"required": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
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

		collection, err := dao.FindCollectionByNameOrId("z5lghx2r3tm45n1")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
