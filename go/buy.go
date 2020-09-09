package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func postBuy(w http.ResponseWriter, r *http.Request) {
	rb := reqBuy{}

	err := json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		outputErrorMsg(w, http.StatusBadRequest, "json decode error")
		return
	}

	if rb.CSRFToken != getCSRFToken(r) {
		outputErrorMsg(w, http.StatusUnprocessableEntity, "csrf token error")

		return
	}

	buyer, errCode, errMsg := getUser(r)
	if errMsg != "" {
		outputErrorMsg(w, errCode, errMsg)
		return
	}

	targetItem := Item{}
	ok := LockItem(rb.ItemID)
	if !ok {
		outputErrorMsg(w, http.StatusNotFound, "item not found")
		return
	}
	defer UnlockItem(rb.ItemID)
	err = dbx.Get(&targetItem, "SELECT * FROM `items` WHERE `id` = ?", rb.ItemID)
	if err == sql.ErrNoRows {
		outputErrorMsg(w, http.StatusNotFound, "item not found")
		return
	}
	if err != nil {
		log.Print(err)

		outputErrorMsg(w, http.StatusInternalServerError, "db error")
		return
	}

	if targetItem.Status != ItemStatusOnSale {
		outputErrorMsg(w, http.StatusForbidden, "item is not for sale")
		return
	}

	if targetItem.SellerID == buyer.ID {
		outputErrorMsg(w, http.StatusForbidden, "自分の商品は買えません")
		return
	}

	seller := User{}
	uc, err := FetchUserCache(targetItem.SellerID)
	if uc == nil {
		outputErrorMsg(w, http.StatusNotFound, "seller not found")
		return
	}
	if err != nil {
		log.Print(err)

		outputErrorMsg(w, http.StatusInternalServerError, "db error")
		return
	}
	seller = uc.BuildUser()

	category, err := getCategoryByID(dbx, targetItem.CategoryID)
	if err != nil {
		log.Print(err)

		outputErrorMsg(w, http.StatusInternalServerError, "category id error")
		return
	}

	chScr := make(chan *APIShipmentCreateRes)
	chErr := make(chan error)
	go func() {
		scr, err := APIShipmentCreate(getShipmentServiceURL(), &APIShipmentCreateReq{
			ToAddress:   buyer.Address,
			ToName:      buyer.AccountName,
			FromAddress: seller.Address,
			FromName:    seller.AccountName,
		})
		chScr <- scr
		chErr <- err
	}()

	tx := dbx.MustBegin()

	result, err := tx.Exec("INSERT INTO `transaction_evidences` (`seller_id`, `buyer_id`, `status`, `item_id`, `item_name`, `item_price`, `item_description`,`item_category_id`,`item_root_category_id`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		targetItem.SellerID,
		buyer.ID,
		TransactionEvidenceStatusWaitShipping,
		targetItem.ID,
		targetItem.Name,
		targetItem.Price,
		targetItem.Description,
		category.ID,
		category.ParentID,
	)
	if err != nil {
		log.Print(err)

		outputErrorMsg(w, http.StatusInternalServerError, "db error")
		tx.Rollback()
		return
	}

	transactionEvidenceID, err := result.LastInsertId()
	if err != nil {
		log.Print(err)

		outputErrorMsg(w, http.StatusInternalServerError, "db error")
		tx.Rollback()
		return
	}

	_, err = tx.Exec("UPDATE `items` SET `buyer_id` = ?, `status` = ?, `updated_at` = ? WHERE `id` = ?",
		buyer.ID,
		ItemStatusTrading,
		time.Now(),
		targetItem.ID,
	)
	if err != nil {
		log.Print(err)

		outputErrorMsg(w, http.StatusInternalServerError, "db error")
		tx.Rollback()
		return
	}

	err = <-chErr
	if err != nil {
		logger.Infow("shipment error", "err", err)
		outputErrorMsg(w, http.StatusInternalServerError, "failed to request to shipment service")
		tx.Rollback()

		return
	}
	scr := <-chScr

	pstr, err := APIPaymentToken(getPaymentServiceURL(), &APIPaymentServiceTokenReq{
		ShopID: PaymentServiceIsucariShopID,
		Token:  rb.Token,
		APIKey: PaymentServiceIsucariAPIKey,
		Price:  targetItem.Price,
	})
	if err != nil {
		log.Print(err)

		outputErrorMsg(w, http.StatusInternalServerError, "payment service is failed")
		tx.Rollback()
		return
	}

	if pstr.Status == "invalid" {
		outputErrorMsg(w, http.StatusBadRequest, "カード情報に誤りがあります")
		tx.Rollback()
		return
	}

	if pstr.Status == "fail" {
		outputErrorMsg(w, http.StatusBadRequest, "カードの残高が足りません")
		tx.Rollback()
		return
	}

	if pstr.Status != "ok" {
		outputErrorMsg(w, http.StatusBadRequest, "想定外のエラー")
		tx.Rollback()
		return
	}

	_, err = tx.Exec("INSERT INTO `shippings` (`transaction_evidence_id`, `status`, `item_name`, `item_id`, `reserve_id`, `reserve_time`, `to_address`, `to_name`, `from_address`, `from_name`, `img_binary`) VALUES (?,?,?,?,?,?,?,?,?,?,?)",
		transactionEvidenceID,
		ShippingsStatusInitial,
		targetItem.Name,
		targetItem.ID,
		scr.ReserveID,
		scr.ReserveTime,
		buyer.Address,
		buyer.AccountName,
		seller.Address,
		seller.AccountName,
		"",
	)
	if err != nil {
		log.Print(err)

		outputErrorMsg(w, http.StatusInternalServerError, "db error")
		tx.Rollback()
		return
	}

	tx.Commit()

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(resBuy{TransactionEvidenceID: transactionEvidenceID})
}
