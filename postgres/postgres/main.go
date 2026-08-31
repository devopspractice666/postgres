package main

import (
	"context"
	. "myapp/connection"
	. "myapp/db_operations"
	. "myapp/postgres"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	conn := GetConn()
	router := mux.NewRouter()
	database := &Postgres{Conn: conn}
	ctx := context.Background()
	InitTable(ctx, conn)
	router.Path("/users").Methods("POST").HandlerFunc(database.Adduser)
	router.Path("/users").Methods("DELETE").HandlerFunc(database.RemoveUserById)
	router.Path("/users/{id}").Methods("GET").HandlerFunc(database.GetUserById)
	router.Path("/users").Methods("PATCH").HandlerFunc(database.ChangeInfoById)
	http.ListenAndServe(":9055", router)

}
