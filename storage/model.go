package storage

type User struct {
	Id       string   `bson:"_id,omitempty" json:"id,omitempty"`
	Username string   `bson:"username"`
	Password string   `bson:"password"`
	Personal Personal `bson:"personal"`
}

type Personal struct {
	Email   string `bson:"email"`
	Country string `bson:"country"`
	Age     int    `bson:"age"`
}
