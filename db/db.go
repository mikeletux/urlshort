package db

// Database is the interface that must be implemented to be used by the handler
type Database interface {
	Insert(shortURL, longURL string) error
	GetFullURL(shortURL string) (longURL string, err error)
}
