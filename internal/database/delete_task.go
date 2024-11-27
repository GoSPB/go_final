package database

import "database/sql"

func DeleteTask(db *sql.DB, id string) error {
	query := "DELETE FROM scheduler WHERE id = :id"
	_, err := db.Query(query, sql.Named("id", id))
	if err != nil {
		return err
	}
	return nil
}
