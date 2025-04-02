package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	sqldb "github.com/saurify/saurify-backend/internal/postgres"
)

func SaveURLToSQL(shortCode, originalURL string, isTemporary bool) error {
	ctx := context.Background()

	query := `INSERT INTO shortlinks (short_code, original_url, is_temporary) VALUES ($1, $2, $3) ON CONFLICT (short_code) DO NOTHING;`
	_, err := sqldb.DB.Exec(ctx, query, shortCode, originalURL, isTemporary)

	if err != nil {
		log.Println("Error saving URL to database")
		return err
	}

	return nil
}

func GetURLFromSQL(shortCode string) (string, bool, error) {
	var originalURL string
	var isTemporary bool

	ctx := context.Background()

	query := `SELECT origin_url FROM shortlinks WHERE short_code = $1;`
	err := sqldb.DB.QueryRow(ctx, query, shortCode).Scan(&originalURL, &isTemporary)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, fmt.Errorf("error fetching URL : %v", err)
	}

	return originalURL, isTemporary, nil
}

func DeleteFromSQL(shortCode string) error {
	ctx := context.Background()

	query := `DELETE FROM shortlinks WHERE short_code = $1`
	_, err := sqldb.DB.Exec(ctx, query, shortCode)
	if err != nil {
		return fmt.Errorf("error deleting URL: %v", err)
	}
	return nil
}
