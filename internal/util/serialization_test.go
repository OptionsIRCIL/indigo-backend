package util

import (
	"fmt"
	"strings"
)

type Student struct {
	Ssn             string `json:"ssn"`
	FirstName       string `json:"firstName" groups:"post"`
	LastName        string `json:"lastName" groups:"post"`
	FavoriteNumbers []int  `json:"favoriteNumbers" groups:"post"`
}

type Course struct {
	Subject          string    `json:"subject" groups:"post"`
	Number           string    `json:"number" groups:"post"` // This is named out of spite
	Name             string    `json:"name" groups:"post"`
	SuperDuperSecret string    `json:"superDuperSecret"`
	Students         []Student `json:"students" groups:"post"`
}

func ExampleDeserialize() {
	courseData := strings.NewReader(`{
	"subject": "CSCI",
	"number": "261",
	"name": "Introduction to Google Docs",
	"students": [
		{
			"firstName": "Carl",
			"lastName": "Weezer",
            "favoriteNumbers": [1, 2, 3]
		},
		{
			"firstName": "Jimmy",
			"lastName": "Neutron"
		}
	]
}`)
	err, deserializedCourse := Deserialize[Course](courseData, "post")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(deserializedCourse)
	// output: {CSCI 261 Introduction to Google Docs }
}
