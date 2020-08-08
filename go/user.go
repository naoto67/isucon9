package main

import (
	"fmt"
	"strings"
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
	fmt.Println("FetchUserSimpleDictByItems: items: ", items)
	for _, v := range items {
		if v.BuyerID != 0 {
			userIDs = append(userIDs, fmt.Sprintf("%d", v.BuyerID))
		}
		userIDs = append(userIDs, fmt.Sprintf("%d", v.SellerID))
	}
	var users []UserSimple
	err := dbx.Select(&users, "SELECT id, account_name, num_sell_items FROM users WHERE id IN (?)", strings.Join(userIDs, ","))
	fmt.Println("FetchUserSimpleDictByItems: users: ", users)
	if err != nil {
		return nil, err
	}
	dict := map[int64]UserSimple{}
	for _, v := range users {
		dict[v.ID] = v
	}

	fmt.Println("FetchUserSimpleDictByItems: dict: ", dict)
	return dict, nil
}
