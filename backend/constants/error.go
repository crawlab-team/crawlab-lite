package constants

import (
	"crawlab-lite/error"
	"net/http"
)

var (
	ErrorMongoError = error.NewSystemOPError(1001, "system error:[mongo]%s", http.StatusInternalServerError)
	//users
	ErrorUserNotFound              = error.NewBusinessError(10001, "user not found.", http.StatusUnauthorized)
	ErrorUsernameOrPasswordInvalid = error.NewBusinessError(11001, "username or password invalid", http.StatusUnauthorized)
)
