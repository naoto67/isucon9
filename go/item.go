package main

func FetchItemIDsFromItems(items []Item) []int64 {
	var res []int64
	for _, v := range items {
		res = append(res, v.ID)
	}
	return res
}

func FetchUserIDsFromItems(items []Item) []int64 {
	var res []int64
	for _, v := range items {
		if v.BuyerID != 0 {
			res = append(res, v.BuyerID)
		}
		res = append(res, v.SellerID)
	}
	return res
}
