package postgres

import (
	"encoding/json"
	"fmt"
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
		fmt.Println("Ошибка при обработке входящего запроса")
		w.WriteHeader(500)
		return
	}
	if u.Info == "" || u.Name == "" {
		fmt.Println("Ошибка при обработке данных входящего запроса")
		w.WriteHeader(400)
		return
	}

	info, err := Insert(ctx, p.Conn, u.Name, u.Info)
	if err != nil {
		fmt.Println("Ошибка при добавлении пользователя в базу данных")
		w.WriteHeader(500)
		return
	}
	fmt.Println(info)
	fmt.Println("В базу данных добавлен пользователь с именем ", u.Name)
}

func (p *Postgres) RemoveUserById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var id map[string]int
	err := json.NewDecoder(r.Body).Decode(&id)
	if err != nil {
		fmt.Println("Неверный формат запроса")
		w.WriteHeader(400)
		return
	}

	info, err := Delete(ctx, p.Conn, id["id"])
	if err == ErrNotFound {
		w.WriteHeader(404)
		fmt.Println("Запрос на несуществующий ресурс")
		return
	}
	if err != nil {
		fmt.Println("Ошибка при удалении пользователя")
		w.WriteHeader(500)
		return
	}
	fmt.Println(info)
	fmt.Println("Из базы данных удален пользователь с id ", id["id"])
}

func (p *Postgres) ChangeInfoById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var data ChangeInfo
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		fmt.Println("Неверный формат запроса")
		w.WriteHeader(400)
		return
	}

	info, err := UpdateByID(ctx, p.Conn, data.Id, data.NewInfo)
	if err == ErrNotFound {
		w.WriteHeader(404)
		fmt.Println("Запрос на несуществующий ресурс")
		return
	}
	if err != nil {
		fmt.Println("Ошибка при обновлении информации пользователя")
		w.WriteHeader(500)
		return
	}
	fmt.Println(info)
	fmt.Println("Изменена информация пользователя с id ", data.Id)

}

func (p *Postgres) GetUserById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data := mux.Vars(r)["id"]

	id, err := strconv.Atoi(data)
	if err != nil {
		fmt.Println("Неверный формат запроса")
		w.WriteHeader(400)
		w.Write([]byte("Неверный формат данных"))
		return
	}
	user, err := SelectByID(ctx, p.Conn, id)
	if err == ErrNotFound {
		w.WriteHeader(404)
		fmt.Println("Отправлен запрос на несуществующий ресурс")
		return
	}
	if err != nil {
		w.WriteHeader(500)
		fmt.Println("Ошибка при форматировании ответа в json")
		return
	}

	response, err := json.MarshalIndent(user, "", " ")
	w.Write([]byte(response))
	fmt.Println("Отправлена информация о юзере с id: ", id)
}
