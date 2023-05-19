package query

import (
	"blazem/pkg/domain/node"

	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

// This is for non-where clause tokens
func checkNest(nestParams []string, getObj map[string]interface{}, docIn *bool) map[string]interface{} {
	for _, nestParam := range nestParams {
		v, exists := getObj[nestParam]
		if !exists {
			*docIn = false
			return getObj
		}
		if reflect.TypeOf(v).String() ==
			"map[string]interface {}" {
			getObj = v.(map[string]interface{})
		}
	}
	return getObj
}

// We want to push the documents that fit the query params
func pushDocs(all bool, whereJson []map[string]interface{}, newMap *[]map[string]interface{},
	fetchKeys []string) []map[string]interface{} {
	if all {
		return whereJson
	}

	for _, doc := range whereJson {
		docIn := true
		newObj := make(map[string]interface{})
		getObj := doc

		for _, fetchKey := range fetchKeys {
			if strings.Contains(fetchKey, ".") {
				nestParams := strings.Split(fetchKey, ".")
				getObj = checkNest(nestParams, getObj, &docIn)
				if v, exists := getObj[nestParams[len(nestParams)-1]]; exists {
					newObj[fetchKey] = v
				} else {
					newObj[fetchKey] = getObj
				}
				continue
			}
			if v, exists := doc[fetchKey]; exists {
				newObj[fetchKey] = v
				continue
			}
			docIn = false
		}
		if docIn {
			*newMap = append(*newMap, newObj)
		}
	}
	return nil
}

// We want to check whether a certain paramater holds for that specific document
func checkParamHolds(ok bool, paramSplit []string,
	getobj map[string]interface{}, mathOp MathOp, holds *int) {
	if !ok {
		*holds = *holds & 0
		return
	}

	whereKey := paramSplit[0]
	whereKey = strings.Trim(whereKey, " ")

	if strings.Contains(whereKey, ".") {
		wherevalue := regexp.MustCompile("(?i)\"[a-zA-Z0-9-_ ]+\"|[0-9]*").FindString(paramSplit[1])
		nestparams := strings.Split(whereKey, ".")

		for _, nestparam := range nestparams {
			if v, exists := getobj[nestparam]; exists {
				if reflect.TypeOf(v).String() == "map[string]interface {}" {
					getobj = v.(map[string]interface{})
				}
			}
		}
		if _, exists := getobj[nestparams[len(nestparams)-1]]; !exists {
			*holds = *holds & 0
		}
		wherevalue = regexp.MustCompile("(?i)[a-zA-Z0-9-_ ]+").FindString(fmt.Sprintf("%v", wherevalue))
		if v, exists := getobj[nestparams[len(nestparams)-1]]; exists {
			checkIfDocHolds(mathOp, v, wherevalue, holds)
			return
		}
	}
	wherevalue := strings.Trim(regexp.MustCompile("(?i)[a-zA-Z0-9-_ ]+").FindString(paramSplit[1]), " ")
	if v, exists := getobj[whereKey]; exists {
		checkIfDocHolds(mathOp, v, wherevalue, holds)
		return
	}
	*holds = *holds & 0
}

// executeQuery is the query chain
func executeQuery(queryType QueryType, whereParams []string,
	fetchKeys []string, jsonData interface{},
	all bool) []map[string]interface{} {

	var newMap []map[string]interface{}
	var whereJson []map[string]interface{}
	newJsonData := jsonData.(sync.Map)

	if len(whereParams) <= 0 {
		newJsonData.Range(func(key, doc any) bool {
			whereJson = append(whereJson, doc.(map[string]interface{}))
			return true
		})
	} else {
		newJsonData.Range(func(key, doc any) bool {
			holds := 1
			getObj := doc.(map[string]interface{})

			for _, param := range whereParams {
				var paramSplit []string
				var mathOp MathOp

				ok := decodeParam(param, &mathOp, &paramSplit)
				checkParamHolds(ok, paramSplit, getObj, mathOp, &holds)
			}
			if holds == 1 {
				whereJson = append(whereJson, getObj)
			}
			return true
		})
	}
	pushed := pushDocs(all, whereJson, &newMap, fetchKeys)
	if pushed == nil {
		return newMap
	}

	if queryType == SELECT {
		return pushed
	}

	if queryType == DELETE {
		for _, doc := range pushed {
			key := doc["key"].(string)
			node.GlobalNode.Data.Delete(key)
		}
		return []map[string]interface{}{}
	}
	return nil
}
