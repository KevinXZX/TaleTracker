package taledb

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/stdlib"
	"taletracker.com/internal/model"
)

type TaleDatabase struct {
	db *sql.DB
}

//go:embed migrations/*
var migrationsFS embed.FS

func OpenDatabase(url string) (tdb *TaleDatabase, err error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}

	tdb = &TaleDatabase{db: db}
	err = tdb.Migrate(url)
	if err != nil {
		return nil, err
	}
	fmt.Println("Database migrated successfully")
	return tdb, nil
}
func (tdb *TaleDatabase) Migrate(dburl string) error {
	data, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", data, dburl)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
func (tdb *TaleDatabase) AddTale(tale model.Tale) error {
	_, err := tdb.db.Exec("INSERT INTO tales (title, author, url,blurb,published,updated,review_score,review_comment) VALUES ($1, $2, $3, $4,$5,$6,$7,$8)", tale.Title, tale.Author, tale.Url, tale.Blurb, tale.Published, tale.Updated, tale.Review.Rating, tale.Review.Description)
	if err != nil {
		return err
	}
	return nil
}
func (tdb *TaleDatabase) GetTales(userID string) ([]model.Tale, error) {
	rows, err := tdb.db.Query("SELECT id,title, author, url,blurb,added,published,updated,review_score,review_comment FROM tales")
	if err != nil {
		return []model.Tale{}, err
	}
	fmt.Println("Query successful")
	defer rows.Close()
	var tales []model.Tale
	for rows.Next() {
		var tale model.Tale
		err = rows.Scan(&tale.ID, &tale.Title, &tale.Author, &tale.Url, &tale.Blurb, &tale.Added, &tale.Published, &tale.Updated, &tale.Review.Rating, &tale.Review.Description)
		if err != nil {
			return []model.Tale{}, err
		}
		tales = append(tales, tale)
	}
	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		return []model.Tale{}, err
	}
	return tales, nil
}
