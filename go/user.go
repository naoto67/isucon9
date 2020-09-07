package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	USER_ID_KEY = "uid:"
)

var (
	DefaultLastBump, _ = time.Parse("1996-01-01 00:00:00", "2000-01-01 00:00:00")
)

func FetchUserSimplesDictFromIDs(userIDs []int64) (map[int64]UserSimple, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	query := "SELECT id, account_name, num_sell_items FROM users WHERE id IN (?)"
	inQuery, inArgs, err := sqlx.In(query, userIDs)
	if err != nil {
		return nil, err
	}
	var users []UserSimple
	err = dbx.Select(&users, inQuery, inArgs...)
	if err != nil {
		return nil, err
	}
	res := map[int64]UserSimple{}
	for _, v := range users {
		res[v.ID] = v
	}
	return res, err
}

func InitUserCache() error {
	var users []UserCache
	err := dbx.Select(&users, "SELECT * FROM users")
	if err != nil {
		return err
	}
	data := map[string][]byte{}
	for _, v := range users {
		d, err := json.Marshal(v)
		if err != nil {
			return err
		}
		key := fmt.Sprintf("%s%d", USER_ID_KEY, v.ID)
		data[key] = d
	}
	return cacheClient.MultiSet(data)
}

func InsertUserCache(user User) error {
	uc := UserCache{
		ID:             user.ID,
		AccountName:    user.AccountName,
		Address:        user.Address,
		HashedPassword: user.HashedPassword,
		NumSellItems:   user.NumSellItems,
		LastBump:       user.LastBump,
		CreatedAt:      user.CreatedAt,
	}

	d, err := json.Marshal(uc)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s%d", USER_ID_KEY, uc.ID)
	return cacheClient.SingleSet(key, d)
}

func FetchUserCache(userID int64) (*UserCache, error) {
	key := fmt.Sprintf("%s%d", USER_ID_KEY, userID)
	data, err := cacheClient.SingleGet(key)
	if err != nil {
		return nil, err
	}
	var res *UserCache
	err = json.Unmarshal(data, &res)
	return res, err
}

func FetchUsersCache(userIDs []int64) ([]*UserCache, error) {
	keys := []string{}
	for _, v := range userIDs {
		key := fmt.Sprintf("%s%d", USER_ID_KEY, v)
		keys = append(keys, key)
	}
	data, err := cacheClient.MultiGet(keys)
	var res []*UserCache
	for i := range data {
		if len(data[i]) == 0 {
			continue
		}
		var uc *UserCache
		err = json.Unmarshal(data[i], &uc)
		if err != nil {
			return nil, err
		}
		res = append(res, uc)
	}
	return res, err
}

func (uc *UserCache) BuildUserSimple() UserSimple {
	return UserSimple{
		ID:           uc.ID,
		AccountName:  uc.AccountName,
		NumSellItems: uc.NumSellItems,
	}
}

func (uc *UserCache) BuildUser() User {
	return User{
		ID:             uc.ID,
		AccountName:    uc.AccountName,
		HashedPassword: uc.HashedPassword,
		Address:        uc.Address,
		NumSellItems:   uc.NumSellItems,
		LastBump:       uc.LastBump,
		CreatedAt:      uc.CreatedAt,
	}
}
