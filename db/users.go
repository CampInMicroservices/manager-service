package db

import (
	"context"
	"time"
)

type User struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	Activated bool      `json:"activated" db:"activated"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateUserParam struct {
	Name      string
	Email     string
	Password  string
	Activated bool
}

type UpdateUserParam struct {
	Name      string
	Email     string
	Password  string
	Activated bool
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
	INSERT INTO "users" ("name", "email", "password", "activated")
	VALUES ($1, $2, $3, $4)
	RETURNING "id", "name", "email", "password", "activated", "created_at"
	`
	row := store.db.QueryRowContext(ctx, query, arg.Name, arg.Email, arg.Password, arg.Activated)

	var user User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Activated,
		&user.CreatedAt,
	)

	return user, err
}
