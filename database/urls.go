package database

import "database/sql"

type UrlRepository struct {
	Db *sql.DB
}

type Url struct {
	OriginalUrl string
	ShortUrl string
}

func (r *UrlRepository) CreateTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS urls (
		originalUrl TEXT PRIMARY KEY UNIQUE,
		shortUrl TEXT UNIQUE
	)`)	
	return err
}

func (r *UrlRepository) Insert(url Url) error {
	_, err := r.Db.Exec("INSERT INTO urls (originalUrl, shortUrl) VALUES (?, ?)",
	url.OriginalUrl, url.ShortUrl)
	return err
}

func (r *UrlRepository) GetByOriginal(oUrl string) Url {
	var url Url
	err := r.Db.
	QueryRow("SELECT originalUrl, shortUrl FROM urls WHERE originalUrl = ?",
	oUrl).Scan(&url.OriginalUrl, &url.ShortUrl)
	if err != nil {
		return Url{}
	}

	return url
}

func (r *UrlRepository) GetByShort(sUrl string) (string, error) {
	var url Url
	err := r.Db.
	QueryRow("SELECT originalUrl, shortUrl FROM urls WHERE shortUrl = ?",
	sUrl).Scan(&url.OriginalUrl, &url.ShortUrl)
	if err != nil {
		return "", err
	}

	return url.OriginalUrl, nil
}
