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

	// get form inputs
	var userId, password, newPassword string

	err := GetUserFormUserId(c, &userId)
	if err != nil {
		return SendError(c, http.StatusBadRequest, err.Error())
	}

	err = GetUserFormPassword(c, &password)
	if err != nil {
		return SendError(c, http.StatusBadRequest, err.Error())
	}

	err = GetUserFormNewPassword(c, &newPassword)
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
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return SendResponse(c, http.StatusOK, data)
}
