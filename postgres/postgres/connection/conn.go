package connection

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

func GetConn() *pgx.Conn {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	bd := os.Getenv("bd")
	conn, err := pgx.Connect(ctx, bd)
	if err != nil {
		panic(err)
	}

	if err := conn.Ping(ctx); err != nil {
		panic(err)
	}

	fmt.Println("Подключение к бд успешно")
	return conn
}
