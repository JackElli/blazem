package query

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
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
	return regexp.MustCompile("(?i)([a-z-_.]*[><=/][0-9]+)|[a-z-_.,]*[a-z-_.,]*([=/~]*\"[a-z-_.0-9 ]+\")*").FindAllString(querystr, 100)
}

// decodeTokens gets each token and decides what it is
func decodeToken(i int, token string, queryType *QueryType, all *bool, where *bool, fetchKeys *[]string, whereParams *[]string) string {
	// switch each index of token
	switch i {
	case 0:
		switch token {
		case "SELECT":
			// set the query type to select
			*queryType = SELECT
		}
	case 1:
		// this is for select like statements
		if token == "all" {
			*all = true
		} else {
			// this is where params are set
			tokensplit := strings.Split(token, ",")
			for _, t := range tokensplit {
				*fetchKeys = append(*fetchKeys, t)
			}
		}

	case 2:
		if token == "WHERE" {
			*where = true
		}
	}
	if *where && i > 2 {
		// this is the where params
		*whereParams = append(*whereParams, token)
	}
	return ""
}

// decodeQuery decodes the query so that we can manipulate it
func decodeQuery(querystr string) ([]map[string]interface{}, []string) {
	// tokenise := strings.Split(querystr, " ")
	// r, _ := regexp.Compile("\"[^\"]*\"|[^,]+")
	// tokenise := regexp.MustCompile("[a-z ]*|'[a-z ]*'").Split("SELECT Title,Duration WHERE Title='key to life'", -1)
	//regex, '[a-z ] a to z with space after *(multiple) or [a-z]many
	//no numbers!

	//want it to be title=asdasd
	//				duration=afsdf
	tokens := tokenise(querystr)
	var queryType QueryType

	var fetchKeys []string
	var whereParams []string
	var where bool
	var all bool
	var errs []string = []string{}
	// for each token
	if len(tokens) < 2 {
		errs = append(errs, "NOT ENOUGH ARGS")
		return nil, errs
	}
	for i, token := range tokens {
		err := decodeToken(i, token, &queryType, &all, &where, &fetchKeys, &whereParams)
		if err != "" {
			errs = append(errs, err)
		}
	}

	jsondata, loaderr := loadTable(jsonLoad)
	if loaderr != "" {
		errs = append(errs, loaderr)
	}
	if len(errs) > 0 {
		return nil, errs
	}
	return executeQuery(queryType, whereParams, fetchKeys, jsondata, all), errs

}

func LoadIntoMemory(filepath string) string {
	// //map of data to return
	var fileType FileType
	if strings.Contains(filepath, "http://") || strings.Contains(filepath, "https://") {
		fileType = U
	} else {
		fileType = F
	}
	UploadFileName = filepath
	// // this is var for getting data from file
	var i interface{}
	var fdata []byte

	if fileType == F {
		fdata, _ = ioutil.ReadFile(filepath)
	} else {
		request, _ := http.Get(filepath)
		fdata, _ = ioutil.ReadAll(request.Body)
	}

	err := json.Unmarshal(fdata, &i)
	if err != nil {

		return "no file uploaded"
	}
	jsonLoad = i
	return ""

}

// loadIntoMemory loads the data in the filepath/url into memory
func loadTable(i interface{}) ([]interface{}, string) {
	// tableSplit := strings.Split(tableName, ".")

	var m map[string]interface{}
	var ok bool = false

	_, ok = i.([]interface{})
	if ok {
		return i.([]interface{}), ""
	}
	m = i.(map[string]interface{})
	for !ok {
		for key, _ := range m {
			if _, gok := m[key].(map[string]interface{}); gok {
				m = m[key].(map[string]interface{})
			}
			if _, fok := m[key].([]interface{}); fok {
				return m[key].([]interface{}), ""
			}

		}
	}
	return i.([]interface{}), ""

}

// decodeParam decodes the paramaters and sets up the next step in the query chain
func decodeParam(param string, mathOp *MathOp, paramsplit *[]string) bool {

	if strings.Contains(param, "=") && !strings.Contains(param, "/") {
		*paramsplit = strings.Split(param, "=")
		*mathOp = EQ
		return true
	} else if strings.Contains(param, ">") {
		*paramsplit = strings.Split(param, ">")
		*mathOp = GT
		return true
	} else if strings.Contains(param, "<") {
		*paramsplit = strings.Split(param, "<")
		*mathOp = LT
		return true
	} else if strings.Contains(param, "/=") {
		*paramsplit = strings.Split(param, "/=")
		*mathOp = NE
		return true
	} else if strings.Contains(param, "~") {
		*paramsplit = strings.Split(param, "~")
		*mathOp = LIKE
		return true
	}
	return false
}

