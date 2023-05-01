package query

import (
	"blazem/domain/global"
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

// Uses regex to split the query string
func tokenise(querystr string) []string {
	regex := "(?i)([a-zA-Z-_.]*[><=/ ]*[0-9]+)|[a-z-_.,]*[a-z-_.,]*([ = /LIKE]*\"[a-z0-9-_.\\[\\]\\* ]+\")*"
	return regexp.MustCompile(regex).FindAllString(querystr, 100)
}

// Execute public and executes the query
func Execute(querystr string, tablename string) ([]global.Document,
	int64, int, []error) {
	var start = time.Now()
	var decodedData, errs = decodeQuery(querystr)
	var elapsed = time.Since(start)

	fmt.Println("------------------")
	fmt.Println(querystr, "executed in", elapsed.Milliseconds(), "ms")
	fmt.Println("------------------")
	return decodedData, elapsed.Milliseconds(), len(decodedData), errs
}

// LoadIntoMemory loads file or API into mem
func LoadIntoMemory(node global.Node) string {
	jsonLoad = node.Data
	return ""
}
