//go:debug x509negativeserial=1
package main

import (
	"log"
	"os"

	"myoptions.info/indigo/backend/internal/entry"
)

func main() {
	// Use a logger with no prefix for program startup
	l := log.New(os.Stderr, "", 0)
	l.Println("Indigo CIL, v0.0.0")

	entry.Entry(l)
}
