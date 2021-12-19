package handler

import (
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

type Handler struct{
   templates  *template.Template 
   db 	      *sqlx.DB
   decoder    *schema.Decoder
   sess       *sessions.CookieStore

}
var sessionName = "storeCookie"

func New(db *sqlx.DB, decoder  *schema.Decoder ,sess  *sessions.CookieStore) *mux.Router{
	h:= &Handler{
		db: db,
		decoder: decoder,
		sess: sess,
	}

	h.parseTemplate()
	r :=mux.NewRouter() 
	l:=r.NewRoute().Subrouter()
	l.HandleFunc("/signUp", h.signUp)
	l.HandleFunc("/signUp/register", h.register)

	l.HandleFunc("/login", h.login)
	l.HandleFunc("/login/check", h.loginCheck)
	l.Use(h.loginmiddleWare)
	r.HandleFunc("/logout", h.logout)
	r.HandleFunc("/", h.Index)

    s:=r.NewRoute().Subrouter()

	s.HandleFunc("/categories/create", h.createCategory)
	s.HandleFunc("/categories/store", h.storeCategory)
	s.HandleFunc("/categories", h.Home)
	s.HandleFunc("/categories/{id}/edit", h.editCategory)
	s.HandleFunc("/categories/{id}/update", h.updateCategory)
	s.HandleFunc("/categories/{id}/delete", h.deleteCategory)
	s.HandleFunc("/categories/search", h.searchCategory)


	s.HandleFunc("/books/create", h.createBook)
	s.HandleFunc("/books/store", h.storeBook)
	
	s.HandleFunc("/books/{id}/edit", h.editBook)
	s.HandleFunc("/books/{id}/update", h.updateBook)
	s.HandleFunc("/books/{id}/delete", h.deleteBook)
	s.HandleFunc("/books/search", h.searchBook)

	s.HandleFunc("/booking/{id}/create", h.createBooking)
	s.HandleFunc("/booking/store", h.storeBooking)
	s.HandleFunc("/booking/{id}/avilable", h.availableBooking)

	s.HandleFunc("/bookings", h.Booking)

	s.PathPrefix("/asset/").Handler(http.StripPrefix("/asset/", http.FileServer(http.Dir("./"))))
	s.Use(h.middleWare)
	

	r.NotFoundHandler = http.HandlerFunc(func (rw http.ResponseWriter, r *http.Request)  {
		if err :=h.templates.ExecuteTemplate(rw,"404.html",nil); err != nil{
			http.Error(rw, err.Error(),http.StatusInternalServerError)
			return
		}
		
	} )

	return r
	
}

func (h *Handler) middleWare(next http.Handler) http.Handler{

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		session, _:= h.sess.Get(r, sessionName)
	
		auth,ok:=session.Values["Authenticated"].(bool)
		if !ok || !auth{
			http.Redirect(rw,r,"/login",http.StatusTemporaryRedirect)
			return
		}

		next.ServeHTTP(rw,r)
	})
}

func (h *Handler) loginmiddleWare(next http.Handler) http.Handler{

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		session, _:= h.sess.Get(r, sessionName)
	
		authUserID :=session.Values["Authenticated"]
		if authUserID != nil{
			http.Redirect(rw,r,"/",http.StatusTemporaryRedirect)
			return			
		}
		next.ServeHTTP(rw,r)
	})
}

func (h *Handler) parseTemplate(){
	h.templates= template.Must(template.ParseFiles(
		"templates/category/create-category.html",
		"templates/category/index-category.html",
		"templates/category/edit-category.html",
		"templates/book/create-book.html",
		"templates/book/index-book.html",
		"templates/book/edit-book.html",
		"templates/booking/create-booking.html",
		"templates/booking/index-booking.html",
		"templates/404.html",
		"templates/login/login.html",
		"templates/registration/signUp.html",
	))
}

