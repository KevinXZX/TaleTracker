package taledb

import (
	"database/sql"
	"embed"
	"encoding/json"
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
	tx, err := tdb.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// Add tale to database
	taleID, err := tdb.addTaleToDB(tx, tale)
	if err != nil {
		fmt.Println("Error adding tale to database", err)
		return err
	}
	// Add or create tags
	for _, tag := range tale.Tags {
		tagID, err := tdb.addOrCreateTaleTag(tx, tag.Name)
		if err != nil {
			return err
		}
		// Add tag to tale
		_, err = tx.Exec("INSERT INTO tales_tags (tale_id,tag_id) VALUES ($1, $2)", taleID, tagID)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (tdb *TaleDatabase) addTaleToDB(tx *sql.Tx, tale model.Tale) (int64, error) {
	// Add tale to database
	var taleID int64
	err := tx.QueryRow("INSERT INTO tales (title, author, url,blurb,published,updated,review_score,review_comment) VALUES ($1, $2, $3, $4,$5,$6,$7,$8) RETURNING id", tale.Title, tale.Author, tale.Url, tale.Blurb, tale.Published, tale.Updated, tale.Review.Rating, tale.Review.Description).Scan(&taleID)
	if err != nil {
		return -1, err
	}

	return taleID, nil
}

func (tdb *TaleDatabase) addOrCreateTaleTag(tx *sql.Tx, tagName string) (int64, error) {
	var tagID int64

	// Check if the tag already exists
	err := tx.QueryRow("SELECT id FROM tags WHERE name = $1", tagName).Scan(&tagID)
	if err == sql.ErrNoRows {
		// Tag doesn't exist, so create it
		err = tx.QueryRow("INSERT INTO tags (name) VALUES ($1) RETURNING id", tagName).Scan(&tagID)
		if err != nil {
			return -1, err
		}
	} else if err != nil {
		return -1, err
	}

	return tagID, nil
}

func (tdb *TaleDatabase) ReviewTale(taleID int, review model.Review) error {
	_, err := tdb.db.Exec("UPDATE tales SET review_score = $1, review_comment = $2 WHERE id = $3", review.Rating, review.Description, taleID)
	if err != nil {
		return err
	}
	return nil
}

func (tdb *TaleDatabase) AddTagToTale(taleID int, tags ...model.Tag) error {
	tx, err := tdb.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, tag := range tags {
		tagID, err := tdb.addOrCreateTaleTag(tx, tag.Name)
		if err != nil {
			return err
		}
		// Add tag to tale
		_, err = tx.Exec("INSERT INTO tales_tags (tale_id,tag_id) VALUES ($1, $2)", taleID, tagID)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (tdb *TaleDatabase) GetTales() ([]model.Tale, error) {
	rows, err := tdb.db.Query(`
		SELECT 
			t.id,
			t.title,
			t.author,
			t.url,
			t.blurb,
			t.added,
			t.published,
			t.updated,
			t.review_score,
			t.review_comment,
			COALESCE(json_agg(DISTINCT tags) FILTER (WHERE tags.id IS NOT NULL), '[]'::json) AS tags
		FROM tales t
		LEFT JOIN tales_tags tt ON t.id = tt.tale_id
		LEFT JOIN tags ON tt.tag_id = tags.id
		GROUP BY t.id
	`)
	if err != nil {
		return []model.Tale{}, err
	}
	defer rows.Close()

	var tales []model.Tale
	for rows.Next() {
		var tale model.Tale
		var tagsJSON string
		err = rows.Scan(&tale.ID, &tale.Title, &tale.Author, &tale.Url, &tale.Blurb, &tale.Added, &tale.Published, &tale.Updated, &tale.Review.Rating, &tale.Review.Description, &tagsJSON)
		if err != nil {
			return []model.Tale{}, err
		}

		// Unmarshal JSON array of tags into []Tag
		var tags []model.Tag
		err = json.Unmarshal([]byte(tagsJSON), &tags)
		if err != nil {
			return []model.Tale{}, err
		}
		tale.Tags = tags

		tales = append(tales, tale)
	}

	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		return []model.Tale{}, err
	}

	return tales, nil
}
