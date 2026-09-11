package postgres

import (
	"encoding/json"
	"log/slog"
	. "myapp/db_operations"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type Postgres struct {
	Conn *pgx.Conn
}

type ChangeInfo struct {
	Id      int    `json:"id"`
	NewInfo string `json:"newInfo"`
}

func (p *Postgres) Adduser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var u *User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		slog.Error("Не удалось разобрать тело запроса", "error", err)
		w.WriteHeader(500)
		return
	}
	if u.Info == "" || u.Name == "" {
		slog.Warn("Пустые поля name или info", "name", u.Name, "info", u.Info)
		w.WriteHeader(400)
		return
	}

	info, err := Insert(ctx, p.Conn, u.Name, u.Info)
	if err != nil {
		slog.Error("Не удалось добавить пользователя в БД", "error", err, "name", u.Name)
		w.WriteHeader(500)
		return
	}
	slog.Info("Пользователь добавлен в БД", "name", u.Name, "result", info)
}

func (p *Postgres) RemoveUserById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var id map[string]int
	err := json.NewDecoder(r.Body).Decode(&id)
	if err != nil {
		slog.Error("Не удалось разобрать тело запроса", "error", err)
		w.WriteHeader(400)
		return
	}

	info, err := Delete(ctx, p.Conn, id["id"])
	if err == ErrNotFound {
		slog.Warn("Пользователь не найден для удаления", "id", id["id"])
		w.WriteHeader(404)
		return
	}
	if err != nil {
		slog.Error("Не удалось удалить пользователя", "error", err, "id", id["id"])
		w.WriteHeader(500)
		return
	}
	slog.Info("Пользователь удалён из БД", "id", id["id"], "result", info)
}

func (p *Postgres) ChangeInfoById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var data ChangeInfo
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		slog.Error("Не удалось разобрать тело запроса", "error", err)
		w.WriteHeader(400)
		return
	}

	info, err := UpdateByID(ctx, p.Conn, data.Id, data.NewInfo)
	if err == ErrNotFound {
		slog.Warn("Пользователь не найден для обновления", "id", data.Id)
		w.WriteHeader(404)
		return
	}
	if err != nil {
		slog.Error("Не удалось обновить информацию пользователя", "error", err, "id", data.Id)
		w.WriteHeader(500)
		return
	}
	slog.Info("Информация о пользователе обновлена", "id", data.Id, "result", info)
}

func (p *Postgres) GetUserById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data := mux.Vars(r)["id"]

	id, err := strconv.Atoi(data)
	if err != nil {
		slog.Warn("Неверный формат id в запросе", "id", data)
		w.WriteHeader(400)
		w.Write([]byte("Неверный формат данных"))
		return
	}
	user, err := SelectByID(ctx, p.Conn, id)
	if err == ErrNotFound {
		slog.Warn("Пользователь не найден", "id", id)
		w.WriteHeader(404)
		return
	}
	if err != nil {
		slog.Error("Не удалось получить пользователя", "error", err, "id", id)
		w.WriteHeader(500)
		return
	}

	response, err := json.MarshalIndent(user, "", " ")
	if err != nil {
		slog.Error("Не удалось сериализовать ответ в JSON", "error", err, "id", id)
		w.WriteHeader(500)
		return
	}
	w.Write([]byte(response))
	slog.Info("Отправлена информация о пользователе", "id", id)
}
