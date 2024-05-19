package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// userDeleteResponse is the response for the user delete handler
type userDeleteResponse struct {
	UserId    string `json:"user_id"`
	DeletedAt string `json:"deleted_at"`
}

// UserDeleteHandler deletes a user
func UserDeleteHandler(c echo.Context) error {

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
		return SendError(c, http.StatusBadRequest, err.Error())
	}

	// delete the user
	err = db.Delete(&user).Error
	// send error if any
	if err != nil {
		return SendError(c, http.StatusBadRequest, err.Error())
	}

	// close the database connection
	defer CloseDB(db)

	// send the response
	data := userDeleteResponse{
		UserId:    userId,
		DeletedAt: time.Now().String(),
	}

	return SendResponse(c, http.StatusOK, data)
}
