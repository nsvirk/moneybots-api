package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	kitesession "github.com/nsvirk/gokitesession"
)

// userLoginResponse is the response for the user login handler
type userLoginResponse struct {
	// Enctoken  string `json:"enctoken"`
	UserId    string `json:"user_id"`
	Enctoken  string `json:"enctoken"`
	Token     string `json:"token"`
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

	// Create a new Kite session instance
	ks := kitesession.New(userId)
	ks.SetDebug(true)

	// Check if the enctoken is valid
	tokenValid, err := ks.CheckEnctokenValid(user.Enctoken)
	if err != nil {
		return SendError(c, http.StatusUnauthorized, err.Error())
	}

	// if enctoken is not valid get a new session
	var session *kitesession.Session
	if !tokenValid {
		// generate totp value
		totpValue, err := ks.GenerateTotpValue(totpSecret)
		if err != nil {
			return SendError(c, http.StatusUnauthorized, err.Error())
		}
		// generate a new kite session
		session, err = ks.GenerateSession(password, totpValue)
		if err != nil {
			return SendError(c, http.StatusUnauthorized, err.Error())
		}

		//check if session is created
		if session.UserId == "" || session.Enctoken == "" {
			return SendError(c, http.StatusInternalServerError, "session not created")
		}

		// update the user enctoken and login_time
		user.Enctoken = session.Enctoken
		user.LoginTime = session.LoginTime
		err = db.Save(&user).Error
		if err != nil {
			return SendError(c, http.StatusInternalServerError, err.Error())
		}
	}

	// close the database connection
	CloseDB(db)

	// create jwt
	jwtToken, err := GenerateJWT(user.UserId, user.Enctoken)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, err.Error())
	}

	// send the response
	data := userLoginResponse{
		UserId:    user.UserId,
		Enctoken:  user.Enctoken,
		Token:     jwtToken,
		LoginTime: user.LoginTime,
	}

	return SendResponse(c, http.StatusOK, data)

}
