package main

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-index"
	"github.com/whosonfirst/go-whosonfirst-index/utils"
	"github.com/whosonfirst/go-whosonfirst-log"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
	"github.com/whosonfirst/go-whosonfirst-sqlite/database"
	"github.com/whosonfirst/go-whosonfirst-sqlite/tables"
	"io"
	"os"
)

func main() {

	mode := flag.String("mode", "files", "The mode to use importing data.")

	geojson := flag.Bool("geojson", false, "Index the 'geojson' table")
	spr := flag.Bool("spr", false, "Index the 'spr' table")
	names := flag.Bool("names", false, "Index the 'names' table")
	all := flag.Bool("all", false, "Index all tables")

	dsn := flag.String("dsn", ":memory:", "")

	flag.Parse()

	logger := log.SimpleWOFLogger()

	db, err := database.NewDB(*dsn)

	if err != nil {
		logger.Fatal("unable to create database (%s) because %s", *dsn, err)
	}

	defer db.Close()

	to_index := make([]sqlite.Table, 0)

	if *geojson || *all {

		gt, err := tables.NewGeoJSONTable()

		if err != nil {
			logger.Fatal("failed to create 'geojson' table because %s", err)
		}

		err = gt.InitializeTable(db)

		if err != nil {
			logger.Fatal("failed to initialize 'geojson' table because %s", err)
		}

		to_index = append(to_index, gt)
	}

	if *spr || *all {

		st, err := tables.NewSPRTable()

		if err != nil {
			logger.Fatal("failed to create 'spr' table because %s", err)
		}

		err = st.InitializeTable(db)

		if err != nil {
			logger.Fatal("failed to initialize 'spr' table because %s", err)
		}

		to_index = append(to_index, st)
	}

	if *names || *all {

		nm, err := tables.NewNamesTable()

		if err != nil {
			logger.Fatal("failed to create 'names' table because %s", err)
		}

		err = nm.InitializeTable(db)

		if err != nil {
			logger.Fatal("failed to initialize 'names' table because %s", err)
		}

		to_index = append(to_index, nm)
	}

	cb := func(fh io.Reader, ctx context.Context, args ...interface{}) error {

		ok, err := utils.IsPrincipalWOFRecord(fh, ctx)

		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		f, err := feature.LoadWOFFeatureFromReader(fh)

		if err != nil {
			return err
		}

		db.Lock()

		defer db.Unlock()

		for _, t := range to_index {

			err = t.IndexFeature(db, f)

			if err != nil {
				logger.Status("failed to index feature in '%s' table because %s", t.Name(), err)
				return err
			}
		}

		return nil
	}

	indexer, err := index.NewIndexer(*mode, cb)

	if err != nil {
		logger.Fatal("Failed to create new indexer because: %s", err)
	}

	err = indexer.IndexPaths(flag.Args())

	if err != nil {
		logger.Fatal("Failed to index paths in %s mode because: %s", *mode, err)
	}

	os.Exit(0)
}
