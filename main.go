package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

type Person struct {
	gorm.Model

	Name  string
	Email string `gorm:"typevarchar(100);unique_index"`
	Books []Book
}

type Book struct {
	gorm.Model

	Title      string
	Author     string
	CallNumber int
	PersonID   int
}

var db *gorm.DB
var err error

func main() {
	// Loading enviroment variables
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbname := os.Getenv("NAME")
	dbpassword := os.Getenv("PASSWORD")

	// Database connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbname, dbpassword, dbPort)

	// Openning connection to database
	db, err = gorm.Open(dialect, dbURI)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connected to database successfully")
	}

	defer db.Close()

	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Book{})

	router := mux.NewRouter()

	router.HandleFunc("/books", GetBooks).Methods("GET")
	router.HandleFunc("/book/{id}", GetBook).Methods("GET")
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/person/{id}", GetPerson).Methods("GET")

	router.HandleFunc("/create/person", CreatePerson).Methods("POST")
	router.HandleFunc("/create/book", CreateBook).Methods("POST")

	router.HandleFunc("/delete/person/{id}", DeletePerson).Methods("DELETE")
	router.HandleFunc("/delete/book/{id}", DeleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	var books []Book

	db.First(&person, params["id"])
	db.Model(&person).Related(&books)

	person.Books = books

	json.NewEncoder(w).Encode(&person)
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	var people []Person

	db.Find(&people)

	json.NewEncoder(w).Encode(&people)
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var person Person
	json.NewDecoder(r.Body).Decode(&person)

	createdPerson := db.Create(&person)
	err = createdPerson.Error
	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(&createdPerson)
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person Person

	db.First(&person, params["id"])
	db.Delete(&person)

	json.NewEncoder(w).Encode(&person)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var book Book

	db.First(&book, params["id"])

	json.NewEncoder(w).Encode(&book)
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book

	db.Find(&books)

	json.NewEncoder(w).Encode(&books)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	createdBook := db.Create(&book)
	err = createdBook.Error
	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(&createdBook)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var book Book

	db.First(&book, params["id"])
	db.Delete(&book)

	json.NewEncoder(w).Encode(&book)
}
