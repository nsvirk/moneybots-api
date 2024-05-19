package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// userLogoutResponse is the response for the user logout handler
type userLogoutResponse struct {
	UserId     string `json:"user_id"`
	LogoutTime string `json:"logout_time"`
}

// UserLogutHandler logs out a user
func UserLogutHandler(c echo.Context) error {

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
		// send error, if record not found
		if err.Error() == "record not found" {
			err = fmt.Errorf("user `%s` not found", userId)
		}
		return SendError(c, http.StatusUnauthorized, err.Error())
	}

	// update the enctoken to blank
	if user.Enctoken == "" {
		err = fmt.Errorf("user `%s` not logged in", userId)
		return SendError(c, http.StatusBadRequest, err.Error())
	} else {
		user.Enctoken = ""
		err = db.Save(&user).Error
		// send error if any
		if err != nil {
			return SendError(c, http.StatusBadRequest, err.Error())
		}
	}

	// close the database connection
	defer CloseDB(db)

	// send the response
	data := userLogoutResponse{
		UserId:     user.UserId,
		LogoutTime: user.UpdatedAt.String(),
	}

	return SendResponse(c, http.StatusOK, data)
}
