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
			"id": "7emu7fnnwuim2ua",
			"created": "2023-08-10 21:19:39.520Z",
			"updated": "2023-08-10 21:19:39.520Z",
			"name": "permissions",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "tzxtuykz",
					"name": "user",
					"type": "relation",
					"required": false,
					"unique": false,
					"options": {
						"collectionId": "27do0wbcuyfmbmx",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": 1,
						"displayFields": []
					}
				},
				{
					"system": false,
					"id": "nd1ekexz",
					"name": "create",
					"type": "bool",
					"required": false,
					"unique": false,
					"options": {}
				},
				{
					"system": false,
					"id": "z79gcpy3",
					"name": "read",
					"type": "relation",
					"required": false,
					"unique": false,
					"options": {
						"collectionId": "z5lghx2r3tm45n1",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": null,
						"displayFields": [
							"name"
						]
					}
				},
				{
					"system": false,
					"id": "t7wwbk2c",
					"name": "update",
					"type": "relation",
					"required": false,
					"unique": false,
					"options": {
						"collectionId": "z5lghx2r3tm45n1",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": null,
						"displayFields": [
							"name"
						]
					}
				},
				{
					"system": false,
					"id": "b4u0gqa7",
					"name": "delete",
					"type": "relation",
					"required": false,
					"unique": false,
					"options": {
						"collectionId": "z5lghx2r3tm45n1",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": null,
						"displayFields": [
							"name"
						]
					}
				},
				{
					"system": false,
					"id": "nyeqt3fn",
					"name": "power",
					"type": "relation",
					"required": false,
					"unique": false,
					"options": {
						"collectionId": "z5lghx2r3tm45n1",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": null,
						"displayFields": []
					}
				}
			],
			"indexes": [
				"CREATE UNIQUE INDEX ` + "`" + `idx_2GIk1Fo` + "`" + ` ON ` + "`" + `permissions` + "`" + ` (` + "`" + `user` + "`" + `)"
			],
			"listRule": "@request.auth.id = user.id",
			"viewRule": "@request.auth.id = user.id",
			"createRule": null,
			"updateRule": null,
			"deleteRule": "@request.auth.id = user.id",
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("7emu7fnnwuim2ua")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
