package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/learningspoons-go/groupchat/db"
	"github.com/learningspoons-go/groupchat/token"
	"github.com/learningspoons-go/groupchat/ws"
)

// Handler api handler
func Handler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/signup", signup).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/signin", signin).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/room", createRoom).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/rooms", getRooms).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/room/{id}", connectToRoom).Methods(http.MethodGet, http.MethodOptions)

	router.Use(handlePanic)
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(w, r)
		})
	})
	router.Use(mux.CORSMethodMiddleware(router))

	return router
}

type usernamePassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func signup(w http.ResponseWriter, r *http.Request) {
	var req usernamePassword
	parseJSON(r.Body, &req)
	id, err := db.CreateUser(req.Username, req.Password)
	must(err)

	t, err := token.New(id)
	must(err)
	writeJSON(w, struct {
		Token string `json:"token"`
	}{t})
}

func signin(w http.ResponseWriter, r *http.Request) {
	var req usernamePassword
	parseJSON(r.Body, &req)
	id, err := db.FindUser(req.Username, req.Password)
	must(err)

	t, err := token.New(id)
	must(err)
	writeJSON(w, struct {
		Token string `json:"token"`
	}{t})
}

func createRoom(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	parseJSON(r.Body, &req)

	id, err := db.CreateRoom(req.Name)
	must(err)
	writeJSON(w, struct {
		ID int `json:"id"`
	}{id})
}

func getRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := db.GetRooms()
	must(err)
	writeJSON(w, rooms)
}

func connectToRoom(w http.ResponseWriter, r *http.Request) {
	uid := userID(r)
	roomID := parseIntParam(r, "id")
	exists, err := db.RoomExists(roomID)
	must(err)

	if !exists {
		panic(notFoundError)
	}

	ws.ChatHandler(roomID, uid).ServeHTTP(w, r)
}
