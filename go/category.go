package main

var (
	CategoryDict       map[int]Category
	Categories         []Category
	CategoriesInParent map[int][]int
)

func InitCategory() {
	CategoryDict = map[int]Category{}
	Categories = []Category{}
	CategoriesInParent = map[int][]int{}
	var categories []Category
	err := dbx.Select(&categories, "SELECT c1.*, c2.category_name as parent_name FROM categories c1 LEFT JOIN categories c2 ON c1.parent_id = c2.id")
	if err != nil {
		logger.Error(err)
		return
	}
	for _, v := range categories {
		CategoryDict[v.ID] = v
		CategoriesInParent[v.ParentID] = append(CategoriesInParent[v.ParentID], v.ID)
	}
	Categories = categories
	return
}
