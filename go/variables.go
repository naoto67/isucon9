package main

import (
	crand "crypto/rand"
	"fmt"
	"html/template"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

const (
	sessionName = "session_isucari"

	DefaultPaymentServiceURL  = "http://localhost:5555"
	DefaultShipmentServiceURL = "http://localhost:7000"

	ItemMinPrice    = 100
	ItemMaxPrice    = 1000000
	ItemPriceErrMsg = "商品価格は100ｲｽｺｲﾝ以上、1,000,000ｲｽｺｲﾝ以下にしてください"

	ItemStatusOnSale  = "on_sale"
	ItemStatusTrading = "trading"
	ItemStatusSoldOut = "sold_out"
	ItemStatusStop    = "stop"
	ItemStatusCancel  = "cancel"

	PaymentServiceIsucariAPIKey = "a15400e46c83635eb181-946abb51ff26a868317c"
	PaymentServiceIsucariShopID = "11"

	TransactionEvidenceStatusWaitShipping = "wait_shipping"
	TransactionEvidenceStatusWaitDone     = "wait_done"
	TransactionEvidenceStatusDone         = "done"

	ShippingsStatusInitial    = "initial"
	ShippingsStatusWaitPickup = "wait_pickup"
	ShippingsStatusShipping   = "shipping"
	ShippingsStatusDone       = "done"

	BumpChargeSeconds = 3 * time.Second

	ItemsPerPage        = 48
	TransactionsPerPage = 10

	BcryptCost = 10

	USERS_KEY         = "USERS"
	USER_FIELD_PREFIX = "USER-FIELD"

	TransactionEvidenceFields = "id, seller_id, buyer_id, trans_status, item_id, item_name, item_price, item_description, item_category_id, item_root_category_id, created_at, trans_updated_at"
	ShippingFields            = "id, ship_status, item_id, item_name, reserve_id, reserve_time, to_address, to_name, from_address, from_name, img_binary, created_at, updated_at"
)

var (
	templates *template.Template
	dbx       *sqlx.DB
	store     sessions.Store
)

func getPaymentServiceURL() string {
	val, _ := getConfigByName("payment_service_url")
	if val == "" {
		return DefaultPaymentServiceURL
	}
	return val
}

func getShipmentServiceURL() string {
	val, _ := getConfigByName("shipment_service_url")
	if val == "" {
		return DefaultShipmentServiceURL
	}
	return val
}

func secureRandomStr(b int) string {
	k := make([]byte, b)
	if _, err := crand.Read(k); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", k)
}

func getImageURL(imageName string) string {
	return fmt.Sprintf("/upload/%s", imageName)
}
