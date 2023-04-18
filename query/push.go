package query

import (
	"blazem/global"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

// This is for non-where clause tokens
func checkNest(nestparams []string, getobj global.Document,
	docin *bool) global.Document {
	for _, nestparam := range nestparams {
		if v, exists := getobj[nestparam]; exists {
			if reflect.TypeOf(v).String() ==
				"map[string]interface {}" {
				getobj = v.(global.Document)
			}
		} else {
			*docin = false
		}
	}
	return getobj
}

// We want to push the documents that fit the query params
func pushDocs(all bool, wherejson []global.Document,
	newmap *[]global.Document,
	fetchKeys []string) []global.Document {
	if all {
		return wherejson
	}

	for _, doc := range wherejson {
		var docin = true
		var newobj = make(global.Document)
		var getobj = doc

		for _, fetchkey := range fetchKeys {
			if strings.Contains(fetchkey, ".") {
				var nestparams = strings.Split(fetchkey, ".")
				getobj = checkNest(nestparams, getobj, &docin)
				if v, exists :=
					getobj[nestparams[len(nestparams)-1]]; exists {
					newobj[fetchkey] = v
				} else {
					newobj[fetchkey] = getobj
				}
				continue
			}
			if v, exists := doc[fetchkey]; exists {
				newobj[fetchkey] = v
				continue
			}
			docin = false
		}
		if docin {
			*newmap = append(*newmap, newobj)
		}
	}
	return nil
}

// We want to check whether a certain paramater holds for that specific document
func checkParamHolds(ok bool, paramsplit []string,
	getobj global.Document, mathOp MathOp, holds *int) {
	if !ok {
		*holds = *holds & 0
		return
	}

	var wherekey = paramsplit[0]
	wherekey = strings.Trim(wherekey, " ")

	if strings.Contains(wherekey, ".") {
		var wherevalue = regexp.MustCompile("(?i)\"[a-zA-Z0-9-_ ]+\"|[0-9]*").FindString(paramsplit[1])
		var nestparams = strings.Split(wherekey, ".")

		for _, nestparam := range nestparams {
			if v, exists := getobj[nestparam]; exists {
				if reflect.TypeOf(v).String() == "map[string]interface {}" {
					getobj = v.(global.Document)
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
	var wherevalue = strings.Trim(regexp.MustCompile("(?i)[a-zA-Z0-9-_ ]+").FindString(paramsplit[1]), " ")
	if v, exists := getobj[wherekey]; exists {
		checkIfDocHolds(mathOp, v, wherevalue, holds)
		return
	}
	*holds = *holds & 0
}

// executeQuery is the query chain
func executeQuery(queryType QueryType, whereParams []string,
	fetchKeys []string, jsondata interface{},
	all bool) []global.Document {

	var newjsondata = jsondata.(sync.Map)
	var newmap []global.Document
	var wherejson []global.Document

	if len(whereParams) > 0 {
		newjsondata.Range(func(key, doc any) bool {
			var holds = 1
			var getobj = doc.(global.Document)

			for _, param := range whereParams {
				var paramsplit []string
				var mathOp MathOp

				var ok = decodeParam(param, &mathOp, &paramsplit)
				checkParamHolds(ok, paramsplit, getobj, mathOp, &holds)
			}
			if holds == 1 {
				wherejson = append(wherejson, getobj)
			}
			return true
		})
	} else {
		newjsondata.Range(func(key, doc any) bool {
			wherejson = append(wherejson, doc.(global.Document))
			return true
		})
	}
	var pushed = pushDocs(all, wherejson, &newmap, fetchKeys)
	if pushed != nil {
		if queryType == SELECT {
			return pushed
		}
		if queryType == DELETE {
			for _, doc := range pushed {
				key := doc["key"].(string)
				global.GlobalNode.Data.Delete(key)
			}
			return []global.Document{}
		}
	}
	return newmap
}
