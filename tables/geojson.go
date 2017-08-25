package tables

import (
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
)

type GeoJSONTable struct {
	sqlite.Table
	name string
}

func NewGeoJSONTable() (*GeoJSONTable, error) {

	t := GeoJSONTable{
		name: "geojson",
	}

	return &t, nil
}

func (t *GeoJSONTable) Name() string {
	return t.name
}

func (t *GeoJSONTable) Schema() string {
	return fmt.Sprintf("CREATE TABLE %s (id INTEGER NOT NULL PRIMARY KEY, body TEXT)", t.Name())
}

func (t *GeoJSONTable) InitializeTable(db sqlite.Database) error {
	return nil
}

func (t *GeoJSONTable) IndexFeature(db sqlite.Database, f geojson.Feature) error {

	conn, err := db.Conn()

	if err != nil {
		return err
	}

	str_id := f.Id()
	body := f.Bytes()

	tx, err := conn.Begin()

	if err != nil {
		return err
	}

	sql := fmt.Sprintf("INSERT INTO %s (id, body) values(?, ?)", t.Name())

	stmt, err := tx.Prepare(sql)

	if err != nil {
		return err
	}

	defer stmt.Close()

	str_body := string(body)

	_, err = stmt.Exec(str_id, str_body)

	if err != nil {
		return err
	}

	return tx.Commit()
}