package query

import (
	"fmt"
	"regexp"
	"time"
)

type QueryType int
type MathOp int

type FileType string

const (
	F FileType = "file"
	U FileType = "url"
)

const (
	SELECT QueryType = 0
	DELETE QueryType = 1
)

const (
	GT   MathOp = 0
	LT   MathOp = 1
	EQ   MathOp = 2
	NE   MathOp = 3
	LIKE MathOp = 4
)

var jsonLoad interface{}
var UploadFileName string

func contains(slice []string, needle string) bool {
	// Checks where a string is in an array
	for _, v := range slice {
		if v == needle {
			return true
		}
	}
	return false
}

func tokenise(querystr string) []string {
	// Uses regex to split the query string
	regex := "(?i)([a-zA-Z-_.]*[><=/ ]*[0-9]+)|[a-z-_.,]*[a-z-_.,]*([ = /LIKE]*\"[a-z0-9-_.\\[\\]\\* ]+\")*"
	return regexp.MustCompile(regex).FindAllString(querystr, 100)
}

func Execute(querystr string, tablename string) ([]map[string]interface{},
	int64, int, []error) {

	// Execute public and executes the query
	var start = time.Now()
	var decodedData, errs = decodeQuery(querystr)
	var elapsed = time.Since(start)

	fmt.Println("------------------")
	fmt.Println(querystr, "executed in", elapsed.Milliseconds(), "ms")
	fmt.Println("------------------")
	return decodedData, elapsed.Milliseconds(), len(decodedData), errs
}
