package api

import "net/http"

func login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Login"))
}
