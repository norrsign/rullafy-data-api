package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/norrsign/rullafy-data-api/db"
	"github.com/vanern/goapi/framework"
	"github.com/vanern/goapi/typs"
)

func Init() {

	// AUTH but no roles
	framework.GET("/me", func(w http.ResponseWriter, r *http.Request) {
		u, ok := typs.UserFromContext(r.Context())
		if !ok {
			http.Error(w, "no user in context", http.StatusUnauthorized)
			return
		}

		user, err := db.Qrs.GetUser(context.Background(), u.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// serialize to JSON and write to response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	}, typs.AuthOnly)

	// AUTH but no roles
	framework.POST("/register", func(w http.ResponseWriter, r *http.Request) {
		u, ok := typs.UserFromContext(r.Context())

		if !ok {
			http.Error(w, "no user in context", http.StatusUnauthorized)
			return
		}
		// get the payload from the request body
		var user db.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "invalid request payload", http.StatusBadRequest)
			return
		}

		fmt.Println(u.ID)
		user.ID = u.ID
		fmt.Println(user.ID, user.Job)
		if _, err := db.Qrs.CreateUser(context.Background(), user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// serialize to JSON and write to response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	}, typs.AuthOnly)

}
