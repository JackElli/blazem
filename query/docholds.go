package query

import (
	"strconv"
	"strings"
)

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
