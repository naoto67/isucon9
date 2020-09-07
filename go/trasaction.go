package main

import "github.com/jmoiron/sqlx"

// res key itemID
func FetchTransactionDictFromItemIDs(itemIDs []int64) (map[int64]TS, error) {
	query := "SELECT * FROM `transaction_evidences` t INNER JOIN shippings s ON s.transaction_evidence_id = t.id WHERE t.`item_id` IN (?)"
	inQuery, inArgs, err := sqlx.In(query, itemIDs)
	if err != nil {
		return nil, err
	}
	rows, err := dbx.Query(inQuery, inArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := map[int64]TS{}
	for rows.Next() {
		var t TransactionEvidence
		var s Shipping
		err = rows.Scan(
			&t.ID,
			&t.SellerID,
			&t.BuyerID,
			&t.Status,
			&t.ItemID,
			&t.ItemName,
			&t.ItemPrice,
			&t.ItemDescription,
			&t.ItemCategoryID,
			&t.ItemRootCategoryID,
			&t.CreatedAt,
			&t.UpdatedAt,
			&s.TransactionEvidenceID,
			&s.Status,
			&s.ItemName,
			&s.ItemID,
			&s.ReserveID,
			&s.ReserveTime,
			&s.ToAddress,
			&s.ToName,
			&s.FromAddress,
			&s.FromName,
			&s.ImgBinary,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		res[t.ItemID] = TS{t, s}
	}
	return res, nil
}

// func FetchTransactionIDsFromTE(transactionEvidences []TransactionEvidence) []int64 {
// 	var res []int64
// 	for _, v := range transactionEvidences {
// 		res = append(res, v.ID)
// 	}
// 	return res
// }
//
// func FetchShippingDictFromTransactionIDs(tIDs []int64) (map[int64]Shipping, error) {
// 	var t []Shipping
// 	query := "SELECT * FROM `shippings` WHERE `transaction_evidence_id` IN (?)"
// 	inQuery, inArgs, err := sqlx.In(query, tIDs)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = dbx.Select(&t, inQuery, inArgs...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	res := map[int64]Shipping{}
// 	for _, v := range t {
// 		res[v.ItemID] = v
// 	}
// 	return res, nil
// }
