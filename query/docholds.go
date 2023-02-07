package query

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// checkIfDocHolds checks if the document matches the query
func checkIfDocHolds(mathOp MathOp, v interface{}, wherevalue interface{}, holds *int) {

	// v is currently in map
	// wherevalue is what user is searching

	// REMEMBER holds is 1 by default
	// we need to prove that it
	// does not hold

	// definitely change
	// this to factory

	var doesWhereCast bool
	opType := fmt.Sprintf("%T", v)
	switch opType {
	case "time.Time":
		_, doesWhereCast = wherevalue.(time.Time)
	case "float64":
		_, doesWhereCast = wherevalue.(float64)
		_, doesFloatParse := strconv.ParseFloat(wherevalue.(string), 64)
		doesWhereCast = (doesFloatParse == nil) || doesWhereCast
	case "string":
		_, doesWhereCast = wherevalue.(string)
	}

	// These are the opposites
	// check if doesnt cast
	if (mathOp == EQ || mathOp == NE) && !doesWhereCast {
		*holds = *holds & 0
		return
	}

	// make wherevalue cast to correct
	// type
	if doesWhereCast {
		// where value will always be string
		// user input
		switch opType {
		case "float64":
			wherevalue, _ = strconv.ParseFloat(wherevalue.(string), 64)
		}
	}

	if mathOp == EQ && v != wherevalue {
		*holds = *holds & 0
		return
	}

	// Not equals
	if mathOp == NE && v == wherevalue {
		*holds = *holds & 0
		return
	}

	// Greater than
	// must be float (check doesWhereCast)
	if mathOp == GT && doesWhereCast && v.(float64) <= wherevalue.(float64) {
		*holds = *holds & 0
		return
	}

	// Less than
	// must be float (check doesWhereCast)
	if mathOp == LT && doesWhereCast && v.(float64) >= wherevalue.(float64) {
		*holds = *holds & 0
		return
	}

	if mathOp == LIKE {
		switch opType {
		case "time.Time":
			if !strings.Contains(v.(time.Time).String(), wherevalue.(string)) {
				*holds = *holds & 0
				return
			}
		case "float64":
			if !strings.Contains(fmt.Sprintf("%g", v.(float64)), wherevalue.(string)) {
				*holds = *holds & 0
				return
			}
		case "string":
			reg, err := regexp.Compile(wherevalue.(string))

			if err != nil || !reg.MatchString(v.(string)) {
				*holds = *holds & 0
				return
			}
		}
	}
}
