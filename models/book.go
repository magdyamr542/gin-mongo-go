package models

type Book struct {
	Name  string `bson:"name"`
	Price int    `bson:"price"`
	_id   string
}
