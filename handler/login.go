package handler

import (
	"log"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

type LoginForm struct {
	UserName string
	Password string
	Errors   map[string]string
	Messages string
}

func (l *LoginForm) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.UserName,
			validation.Required.Error("This filed cannot be null"),
			validation.Length(3, 15).Error("The User name length must be between 3 and 15"),
		),
		validation.Field(&l.Password,
			validation.Required.Error("This filed cannot be null"),
			validation.Length(6, 25).Error("The Password length must be between 6 and 25"),
		),
	)
}

func (h *Handler) login(rw http.ResponseWriter, r *http.Request) {
	sessions, err :=h.sess.Get(r,sessionName )
	if err!=nil{
		log.Fatal(err)
	}
	form := LoginForm{}
	if flashes := sessions.Flashes(); len(flashes) > 0{
		if val, ok:=flashes[0].(string); ok{
			form.Messages = val
		}
	
		if err:=sessions.Save(r, rw); err!=nil{
			log.Fatal(err)
		}
	}
	if err := h.templates.ExecuteTemplate(rw, "login.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) loginCheck(rw http.ResponseWriter, r *http.Request) {
	
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}
	form := LoginForm{}
	if err := h.decoder.Decode(&form, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := form.Validate(); err != nil {
		VErrors, ok := err.(validation.Errors)

		if ok {
			Errs := map[string]string{}
			for key, value := range VErrors {
				Errs[strings.Title(key)] = value.Error()
			}

			h.loadLoginForm(rw, form.UserName, Errs)
			return
		}
	}

	userCheck := `SELECT * FROM users WHERE username =$1`
	var user SignUpForm
	h.db.Get(&user, userCheck, form.UserName)
	if user.UserName == "" {
		Errs := map[string]string{"UserName": "The User name does not match"}

		h.loadLoginForm(rw, form.UserName, Errs)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		Errs := map[string]string{"Password": "The password does not match"}

		h.loadLoginForm(rw, form.UserName, Errs)
		return
	}

	session, err := h.sess.Get(r, sessionName)
	if err != nil {
		log.Fatal(err)
	}
	session.Options.HttpOnly = true
	session.Values["Authenticated"] = true

	session.Save(r, rw)

	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handler) loadLoginForm(rw http.ResponseWriter, userName string, errors map[string]string) {

	form := LoginForm{
		UserName: userName,
		Errors:   errors,
	}

	if err := h.templates.ExecuteTemplate(rw, "login.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) logout(rw http.ResponseWriter, r *http.Request) {

	session, err :=h.sess.Get(r,sessionName )
	if err!=nil{
		log.Fatal(err)
	}
	// session.Options.HttpOnly = false
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1
	session.Save(r, rw)

	http.Redirect(rw, r, "/login", http.StatusTemporaryRedirect)

}
