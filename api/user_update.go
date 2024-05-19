package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// userUpdateResponse is the response for the user update handler
type userUpdateResponse struct {
	UserId    string `json:"user_id"`
	UpdatedAt string `json:"updated_at"`
}

// UserUpdateHandler updates a user
func UserUpdateHandler(c echo.Context) error {

	// get the form inputs
	userId := c.FormValue("user_id")
	if userId == "" {
		return SendError(c, http.StatusBadRequest, "`user_id` is required")
	}

	password := c.FormValue("password")
	if password == "" {
		return SendError(c, http.StatusBadRequest, "`password` is required")
	}

	newPassword := c.FormValue("new_password")
	if newPassword == "" {
		return SendError(c, http.StatusBadRequest, "`new_password` is required")
	}

	// initialize the database connection
	db := ConnectToDB()

	// check if the user exists or not
	var user = UserModel{}
	var passwordHash = generateMD5Hash(password)
	err := db.Where("user_id = ? and password_hash = ? ", userId, passwordHash).First(&user).Error

	if err != nil {
		// send error, if record not found
		if err.Error() == "record not found" {
			err = fmt.Errorf("user `%s` not found", userId)
		}
		return SendError(c, http.StatusBadRequest, err.Error())
	}

	// update the user password
	user.PasswordHash = generateMD5Hash(newPassword)
	err = db.Save(&user).Error

	// send error if any
	if err != nil {
		return SendError(c, http.StatusBadRequest, err.Error())
	}

	// close the database connection
	defer CloseDB(db)

	// send the response
	data := userUpdateResponse{
		UserId:    user.UserId,
		UpdatedAt: user.UpdatedAt.String(),
	}

	return SendResponse(c, http.StatusOK, data)
}
