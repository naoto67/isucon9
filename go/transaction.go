package main

import (
	"fmt"
	"strings"
)

type TS struct {
	TransactionEvidence TransactionEvidence
	Shipping            Shipping
}

func FetchTransactionDictByItems(items []Item) (map[int64]TS, error) {
	var itemIDs []string
	for _, v := range items {
		itemIDs = append(itemIDs, fmt.Sprintf("%d", v.ID))
	}
	fmt.Println("FetchTransactionDictByItems: itemIDs: ", itemIDs)
	query := fmt.Sprintf("SELECT * FROM `transaction_evidences` t INNER JOIN `shippings` s ON t.id = s.transaction_evidence_id WHERE t.`item_id` IN (%s)", strings.Join(itemIDs, ","))
	fmt.Println("FetchTransactionDictByItems: query: ", query)
	rows, err := dbx.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dict := map[int64]TS{}

	for rows.Next() {
		var t TransactionEvidence
		var s Shipping
		if err = rows.Scan(
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
		); err != nil {
			fmt.Println(err)
			return nil, err
		}

		dict[t.ItemID] = TS{
			TransactionEvidence: t,
			Shipping:            s,
		}
	}
	fmt.Println("FetchTransactionDictByItems: dict: ", dict)
	fmt.Println("FetchTransactionDictByItems: len(dict): ", len(dict))

	return dict, nil
}
