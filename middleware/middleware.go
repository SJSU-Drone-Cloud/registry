package middleware

import (
	"net/http"

	"github.com/QianMason/drone-backend/sessions"
)

func AuthRequired(handler http.HandlerFunc) http.HandlerFunc { //middleware func
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessions.Store.Get(r, "session")
		_, ok := session.Values["user_id"]
		if !ok {
			http.Redirect(w, r, "/login", 302)
			return
		}
		handler.ServeHTTP(w, r)
	}
}
