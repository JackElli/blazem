package query

import (
	"fmt"
	"regexp"
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
	case "string":
		_, doesWhereCast = wherevalue.(string)
	}

	// check if doesnt cast
	if (mathOp == EQ || mathOp == NE) && !doesWhereCast {
		*holds = *holds & 0
		return
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

	// TODO add > and <

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
