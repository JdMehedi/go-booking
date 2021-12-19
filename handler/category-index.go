package handler

import (
	"net/http"
)

type IndexCategory struct{
   Category []Category
}

func (h *Handler) Home (rw http.ResponseWriter, r *http.Request) {
	categories := []Category{}

    h.db.Select(&categories, "SELECT * FROM categories")

	lt := IndexCategory{
		Category:categories,
	}

	if err:= h.templates.ExecuteTemplate(rw,"index-category.html", lt); err !=nil{
		http.Error(rw, err.Error(),http.StatusInternalServerError)
		return
	}
}
func (h *Handler)searchCategory(rw http.ResponseWriter, r *http.Request){

	if err:=r.ParseForm(); err !=nil{
		http.Error(rw, err.Error(),http.StatusInternalServerError)
		return
	}

	res:=r.FormValue("search")

	const searchValue = `Select * FROM categories WHERE name ILIKE '%%' || $1 || '%%'`
	var cat []Category
	h.db.Select(&cat,searchValue,res)
	lt := IndexCategory{
		Category:cat,
	}
	// fmt.Println(lt)
	if err:= h.templates.ExecuteTemplate(rw,"index-category.html", lt); err !=nil{
		http.Error(rw, err.Error(),http.StatusInternalServerError)
		return
	}



}
