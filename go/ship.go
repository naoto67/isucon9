package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func postShipDone(w http.ResponseWriter, r *http.Request) {
	reqpsd := reqPostShipDone{}

	err := json.NewDecoder(r.Body).Decode(&reqpsd)
	if err != nil {
		outputErrorMsg(w, http.StatusBadRequest, "json decode error")
		return
	}

	csrfToken := reqpsd.CSRFToken
	itemID := reqpsd.ItemID

	if csrfToken != getCSRFToken(r) {
		outputErrorMsg(w, http.StatusUnprocessableEntity, "csrf token error")

		return
	}

	var t TransactionEvidence
	var s Shipping
	err = dbx.QueryRow("SELECT * FROM `transaction_evidences` t INNER JOIN shippings s ON t.id = s.transaction_evidence_id WHERE t.item_id = ?", itemID).Scan(
		&t.ID, &t.SellerID, &t.BuyerID, &t.Status, &t.ItemID, &t.ItemName, &t.ItemPrice, &t.ItemDescription, &t.ItemCategoryID, &t.ItemRootCategoryID, &t.CreatedAt, &t.UpdatedAt,
		&s.TransactionEvidenceID, &s.Status, &s.ItemName, &s.ItemID, &s.ReserveID, &s.ReserveTime, &s.ToAddress, &s.ToName, &s.FromAddress, &s.FromName, &s.ImgBinary, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		outputErrorMsg(w, http.StatusNotFound, "transaction_evidence not found")
		return
	}
	ssrStatus := make(chan string)
	ssrErr := make(chan error)
	go func() {
		ssr, err := APIShipmentStatus(getShipmentServiceURL(), &APIShipmentStatusReq{
			ReserveID: s.ReserveID,
		})
		ssrStatus <- ssr.Status
		ssrErr <- err
	}()

	transactionEvidence := t

	seller, errCode, errMsg := getUser(r)
	if errMsg != "" {
		outputErrorMsg(w, errCode, errMsg)
		return
	}

	if transactionEvidence.SellerID != seller.ID {
		outputErrorMsg(w, http.StatusForbidden, "権限がありません")
		return
	}

	tx := dbx.MustBegin()

	item := Item{}
	err = tx.Get(&item, "SELECT * FROM `items` WHERE `id` = ? ", itemID)
	if err == sql.ErrNoRows {
		outputErrorMsg(w, http.StatusNotFound, "items not found")
		tx.Rollback()
		return
	}
	if err != nil {
		log.Print(err)
		outputErrorMsg(w, http.StatusInternalServerError, "db error")
		tx.Rollback()
		return
	}

	if item.Status != ItemStatusTrading {
		outputErrorMsg(w, http.StatusForbidden, "商品が取引中ではありません")
		tx.Rollback()
		return
	}

	// 元の場所のまま
	if transactionEvidence.Status != TransactionEvidenceStatusWaitShipping {
		outputErrorMsg(w, http.StatusForbidden, "準備ができていません")
		tx.Rollback()
		return
	}

	status := <-ssrStatus
	err = <-ssrErr
	if err != nil {
		log.Print(err)
		outputErrorMsg(w, http.StatusInternalServerError, "failed to request to shipment service")
		tx.Rollback()

		return
	}
	if !(status == ShippingsStatusShipping || status == ShippingsStatusDone) {
		outputErrorMsg(w, http.StatusForbidden, "shipment service側で配送中か配送完了になっていません")
		tx.Rollback()
		return
	}

	_, err = tx.Exec("UPDATE `shippings` SET `status` = ?, `updated_at` = ? WHERE `transaction_evidence_id` = ?",
		status,
		time.Now(),
		transactionEvidence.ID,
	)
	if err != nil {
		log.Print(err)

		outputErrorMsg(w, http.StatusInternalServerError, "db error")
		tx.Rollback()
		return
	}

	_, err = tx.Exec("UPDATE `transaction_evidences` SET `status` = ?, `updated_at` = ? WHERE `id` = ?",
		TransactionEvidenceStatusWaitDone,
		time.Now(),
		transactionEvidence.ID,
	)
	if err != nil {
		log.Print(err)

		outputErrorMsg(w, http.StatusInternalServerError, "db error")
		tx.Rollback()
		return
	}

	tx.Commit()

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(resBuy{TransactionEvidenceID: transactionEvidence.ID})
}
