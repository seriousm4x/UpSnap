package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("27do0wbcuyfmbmx")
		if err != nil {
			return err
		}

		options := map[string]any{}
		json.Unmarshal([]byte(`{
			"allowEmailAuth": true,
			"allowOAuth2Auth": true,
			"allowUsernameAuth": true,
			"exceptEmailDomains": [],
			"manageRule": null,
			"minPasswordLength": 5,
			"onlyEmailDomains": [],
			"requireEmail": false
		}`), &options)
		collection.SetOptions(options)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("27do0wbcuyfmbmx")
		if err != nil {
			return err
		}

		options := map[string]any{}
		json.Unmarshal([]byte(`{
			"allowEmailAuth": true,
			"allowOAuth2Auth": true,
			"allowUsernameAuth": true,
			"exceptEmailDomains": [],
			"manageRule": null,
			"minPasswordLength": 8,
			"onlyEmailDomains": [],
			"requireEmail": false
		}`), &options)
		collection.SetOptions(options)

		return dao.SaveCollection(collection)
	})
}
