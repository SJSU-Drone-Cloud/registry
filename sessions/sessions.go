package sessions

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("top-secret"))

func GetSession(w http.ResponseWriter, r *http.Request, name string, userID int64) {
	session, _ := Store.Get(r, name)
	session.Values["user_id"] = userID
	session.Save(r, w)
}

func EndSession(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session")
	delete(session.Values, "user_id")
	session.Save(r, w)
}
