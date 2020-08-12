package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	USER_KEY = "users"
)

func FetchUserDictByItems(items []Item) (map[int64]User, error) {
	var userIDs []string
	for _, v := range items {
		if v.BuyerID != 0 {
			userIDs = append(userIDs, fmt.Sprintf("%d", v.BuyerID))
		}
		userIDs = append(userIDs, fmt.Sprintf("%d", v.SellerID))
	}
	var users []User
	err := dbx.Select(&users, "SELECT * FROM users WHERE id IN (?)", strings.Join(userIDs, ","))
	if err != nil {
		return nil, err
	}
	dict := map[int64]User{}
	for _, v := range users {
		dict[v.ID] = v
	}

	return dict, nil
}

func FetchUserSimpleDictByItems(items []Item) (map[int64]UserSimple, error) {
	var userIDs []string
	for _, v := range items {
		if v.BuyerID != 0 {
			userIDs = append(userIDs, fmt.Sprintf("%d", v.BuyerID))
		}
		userIDs = append(userIDs, fmt.Sprintf("%d", v.SellerID))
	}
	var users []UserSimple
	query := fmt.Sprintf("SELECT id, account_name, num_sell_items FROM users WHERE id IN (%s)", strings.Join(userIDs, ","))
	err := dbx.Select(&users, query)
	if err != nil {
		return nil, err
	}
	dict := map[int64]UserSimple{}
	for _, v := range users {
		dict[v.ID] = v
	}

	return dict, nil
}

func InitUsersCache() error {
	var users []User
	err := dbx.Select(&users, "SELECT * FROM users")
	if err != nil {
		return err
	}

	var values []interface{}

	for _, v := range users {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}

		values = append(values, v.ID)
		values = append(values, b)
	}

	err = redisClient.HMSET(USER_KEY, values)
	return err
}
