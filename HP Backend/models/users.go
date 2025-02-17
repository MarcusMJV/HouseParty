package models

import (
	"database/sql"
	"errors"

	"houseparty.com/storage"
	"houseparty.com/utils"
)

type User struct {
	Id       int64  `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

type UserResponse struct {
    Id       int64  `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

func (u *User) Save() error {
	stmt, err := storage.DB.Prepare(storage.SaveUserQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	HashPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Email, HashPassword, u.Username)
	if err != nil {
		return err
	}

	u.Id, _ = result.LastInsertId()

	return nil
}

func (u *User) GetUserById(id int64) error {
	row := storage.DB.QueryRow(storage.GetUserByIdQuery, id)

	err := row.Scan(&u.Id, &u.Username, &u.Email)

	if err == sql.ErrNoRows {
		return errors.New("user not found")
	} else if err != nil {
		return err
	}

	return nil
}

func (u *User) RetrieveHashPassword() (string, error) {
	row := storage.DB.QueryRow(storage.RetrieveHashPasswordQuery, u.Email, u.Username)

	var HashPassword string
	err := row.Scan(&u.Id, &u.Username, &u.Email, &HashPassword)

	if err == sql.ErrNoRows{
		return "", errors.New("there is no user with these credentials")
	} else if err != nil {
		return "", err
	}

	return HashPassword, nil
}

func (u *User) ToUserResponse() *UserResponse {
    return &UserResponse{
        Id:       u.Id,
        Username: u.Username,
        Email:    u.Email,
    }
}
