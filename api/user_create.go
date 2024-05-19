package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// userCreateResponse is the response for the user create handler
type userCreateResponse struct {
	UserId    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

// UserCreateHandler creates a new user
func UserCreateHandler(c echo.Context) error {

	// get form inputs
	var userId, password string

	err := GetUserFormUserId(c, &userId)
	if err != nil {
		return SendError(c, http.StatusBadRequest, err.Error())
	}

	err = GetUserFormPassword(c, &password)
	if err != nil {
		return SendError(c, http.StatusBadRequest, err.Error())
	}

	// initialize the database connection
	db := ConnectToDB()

	// check if the user exists or not
	var user = UserModel{}
	var passwordHash = generateMD5Hash(password)
	err = db.Where("user_id = ? and password_hash = ? ", userId, passwordHash).First(&user).Error
	if err != nil {
		// send error only if err is not qual to "record not found"
		if err.Error() != "record not found" {
			return SendError(c, http.StatusBadRequest, err.Error())
		}
	}

	// if user already exists, send error
	if user.UserId == userId {
		err = fmt.Errorf("user `%s` already exists", userId)
		return SendError(c, http.StatusBadRequest, err.Error())
	}

	// create a new user
	user.UserId = userId
	user.PasswordHash = generateMD5Hash(password)
	err = db.Create(&user).Error

	// send error if any
	if err != nil {
		return SendError(c, http.StatusBadRequest, err.Error())
	}

	// send the response
	data := userCreateResponse{
		UserId:    user.UserId,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return SendResponse(c, http.StatusCreated, data)
}
