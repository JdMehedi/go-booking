package handler

import (

	"log"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

type SignUpFormData struct {
	SignUpForm SignUpForm
	Errors   map[string]string
}

type SignUpForm struct {
	ID 				int    `db:"id"`
	FirstName 		string `db:"first_name"`
	LastName 		string	`db:"last_name"`
	UserName 		string	`db:"username"`
	Password 		string	`db:"password"`
	Email    		string	`db:"email"`
	ConfirmPassword string	
}

func (s *SignUpForm) validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.FirstName,
			validation.Required.Error("This filed cannot be null"),
			validation.Length(3, 15).Error("The Firstname length must be between 3 and 15"),
		),
		validation.Field(&s.LastName,
			validation.Required.Error("This filed cannot be null"),
			validation.Length(3, 15).Error("The Lastname length must be between 3 and 15"),
		),
		validation.Field(&s.UserName,
			validation.Required.Error("This filed cannot be null"),
			validation.Length(3, 10).Error("The Username length must be between 3 and 15"),
		),
		validation.Field(&s.Email,
			validation.Required.Error("This filed cannot be null"),
		),
	)
}

func (h *Handler) signUp(rw http.ResponseWriter, r *http.Request) {
	form := SignUpForm{}
	Errors:=map[string]string{}
	h.loadregisterForm(rw,form,Errors)

}

func (h *Handler) register(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}
	form := SignUpForm{}
	if err := h.decoder.Decode(&form, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if form.Password != form.ConfirmPassword{
		Errors:= map[string]string{"Password":"The given password does not match with previous password"}
		h.loadregisterForm(rw,form,Errors)
	}
	if err := form.validate(); err != nil {
		VErrors, ok := err.(validation.Errors)
		if ok {
			Errs := map[string]string{}
			for key, value := range VErrors {
				Errs[strings.Title(key)] = value.Error()
			}

			h.loadregisterForm(rw,form,Errs)
			return
		}
	}

	const registration =`INSERT INTO users (first_name,last_Name,username,email,password) VALUES ($1,$2,$3,$4,$5)`
	pass, err:= bcrypt.GenerateFromPassword([]byte(form.Password),bcrypt.DefaultCost); 
	if err!=nil{
		log.Fatal(err)
	}
	res:=h.db.MustExec(registration, form.FirstName,form.LastName,form.UserName,form.Email,string(pass))
	if ok, err:=res.RowsAffected(); err !=nil || ok==0{
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	sessions, err :=h.sess.Get(r,sessionName )
	if err!=nil{
		log.Fatal(err)
	}
	sessions.AddFlash("Registration Successfully")
	if err:=sessions.Save(r, rw); err!=nil{
		log.Fatal(err)
	}
	http.Redirect(rw,r,"/login", http.StatusTemporaryRedirect)
}

func (h *Handler) loadregisterForm(rw http.ResponseWriter,signupform SignUpForm,errors map[string]string){

	form:= SignUpFormData{
		SignUpForm: signupform,
		Errors: errors,
	}
	if err:= h.templates.ExecuteTemplate(rw,"signUp.html", form); err !=nil{
		http.Error(rw, err.Error(),http.StatusInternalServerError)
		return
	}

}
