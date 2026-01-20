package util

import (
	"fmt"
	"strings"
)

type Course struct {
	Subject          string `json:"subject" groups:"post"`
	Number           string `json:"number" groups:"post"` // This is named out of spite
	Name             string `json:"name" groups:"post"`
	SuperDuperSecret string `json:"-"`
}

func ExampleDeserialize() {
	courseData := strings.NewReader(`{
	"subject": "CSCI",
	"number": "261",
	"name": "Introduction to Google Docs"
}`)
	deserializedCourse := Course{}
	err := Deserialize(courseData, deserializedCourse, "post")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("subject: %s; number: %s; name: %s;", deserializedCourse.Subject, deserializedCourse.Number, deserializedCourse.Name)
	// output: subject: CSCI; number: 261; name: Introduction to Google Docs;
}
