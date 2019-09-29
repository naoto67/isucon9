package main

import (
	"encoding/json"
	"fmt"
	"log"
)

var (
	SHIP_KEY = string("SHIPPINGS")
)

func makeShippingField(transactionEvidenceID int64) string {
	return fmt.Sprintf("%d", transactionEvidenceID)
}

func (r *Redisful) StoreShipping(shipping Shipping) {
	field := makeShippingField(shipping.TransactionEvidenceID)
	r.SetHashToCache(SHIP_KEY, field, shipping)
}

func (r *Redisful) InitShippings() {
	var shippings []Shipping
	err := dbx.Select(&shippings, "SELECT * FROM `shippings`")
	if err != nil {
		log.Println("failed to select shippings")
		return
	}
	v := make([]interface{}, 0, 1000)
	for i, _ := range shippings {
		var field string
		field = makeShippingField(shippings[i].TransactionEvidenceID)
		v = append(v, field)
		data, _ := json.Marshal(shippings[i])
		v = append(v, data)
		if i%1000 == 999 {
			r.SetMultiHashToCache(SHIP_KEY, v)
			v = make([]interface{}, 0, 1000)
		}
	}
	r.SetMultiHashToCache(SHIP_KEY, v)
}
