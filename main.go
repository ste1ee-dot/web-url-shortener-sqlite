package main

import (
	"database/sql"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
	"web-url-shortener/database"

	_ "github.com/mattn/go-sqlite3"
)

// Set your domain and csv path here:
const domain = "localhost"
const port = "8080"

const csvPath = "data.csv"
// Maximum lenght of shortened url path
const keyLength = 6


type Data struct {
	Url string
}

type Link struct {
	ShortUrl string `csv:"shorturl"`
	Url string `csv:"url"`
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateShortKey() string {
	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func addUrl(url string, urlRepository *database.UrlRepository) string {
	emptyUrl := database.Url{}
	onUrl := urlRepository.GetByOriginal(url)
	if onUrl != emptyUrl {
		return onUrl.ShortUrl
	}
	shortUrl := generateShortKey()
	urlRepository.Insert(database.Url{OriginalUrl: url, ShortUrl: shortUrl})

	return shortUrl
}

func main() {

	dbConnection, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		panic(err)
	}

	var urlRepository = &database.UrlRepository{Db: dbConnection}
	err = urlRepository.CreateTable()
	if err != nil {
		panic(err)
	}

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.ExecuteTemplate(w, "index.gohtml", nil)
	})

	router.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		url	:= r.PostFormValue("url")
		shortUrl := addUrl(url, urlRepository)

		tpl.ExecuteTemplate(w, "shorten.gohtml", Data{
			Url: domain + ":" + port + "/" + shortUrl,

		})
	})

	router.HandleFunc("/{shortUrl}", func(w http.ResponseWriter, r *http.Request) {
		shortUrl := r.URL.Path[1:]
		oUrl, err := urlRepository.GetByShort(shortUrl)
		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, oUrl, http.StatusMovedPermanently)
	})

	server := http.Server{
		Addr: ":" + port,
		Handler: router,
	}
	log.Println("Starting server on port :" + port)
	server.ListenAndServe()
}

