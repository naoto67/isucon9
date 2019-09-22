package main

import (
	"database/sql"
	"encoding/json"
	"errors"
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

	tx := dbx.MustBegin()

	targetItem := Item{}
	err = tx.Get(&targetItem, "SELECT * FROM `items` WHERE `id` = ? ", rb.ItemID)
	if err == sql.ErrNoRows {
		outputErrorMsg(w, http.StatusNotFound, "item not found")
		tx.Rollback()
		return
	}
	if err != nil {
		log.Print(err)

		outputErrorMsg(w, http.StatusInternalServerError, "db error")
		tx.Rollback()
		return
	}

	if targetItem.Status != ItemStatusOnSale {
		outputErrorMsg(w, http.StatusForbidden, "item is not for sale")
		tx.Rollback()
		return
	}

	if targetItem.SellerID == buyer.ID {
		outputErrorMsg(w, http.StatusForbidden, "自分の商品は買えません")
		tx.Rollback()
		return
	}

	seller := User{}
	redisful, _ := NewRedisful()
	seller, err = redisful.fetchUserByID(targetItem.SellerID)
	defer redisful.Close()
	if err != nil {
		err = tx.Get(&seller, "SELECT * FROM `users` WHERE `id` = ?", targetItem.SellerID)
		if err == sql.ErrNoRows {
			outputErrorMsg(w, http.StatusNotFound, "seller not found")
			tx.Rollback()
			return
		}
		if err != nil {
			log.Print(err)

			outputErrorMsg(w, http.StatusInternalServerError, "db error")
			tx.Rollback()
			return
		}
	}

	category, err := getCategoryByID(tx, targetItem.CategoryID)
	if err != nil {
		log.Print(err)

		outputErrorMsg(w, http.StatusInternalServerError, "category id error")
		tx.Rollback()
		return
	}

	errChan := make(chan error)
	shipmentCreateRes := make(chan APIShipmentCreateRes)

	go func() {
		scr, err := APIShipmentCreate(getShipmentServiceURL(), &APIShipmentCreateReq{
			ToAddress:   buyer.Address,
			ToName:      buyer.AccountName,
			FromAddress: seller.Address,
			FromName:    seller.AccountName,
		})

		if err != nil {
			errChan <- err
		}
		shipmentCreateRes <- *scr
	}()

	go func() {
		pstr, err := APIPaymentToken(getPaymentServiceURL(), &APIPaymentServiceTokenReq{
			ShopID: PaymentServiceIsucariShopID,
			Token:  rb.Token,
			APIKey: PaymentServiceIsucariAPIKey,
			Price:  targetItem.Price,
		})
		if err != nil {
			errChan <- err
			return
		}

		if pstr.Status == "invalid" {
			errChan <- errors.New("カード情報に誤りがあります")
			return
		}

		if pstr.Status == "fail" {
			errChan <- errors.New("カードの残高が足りません")
			return
		}

		if pstr.Status != "ok" {
			errChan <- errors.New("想定外のエラー")
			return
		}
	}()

	done := make(chan int64)

	go func() {
		_, err = tx.Exec("UPDATE `items` SET `buyer_id` = ?, `status` = ?, `updated_at` = ? WHERE `id` = ?",
			buyer.ID,
			ItemStatusTrading,
			time.Now(),
			targetItem.ID,
		)
		if err != nil {
			errChan <- err
			return
		}

		scr := <-shipmentCreateRes

		result, err := tx.Exec("INSERT INTO `transactions` (`seller_id`, `buyer_id`, `trans_status`, `item_id`, `item_name`, `item_price`, `item_description`,`item_category_id`,`item_root_category_id`, `ship_status`, `reserve_id`, `reserve_time`, `to_address`, `to_name`, `from_address`, `from_name`, `img_binary`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			targetItem.SellerID,
			buyer.ID,
			TransactionEvidenceStatusWaitShipping,
			targetItem.ID,
			targetItem.Name,
			targetItem.Price,
			targetItem.Description,
			category.ID,
			category.ParentID,

			ShippingsStatusInitial,
			scr.ReserveID,
			scr.ReserveTime,
			buyer.Address,
			buyer.AccountName,
			seller.Address,
			seller.AccountName,
			"",
		)

		if err != nil {
			errChan <- errors.New("db error")
			return
		}

		transactionEvidenceID, err := result.LastInsertId()
		if err != nil {
			errChan <- errors.New("db error")
			return
		}
		done <- transactionEvidenceID
	}()

	var transactionEvidenceID int64
Label:
	for {
		select {
		case err := <-errChan:
			tx.Rollback()
			outputErrorMsg(w, http.StatusBadRequest, err.Error())
			return
		case transactionEvidenceID = <-done:
			tx.Commit()
			break Label
		}
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(resBuy{TransactionEvidenceID: transactionEvidenceID})
}
