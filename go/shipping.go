package main

import (
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
	var field string
	for i, _ := range shippings {
		field = makeShippingField(shippings[i].TransactionEvidenceID)
		r.SetHashToCache(SHIP_KEY, field, shippings[i])
	}
}
