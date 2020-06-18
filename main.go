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

func init() {
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
	fmt.Println("Database Connection Sucess")
}

func main() {
	defer db.Close()
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
	// insertQuery := fmt.Sprintf("")
	// fmt.Println(insertQuery)
	_, err = db.Query("INSERT INTO konten(title, description, author) VALUES (?, ?, ?)", article.Title, article.Desc, article.Author)
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	fmt.Fprintf(w, "Insert Data Success")
}

func editArticle(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	kontenID := urlParams["id"]
	fmt.Fprintf(w, "id: %v\n", urlParams["id"])
	var article constant.Article
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = db.Query("UPDATE konten SET title=?,description=?,author=? WHERE id=?", article.Title, article.Desc, article.Author, kontenID)
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	fmt.Fprintf(w, "Update Data Success")
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	kontenID := urlParams["id"]
	fmt.Fprintf(w, "id: %v\n", urlParams["id"])
	_, err = db.Query("DELETE FROM konten WHERE id = ?", kontenID)
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	fmt.Fprintf(w, "Delete Data Success")
}

// Routing
func handleRequest() {
	router := mux.NewRouter()
	router.HandleFunc("/", home)
	router.HandleFunc("/article", articles)
	router.HandleFunc("/article/add", addArticle).Methods("POST")
	router.HandleFunc("/article/edit/{id:[0-9]+}", editArticle).Methods("PUT")
	router.HandleFunc("/article/delete/{id:[0-9]+}", deleteArticle).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8081", router))
}
