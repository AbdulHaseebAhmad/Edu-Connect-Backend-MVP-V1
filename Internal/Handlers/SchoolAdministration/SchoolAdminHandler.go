package SchoolAdminHandler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage"
	response "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Responses"
)

func LinkValidation(storage Storage.SchoolAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			slog.Info("request Error", "error", "token is missing in request")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("token is invalid")))
			return
		}

		status, err := storage.ValidateLink(r.Context(), token)
		if err != nil {
			slog.Info("Authorization Error", "error", "token is invalid/expired")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusAccepted, status)
	}
}
