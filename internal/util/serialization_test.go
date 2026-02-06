package util

import (
	"fmt"
	"strings"
	"testing"
)

type Student struct {
	Ssn       string
	FirstName string `groups:"post,get"`
	LastName  string `groups:"post,get"`
}

type Course struct {
	Subject    string `groups:"post"`
	Number     string `groups:"post"` // Cry about it
	Name       string `groups:"post"`
	Sections   []int  `groups:"post"`
	EnrollCode string
	Students   []Student `groups:"post"`
}

func ExampleDeserialize() {
	goodStudentData := strings.NewReader(`{ "firstName": "Jimmy", "lastName": "Neutron" }`)
	badStudentData := strings.NewReader(`{ "firstName": "Carl", "lastName": "Wheezer", "ssn": "123123123" }`)

	goodErr, goodStudent := Deserialize[Student](goodStudentData, []string{"post"})
	badErr, badStudent := Deserialize[Student](badStudentData, []string{"post"})

	if goodErr == nil {
		fmt.Printf("goodStudent Successfully Deserialized! Data: %s\n", goodStudent)
	} else {
		fmt.Printf("goodStudent Failed Deserialization. Error: %s\n", goodErr)
	}

	if badErr == nil {
		fmt.Printf("badStudent Successfully Deserialized! Data: %s\n", badStudent)
	} else {
		fmt.Printf("badStudent Failed Deserialization. Error: %s\n", badErr)
	}

	// output:
	// goodStudent Successfully Deserialized! Data: { Jimmy Neutron}
	// badStudent Failed Deserialization. Error: json: unknown field "ssn"
}

func TestDeserialize(t *testing.T) {
	// Case 1 - Basic struct
	case1Error, _ := Deserialize[Student](
		strings.NewReader(`{
			"firstName": "Jimmy",
			"lastName": "Neutron"
		}`),
		[]string{"post"},
	)
	if case1Error != nil {
		t.Error("Deserialization unexpectedly failed for case 1")
	}

	// Case 2 - Basic struct, but masked property passed
	case2Error, _ := Deserialize[Course](
		strings.NewReader(`{
			"firstName": "Jimmy",
			"lastName": "Neutron",
			"ssn": "123123123"
		}`),
		[]string{"post"},
	)
	if case2Error == nil {
		t.Error("Deserialization unexpectedly succeeded for case 2")
	}

	// Case 3 - Complex struct
	case3Error, _ := Deserialize[Course](
		strings.NewReader(`{
			"subject": "CSCI",
			"number": "261",
			"name": "Introduction to Google Docs",
			"students": [
				{
					"firstName": "Carl",
					"lastName": "Wheezer"
				},
				{
					"firstName": "Jimmy",
					"lastName": "Neutron"
				}
			],
			"sections": [123, 456, 789]
		}`),
		[]string{"post"},
	)
	if case3Error != nil {
		t.Error("Deserialization unexpectedly failed for case 3")
	}

	// Case 4 - Complex struct, bad nested property
	case4Error, _ := Deserialize[Course](
		strings.NewReader(`{
			"subject": "CSCI",
			"number": "261",
			"name": "Introduction to Google Docs",
			"students": [
				{
					"firstName": "Carl",
					"lastName": "Wheezer"
				},
				{
					"firstName": "Jimmy",
					"lastName": "Neutron",
					"ssn": "123123123"
				}
			],
			"sections": [123, 456, 789]
		}`),
		[]string{"post"},
	)
	if case4Error == nil {
		t.Error("Deserialization unexpectedly succeeded for case 4")
	}
}

func ExampleSerialize() {
	student := Student{
		Ssn:       "123123123",
		FirstName: "Mercedes",
		LastName:  "Benz",
	}

	serializedStudent, _ := Serialize(student, []string{"get"})
	fmt.Println(string(serializedStudent))

	// output: {"firstName":"Mercedes","lastName":"Benz"}
}
