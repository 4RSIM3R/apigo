package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/4RSIM3R/belajar_golang/constant"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var db *sql.DB
var err error

func main() {
	godotErr := godotenv.Load()
	if godotErr != nil {
		log.Fatal("Error loading .env file")
	}

	dbConnection := os.Getenv("DB_CONNECTION")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbDatabase := os.Getenv("DB_DATABASE")
	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	fmt.Println(dbConnection)
	fmt.Println(dbHost)
	fmt.Println(dbPort)
	fmt.Println(dbDatabase)
	fmt.Println(dbUser)
	fmt.Println(dbPassword)

	dbInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbDatabase)
	db, err = sql.Open(dbConnection, dbInfo)
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	fmt.Println("Database Connection Sucess")

	handleRequest()
}

// Controller
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from index")
}

func articles(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var results []constant.Article

	selectQuery := "SELECT title, description, author FROM konten"

	result, err := db.Query(selectQuery)

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	for result.Next() {
		var article constant.Article
		err := result.Scan(&article.Title, &article.Desc, &article.Author)
		if err != nil {
			panic(err.Error())
		}
		results = append(results, article)
	}

	json.NewEncoder(w).Encode(results)

	// contents := constant.Articles{
	// 	constant.Article{Title: "Title", Desc: "Desc", Author: "Author"},
	// 	constant.Article{Title: "Judul", Desc: "Deskripis", Author: "Pemilik"},
	// }
	// json.NewEncoder(w).Encode(contents)
}

func addArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST endpoint test")

	var article constant.Article
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Author : ", article.Author)
	insertQuery := fmt.Sprintf("INSERT INTO konten(title, description, author) VALUES ('%s', '%s', '%s')", article.Title, article.Desc, article.Author)
	fmt.Println(insertQuery)
	insert, err := db.Query(insertQuery)
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	defer insert.Close()
}

// Routing
func handleRequest() {
	router := mux.NewRouter()
	router.HandleFunc("/", home)
	router.HandleFunc("/article", articles)
	router.HandleFunc("/article/add", addArticle).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", router))
}
