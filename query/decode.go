package query

import (
	"fmt"
	"strings"
)

//	index of token: {
//			SPECIAL TOKEN: {
//					func to complete
//			}
//	}
//
// asdsas
type DecodeFunc func(string, *QueryType,
	*bool, *bool, *[]string, *[]string) error

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

// ANY being any token
var decodeTokenTable = map[int]map[string]DecodeFunc{
	0: {
		"SELECT": func(token string, queryType *QueryType,
			all *bool, where *bool, fetchKeys *[]string, whereParams *[]string) error {
			*queryType = SELECT
			return nil
		},
		"DELETE": func(token string, queryType *QueryType,
			all *bool, where *bool, fetchKeys *[]string, whereParams *[]string) error {
			*queryType = DELETE
			return nil
		},
	},
	1: {
		"ANY": func(token string, queryType *QueryType,
			all *bool, where *bool, fetchKeys *[]string, whereParams *[]string) error {
			if token == "none" {
				return fmt.Errorf("NONE")
			}
			if token != "all" {
				// this is where params are set
				tokensplit := strings.Split(token, ",")
				for _, t := range tokensplit {
					*fetchKeys = append(*fetchKeys, t)
				}
				return nil
			}
			*all = true
			return nil
		},
	},
	2: {
		"WHERE": func(token string, queryType *QueryType,
			all *bool, where *bool, fetchKeys *[]string, whereParams *[]string) error {
			*where = true
			return nil
		},
	},
}

var specialTokens = map[string]bool{
	"SELECT": true,
	"WHERE":  true,
	"DELETE": true,
}

// decodeTokens gets each token and decides what it is
func decodeToken(i int, token string, queryType *QueryType,
	all *bool, where *bool, fetchKeys *[]string, whereParams *[]string) error {
	// If there is a where clause
	// get all predicates after the where
	// do this first
	if *where && i > 2 {
		*whereParams = append(*whereParams, token)
		return nil
	}
	// this is for the decoder
	// before the where clause
	var decoderError error
	if _, ok := specialTokens[token]; ok {
		decoderError = decodeTokenTable[i][token](token, queryType, all,
			where, fetchKeys, whereParams)
	} else {
		decoderError = decodeTokenTable[i]["ANY"](token, queryType, all,
			where, fetchKeys, whereParams)
	}

	return decoderError
}

// decodeQuery decodes the query so that we can manipulate it
func decodeQuery(querystr string) ([]map[string]interface{}, []error) {
	var tokens = tokenise(querystr)
	var queryType QueryType

	var fetchKeys []string
	var whereParams []string
	var where bool
	var all bool
	var errs []error = make([]error, 0)

	if len(tokens) < 2 {
		errs = append(errs, fmt.Errorf("NOT ENOUGH ARGS"))
		return nil, errs
	}

	for i, token := range tokens {
		err := decodeToken(i, token, &queryType, &all, &where, &fetchKeys, &whereParams)
		if err != nil {
			errs = append(errs, err)
		}
	}

	jsondata, loaderr := loadTable(jsonLoad)
	if loaderr != nil {
		errs = append(errs, loaderr)
	}
	if len(errs) > 0 {
		return nil, errs
	}
	return executeQuery(queryType, whereParams, fetchKeys, jsondata, all), errs

}
