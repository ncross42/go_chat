package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/learningspoons-go/groupchat/db"
	"github.com/learningspoons-go/groupchat/token"
)

func must(err error) {
	if err == db.ErrUnauthorized {
		panic(unauthorizedError)
	} else if err != nil {
		log.Println("internal error:", err)
		panic(internalError)
	}
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	must(json.NewEncoder(w).Encode(v))
}

func parseJSON(r io.Reader, v interface{}) {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		log.Println("parsing json body:", err)
		panic(malformedInputError)
	}
}

func parseIntParam(r *http.Request, key string) int {
	vars := mux.Vars(r)
	if v, ok := vars[key]; ok {
		i, err := strconv.Atoi(v)
		if err == nil {
			return i
		}
	}
	panic(malformedInputError)
}

func userID(r *http.Request) int {
	t := r.URL.Query().Get("token")
	id, err := token.Parse(t)
	if err != nil {
		log.Println(err)
		panic(unauthorizedError)
	}
	return id
}
