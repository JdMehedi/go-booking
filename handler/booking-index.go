package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

type IndexBooking struct{
   Booking []Booking
}

func (h *Handler) Booking (rw http.ResponseWriter, r *http.Request) {
	bookings := []Booking{}
	

    h.db.Select(&bookings, "SELECT * FROM bookings ORDER BY id DESC")
	for key, value := range bookings {
		const getBook = `SELECT name FROM books WHERE id=$1`
		var book Book
		h.db.Get(&book, getBook, value.Book_id)
		start_time:= value.Start_time.Format("Mon Jan _2 2006 15:04 AM")
		end_time:= value.End_time.Format("Mon Jan _2 2006 15:04 AM")
		bookings[key].Book_name = book.Name
		bookings[key].S_time = start_time
		bookings[key].E_time = end_time
		
	}
	lt := IndexBooking{
		Booking:bookings,
	}
	if err:= h.templates.ExecuteTemplate(rw,"index-booking.html", lt); err !=nil{
		http.Error(rw, err.Error(),http.StatusInternalServerError)
		return
	}
}

func (h *Handler) availableBooking(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	const availableStatus = `UPDATE books SET status = false WHERE id = $1`
	h.db.MustExec( availableStatus, id)

	http.Redirect(rw,r,"/",http.StatusTemporaryRedirect)
}
