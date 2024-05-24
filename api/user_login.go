package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	kiteauth "github.com/nsvirk/gokiteauth"
)

// userLoginResponse is the response for the user login handler
type userLoginResponse struct {
	UserId    string `json:"user_id"`
	Enctoken  string `json:"enctoken"`
	LoginTime string `json:"login_time"`
}

// UserLoginHandler logs in a user
func UserLoginHandler(c echo.Context) error {

	// get form inputs
	var userId, password, totpSecret string

	err := GetUserFormUserId(c, &userId)
	if err != nil {
		return SendError(c, http.StatusBadRequest, err.Error())
	}

	err = GetUserFormPassword(c, &password)
	if err != nil {
		return SendError(c, http.StatusBadRequest, err.Error())
	}

	err = GetUserFormTotpSecret(c, &totpSecret)
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

	// Create a new Kite connect instance
	ka := kiteauth.New(userId)
	// kc.SetDebug(true)

	// Check if the enctoken is valid
	ka.SetBaseURI("https://kite.zerodha.com/oms")
	tokenValid, err := ka.CheckEnctokenValid(user.Enctoken)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, err.Error())
	}

	if tokenValid {
		// enctoken is valid send response
		data := userLoginResponse{
			UserId:    user.UserId,
			Enctoken:  user.Enctoken,
			LoginTime: user.LoginTime,
		}
		return SendResponse(c, http.StatusOK, data)

		// enctoken is not valid sget a new session
	} else {
		// generate a new kite session
		ka.SetBaseURI("https://kite.zerodha.com")
		session, err := ka.GenerateSession(password, totpSecret)
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
		// send the response
		data := userLoginResponse{
			UserId:    user.UserId,
			Enctoken:  user.Enctoken,
			LoginTime: user.LoginTime,
		}
		return SendResponse(c, http.StatusOK, data)
	}

}
