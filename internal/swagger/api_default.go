/*
 * People info
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 0.0.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/ploschka/golang_task/internal/logger"
	m "github.com/ploschka/golang_task/internal/model"
	"github.com/ploschka/golang_task/internal/server"
	"gorm.io/gorm"
)

func InfoDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func InfoGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	serie := r.URL.Query().Get("PassSerie")
	number := r.URL.Query().Get("PassNum")

	log.Debug(serie, number)

	if len(serie) != 4 {
		server.BadRequest(w, server.ErrInvalidPassSerie)
		return
	}

	if len(number) != 6 {
		server.BadRequest(w, server.ErrInvalidPassNum)
		return
	}

	serie_int, err := strconv.ParseUint(serie, 10, 32)
	if err != nil {
		server.BadRequest(w, errors.Join(server.ErrInvalidPassSerie, err))
		return
	}

	number_int, err := strconv.ParseUint(number, 10, 32)
	if err != nil {
		server.BadRequest(w, errors.Join(server.ErrInvalidPassNum, err))
		return
	}

	user := m.User{PassSerie: uint32(serie_int), PassNumber: uint32(number_int)}

	db := m.GetDB()
	q := func(tx *gorm.DB) *gorm.DB {
		return tx.First(&user)
	}
	log.Debug(db.ToSQL(q))

	result := q(db)
	if result.Error != nil {
		server.BadRequest(w, errors.Join(server.ErrDatabase, result.Error))
		return
	}
	log.Debug(user)

	ret := People{
		Surname:    user.Surname,
		Name:       user.Name,
		Patronymic: user.Patronimic,
		Address:    user.Address,
	}
	str, err := json.Marshal(ret)
	if err != nil {
		server.InternalError(w, errors.Join(server.ErrJson, err))
		return
	}
	server.OK(w)
	w.Write(str)
}

func InfoListGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// body, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	server.InternalError(w, errors.Join(server.ErrBodyRead, err))
	// 	return
	// }
	user_query := People{}

	user_query.Name = r.URL.Query().Get("name")
	user_query.Surname = r.URL.Query().Get("surname")
	user_query.Patronymic = r.URL.Query().Get("patronymic")
	user_query.Address = r.URL.Query().Get("address")

	page, err := strconv.ParseUint(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		server.BadRequest(w, errors.Join(server.ErrInvalidPage, err))
		return
	}
	length, err := strconv.ParseUint(r.URL.Query().Get("len"), 10, 64)
	if err != nil {
		server.BadRequest(w, errors.Join(server.ErrInvalidLen, err))
		return
	}

	log.Debug(int(page), int(length))

	users := []m.User{}

	db := m.GetDB()
	q := func(tx *gorm.DB) *gorm.DB {
		return tx.Where(&user_query).Limit(int(length)).Offset(int((page - 1) * length)).Find(&users)
	}
	log.Debug(db.ToSQL(q))

	result := q(db)
	if result.Error != nil {
		server.BadRequest(w, errors.Join(server.ErrDatabase, result.Error))
		return
	}

	ret_users := []PeopleWithNumber{}

	for _, u := range users {
		ret_users = append(ret_users, PeopleWithNumber{
			Surname:    u.Surname,
			Name:       u.Name,
			Patronymic: u.Patronimic,
			Address:    u.Address,
			Passport: &Passport{
				PassSerie: int32(u.PassSerie),
				PassNum:   int32(u.PassNumber),
			},
		})
	}

	str, err := json.Marshal(ret_users)
	if err != nil {
		server.InternalError(w, errors.Join(server.ErrJson, err))
		return
	}
	server.OK(w)
	w.Write(str)
}

func InfoPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		server.BadRequest(w, errors.Join(server.ErrBodyRead, err))
		return
	}

	pass := PassportSerie{}
	err = json.Unmarshal(body, &pass)
	if err != nil {
		server.BadRequest(w, errors.Join(server.ErrJson, err))
		return
	}
	log.Debug(pass.PassportNumber)

	strs := strings.Split(pass.PassportNumber, " ")

	if len(strs[0]) != 4 {
		server.BadRequest(w, server.ErrInvalidPassSerie)
		return
	}

	if len(strs[1]) != 6 {
		server.BadRequest(w, server.ErrInvalidPassNum)
		return
	}

	user := m.User{}

	serie, err := strconv.ParseUint(strs[0], 10, 32)
	if err != nil {
		server.BadRequest(w, errors.Join(server.ErrInvalidPassSerie, err))
		return
	}

	number, err := strconv.ParseUint(strs[1], 10, 32)
	if err != nil {
		server.BadRequest(w, errors.Join(server.ErrInvalidPassNum, err))
		return
	}

	user.PassSerie = uint32(serie)
	user.PassNumber = uint32(number)

	db := m.GetDB()
	q := func(tx *gorm.DB) *gorm.DB {
		return tx.Create(&user)
	}
	log.Debug(db.ToSQL(q))

	result := q(db)
	if result.Error != nil {
		server.BadRequest(w, errors.Join(server.ErrDatabase, result.Error))
		return
	}
	server.OK(w)
}

func InfoPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		server.BadRequest(w, errors.Join(server.ErrBodyRead, err))
		return
	}

	user := PeopleWithNumber{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		server.BadRequest(w, errors.Join(server.ErrJson, err))
		return
	}
	log.Debug(user, user.Passport)

	query_user := m.User{}
	query_user.PassSerie = uint32(user.Passport.PassSerie)
	query_user.PassNumber = uint32(user.Passport.PassNum)

	log.Debug(query_user)

	db := m.GetDB()
	q1 := func(tx *gorm.DB) *gorm.DB {
		return tx.Where(&query_user).First(&query_user)
	}
	q2 := func(tx *gorm.DB) *gorm.DB {
		return tx.Save(&query_user)
	}
	log.Debug(db.ToSQL(q1))

	result := q1(db)
	if result.Error != nil {
		server.BadRequest(w, errors.Join(server.ErrDatabase, result.Error))
		return
	}
	log.Debug(query_user)

	query_user.Address = user.Address
	query_user.Name = user.Name
	query_user.Surname = user.Surname
	query_user.Patronimic = user.Patronymic

	log.Debug(db.ToSQL(q2))

	result = q2(db)
	if result.Error != nil {
		server.InternalError(w, errors.Join(server.ErrDatabase, result.Error))
		return
	}

	server.OK(w)
}

func TimeEndPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		server.BadRequest(w, errors.Join(server.ErrBodyRead, err))
		return
	}

	usertime := UserTime{}
	err = json.Unmarshal(body, &usertime)
	if err != nil {
		server.BadRequest(w, errors.Join(server.ErrJson, err))
		return
	}
	log.Debug(usertime, usertime.Passport)

	query_user := m.User{}
	query_user.PassSerie = uint32(usertime.Passport.PassSerie)
	query_user.PassNumber = uint32(usertime.Passport.PassNum)
	log.Debug(query_user)

	task := m.Task{}
	task.Id = uint(usertime.TaskId)

	db := m.GetDB()
	q1 := func(tx *gorm.DB) *gorm.DB {
		return tx.Where(&query_user).First(&query_user)
	}
	q2 := func(tx *gorm.DB) *gorm.DB {
		return tx.First(&task)
	}
	q3 := func(tx *gorm.DB) *gorm.DB {
		return tx.Save(&task)
	}
	log.Debug(db.ToSQL(q1))

	result := q1(db)
	if result.Error != nil {
		server.BadRequest(w, errors.Join(server.ErrDatabase, result.Error))
		return
	}
	log.Debug(query_user)

	log.Debug(db.ToSQL(q2))
	result = q2(db)
	if result.Error != nil {
		server.BadRequest(w, errors.Join(server.ErrDatabase, result.Error))
		return
	}
	log.Debug(task)

	if task.TimeStart == nil {
		server.BadRequest(w, errors.Join(server.ErrUnstartedTask))
	}
	if task.UserId != query_user.Id {
		server.BadRequest(w, errors.Join(server.ErrIncorrectUser))
	}

	currtime := time.Now()

	task.TimeEnd = &currtime

	log.Debug(db.ToSQL(q3))

	result = q3(db)
	if result.Error != nil {
		server.BadRequest(w, errors.Join(server.ErrDatabase, result.Error))
		return
	}
	server.OK(w)
}

func TimeGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	server.OK(w)
}

func TimeStartPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		server.BadRequest(w, errors.Join(server.ErrBodyRead, err))
		return
	}

	usertime := UserTime{}
	err = json.Unmarshal(body, &usertime)
	if err != nil {
		server.BadRequest(w, errors.Join(server.ErrJson, err))
		return
	}
	log.Debug(usertime, usertime.Passport)

	query_user := m.User{}
	query_user.PassSerie = uint32(usertime.Passport.PassSerie)
	query_user.PassNumber = uint32(usertime.Passport.PassNum)
	log.Debug(query_user)

	task := m.Task{}
	task.Id = uint(usertime.TaskId)

	db := m.GetDB()
	q1 := func(tx *gorm.DB) *gorm.DB {
		return tx.Where(&query_user).First(&query_user)
	}
	q2 := func(tx *gorm.DB) *gorm.DB {
		return tx.First(&task)
	}
	q3 := func(tx *gorm.DB) *gorm.DB {
		return tx.Save(&task)
	}
	log.Debug(db.ToSQL(q1))

	result := q1(db)
	if result.Error != nil {
		server.BadRequest(w, errors.Join(server.ErrDatabase, result.Error))
		return
	}
	log.Debug(query_user)

	log.Debug(db.ToSQL(q2))
	result = q2(db)
	if result.Error != nil {
		server.BadRequest(w, errors.Join(server.ErrDatabase, result.Error))
		return
	}
	log.Debug(task)

	currtime := time.Now()

	task.TimeStart = &currtime
	task.TimeEnd = nil
	task.UserId = query_user.Id

	log.Debug(db.ToSQL(q3))

	result = q3(db)
	if result.Error != nil {
		server.BadRequest(w, errors.Join(server.ErrDatabase, result.Error))
		return
	}
	server.OK(w)
}
