package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	kiteconnect "github.com/nsvirk/gokiteconnect/v4"
)

// UserProfileHandler gets the user profile
func UserProfileHandler(c echo.Context) error {

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
		return SendError(c, http.StatusUnauthorized, err.Error())
	}

	// close the database connection
	defer CloseDB(db)

	// check if user logged in
	if user.Enctoken == "" {
		err = fmt.Errorf("user `%s` not logged in", userId)
		return SendError(c, http.StatusBadRequest, err.Error())
	}

	// get the user profile
	// Create a new Kite connect instance
	kc := kiteconnect.New(userId)
	// kc.SetDebug(true)

	// get profile
	kc.SetBaseURI("https://kite.zerodha.com/oms")
	kc.SetEnctoken(user.Enctoken)

	profile, err := kc.GetUserProfile()
	// send error if any
	if err != nil {
		return SendError(c, http.StatusBadRequest, err.Error())
	}

	// send the response
	return SendResponse(c, http.StatusOK, profile)
}
