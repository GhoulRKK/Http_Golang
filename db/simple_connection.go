package db

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Simple_connection(ctx context.Context) (*pgxpool.Pool, error) {
	conString := os.Getenv("CONN_STRING")
	pool, err := pgxpool.New(ctx, conString)
	return pool, err
}

func Create_DB(pool *pgxpool.Pool, ctx context.Context) error {
	sql_string := `
	CREATE TABLE IF NOT EXISTS notes(
	ID SERIAL PRIMARY KEY,
	TITLE TEXT NOT NULL,
	CONTENT TEXT NOT NULL,
	CREATED_AT TIMESTAMP NOT NULL);
	`
	_, err := pool.Exec(ctx, sql_string)
	return err
}

type Notes struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func Create_notes(pool *pgxpool.Pool, ctx context.Context, note Notes) error {
	sql_string := `
	INSERT INTO notes(TITLE, CONTENT, CREATED_AT) VALUES ($1, $2, $3);
	`
	_, err := pool.Exec(ctx, sql_string, note.Title, note.Content, time.Now())
	return err
}

func CheclAll(pool *pgxpool.Pool, ctx context.Context) ([]Notes, error) {
	array := make([]Notes, 0, 10)
	sqlString := `
	SELECT ID, TITLE, CONTENT, CREATED_AT FROM notes;
	`
	rows, err := pool.Query(ctx, sqlString)
	if err != nil {
		return array, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var title, content string
		var createdAT time.Time
		if err := rows.Scan(&id, &title, &content, &createdAT); err != nil {
			return array, err
		}
		array = append(array, Notes{id, title, content, createdAT})
	}
	return array, rows.Err()
}

func GetByID(pool *pgxpool.Pool, ctx context.Context, idNote int) (Notes, error) {
	sqlString := `
	SELECT ID, TITLE, CONTENT, CREATED_AT FROM notes
	WHERE ID = $1;
	`
	row := pool.QueryRow(ctx, sqlString, idNote)
	var id int
	var title, content string
	var createdAT time.Time
	err := row.Scan(&id, &title, &content, &createdAT)
	return Notes{id, title, content, createdAT}, err
}

func UpdateByID(pool *pgxpool.Pool, ctx context.Context, updateNote Notes) error {
	sqlString := `
	UPDATE notes
	SET TITLE = $1, CONTENT = $2, CREATED_AT = $3
	WHERE ID = $4;
	`
	_, err := pool.Exec(ctx, sqlString, updateNote.Title, updateNote.Content, time.Now(), updateNote.Id)
	return err
}

func DeleteById(pool *pgxpool.Pool, ctx context.Context, id int) error {
	sqlString := `
	DELETE FROM notes
	WHERE ID = $1;
	`
	_, err := pool.Exec(ctx, sqlString, id)
	return err
}
