package crypto

import (
	"fmt"
)

func ExampleVerify() {
	hash, _ := Hash("password123")

	goodMatches := Verify(hash, "password123")
	badMatches := Verify(hash, "password987")

	fmt.Println(goodMatches, badMatches)
	// output: true false
}
