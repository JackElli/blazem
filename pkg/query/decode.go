package query

import (
	"fmt"
	"regexp"
	"strings"
)

type DecodeFunc func(string, *QueryType,
	*bool, *bool, *[]string, *[]string) error

var specialTokens = map[string]bool{
	"SELECT": true,
	"WHERE":  true,
	"DELETE": true,
	"LIKE":   true,
}

// DecodeParam decodes the paramaters and sets up the next step in the query chain
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
	} else if strings.Contains(param, "LIKE") {
		*paramsplit = strings.Split(param, "LIKE")
		*mathOp = LIKE
		return true
	}
	return false
}

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
		// this is for the values fetched
		"ANY": func(token string, queryType *QueryType,
			all *bool, where *bool, fetchKeys *[]string, whereParams *[]string) error {
			if token == "none" {
				return fmt.Errorf("NONE")
			}
			if token != "all" {
				// this is where params are set
				tokensplit := strings.Split(token, ",")
				// if there's no key entered in the query
				// add it so we know which doc to retrieve
				// and send back
				if !strings.Contains(token, "key") {
					*fetchKeys = append(*fetchKeys, "key")
				}
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

// DecodeTokens gets each token and decides what it is
func decodeToken(i int, token string, queryType *QueryType,
	all *bool, where *bool, fetchKeys *[]string, whereParams *[]string) error {
	var decoderError error
	if *where && i > 2 {
		noWhiteSpaceReg := regexp.MustCompile("[a-zA-Z-_.= ]*")
		findSection := noWhiteSpaceReg.FindString(token)
		findSectionNoWhiteSpace := strings.ReplaceAll(findSection, " ", "")
		trimmedToken := strings.ReplaceAll(token, findSection, findSectionNoWhiteSpace)
		*whereParams = append(*whereParams, string(trimmedToken))
		return nil
	}
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
func (query *Query) decodeQuery(queryStr string) ([]map[string]interface{}, []error) {

	var queryType QueryType
	var fetchKeys []string
	var whereParams []string
	var where bool
	var all bool

	errs := make([]error, 0)
	tokens := query.tokenise(queryStr)

	if len(tokens) < 2 {
		errs = append(errs, fmt.Errorf("Not enough query arguments passed"))
		return nil, errs
	}

	for i, token := range tokens {
		err := decodeToken(i, token, &queryType, &all, &where, &fetchKeys, &whereParams)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return executeQuery(queryType, whereParams, fetchKeys, query.JsonLoad, all), nil
}
