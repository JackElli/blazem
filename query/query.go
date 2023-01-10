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

// contains checks where a string is in an array
func contains(slice []string, needle string) bool {
	for _, v := range slice {
		if v == needle {
			return true
		}
	}
	return false
}

// printResults prints all the documents inside a query result
func printResults(resultMap []map[string]interface{}) {
	for _, doc := range resultMap {
		for key, value := range doc {
			fmt.Println(key, " ", value)
		}
		fmt.Println("---")
	}
}

// tonise uses regex to split the query string
func tokenise(querystr string) []string {
	regex := "(?i)([a-z-_.]*[><=/][0-9]+)|[a-z-_.,]*[a-z-_.,]*([=/~]*\"[a-z-_.0-9 ]+\")*"
	return regexp.MustCompile(regex).FindAllString(querystr, 100)
}

// Execute public and executes the query
func Execute(querystr string, tablename string) ([]map[string]interface{},
	int64, int, []error) {

	// set start time to now
	start := time.Now()
	// decode the query
	decodedData, errs := decodeQuery(querystr)
	elapsed := time.Since(start)
	fmt.Println("------------------")
	fmt.Println(querystr, "executed in", elapsed.Milliseconds(), "ms")
	fmt.Println("------------------")
	return decodedData, elapsed.Milliseconds(), len(decodedData), errs
	// calc the time taken
	// fmt.Println("Program took ", elapsed, " to get ", len(decodedData), " docs")
}
