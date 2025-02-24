package services

import (
	"database/sql"
	"errors"

	"houseparty.com/models"
	"houseparty.com/storage"
	"houseparty.com/utils"
)

func CreateNewUser(user *models.User) (string, error) {

	isInUse, errorString := checkInUse(user.Username, user.Email)

	if isInUse {
		err := errors.New(errorString)
		return "", err
	}

	err := user.Save()
	if err != nil {
		return "", err
	}

	return utils.GenerateToken(user.Email, user.Username, user.Id)
}

func ValidateCredentials(user *models.User) (string, error){
	hashPassword, err :=  user.RetrieveHashPassword()
	if err != nil{
		return "", err
	}
	isValidPassword := utils.CheckPasswordHash(user.Password, hashPassword)

	if !isValidPassword{
		return "", errors.New("password that is provided is incorrect")
	}

	err = user.GetUserById(user.Id)
	if err != nil{
		return "", err
	}

	return utils.GenerateToken(user.Email, user.Username, user.Id)
}

func checkInUse(username, email string) (bool, string) { 
	var conflictField string
	err := storage.DB.QueryRow(storage.CheckInUseQuery, username, email, username, email).Scan(&conflictField)

	if err != nil && err != sql.ErrNoRows {
		return true, "Database Error"
	} else if conflictField == "email"{
		return true, "Email already in use."
	}else if conflictField == "username"{
		return true, "Username already in use."
	}

	return false, ""
}

