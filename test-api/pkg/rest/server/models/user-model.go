package models

type User struct {
	ID string `json:"id,omitempty" bson:"_id,omitempty"`

	Age int8 `json:"age,omitempty" bson:"age,omitempty"`

	Name string `json:"name,omitempty" bson:"name,omitempty"`
}
