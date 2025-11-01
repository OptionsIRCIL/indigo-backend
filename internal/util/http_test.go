package util

import (
	"bytes"
	"fmt"
	"log"
	"net/http/httptest"
)

type basicError struct {
	Msg string
}

func (e *basicError) Error() string {
	return e.Msg
}

type examplePayload struct {
	Hostname     string `json:"hostname"`
	Ip           string `json:"ip"`
	PackageCount int    `json:"package_count"`
}

func ExampleReturnSerialized() {
	w := httptest.NewRecorder()
	payload := examplePayload{
		"ORGSERVERPROD1",
		"172.16.100.67",
		175,
	}

	ReturnSerialized(w, 200, &payload)
	fmt.Printf(
		"Body=%s, Code=%d, Content-Type=%s",
		w.Body,
		w.Code,
		w.Header().Get("Content-Type"),
	)
	// Output: Body={"hostname":"ORGSERVERPROD1","ip":"172.16.100.67","package_count":175}, Code=200, Content-Type=application/json
}

func ExampleThrowHttpUnhandled() {
	w := httptest.NewRecorder()
	buf := bytes.Buffer{}

	log.SetOutput(&buf)
	ThrowHttpUnhandled(w, &basicError{"I WARNED YOU ABOUT STAIRS BRO!!!! I TOLD YOU DOG!"})
	fmt.Printf("Log=%s, Body=%s", buf.String()[20:buf.Len()-1], w.Body)
	// Output: Log=Unhandled error: I WARNED YOU ABOUT STAIRS BRO!!!! I TOLD YOU DOG!, Body={"status":500,"message":"Internal server error"}
}

func ExampleThrowHttpError() {
	w := httptest.NewRecorder()

	ThrowHttpError(w, 419, "I'm a teapot")

	fmt.Printf(
		"Body=%s, Code=%d, Content-Type=%s",
		w.Body,
		w.Code,
		w.Header().Get("Content-Type"),
	)
	// Output: Body={"status":419,"message":"I'm a teapot"}, Code=419, Content-Type=application/json
}

func ExampleThrowHttpStatus() {
	w := httptest.NewRecorder()

	ThrowHttpStatus(w, 500)

	fmt.Printf(
		"Body=%s, Code=%d, Content-Type=%s",
		w.Body,
		w.Code,
		w.Header().Get("Content-Type"),
	)
	// Output: Body={"status":500,"message":"Internal server error"}, Code=500, Content-Type=application/json
}
