package main

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func init() {
	for _, v := range categories {
		categoryDict[v.ID] = v
	}
}

var (
	categoryDict = map[int]Category{}
	categories   = []Category{
		{ID: 1, ParentID: 0, CategoryName: "ソファー"},
		{ID: 2, ParentID: 1, CategoryName: "一人掛けソファー", ParentCategoryName: "ソファー"},
		{ID: 3, ParentID: 1, CategoryName: "二人掛けソファー", ParentCategoryName: "ソファー"},
		{ID: 4, ParentID: 1, CategoryName: "コーナーソファー", ParentCategoryName: "ソファー"},
		{ID: 5, ParentID: 1, CategoryName: "二段ソファー", ParentCategoryName: "ソファー"},
		{ID: 6, ParentID: 1, CategoryName: "ソファーベッド", ParentCategoryName: "ソファー"},
		{ID: 10, ParentID: 0, CategoryName: "家庭用チェア"},
		{ID: 11, ParentID: 10, CategoryName: "スツール", ParentCategoryName: "家庭用チェア"},
		{ID: 12, ParentID: 10, CategoryName: "クッションスツール", ParentCategoryName: "家庭用チェア"},
		{ID: 13, ParentID: 10, CategoryName: "ダイニングチェア", ParentCategoryName: "家庭用チェア"},
		{ID: 14, ParentID: 10, CategoryName: "リビングチェア", ParentCategoryName: "家庭用チェア"},
		{ID: 15, ParentID: 10, CategoryName: "カウンターチェア", ParentCategoryName: "家庭用チェア"},
		{ID: 20, ParentID: 0, CategoryName: "キッズチェア"},
		{ID: 21, ParentID: 20, CategoryName: "学習チェア", ParentCategoryName: "キッズチェア"},
		{ID: 22, ParentID: 20, CategoryName: "ベビーソファ", ParentCategoryName: "キッズチェア"},
		{ID: 23, ParentID: 20, CategoryName: "キッズハイチェア", ParentCategoryName: "キッズチェア"},
		{ID: 24, ParentID: 20, CategoryName: "テーブルチェア", ParentCategoryName: "キッズチェア"},
		{ID: 30, ParentID: 0, CategoryName: "オフィスチェア"},
		{ID: 31, ParentID: 30, CategoryName: "デスクチェア", ParentCategoryName: "オフィスチェア"},
		{ID: 32, ParentID: 30, CategoryName: "ビジネスチェア", ParentCategoryName: "オフィスチェア"},
		{ID: 33, ParentID: 30, CategoryName: "回転チェア", ParentCategoryName: "オフィスチェア"},
		{ID: 34, ParentID: 30, CategoryName: "リクライニングチェア", ParentCategoryName: "オフィスチェア"},
		{ID: 35, ParentID: 30, CategoryName: "投擲用椅子", ParentCategoryName: "オフィスチェア"},
		{ID: 40, ParentID: 0, CategoryName: "折りたたみ椅子"},
		{ID: 41, ParentID: 40, CategoryName: "パイプ椅子", ParentCategoryName: "折りたたみ椅子"},
		{ID: 42, ParentID: 40, CategoryName: "木製折りたたみ椅子", ParentCategoryName: "折りたたみ椅子"},
		{ID: 43, ParentID: 40, CategoryName: "キッチンチェア", ParentCategoryName: "折りたたみ椅子"},
		{ID: 44, ParentID: 40, CategoryName: "アウトドアチェア", ParentCategoryName: "折りたたみ椅子"},
		{ID: 45, ParentID: 40, CategoryName: "作業椅子", ParentCategoryName: "折りたたみ椅子"},
		{ID: 50, ParentID: 0, CategoryName: "ベンチ"},
		{ID: 51, ParentID: 50, CategoryName: "一人掛けベンチ", ParentCategoryName: "ベンチ"},
		{ID: 52, ParentID: 50, CategoryName: "二人掛けベンチ", ParentCategoryName: "ベンチ"},
		{ID: 53, ParentID: 50, CategoryName: "アウトドア用ベンチ", ParentCategoryName: "ベンチ"},
		{ID: 54, ParentID: 50, CategoryName: "収納付きベンチ", ParentCategoryName: "ベンチ"},
		{ID: 55, ParentID: 50, CategoryName: "背もたれ付きベンチ", ParentCategoryName: "ベンチ"},
		{ID: 56, ParentID: 50, CategoryName: "ベンチマーク", ParentCategoryName: "ベンチ"},
		{ID: 60, ParentID: 0, CategoryName: "座椅子"},
		{ID: 61, ParentID: 60, CategoryName: "和風座椅子", ParentCategoryName: "座椅子"},
		{ID: 62, ParentID: 60, CategoryName: "高座椅子", ParentCategoryName: "座椅子"},
		{ID: 63, ParentID: 60, CategoryName: "ゲーミング座椅子", ParentCategoryName: "座椅子"},
		{ID: 64, ParentID: 60, CategoryName: "ロッキングチェア", ParentCategoryName: "座椅子"},
		{ID: 65, ParentID: 60, CategoryName: "座布団", ParentCategoryName: "座椅子"},
		{ID: 66, ParentID: 60, CategoryName: "空気椅子", ParentCategoryName: "座椅子"},
	}
)

func getCategoryByID(q sqlx.Queryer, categoryID int) (category Category, err error) {
	category, ok := categoryDict[categoryID]
	if ok {
		return category, nil
	}
	return Category{}, sql.ErrNoRows
}
