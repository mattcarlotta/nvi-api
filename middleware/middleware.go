package middleware

import (
	"context"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/mattcarlotta/nvi-api/utils"
)

func CookieSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("SESSION_TOKEN")
		if err != nil {
			utils.SendErrorResponse(res, http.StatusUnauthorized, "You must be logged in order to do that!")
			return
		} else if len(cookie.Value) == 0 {
			utils.SendErrorResponse(res, http.StatusUnauthorized, "You must be logged in order to do that!")
			return
		}

		data, err := utils.ValidateSessionToken(cookie.Value)
		if err != nil {
			utils.SendErrorResponse(res, http.StatusUnauthorized, err.Error())
			return
		}
		ctx := context.WithValue(req.Context(), "userSessionId", data.UserId)

		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		origin := req.Header.Get("Origin")
		var CLIENT_ORIGIN = utils.GetEnv("CLIENT_HOST")
		if origin == CLIENT_ORIGIN {
			res.Header().Set("Access-Control-Allow-Origin", CLIENT_ORIGIN)
			res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			res.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			res.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		next.ServeHTTP(res, req)
	})
}

func Logging(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}
