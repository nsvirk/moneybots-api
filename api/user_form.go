package api

import (
	"errors"

	"github.com/labstack/echo/v4"
)

// GetUserFormUserId gets the user_id from the form
func GetUserFormUserId(c echo.Context, userId *string) error {
	*userId = c.FormValue("user_id")
	if *userId == "" {
		return errors.New("`user_id` is required")
	}
	if len(*userId) != 6 {
		return errors.New("`user_id` should be 6 characters")
	}
	return nil
}

// GetUserFormPassword gets the password from the form
func GetUserFormPassword(c echo.Context, password *string) error {
	*password = c.FormValue("password")
	if *password == "" {
		return errors.New("`password` is required")
	}
	return nil
}

// GetUserFormTotpSecret gets the totp_secret from the form
func GetUserFormTotpSecret(c echo.Context, totpSecret *string) error {
	*totpSecret = c.FormValue("totp_secret")
	if *totpSecret == "" {
		return errors.New("`totp_secret` is required")
	}
	if len(*totpSecret) != 32 {
		return errors.New("`totp_secret` should be 32 characters")
	}
	return nil
}

// GetUserFormNewPassword gets the password from the form
func GetUserFormNewPassword(c echo.Context, newPassword *string) error {
	*newPassword = c.FormValue("new_password")
	if *newPassword == "" {
		return errors.New("`new_password` is required")
	}
	return nil
}
