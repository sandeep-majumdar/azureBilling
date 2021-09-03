package utils

import (
	"log"
	"os"
)

func Quiet() func() {
	n, _ := os.Open(os.DevNull)
	sout := os.Stdout
	serr := os.Stderr
	os.Stdout = n
	os.Stderr = n
	log.SetOutput(n)
	return func() {
		defer n.Close()
		os.Stdout = sout
		os.Stderr = serr
		log.SetOutput(os.Stderr)
	}
}
