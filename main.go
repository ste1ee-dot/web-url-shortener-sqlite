package main

import (
	"database/sql"
	"html/template"
	"log"
	"math/rand"
	"net/http"

	"web-url-shortener/database"

	_ "github.com/mattn/go-sqlite3"
)

const (
	// Set your domain and DB path here:
	domain = "localhost"
	port   = "8080"

	dbPath = "./data.db"

	// Maximum lenght of shortened url path
	keyLength = 6
)

type Data struct {
	Url string
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateShortKey() string {
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

func addUrl(url string, shortUrl string, urlRepository *database.UrlRepository) {

	err := urlRepository.Insert(database.Url{OriginalUrl: url, ShortUrl: shortUrl})
	if err != nil {
		panic(err)
	}

}

func main() {

	dbConnection, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	var urlRepository = &database.UrlRepository{Db: dbConnection}
	urlRepository.CreateTable()

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.ExecuteTemplate(w, "index.gohtml", nil)
	})

	router.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		var shortUrl string
		originalUrl := r.PostFormValue("url")
		var emptyUrl database.Url
		urltype := urlRepository.GetByOriginal(originalUrl)

		if urltype == emptyUrl {
			shortUrl = generateShortKey()
			addUrl(originalUrl, shortUrl, urlRepository)

		} else {
			shortUrl = urltype.ShortUrl
		}

		tpl.ExecuteTemplate(w, "shorten.gohtml", Data{
			Url: domain + ":" + port + "/" + shortUrl,
		})
	})

	router.HandleFunc("/{shortUrl}", func(w http.ResponseWriter, r *http.Request) {
		shortURL := r.PathValue("shortUrl")

		oUrl, err := urlRepository.GetByShort(shortURL)
		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, oUrl, http.StatusMovedPermanently)
	})

	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	log.Println("Starting server on port :" + port)
	server.ListenAndServe()
}
