package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	repo "github.com/moemoe89/go-unit-test-sql/repository"

	_ "github.com/go-sql-driver/mysql"
)

type repository struct {
	db *sql.DB
}

var temp = &repo.UserModel{
	ID:    uuid.New().String(),
	Name:  "sanjai",
	Email: "sanjai@mail.com",
	Phone: "9488900582",
}

// NewRepository will create a variable that represent the Repository struct
func NewRepository(dialect, dsn string, idleConn, maxConn int) (repo.Repository, error) {
	db, err := sql.Open(dialect, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(idleConn)
	db.SetMaxOpenConns(maxConn)

	return &repository{db}, nil
}

// Close the connection
func (r *repository) Close() {
	r.db.Close()
}

// FindByID attaches the user repository and find data based on id
func (r *repository) FindByID(id string) (*repo.UserModel, error) {
	user := new(repo.UserModel)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, "SELECT id, name, email, phone FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email, &user.Phone)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Find attaches the user repository and find all data
func (r *repository) Find() ([]*repo.UserModel, error) {
	users := make([]*repo.UserModel, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 5%time.Second)
	defer cancel()
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, email, phone FROM users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		user := new(repo.UserModel)
		rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
		)

		users = append(users, user)
	}

	return users, nil
}

// Create attaches the user repository and creating the data
func (r *repository) Create(user *repo.UserModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO users (id, name, email, phone) VALUES (?, ?, ?, ?)"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.ID, user.Name, user.Email, user.Phone)
	return err
}

// Update attaches the user repository and update data based on id
func (r *repository) Update(user *repo.UserModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE users SET name = ?, email = ?, phone = ? WHERE id = ?"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.Name, user.Email, user.Phone, user.ID)
	return err
}

// Delete attaches the user repository and delete data based on id
func (r *repository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM users WHERE id = ?"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		err = errors.New("querry mismatch error")
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	fmt.Println(err)
	// err = errors.New("some error")
	return err
}

func main() {
	var repos *repository
	// repos.Delete("1")
	err := repos.Delete(temp.ID)
	fmt.Println(err)
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// query := "DELETE FROM users WHRE id = ?"
	// stmt, err := repo.db.PrepareContext(ctx, query)
	// if err != nil {
	// 	return err
	// }
	// defer stmt.Close()

	// _, err = stmt.ExecContext(ctx, temp.ID)
	// fmt.Println(err)
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}
