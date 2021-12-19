package handler

import (
	"net/http"
)

type IndexBook struct{
   Book []Book
   Booking []Booking
   Category []Category
}

func (h *Handler) Index (rw http.ResponseWriter, r *http.Request) {
	books := []Book{}

    h.db.Select(&books, "SELECT * FROM books")
	for key, value := range books {
		const getCategory = `SELECT name FROM categories WHERE id=$1`
		var category Category
		h.db.Get(&category, getCategory, value.Category_id)
		books[key].Category_Name = category.Name
	}
	
	categories := []Category{}

    h.db.Select(&categories, "SELECT * FROM categories")

	lt := IndexBook{
		Book:books,
		Category:categories,
	}
	//  fmt.Println(lt)
	if err:= h.templates.ExecuteTemplate(rw,"index-book.html", lt); err !=nil{
		http.Error(rw, err.Error(),http.StatusInternalServerError)
		return
	}
}

func (h *Handler)searchBook(rw http.ResponseWriter, r *http.Request){

	if err:=r.ParseForm(); err !=nil{
		http.Error(rw, err.Error(),http.StatusInternalServerError)
		return
	}
	res:=r.FormValue("search")

	const searchValue = `Select * FROM books WHERE name ILIKE '%%' || $1 || '%%'`
	var book []Book
	h.db.Select(&book,searchValue,res)
	
	for key, value := range book {
		const getCategory = `SELECT name FROM categories WHERE id=$1`
		var category Category
		h.db.Get(&category, getCategory, value.Category_id)
		book[key].Category_Name = category.Name
	}
	lt := IndexBook{
		Book:book,
	}
	if err:= h.templates.ExecuteTemplate(rw,"index-book.html", lt); err !=nil{
		http.Error(rw, err.Error(),http.StatusInternalServerError)
		return
	}

}
