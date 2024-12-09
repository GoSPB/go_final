package database

import (
	"database/sql"
	"errors"
)

func DeleteTask(db *sql.DB, id string) error {
	query := "DELETE FROM scheduler WHERE id = ?"
	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("нет такой задачи")
	}
	return nil
}
