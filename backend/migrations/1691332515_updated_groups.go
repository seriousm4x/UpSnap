package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("djqp3uxrac2ores")
		if err != nil {
			return err
		}

		collection.Name = "device_groups"

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("djqp3uxrac2ores")
		if err != nil {
			return err
		}

		collection.Name = "groups"

		return dao.SaveCollection(collection)
	})
}
