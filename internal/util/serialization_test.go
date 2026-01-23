package util

import (
	"fmt"
	"strings"
)

type Course struct {
	Subject          string `json:"subject" groups:"post"`
	Number           string `json:"number" groups:"post"` // This is named out of spite
	Name             string `json:"name" groups:"post"`
	SuperDuperSecret string `json:"superDuperSecret"`
}

func ExampleDeserialize() {
	courseData := strings.NewReader(`{
	"subject": "CSCI",
	"number": "261",
	"name": "Introduction to Google Docs"
}`)
	err, deserializedCourse := Deserialize[Course](courseData, "post")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(deserializedCourse)
	// output: {CSCI 261 Introduction to Google Docs }
}
