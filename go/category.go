package main

import "github.com/jmoiron/sqlx"

func getCategoryByID(q sqlx.Queryer, categoryID int) (category Category, err error) {
	err = sqlx.Get(q, &category, "SELECT * FROM `categories` WHERE `id` = ?", categoryID)
	if category.ParentID != 0 {
		parentCategory, err := getCategoryByID(q, category.ParentID)
		if err != nil {
			return category, err
		}
		category.ParentCategoryName = parentCategory.CategoryName
	}
	return category, err
}
