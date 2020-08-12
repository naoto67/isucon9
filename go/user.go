package main

import (
	"encoding/json"
	"fmt"
)

const (
	USER_KEY_PREFIX         = "u:"
	USER_ACCOUNT_KEY_PREFIX = "ua:"
	USER_KEY                = "users"
	USER_COUNT_KEY          = "uc"
)

func FetchUserDictByItems(items []Item) (map[int64]User, error) {
	var userIDs []string
	for _, v := range items {
		if v.BuyerID != 0 {
			userIDs = append(userIDs, fmt.Sprintf("%s%d", USER_KEY_PREFIX, v.BuyerID))
		}
		userIDs = append(userIDs, fmt.Sprintf("%s%d", USER_KEY_PREFIX, v.SellerID))
	}
	b, err := redisClient.MGET(userIDs)
	if err != nil {
		return nil, err
	}
	dict := map[int64]User{}
	for _, v := range b {
		var u CacheUser
		err = json.Unmarshal(v, &u)
		if err != nil {
			return nil, err
		}
		dict[u.ID] = User{
			ID:             u.ID,
			AccountName:    u.AccountName,
			HashedPassword: u.HashedPassword,
			Address:        u.Address,
			NumSellItems:   u.NumSellItems,
			LastBump:       u.LastBump,
			CreatedAt:      u.CreatedAt,
		}

	}

	return dict, nil
}

func FetchUserSimpleDictByItems(items []Item) (map[int64]UserSimple, error) {
	var userIDs []string
	for _, v := range items {
		if v.BuyerID != 0 {
			userIDs = append(userIDs, fmt.Sprintf("%s%d", USER_KEY_PREFIX, v.BuyerID))
		}
		userIDs = append(userIDs, fmt.Sprintf("%s%d", USER_KEY_PREFIX, v.SellerID))
	}
	b, err := redisClient.MGET(userIDs)
	if err != nil {
		return nil, err
	}
	dict := map[int64]UserSimple{}
	for _, v := range b {
		var u User
		err = json.Unmarshal(v, &u)
		if err != nil {
			return nil, err
		}
		dict[u.ID] = UserSimple{
			ID:           u.ID,
			AccountName:  u.AccountName,
			NumSellItems: u.NumSellItems,
		}

	}

	return dict, nil
}

func InitUsersCache() error {
	var users []CacheUser
	err := dbx.Select(&users, "SELECT * FROM users")
	if err != nil {
		return err
	}

	m := make(map[string][]byte)
	ua := make(map[string][]byte)

	for _, v := range users {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		m[fmt.Sprintf("%s%d", USER_KEY_PREFIX, v.ID)] = b

		b, err = json.Marshal(v.ID)
		if err != nil {
			return err
		}
		m[fmt.Sprintf("%s%s", USER_ACCOUNT_KEY_PREFIX, v.AccountName)] = b
	}

	err = redisClient.MSET(m)
	if err != nil {
		return err
	}
	err = redisClient.MSET(ua)
	if err != nil {
		return err
	}
	err = redisClient.SET(USER_COUNT_KEY, len(users))
	return err
}

func GetUserCacheByID(id int64) (*User, error) {
	b, err := redisClient.GET(fmt.Sprintf("%s%d", USER_KEY_PREFIX, id))
	if err != nil {
		return nil, err
	}
	var user CacheUser
	err = json.Unmarshal(b, &user)
	u := User{
		ID:             user.ID,
		AccountName:    user.AccountName,
		HashedPassword: user.HashedPassword,
		Address:        user.Address,
		NumSellItems:   user.NumSellItems,
		LastBump:       user.LastBump,
		CreatedAt:      user.CreatedAt,
	}
	return &u, err
}

func UpdateUserCache(user User) error {
	u := CacheUser{
		ID:             user.ID,
		AccountName:    user.AccountName,
		HashedPassword: user.HashedPassword,
		Address:        user.Address,
		NumSellItems:   user.NumSellItems,
		LastBump:       user.LastBump,
		CreatedAt:      user.CreatedAt,
	}
	b, err := json.Marshal(u)
	if err != nil {
		return err
	}
	err = redisClient.SET(fmt.Sprintf("%s%d", USER_KEY_PREFIX, user.ID), b)
	return err
}

func GetUserCacheByAccountName(accountName string) (*User, error) {
	b, err := redisClient.GET(fmt.Sprintf("%s%s", USER_ACCOUNT_KEY_PREFIX, accountName))
	if err != nil {
		return nil, err
	}
	var userID int64
	err = json.Unmarshal(b, &userID)
	if err != nil {
		return nil, err
	}

	return GetUserCacheByID(userID)
}

func GetUserCount() (int, error) {
	b, err := redisClient.GET(USER_COUNT_KEY)
	if err != nil {
		return 0, nil
	}
	var cnt int
	err = json.Unmarshal(b, &cnt)
	return cnt, err
}

func IncrUserCount() error {
	return redisClient.INCR(USER_COUNT_KEY)
}
