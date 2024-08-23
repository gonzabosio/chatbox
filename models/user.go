package models

type User struct {
	Id       string   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string   `bson:"name" json:"name"`
	Password string   `bson:"password" json:"password"`
	Personal Personal `bson:"personal" json:"personal"`
}

type Personal struct {
	Email   string `bson:"email"`
	Country string `bson:"country"`
	Age     int    `bson:"age"`
}
