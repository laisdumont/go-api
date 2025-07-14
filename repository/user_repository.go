package repository

import (
	"database/sql"
	"go-api/model"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *model.User) error {
	_, err := r.DB.Exec("INSERT INTO users(name, password) VALUES(?, ?)", user.Name, user.Password)
	return err
}

func (r *UserRepository) GetAll() ([]model.User, error) {
	rows, err := r.DB.Query("SELECT id, name FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		rows.Scan(&u.ID, &u.Name)
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) Update(user *model.User) error {
	_, err := r.DB.Exec("UPDATE users SET name = ? WHERE id = ?", user.Name, user.ID)
	return err
}

func (r *UserRepository) Delete(id int) error {
	_, err := r.DB.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func (r *UserRepository) FindByName(name string) (*model.User, error) {
	row := r.DB.QueryRow("SELECT id, name, password FROM users WHERE name = ?", name)
	var u model.User
	err := row.Scan(&u.ID, &u.Name, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
