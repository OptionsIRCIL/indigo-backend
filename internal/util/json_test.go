package util

import (
	"fmt"
	"net/http/httptest"
	"strings"
)

type inventoryItem struct {
	Item  string `json:"item"`
	Count int    `json:"count"`
}

func ExampleDecodeJsonBody() {
	r := httptest.NewRequest(
		"POST",
		"http://example.com/api/v1/demo",
		strings.NewReader(`{"item": "Banana", "count": 7}`),
	)
	r.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	item := inventoryItem{}
	err := DecodeJsonBody(w, r, &item)

	if err == nil {
		fmt.Printf("Decode succeeded! Item=%s, Count=%d\n", item.Item, item.Count)
	} else {
		fmt.Println("Oh no, decode failed!")
	}

	// Output: Decode succeeded! Item=Banana, Count=7
}
