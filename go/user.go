package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func getUserSimpleByID(q sqlx.Queryer, userID int64) (userSimple UserSimple, err error) {
	user := User{}
	err = sqlx.Get(q, &user, "SELECT * FROM `users` WHERE `id` = ?", userID)
	if err != nil {
		return userSimple, err
	}
	userSimple.ID = user.ID
	userSimple.AccountName = user.AccountName
	userSimple.NumSellItems = user.NumSellItems
	return userSimple, err
}

func getUser(r *http.Request) (user User, errCode int, errMsg string) {
	session := getSession(r)
	userID, ok := session.Values["user_id"]
	if !ok {
		return user, http.StatusNotFound, "no session"
	}

	redisful, _ := NewRedisful()
	defer redisful.Close()
	var err error
	user, err = redisful.fetchUserByID(userID.(int64))
	if err != nil {
		err := dbx.Get(&user, "SELECT * FROM `users` WHERE `id` = ?", userID)
		if err == sql.ErrNoRows {
			return user, http.StatusNotFound, "user not found"
		}
		if err != nil {
			log.Print(err)
			return user, http.StatusInternalServerError, "db error"
		}

	}

	return user, http.StatusOK, ""
}

func postRegister(w http.ResponseWriter, r *http.Request) {
	rr := reqRegister{}
	err := json.NewDecoder(r.Body).Decode(&rr)
	if err != nil {
		outputErrorMsg(w, http.StatusBadRequest, "json decode error")
		return
	}

	accountName := rr.AccountName
	address := rr.Address
	password := rr.Password

	if accountName == "" || password == "" || address == "" {
		outputErrorMsg(w, http.StatusBadRequest, "all parameters are required")

		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		log.Print(err)

		outputErrorMsg(w, http.StatusInternalServerError, "error")
		return
	}

	timeNow := time.Now()
	result, err := dbx.Exec("INSERT INTO `users` (`account_name`, `hashed_password`, `address`, `created_at`) VALUES (?, ?, ?, ?)",
		accountName,
		hashedPassword,
		address,
		timeNow,
	)
	if err != nil {
		log.Print(err)

		outputErrorMsg(w, http.StatusInternalServerError, "db error")
		return
	}

	userID, err := result.LastInsertId()

	if err != nil {
		log.Print(err)

		outputErrorMsg(w, http.StatusInternalServerError, "db error")
		return
	}

	redisful, _ := NewRedisful()

	redisful.StoreUserCache(User{
		ID:             userID,
		AccountName:    accountName,
		HashedPassword: hashedPassword,
		Address:        address,
		NumSellItems:   0,
		LastBump:       time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		CreatedAt:      timeNow,
	})

	u := User{
		ID:          userID,
		AccountName: accountName,
		Address:     address,
	}

	session := getSession(r)
	session.Values["user_id"] = u.ID
	session.Values["csrf_token"] = secureRandomStr(20)
	if err = session.Save(r, w); err != nil {
		log.Print(err)
		outputErrorMsg(w, http.StatusInternalServerError, "session error")
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(u)
}

func (r *Redisful) InitUsersCache() error {
	users := []User{}
	err := dbx.Select(&users, "SELECT * FROM users")
	if err != nil {
		return err
	}

	for i := range users {
		err = r.SetHashToCache(USERS_KEY, makeUsersField(users[i].ID), users[i].toJson())
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Redisful) StoreUserCache(user User) error {
	err := r.SetHashToCache(USERS_KEY, makeUsersField(user.ID), user.toJson())
	return err
}

func (r *Redisful) fetchUserByID(userID int64) (User, error) {
	user := User{}
	err := r.GetHashFromCache(USERS_KEY, makeUsersField(userID), &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *Redisful) fetchUserSimpleByID(userID int64) (UserSimple, error) {
	user := UserSimple{}
	err := r.GetHashFromCache(USERS_KEY, makeUsersField(userID), &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (user User) toJson() map[string]interface{} {
	res := make(map[string]interface{})
	res["id"] = user.ID
	res["account_name"] = user.AccountName
	res["hashed_password"] = user.HashedPassword
	res["address"] = user.Address
	res["num_sell_items"] = user.NumSellItems
	res["last_bump"] = user.LastBump
	res["created_at"] = user.CreatedAt
	return res
}

func makeUsersField(userID int64) string {
	return fmt.Sprintf("%s-%d", USER_FIELD_PREFIX, userID)
}
