package cli

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func getBodyFromFlagInputs(data, file string) (buf *strings.Builder, err error) {

	var f *os.File

	// handle any missing args
	switch {
	case data == "" && file == "":
		err = errors.New("Missing data - please provide the data that you would like to create")
		fmt.Fprintln(os.Stderr, fmt.Sprintf("%s", err.Error()))
		fmt.Fprintln(os.Stderr, "Must specify one of:")
		fmt.Fprintln(os.Stderr, "-d '<string>'")
		fmt.Fprintln(os.Stderr, "-d -")
		fmt.Fprintln(os.Stderr, "-f <filename>")
	}

	if err == nil {

		// use stdin as data to send if "-" is specified on command line `-d -`
		if data == "-" {
			body = os.Stdin
		} else {
			body = bytes.NewBuffer([]byte(data))
		}

		// if file is specified, use that instead of any `-d`
		if file != "" {
			f, err = os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to open file to read - %s\n", err)
			}
			defer f.Close()
			body = f
		}
	}

	if err == nil {
		buf = new(strings.Builder)
		_, err = io.Copy(buf, body)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input data: %s\n", err.Error())
	}

	return buf, err
}

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
