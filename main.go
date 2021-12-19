package main

import (
	"log"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"main.go/handler"
)


var sessionName= "book-session"

func main(){
	var createTable = `
	CREATE TABLE IF NOT EXISTS users (
		id serial,
		first_name 			text,
		last_name 			text,
		username 			text,
		email 				text,
		password 			text,

		primary key(id)
	);
	CREATE TABLE IF NOT EXISTS categories (
		id serial,
		name text,
		image text,

		primary key(id)
	);
	CREATE TABLE IF NOT EXISTS books (
		id serial,
		category_id integer,
		name text,
		status boolean,
		
		primary key(id)
	
		
	);
	CREATE TABLE IF NOT EXISTS bookings (
		id serial,
		user_id integer,
		book_id integer,
		start_time timestamp,
		end_time   timestamp,


		primary key(id)	
	)`

	db, err := sqlx.Connect("postgres", "user=postgres password=Passw0rd dbname=book sslmode=disable")
    if err != nil {
        log.Fatalln(err)
    }

	db.MustExec(createTable)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	var store = sessions.NewCookieStore([]byte(securecookie.GenerateRandomKey(32)))
	r:=handler.New(db, decoder, store)


     log.Println("Server Starting.............")

	if err := http.ListenAndServe("127.0.0.1:3000",r); err !=nil{
		log.Fatal(err)
	}

}