package main

import (
	"context"
	. "myapp/connection"
	. "myapp/db_operations"
	. "myapp/postgres"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	conn := GetConn()
	router := mux.NewRouter()
	database := &Postgres{Conn: conn}
	ctx := context.Background()
	InitTable(ctx, conn)

	router.Path("/users").Methods("POST").Handler(
		InstrumentHandler(database.Adduser, "POST", "/users"),
	)
	router.Path("/users").Methods("DELETE").Handler(
		InstrumentHandler(database.RemoveUserById, "DELETE", "/users"),
	)
	router.Path("/users/{id}").Methods("GET").Handler(
		InstrumentHandler(database.GetUserById, "GET", "/users/{id}"),
	)
	router.Path("/users").Methods("PATCH").Handler(
		InstrumentHandler(database.ChangeInfoById, "PATCH", "/users"),
	)

	router.Path("/metrics").Methods("GET").Handler(promhttp.Handler())

	http.ListenAndServe(":9055", router)
}
