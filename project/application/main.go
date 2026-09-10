package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	. "myapp/connection"
	. "myapp/db_operations"
	. "myapp/postgres"

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

	router.Path("/health").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := conn.Ping(r.Context()); err != nil {
			w.WriteHeader(503)
			w.Write([]byte("БД НЕДОСТУПНА!"))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	srv := &http.Server{
		Addr:    ":9055",
		Handler: router,
	}

	go func() {
		log.Println("Запуск сервера на порту :9055")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Остановка сервера...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Сервер принудительно остановлен: %v", err)
	}

	if err := conn.Close(shutdownCtx); err != nil {
		log.Printf("Ошибка закрытия соединения с БД: %v", err)
	}

	log.Println("Сервер завершил работу корректно")
}
