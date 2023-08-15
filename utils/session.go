package utils

import "net/http"

func GetUserSessionId(res http.ResponseWriter, req *http.Request) string {
	userSessionId, ok := req.Context().Value("userSessionId").(string)
	if !ok {
		SendErrorResponse(
			res,
			http.StatusInternalServerError,
			"Encountered an unexpected error when attempting to access the user session id.",
		)
	}

	return userSessionId
}
