package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

const (
	USER_CACHE_PREFIX         = "uid:"
	USER_ACCOUNT_CACHE_PREFIX = "ua:"
)

var (
	DefaultLastBump, _ = time.Parse("2006-01-02 00:00:00", "2000-01-01 00:00:00")
)

type UserCache struct {
	ID             int64     `json:"id" db:"id"`
	AccountName    string    `json:"account_name" db:"account_name"`
	HashedPassword []byte    `json:"hashed_password" db:"hashed_password"`
	Address        string    `json:"address,omitempty" db:"address"`
	NumSellItems   int       `json:"num_sell_items" db:"num_sell_items"`
	LastBump       time.Time `json:"last_bump" db:"last_bump"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

func (user User) toUserCache() UserCache {
	return UserCache{
		ID:             user.ID,
		AccountName:    user.AccountName,
		HashedPassword: user.HashedPassword,
		Address:        user.Address,
		NumSellItems:   user.NumSellItems,
		LastBump:       user.LastBump,
		CreatedAt:      user.CreatedAt,
	}
}

func (user UserCache) toUser() User {
	return User{
		ID:             user.ID,
		AccountName:    user.AccountName,
		HashedPassword: user.HashedPassword,
		Address:        user.Address,
		NumSellItems:   user.NumSellItems,
		LastBump:       user.LastBump,
		CreatedAt:      user.CreatedAt,
	}
}

func StoreUserCache(user User) error {
	data, err := json.Marshal(user.toUserCache())
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s%d", USER_CACHE_PREFIX, user.ID)
	err = cacheClient.SingleSet(key, data)
	if err != nil {
		return err
	}
	key = fmt.Sprintf("%s%s", USER_ACCOUNT_CACHE_PREFIX, user.AccountName)
	return cacheClient.SingleSet(key, data)
}

func InitUserCache() error {
	var users []UserCache
	err := dbx.Select(&users, "SELECT * FROM users")
	if err != nil {
		return err
	}
	userMap := map[string][]byte{}
	userAccountMap := map[string][]byte{}
	for _, user := range users {
		data, err := json.Marshal(user)
		if err != nil {
			return err
		}
		key := fmt.Sprintf("%s%d", USER_CACHE_PREFIX, user.ID)
		userMap[key] = data
		key = fmt.Sprintf("%s%s", USER_ACCOUNT_CACHE_PREFIX, user.AccountName)
		userAccountMap[key] = data
	}
	err = cacheClient.MultiSet(userMap)
	if err != nil {
		return err
	}
	return cacheClient.MultiSet(userAccountMap)
}

func FetchUserCache(userID int64) (*User, error) {
	key := fmt.Sprintf("%s%d", USER_CACHE_PREFIX, userID)
	data, err := cacheClient.SingleGet(key)
	if err != nil {
		return nil, err
	}
	var user UserCache
	err = json.Unmarshal(data, &user)
	res := user.toUser()
	return &res, err
}

func FetchUserCacheByAccountName(accountName string) (*User, error) {
	key := fmt.Sprintf("%s%s", USER_ACCOUNT_CACHE_PREFIX, accountName)
	data, err := cacheClient.SingleGet(key)
	if err != nil {
		return nil, err
	}
	var user UserCache
	err = json.Unmarshal(data, &user)
	res := user.toUser()
	return &res, err
}

func FetchUserSimpleDictFromCache(items []Item) (map[int64]UserSimple, error) {
	var userIDs []string
	for _, v := range items {
		if v.BuyerID != 0 {
			userIDs = append(userIDs, fmt.Sprintf("%s%d", USER_CACHE_PREFIX, v.BuyerID))
		}
		userIDs = append(userIDs, fmt.Sprintf("%s%d", USER_CACHE_PREFIX, v.SellerID))
	}
	data, err := cacheClient.MultiGet(userIDs)
	if err != nil {
		return nil, err
	}
	dict := map[int64]UserSimple{}
	for i := range data {
		if data[i] == nil {
			log.Println("user not found")
			continue
		}
		var user User
		err = json.Unmarshal(data[i], &user)
		if err != nil {
			return nil, err
		}
		userSimple := UserSimple{
			ID:           user.ID,
			AccountName:  user.AccountName,
			NumSellItems: user.NumSellItems,
		}
		dict[user.ID] = userSimple
	}

	return dict, nil
}
