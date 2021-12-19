package handler

import (
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

type formBookData struct{
	Book Book
	Category []Category
	Errors map[string]string
}

type Book struct{
	ID int `db:"id"`
	Category_id int `db:"category_id"`
	Name string `db:"name"`
	Status bool `db:"status"`
	Category_Name string
}

func (b *Book) Validate() error{
	return validation.ValidateStruct(b,
		validation.Field(&b.Name, validation.Required.Error("This filed cannot be null"),
		validation.Length(3,20).Error("The Category name length must be between 3 and 20")),
	)
}



func (h *Handler) createBook( rw http.ResponseWriter, r *http.Request){
	Book := Book{}
	Errors := map[string]string{}
	
	category:=[]Category{}
	h.db.Select(&category,"Select * from categories")
	h.loadCreatedBookForm(rw,category,Book,Errors)
}

func (h *Handler) loadCreatedBookForm(rw http.ResponseWriter,cat []Category, books Book, errs map[string]string){

	form:=formBookData{
		Book: books,
			Errors: errs,
			Category:cat,
		}

		if err:= h.templates.ExecuteTemplate(rw,"create-book.html", form); err !=nil{
			http.Error(rw, err.Error(),http.StatusInternalServerError)
			return
		}

}

func (h *Handler) storeBook( rw http.ResponseWriter, r *http.Request){
	// fmt.Println("done")

	if err:=r.ParseForm(); err != nil{
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var book Book
	if err :=h.decoder.Decode(&book, r.PostForm); err != nil{

		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	category:=[]Category{}
	h.db.Select(&category,"select * from categories")

	if err:= book.Validate(); err !=nil{
		vErrors, ok := err.(validation.Errors)
		if ok {
			vErrs:=make(map[string]string)
			for key, value := range vErrors{
				vErrs[strings.Title(key)]=value.Error()

			}
			h.loadCreatedBookForm(rw,category,book,vErrs)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
		
	}

	const insertBook =`INSERT INTO books (name,category_id,status) VALUES ($1,$2,$3)`

		res:=h.db.MustExec(insertBook, book.Name, book.Category_id,false)

		if ok, err:=res.RowsAffected(); err !=nil || ok==0{
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		
		http.Redirect(rw,r,"/books", http.StatusTemporaryRedirect)
 	
}



func (h *Handler) editBook(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(rw, "invalid ", http.StatusTemporaryRedirect)
		return
	}
	const getBook = `SELECT * FROM books WHERE id = $1`
	var book Book
	h.db.Get(&book,getBook,id)

	Category:=[]Category{}
	h.db.Select(&Category,"select * from categories")

	if book.ID == 0 {
		http.Error(rw, "invalid URL", http.StatusInternalServerError)
		return
	}

	h.loadUpdateBookForm(rw,Category,book,map[string]string{})

}

func (h *Handler) updateBook(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(rw, "invalid update", http.StatusTemporaryRedirect)
		return
	}

	const getBook = `SELECT * FROM books WHERE id = $1`
	var book Book
	h.db.Get(&book, getBook, id )

	if book.ID == 0 {
		http.Error(rw, "invalid URL", http.StatusInternalServerError)
		return
	}

	if err :=r.ParseForm(); err !=nil{
		http.Error(rw, err.Error(),http.StatusInternalServerError)
		 }

		 category:=[]Category{}
	h.db.Select(&category,"select * from categories")

		if err :=h.decoder.Decode(&book, r.PostForm); err != nil{
	
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

	
		if err:= book.Validate(); err !=nil{
			vErrors, ok := err.(validation.Errors)
			if ok {
				vErrs:=make(map[string]string)
				for key, value := range vErrors{
					vErrs[strings.Title(key)]=value.Error()
	
				}
				h.loadUpdateBookForm(rw,category,book,vErrs)
				return
			}
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
			
		}

		 if book.Name == ""{
			Errors := map[string]string{
				"Name":"This filed cannot be null",
			}
			h.loadUpdateBookForm(rw,category,book,Errors)
			   return
		   }
	
		if len(book.Name) <3 {
	
			Errors := map[string]string{
				"Name":"This filed must be greater than or equals 3",
			}
			h.loadCreatedBookForm(rw,category,book,Errors)	
			return
		} 
	
		const completedBook = `UPDATE books SET name = $2, category_id=$3 WHERE id = $1`
		res:= h.db.MustExec( completedBook, id, book.Name,book.Category_id)

		if ok, err:= res.RowsAffected(); err != nil || ok == 0 {
			http.Error(rw, err.Error(),http.StatusInternalServerError)
	
			return
		}

	http.Redirect(rw,r, "/", http.StatusTemporaryRedirect)
}


func (h *Handler) deleteBook(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	if id == "" {
		http.Error(rw, "invalid update", http.StatusTemporaryRedirect)
		return
	}

	const getBook = `SELECT * FROM books WHERE id = $1`
	var book Book
	h.db.Get(&book, getBook, id )

	if book.ID == 0 {
		http.Error(rw, "invalid URL", http.StatusInternalServerError)
		return
	}

	const deleteBook =`DELETE FROM books WHERE id =$1`

	res:= h.db.MustExec( deleteBook, id)

	if ok, err:= res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(),http.StatusInternalServerError)

		return
	}


	http.Redirect(rw,r, "/books", http.StatusTemporaryRedirect)
}



func (h *Handler) loadUpdateBookForm(rw http.ResponseWriter,cat []Category, books Book, errs map[string]string){

	form:=formBookData{
		Book: books,
			Errors: errs,
			Category:cat,
		}

		if err:= h.templates.ExecuteTemplate(rw,"edit-book.html", form); err !=nil{
			http.Error(rw, err.Error(),http.StatusInternalServerError)
			return
		}

}