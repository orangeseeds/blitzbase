package store

import (
	"fmt"

	dbx "github.com/go-ozzo/ozzo-dbx"
	model "github.com/orangeseeds/blitzbase/refactored/models"
)

func (s *BaseStore) TableExists(tableName string) bool {
	var exists bool
	err := s.DB().Select("count(*)").
		From("sqlite_schema").
		AndWhere(dbx.HashExp{"type": []any{"table", "view"}}).
		AndWhere(dbx.NewExp("LOWER([[name]])=LOWER({:tableName})", dbx.Params{"tableName": tableName})).
		Limit(1).
		Row(&exists)

	return err == nil && exists
}

func (s *BaseStore) createTable(name string, defn map[string]string) error {
	query := s.DB().CreateTable(name, defn)
	if _, err := query.Execute(); err != nil {
		return err
	}

	// Possible place to add migrations
	return nil
}

// Make the default fields fit in better
func (s *BaseStore) CreateCollectionTable(c *model.Collection) error {
	if s.TableExists(c.GetName()) {
		return fmt.Errorf("Table with name %s already exists.", c.GetName())
	}

	json, err := c.Schema.MarshalJSON()
	if err != nil {
		return err
	}

	_, err = s.DB().Insert(c.TableName(), dbx.Params{
		"Id":     c.GetID(),
		"Name":   c.GetName(),
		"Type":   c.Type,
		"Schema": string(json),
		"Rule":   c.Rule,
	}).Execute()
	if err != nil {
		return err
	}

	err = s.createTable(c.GetName(), c.DataDefn())
	if err != nil {
		return err
	}

	return nil
}

// make sure if you are running this check if previous admin tables are present
func (s *BaseStore) CreateAdminTable() error {
	var a *model.Admin

	if s.TableExists(a.TableName()) {
		return fmt.Errorf("Table with name %s already exists.", a.TableName())
	}

	err := s.createTable(a.TableName(), a.ExportForTable())
	if err != nil {
		return err
	}
	return nil
}

func (s *BaseStore) CreateCollectionMetaTable() error {
	var col model.Collection

	if s.TableExists(col.TableName()) {
		return fmt.Errorf("Table with name %s already exists.", col.TableName())
	}

	err := s.createTable(col.TableName(), col.MetaDataDefn())
	if err != nil {
		return err
	}
	return nil
}
