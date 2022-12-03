package db

import (
	"context"
	"time"
)

type User struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateUserParam struct {
	Name string
}

type UpdateUserParam struct {
	Name string
}

type ListUserParam struct {
	Offset int32
	Limit  int32
}

func (store *Store) GetUserByID(ctx context.Context, id int64) (user User, err error) {

	const query = `SELECT * FROM "users" WHERE "id" = $1`
	err = store.db.GetContext(ctx, &user, query, id)

	return
}

func (store *Store) CreateUser(ctx context.Context, arg CreateUserParam) (User, error) {

	const query = `
	INSERT INTO "users" ("name") 
	VALUES ($1)
	RETURNING "id", "name", "created_at"
	`
	row := store.db.QueryRowContext(ctx, query, arg.Name)

	var user User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.CreatedAt,
	)

	return user, err
}
