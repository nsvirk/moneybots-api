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

	// get the form inputs
	userId := c.FormValue("user_id")
	if userId == "" {
		return SendError(c, http.StatusBadRequest, "`user_id` is required")
	}

	password := c.FormValue("password")
	if password == "" {
		return SendError(c, http.StatusBadRequest, "`password` is required")
	}

	// initialize the database connection
	db := ConnectToDB()

	// check if the user exists or not
	var user = UserModel{}
	var passwordHash = generateMD5Hash(password)
	err := db.Where("user_id = ? and password_hash = ? ", userId, passwordHash).First(&user).Error
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
		CreatedAt: user.CreatedAt.String(),
	}

	return SendResponse(c, http.StatusCreated, data)
}
