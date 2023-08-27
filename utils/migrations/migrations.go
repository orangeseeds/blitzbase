package migrations

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/orangeseeds/blitzbase/store"
)

func CreateNewTable(store *store.Storage, c store.Collection) error {
	q := store.DB.CreateTable(c.TableName(), c.TableSchema())
	if _, err := q.Execute(); err != nil {
		return err
	}

	upSQL := fmt.Sprintf("%s;", q.SQL())
	downSQL := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, c.TableName())
	if err := CreateNewMigration(upSQL, downSQL, store.MigrationsPath); err != nil {
		return err
	}

	return nil
}

func AddCollectionRecord(s *store.Storage, c store.Collection) error {

	cType := store.BaseType
	if c.IsAuth() {
		cType = store.AuthType
	}

	schemaJson, _ := json.Marshal(c.RawSchema())

	q := s.DB.Insert("_collections", dbx.Params{
		"name":   c.TableName(),
		"type":   cType,
		"schema": schemaJson,
	})

	_, err := q.Execute()

	if err != nil {
		return err
	}
	return nil
}

func CreateNewMigration(upSQL, downSQL, path string) error {

	currVersion, err := latestVersion(path)
	if err != nil {
		return err
	}

	err = genMigrationFile(currVersion+1, downSQL, "down", path)
	if err != nil {
		return err
	}

	err = genMigrationFile(currVersion+1, upSQL, "up", path)
	if err != nil {
		return err
	}
	return nil
}

func genMigrationFile(version int, query string, upDown string, path string) error {

	fsUp, err := os.Create(fmt.Sprintf("%s/%d_install.%s.sql", path, version, upDown))
	if err != nil {
		return err
	}
	defer fsUp.Close()

	_, err = fsUp.Write([]byte(query))
	if err != nil {
		return err
	}

	return nil
}

func latestVersion(path string) (int, error) {

	var largest int
	dir, err := os.Open(path)

	if err != nil && !os.IsNotExist(err) {
		return 0, err
	}
	if os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return 0, err
		}

		dir, err = os.Open(path)
		if err != nil {
			return 0, err
		}
		defer dir.Close()
	}
	defer dir.Close()

	// Read the directory contents
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return 0, err
	}

	// Extract file names
	var fileNames []string
	for _, fileInfo := range fileInfos {
		if fileInfo.Mode().IsRegular() { // Only include regular files, not directories
			fileNames = append(fileNames, fileInfo.Name())
		}
	}

	for _, m := range fileNames {
		latest := strings.Split(m, "_")[0]
		n, err := strconv.Atoi(latest)
		if err != nil {
			continue
		}

		if n > largest {
			largest = n
		}
	}
	return largest, nil
}
