package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// Up migration: Change `url` field type from URL to string
		collectionName := "devices"

		return app.RunInTransaction(func(txApp core.App) error {
			// Step 1: Fetch the collection
			collection, err := txApp.FindCollectionByNameOrId(collectionName)
			if err != nil {
				return err
			}

			// Step 2: Add a new text field
			newField := &core.TextField{
				Name:     "new_field",
				Required: false,
			}
			collection.Fields.AddAt(6, newField)

			if err := txApp.Save(collection); err != nil {
				return err
			}

			// Step 3: Migrate data from `url` to `new_field`
			records, err := txApp.FindAllRecords(collectionName)
			if err != nil {
				return err
			}

			for _, record := range records {
				url := record.GetString("link")
				record.Set("new_field", url)
				if err := txApp.Save(record); err != nil {
					return err
				}
			}

			// Step 4: Remove the old `url` field
			collection.Fields.RemoveByName("link")
			if err := txApp.Save(collection); err != nil {
				return err
			}

			// Step 5: Rename `new_field` to `url`
			renamedField := collection.Fields.GetByName("new_field")
			renamedField.SetName("link")

			return txApp.Save(collection)
		})
	}, func(app core.App) error {
		// Down migration: Restore the `url` field
		collectionName := "devices"

		return app.RunInTransaction(func(txApp core.App) error {
			// Step 1: Fetch the collection
			collection, err := txApp.FindCollectionByNameOrId(collectionName)
			if err != nil {
				return err
			}

			// Step 2: Add a new `url` field
			newField := &core.URLField{
				Name:     "new_field",
				Required: false,
			}
			collection.Fields.AddAt(6, newField)

			if err := txApp.Save(collection); err != nil {
				return err
			}

			// Step 3: Migrate data from `url` to `new_field`
			records, err := txApp.FindAllRecords(collectionName)
			if err != nil {
				return err
			}

			for _, record := range records {
				text := record.GetString("link")
				record.Set("new_field", text)
				if err := txApp.Save(record); err != nil {
					return err
				}
			}

			// Step 4: Remove the old `url` field
			collection.Fields.RemoveByName("link")
			if err := txApp.Save(collection); err != nil {
				return err
			}

			// Step 5: Rename `new_field` to `url`
			renamedField := collection.Fields.GetByName("new_field")
			renamedField.SetName("link")

			return txApp.Save(collection)
		})
	})
}
