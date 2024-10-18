package models

type User struct {
	ID       string   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string   `bson:"name" json:"name" validate:"required,max=12"`
	Password string   `bson:"password" json:"password" validate:"required,max=20"`
	Personal Personal `bson:"personal" json:"personal"`
}

type UserDataResponse struct {
	ID       string   `bson:"_id" json:"id"`
	Name     string   `bson:"name" json:"name"`
	Personal Personal `bson:"personal" json:"personal"`
}

type Personal struct {
	Email   string `bson:"email" json:"email"`
	Country string `bson:"country" json:"country"`
	Age     int    `bson:"age" json:"age"`
}
