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
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	}, func(db dbx.Builder) error {
		jsonData := `{
			"id": "_pb_users_auth_",
			"created": "2023-01-27 21:08:42.303Z",
			"updated": "2023-01-27 21:08:42.304Z",
			"name": "users",
			"type": "auth",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "users_name",
					"name": "name",
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
					"id": "users_avatar",
					"name": "avatar",
					"type": "file",
					"required": false,
					"unique": false,
					"options": {
						"maxSelect": 1,
						"maxSize": 5242880,
						"mimeTypes": [
							"image/jpg",
							"image/jpeg",
							"image/png",
							"image/svg+xml",
							"image/gif",
							"image/webp"
						],
						"thumbs": null
					}
				}
			],
			"listRule": "id = @request.auth.id",
			"viewRule": "id = @request.auth.id",
			"createRule": "",
			"updateRule": "id = @request.auth.id",
			"deleteRule": "id = @request.auth.id",
			"options": {
				"allowEmailAuth": true,
				"allowOAuth2Auth": true,
				"allowUsernameAuth": true,
				"exceptEmailDomains": null,
				"manageRule": null,
				"minPasswordLength": 8,
				"onlyEmailDomains": null,
				"requireEmail": false
			}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	})
}