// checkIfDocHolds checks if the document matches the query
func checkIfDocHolds(mathOp MathOp, v interface{}, wherevalue string, holds *int) {
	// parse float, will return string if its not number

	//NEED TO CHECK FOR MULTIPLE TYPES
	intVar, _ := strconv.ParseFloat(wherevalue, 64)
	if intVar == 0 { //not an number
		if mathOp == EQ {
			if v != wherevalue {
				*holds = *holds & 0
			}
		} else if mathOp == NE {
			if v == wherevalue {
				*holds = *holds & 0
			}
		} else {
			//like
			if !strings.Contains(v.(string), wherevalue) {
				*holds = *holds & 0
			}
		}

	} else { //must be a number
		var vNumVar float64
		if _, ok := v.(string); ok {
			vNumVar, _ = strconv.ParseFloat(v.(string), 64)
		} else {
			//possible a float
			vNumVar = v.(float64)
		}

		switch mathOp {
		case EQ:
			if vNumVar == intVar {
				*holds = *holds & 1
			} else {
				*holds = *holds & 0
			}
		case LT:
			if vNumVar < intVar {
				*holds = *holds & 1
			} else {
				*holds = *holds & 0
			}
		case GT:
			if vNumVar > intVar {
				*holds = *holds & 1
			} else {
				*holds = *holds & 0
			}
		}
	}
}
func pushDocs(all bool, wherejson []map[string]interface{}, newmap *[]map[string]interface{}, fetchKeys []string) []map[string]interface{} {
	// if the select param is all return all data
	if all {
		return wherejson
	} else {
		// this is essentially a Primary Index
		// for each doc in filtered query
		for _, doc := range wherejson {
			docin := 1
			// allocate memory for the obj
			newobj := make(map[string]interface{})
			// for each key in the token fetch keys
			for _, fetchkey := range fetchKeys {
				// if the key exists in the data
				if strings.Contains(fetchkey, ".") {
					//genius
					nestparams := strings.Split(fetchkey, ".")
					getobj := doc
					for _, nestparam := range nestparams {
						if v, exists := getobj[nestparam]; exists {
							if reflect.TypeOf(v).String() == "map[string]interface {}" {
								getobj = v.(map[string]interface{})
							}
						} else {
							docin = 0
						}
					}
					if v, exists := getobj[nestparams[len(nestparams)-1]]; exists {
						newobj[fetchkey] = v
					} else {
						newobj[fetchkey] = getobj
					}

				} else {
					if v, exists := doc[fetchkey]; exists {
						newobj[fetchkey] = v
					} else {
						docin = 0
					}
				}
			}
			// append this data to the new map to return
			if docin == 1 {
				*newmap = append(*newmap, newobj)
			}
		}
	}
	return nil
}

// executeQuery is the query chain
func executeQuery(queryType QueryType, whereParams []string, fetchKeys []string, jsondata []interface{}, all bool) []map[string]interface{} {
	var newmap []map[string]interface{}

	switch queryType {
	case SELECT:
		var wherejson []map[string]interface{}
		// if there are whereParams
		//filter
		if len(whereParams) > 0 {
			// for each document
			for _, doc := range jsondata {
				// this is seeing if the doc matches the query
				var holds int = 1
				//could be better
				getobj := doc.(map[string]interface{})

				// for each param
				for _, param := range whereParams {
					// split the paramstring
					var paramsplit []string
					// math operator <>=
					var mathOp MathOp
					//decodes the paramater, splitting int paramsplit and mathop
					ok := decodeParam(param, &mathOp, &paramsplit)
					if ok {
						// this is the where key
						wherekey := paramsplit[0]
						// where value, value of the where param
						wherevalue := regexp.MustCompile("(?i)[a-z0-9-_. ]+|[0-9-_. ]+").FindString(paramsplit[1])

						// fmt.Println(getobj)
						if strings.Contains(wherekey, ".") {
							//genius
							nestparams := strings.Split(wherekey, ".")
							for _, nestparam := range nestparams {
								if v, exists := getobj[nestparam]; exists {
									if reflect.TypeOf(v).String() == "map[string]interface {}" {
										getobj = v.(map[string]interface{})
									}
								}
							}
							if _, exists := getobj[nestparams[len(nestparams)-1]]; !exists {
								holds = holds & 0
							}
							// objval := getobj[nestparams[len(nestparams)-1]]
							wherevalue = regexp.MustCompile("(?i)[a-z ]+|[0-9]+").FindString(fmt.Sprintf("%v", wherevalue))
							if v, exists := getobj[nestparams[len(nestparams)-1]]; exists {
								checkIfDocHolds(mathOp, v, wherevalue, &holds)
							}
						}
						if v, exists := getobj[wherekey]; exists {
							checkIfDocHolds(mathOp, v, wherevalue, &holds)
						}
					} else {
						holds = holds & 0
					}
				}
				// if the doc matches the query
				if holds == 1 {
					wherejson = append(wherejson, getobj)
				}
			}
		} else {
			for _, doc := range jsondata {
				wherejson = append(wherejson, doc.(map[string]interface{}))
			}

		}
		//push docs
		pushed := pushDocs(all, wherejson, &newmap, fetchKeys)
		if pushed != nil {
			return pushed
		}
	}
	return newmap
}

// createIndex creates an index so that queries run faster
func createIndex() {

}

// Execute public and executes the query
func Execute(querystr string, tablename string) ([]map[string]interface{}, int64, int, []string) {

	// set start time to now
	start := time.Now()
	// decode the query
	decodedData, err := decodeQuery(querystr)
	elapsed := time.Since(start)
	fmt.Println("------------------")
	fmt.Println(querystr, "executed in", elapsed.Milliseconds(), "ms")
	fmt.Println("------------------")
	return decodedData, elapsed.Milliseconds(), len(decodedData), err
	// calc the time taken
	// fmt.Println("Program took ", elapsed, " to get ", len(decodedData), " docs")
}
