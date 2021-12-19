package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

type formData struct{
	Category Category
	Errors map[string]string
}

type Category struct{
	ID int `db:"id"`
	Name string `db:"name"`
	Image string `db:"image"` 
}

func (c *Category) validate() error{

	return validation.ValidateStruct(c,
		validation.Field(&c.Name,
			 validation.Required.Error("This filed cannot be null"),
			 validation.Length(3,30).Error("The Category name length must be between 3 and 30"),
			),
	)
}


func (h *Handler) createCategory( rw http.ResponseWriter, r *http.Request){
	Category := Category{}
	Errors := map[string]string{}
	
	h.loadCreatedCategoryForm(rw,Category,Errors)
}
func (h *Handler) storeCategory( rw http.ResponseWriter, r *http.Request){

	if err:=r.ParseMultipartForm(10 << 20); err != nil{
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var category Category
	if err :=h.decoder.Decode(&category, r.PostForm); err != nil{

		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

    file, _, err := r.FormFile("Image")
	fmt.Println(file)
    if err != nil {
        fmt.Println("Error Retrieving the File")
        fmt.Println(err)
        return
    }
    defer file.Close()
   
	var img = "upload-*.png"
    tempFile, err := ioutil.TempFile("assets/images", img)
	fmt.Printf("%T",img)
    if err != nil {
        fmt.Println(err)
    }
    defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
        fmt.Println(err)
    }
	
    tempFile.Write(fileBytes)
	a := tempFile.Name()



	if err:= category.validate(); err !=nil{  
		vErrors, ok := err.(validation.Errors)
		if ok {
			vErrs:=make(map[string]string)
			for key, value := range vErrors{
				vErrs[strings.Title(key)]=value.Error()

			}
			h.loadCreatedCategoryForm(rw,category,vErrs)
			return
		}

		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
		
	}

	const insertCategory =`INSERT INTO categories (name,image) VALUES ($1,$2)`

		res:=h.db.MustExec(insertCategory, category.Name,a)

		if ok, err:=res.RowsAffected(); err !=nil || ok==0{
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		
		http.Redirect(rw,r,"/categories", http.StatusTemporaryRedirect)
 	
}

func (h *Handler) loadCreatedCategoryForm(rw http.ResponseWriter, categories Category, errs map[string]string){

	form:=formData{
		Category: categories,
			Errors: errs,
		}
		if err:= h.templates.ExecuteTemplate(rw,"create-category.html", form); err !=nil{
			http.Error(rw, err.Error(),http.StatusInternalServerError)
			return
		}

}

func (h *Handler) editCategory(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(rw, "invalid ", http.StatusTemporaryRedirect)
		return
	}
	const getCategory = `SELECT * FROM categories WHERE id = $1`
	fmt.Println(getCategory)
	var category Category
	h.db.Get(&category, getCategory, id )
	fmt.Println(category)

	if category.ID == 0 {
		http.Error(rw, "invalid URL", http.StatusInternalServerError)
		return
	}

	h.loadUpdateCategoryForm(rw,category,map[string]string{})

}

func (h *Handler) updateCategory(rw http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(rw, "invalid update", http.StatusTemporaryRedirect)
		return
	}

	const getCategory = `SELECT * FROM categories WHERE id = $1`

	var category Category

	h.db.Get(&category, getCategory, id )

	if category.ID == 0 {
		http.Error(rw, "invalid URL", http.StatusInternalServerError)
		return
	}


	if err :=r.ParseForm(); err !=nil{
		http.Error(rw, err.Error(),http.StatusInternalServerError)
		 }
		 var form Category

		 if err :=h.decoder.Decode(&form, r.PostForm); err != nil{
	 
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		//  var cat Category
		
		 fmt.Printf("%#v",form)

		 if err:= form.validate(); err !=nil{
			vErrors, ok := err.(validation.Errors)
			if ok {
				vErrs:=make(map[string]string)
				for key, value := range vErrors{
					vErrs[strings.Title(key)]=value.Error()
	
				}
				h.loadUpdateCategoryForm(rw,category,vErrs)
				return
			}
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
			
		}
	
		const completedCategory = `UPDATE categories SET name = $2 WHERE id = $1`
		res:= h.db.MustExec( completedCategory, id, form.Name)

		if ok, err:= res.RowsAffected(); err != nil || ok == 0 {
			http.Error(rw, err.Error(),http.StatusInternalServerError)
	
			return
		}

	http.Redirect(rw,r, "/categories", http.StatusTemporaryRedirect)
}


func (h *Handler) deleteCategory(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// fmt.Println(vars)
	id := vars["id"]
	
	if id == "" {
		http.Error(rw, "invalid update", http.StatusTemporaryRedirect)
		return
	}

	const getcategory = `SELECT * FROM categories WHERE id = $1`
	var category Category
	h.db.Get(&category, getcategory, id )

	if category.ID == 0 {
		http.Error(rw, "invalid URL", http.StatusInternalServerError)
		return
	}

	const deleteCategory =`DELETE FROM categories WHERE id =$1`

	res:= h.db.MustExec( deleteCategory, id)

	if ok, err:= res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(),http.StatusInternalServerError)

		return
	}


	http.Redirect(rw,r, "/categories", http.StatusTemporaryRedirect)
}




func (h *Handler) loadUpdateCategoryForm(rw http.ResponseWriter, categories Category, errs map[string]string){

	form:=formData{
		Category: categories,
			Errors: errs,
		}
		fmt.Println(form)
  
		if err:= h.templates.ExecuteTemplate(rw,"edit-category.html", nil); err !=nil{
			http.Error(rw, err.Error(),http.StatusInternalServerError)
			return
		}

}