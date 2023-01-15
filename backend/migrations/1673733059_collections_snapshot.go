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
		jsonData := `[
			{
				"id": "z5lghx2r3tm45n1",
				"created": "2023-01-14 21:50:42.797Z",
				"updated": "2023-01-14 21:50:42.797Z",
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
						"id": "ilrwvlev",
						"name": "port",
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
					}
				],
				"listRule": "",
				"viewRule": "",
				"createRule": "",
				"updateRule": "",
				"deleteRule": "",
				"options": {}
			},
			{
				"id": "cti4l8f4mz8df3r",
				"created": "2023-01-14 21:50:42.797Z",
				"updated": "2023-01-14 21:50:42.797Z",
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
					}
				],
				"listRule": "",
				"viewRule": "",
				"createRule": "",
				"updateRule": "",
				"deleteRule": "",
				"options": {}
			},
			{
				"id": "nmj3ko20gzkg8n3",
				"created": "2023-01-14 21:50:42.797Z",
				"updated": "2023-01-14 21:50:42.797Z",
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
			}
		]`

		collections := []*models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collections); err != nil {
			return err
		}

		return daos.New(db).ImportCollections(collections, true, nil)
	}, func(db dbx.Builder) error {
		return nil
	})
}
