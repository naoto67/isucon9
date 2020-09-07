package main

import "github.com/jmoiron/sqlx"

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
