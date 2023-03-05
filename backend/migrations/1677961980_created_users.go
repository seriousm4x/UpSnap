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
			"id": "27do0wbcuyfmbmx",
			"created": "2023-03-04 20:33:00.558Z",
			"updated": "2023-03-04 20:33:00.558Z",
			"name": "users",
			"type": "auth",
			"system": false,
			"schema": [],
			"listRule": null,
			"viewRule": null,
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {
				"allowEmailAuth": true,
				"allowOAuth2Auth": false,
				"allowUsernameAuth": true,
				"exceptEmailDomains": [],
				"manageRule": null,
				"minPasswordLength": 8,
				"onlyEmailDomains": [],
				"requireEmail": false
			}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("27do0wbcuyfmbmx")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
