package simple

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
)

var ErrNotFound = errors.New("Указанный ресурс не найден")

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Info string `json:"info"`
}

func InitTable(ctx context.Context, conn *pgx.Conn) error {
	sqlQUery := `
	create table if not exists users(
	 	id SERIAL PRIMARY KEY,
        name Varchar(100) unique,
        info varchar(200)
	)
	`
	_, err := conn.Exec(ctx, sqlQUery)
	if err != nil {
		slog.Error("Не удалось создать таблицу users", "error", err)
		return err
	}
	slog.Info("Таблица users инициализирована успешно")
	return nil
}

func Delete(ctx context.Context, conn *pgx.Conn, id int) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlQuery := `
	DELETE FROM users
	where id=$1;
	`
	reply, err := conn.Exec(ctx, sqlQuery, id)
	if err != nil {
		slog.Error("Ошибка выполнения DELETE", "error", err, "id", id)
		return "", err
	}
	str := reply.String()
	return str, nil
}

func Insert(ctx context.Context, conn *pgx.Conn, name, info string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlQuery := `
	insert into users (name,info)
	values ($1,$2);
	`
	reply, err := conn.Exec(ctx, sqlQuery, name, info)
	if err != nil {
		slog.Error("Ошибка выполнения INSERT", "error", err, "name", name)
		return "ОШИБКА", err
	}
	str := reply.String()
	return str, nil
}

func UpdateByID(ctx context.Context, conn *pgx.Conn, id int, newInfo string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlQuery := `
	UPDATE users
	SET  info = $2
	where id = $1
	`
	reply, err := conn.Exec(ctx, sqlQuery, id, newInfo)
	if err != nil {
		slog.Error("Ошибка выполнения UPDATE", "error", err, "id", id)
		return "", err
	}
	if reply.RowsAffected() == 0 {
		return "", ErrNotFound
	}
	str := reply.String()
	return str, nil
}

func SelectByID(ctx context.Context, conn *pgx.Conn, id int) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var user User
	sqlQuery := `
	Select id, name, info from users
	where id=$1
	`
	rows, err := conn.Query(ctx, sqlQuery, id)
	if err != nil {
		slog.Error("Ошибка выполнения SELECT", "error", err, "id", id)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&user.ID, &user.Name, &user.Info)
	}

	if rows.CommandTag().RowsAffected() == 0 {
		return nil, ErrNotFound
	}
	return &user, nil
}
