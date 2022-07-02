package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	// use go install [package@version] to install
	// use go mod init
	// use go mod tidy
	"github.com/jinzhu/gorm"
	// Not used directly in the code, use _, using for middleware
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	gorm.Model

	Username string
	Email string `gorm:"typevarchar(100);unique_index"`
	Items []Item
}

type Item struct {
	gorm.Model
	
	Name string
	Stat string
	ItemNumber int `gorm:"unique_index"`
	UserID  int
}

var DB *gorm.DB
var ERR error

func main() {
	// Loading environment variables
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbport := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbname := os.Getenv("NAME")
	password:= os.Getenv("PASSWORD")

	// Database connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbname, password, dbport)

	// Openning conection to database
	DB, ERR = gorm.Open(dialect, dbURI)

	if ERR != nil {
		log.Fatal(ERR)
 	} else {
		fmt.Println("Successfully connectected to database!")
	}

	// Close connection to database when the main function finishes
	// defer = do this after the function that is currently running has finished
	defer DB.Close() 

	// Make migrations to the database, if they have not already been created
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Item{})

	router := mux.NewRouter()

	// route for handling get /people
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/user/{id}", getUser).Methods("GET")
	router.HandleFunc("/create/user", createUser).Methods("POST")
	router.HandleFunc("/delete/user/{id}", deleteUser).Methods("DELETE")
	// -- books
	router.HandleFunc("/items", getItems).Methods("GET")
	router.HandleFunc("/item/{id}", getItem).Methods("GET")
	router.HandleFunc("/create/item", createItem).Methods("POST")
	router.HandleFunc("/delete/item/{id}", deleteItem).Methods("DELETE")

	// startup the server
	log.Fatal(http.ListenAndServe(":8080", router ))
}


// API Controllers
func getUsers(w http.ResponseWriter, r *http.Request) {
	var user []User
	
	// Gets all people in the database
	DB.Find(&user)

	// Transforms to JSON 
	json.NewEncoder(w).Encode(&user)

	fmt.Println("Called getUsers")
}

func getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var user User 
  var items []Item
	// only find first object with given input
	DB.First(&user, params["id"])
	// gets all books associated with this person
	DB.Model(&user).Related(&items)
	
	user.Items = items

	json.NewEncoder(w).Encode(&user)
	fmt.Println("Called getUser")
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	
	// parse incoming json payload
	json.NewDecoder(r.Body).Decode(&user)

	createdUser := DB.Create(&user)
	ERR = createdUser.Error

	if ERR != nil {
		json.NewEncoder(w).Encode(ERR)
	} else {
		json.NewEncoder(w).Encode(&user)
	}
	fmt.Println("Called createUser")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	var user User

	DB.First(&user, params["id"])
	DB.Delete(&user)

	json.NewEncoder(w).Encode(&user)
	fmt.Println("Called deleteUser")
}

func getItems(w http.ResponseWriter, r *http.Request) {
	var items []Item

	DB.Find(&items)

	json.NewEncoder(w).Encode(&items)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var item Item
	
	DB.First(&item, params["id"])
	json.NewEncoder(w).Encode(&item)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	
	json.NewDecoder(r.Body).Decode(&item)
	createdItem := DB.Create(&item)

	ERR = createdItem.Error

	if ERR != nil {
		json.NewEncoder(w).Encode(ERR)

	} else {
		json.NewEncoder(w).Encode(ERR)
	}
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var item Item

	DB.First(&item, params["id"])
	DB.Delete(&item)

	json.NewEncoder(w).Encode(&item)
}