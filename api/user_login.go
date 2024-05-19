package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	kiteconnect "github.com/nsvirk/gokiteconnect/v4"
)

// userLoginResponse is the response for the user login handler
type userLoginResponse struct {
	UserId    string `json:"user_id"`
	Enctoken  string `json:"enctoken"`
	LoginTime string `json:"login_time"`
}

// UserLoginHandler logs in a user
func UserLoginHandler(c echo.Context) error {

	// get the form inputs
	userId := c.FormValue("user_id")
	if userId == "" {
		return SendError(c, http.StatusBadRequest, "`user_id` is required")
	}

	password := c.FormValue("password")
	if password == "" {
		return SendError(c, http.StatusBadRequest, "`password` is required")
	}

	totpSecret := c.FormValue("totp_secret")
	if totpSecret == "" {
		return SendError(c, http.StatusBadRequest, "`totp_secret` is required")
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

	// Create a new Kite connect instance
	kc := kiteconnect.New(userId)
	// kc.SetDebug(true)

	// if user exists then check if db session is valid
	kc.SetBaseURI("https://kite.zerodha.com/oms")
	tokenValid := kc.CheckEnctokenValid(user.Enctoken)

	// enctoken not valid, get a new one
	if !tokenValid {
		// generate a new kite session
		kc.SetBaseURI("https://kite.zerodha.com")
		session, err := kc.GenerateSession(password, totpSecret)
		if err != nil {
			return SendError(c, http.StatusUnauthorized, err.Error())
		}
		// update the user enctoken and login_time
		user.Enctoken = session.Enctoken
		user.LoginTime = session.LoginTime
		err = db.Save(&user).Error
		// close the database connection
		defer CloseDB(db)
		if err != nil {
			return SendError(c, http.StatusInternalServerError, err.Error())
		}
	} else {
		// close the database connection
		defer CloseDB(db)

	}

	// send the response
	data := userLoginResponse{
		UserId:    user.UserId,
		Enctoken:  user.Enctoken,
		LoginTime: user.LoginTime,
	}

	return SendResponse(c, http.StatusOK, data)
}
