package main

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	categories = map[int]interface{}{
		1:  Category{ID: 1, ParentID: 0, CategoryName: "ソファー"},
		2:  Category{ID: 2, ParentID: 1, CategoryName: "一人掛けソファー", ParentCategoryName: "ソファー"},
		3:  Category{ID: 3, ParentID: 1, CategoryName: "二人掛けソファー", ParentCategoryName: "ソファー"},
		4:  Category{ID: 4, ParentID: 1, CategoryName: "コーナーソファー", ParentCategoryName: "ソファー"},
		5:  Category{ID: 5, ParentID: 1, CategoryName: "二段ソファー", ParentCategoryName: "ソファー"},
		6:  Category{ID: 6, ParentID: 1, CategoryName: "ソファーベッド", ParentCategoryName: "ソファー"},
		10: Category{ID: 10, ParentID: 0, CategoryName: "家庭用チェア"},
		11: Category{ID: 11, ParentID: 10, CategoryName: "スツール", ParentCategoryName: "家庭用チェア"},
		12: Category{ID: 12, ParentID: 10, CategoryName: "クッションスツール", ParentCategoryName: "家庭用チェア"},
		13: Category{ID: 13, ParentID: 10, CategoryName: "ダイニングチェア", ParentCategoryName: "家庭用チェア"},
		14: Category{ID: 14, ParentID: 10, CategoryName: "リビングチェア", ParentCategoryName: "家庭用チェア"},
		15: Category{ID: 15, ParentID: 10, CategoryName: "カウンターチェア", ParentCategoryName: "家庭用チェア"},
		20: Category{ID: 20, ParentID: 0, CategoryName: "キッズチェア"},
		21: Category{ID: 21, ParentID: 20, CategoryName: "学習チェア", ParentCategoryName: "キッズチェア"},
		22: Category{ID: 22, ParentID: 20, CategoryName: "ベビーソファ", ParentCategoryName: "キッズチェア"},
		23: Category{ID: 23, ParentID: 20, CategoryName: "キッズハイチェア", ParentCategoryName: "キッズチェア"},
		24: Category{ID: 24, ParentID: 20, CategoryName: "テーブルチェア", ParentCategoryName: "キッズチェア"},
		30: Category{ID: 30, ParentID: 0, CategoryName: "オフィスチェア"},
		31: Category{ID: 31, ParentID: 30, CategoryName: "デスクチェア", ParentCategoryName: "オフィスチェア"},
		32: Category{ID: 32, ParentID: 30, CategoryName: "ビジネスチェア", ParentCategoryName: "オフィスチェア"},
		33: Category{ID: 33, ParentID: 30, CategoryName: "回転チェア", ParentCategoryName: "オフィスチェア"},
		34: Category{ID: 34, ParentID: 30, CategoryName: "リクライニングチェア", ParentCategoryName: "オフィスチェア"},
		35: Category{ID: 35, ParentID: 30, CategoryName: "投擲用椅子", ParentCategoryName: "オフィスチェア"},
		40: Category{ID: 40, ParentID: 0, CategoryName: "折りたたみ椅子"},
		41: Category{ID: 41, ParentID: 40, CategoryName: "パイプ椅子", ParentCategoryName: "折りたたみ椅子"},
		42: Category{ID: 42, ParentID: 40, CategoryName: "木製折りたたみ椅子", ParentCategoryName: "折りたたみ椅子"},
		43: Category{ID: 43, ParentID: 40, CategoryName: "キッチンチェア", ParentCategoryName: "折りたたみ椅子"},
		44: Category{ID: 44, ParentID: 40, CategoryName: "アウトドアチェア", ParentCategoryName: "折りたたみ椅子"},
		45: Category{ID: 45, ParentID: 40, CategoryName: "作業椅子", ParentCategoryName: "折りたたみ椅子"},
		50: Category{ID: 50, ParentID: 0, CategoryName: "ベンチ"},
		51: Category{ID: 51, ParentID: 50, CategoryName: "一人掛けベンチ", ParentCategoryName: "ベンチ"},
		52: Category{ID: 52, ParentID: 50, CategoryName: "二人掛けベンチ", ParentCategoryName: "ベンチ"},
		53: Category{ID: 53, ParentID: 50, CategoryName: "アウトドア用ベンチ", ParentCategoryName: "ベンチ"},
		54: Category{ID: 54, ParentID: 50, CategoryName: "収納付きベンチ", ParentCategoryName: "ベンチ"},
		55: Category{ID: 55, ParentID: 50, CategoryName: "背もたれ付きベンチ", ParentCategoryName: "ベンチ"},
		56: Category{ID: 56, ParentID: 50, CategoryName: "ベンチマーク", ParentCategoryName: "ベンチ"},
		60: Category{ID: 60, ParentID: 0, CategoryName: "座椅子"},
		61: Category{ID: 61, ParentID: 60, CategoryName: "和風座椅子", ParentCategoryName: "座椅子"},
		62: Category{ID: 62, ParentID: 60, CategoryName: "高座椅子", ParentCategoryName: "座椅子"},
		63: Category{ID: 63, ParentID: 60, CategoryName: "ゲーミング座椅子", ParentCategoryName: "座椅子"},
		64: Category{ID: 64, ParentID: 60, CategoryName: "ロッキングチェア", ParentCategoryName: "座椅子"},
		65: Category{ID: 65, ParentID: 60, CategoryName: "座布団", ParentCategoryName: "座椅子"},
		66: Category{ID: 66, ParentID: 60, CategoryName: "空気椅子", ParentCategoryName: "座椅子"},
	}
)

func getCategoryByID(q sqlx.Queryer, categoryID int) (category Category, err error) {
	c := categories[categoryID]
	if c == nil {
		return category, errors.New("not found")
	}
	category = c.(Category)
	return category, nil
}
