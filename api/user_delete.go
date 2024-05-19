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
		DeletedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	return SendResponse(c, http.StatusOK, data)
}
