package models

type Author struct {
	Name  string `bson:"name"`
	Age   int    `bson:"age"`
	Id    string `bson:"_id" json:"_id"`
	Books []Book `bson:"books"`
}
