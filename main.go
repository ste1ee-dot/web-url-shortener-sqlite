package main

import (
	"encoding/csv"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
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

func loadUrls(links *[]Link) {
	file, err := os.Open(csvPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil { panic(err) }

	for _, record := range records {
		link := Link{
			ShortUrl: record[0],
			Url: record[1],
		}

		*links = append(*links, link)
	}

}

func addUrl(shortUrl string, url string) {
	file, err := os.OpenFile(csvPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	record := []string{shortUrl, url}

	err = writer.Write(record)
	if err != nil {
		panic(err)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		panic(err)
	}
}

func main() {

	var links []Link
	loadUrls(&links)

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.ExecuteTemplate(w, "index.gohtml", nil)
	})

	var url string
	var shortUrl string

	router.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		url	= r.PostFormValue("url")
		shortUrl = generateShortKey() 
		addUrl(shortUrl, url)

		link := Link {
			ShortUrl: shortUrl,
			Url: url,
		}

		links = append(links, link)

		tpl.ExecuteTemplate(w, "shorten.gohtml", Data{
			Url: domain + ":" + port + "/" + shortUrl,

		})
	})

	router.HandleFunc("/{shortUrl}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	})

	server := http.Server{
		Addr: ":" + port,
		Handler: router,
	}
	log.Println("Starting server on port :" + port)
	server.ListenAndServe()
}

