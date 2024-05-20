package handler

import "net/http"

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("dummy"))
}
