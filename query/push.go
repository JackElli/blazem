package query

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func checkNest(nestparams []string, getobj map[string]interface{},
	docin *bool) {
	for _, nestparam := range nestparams {
		if v, exists := getobj[nestparam]; exists {
			if reflect.TypeOf(v).String() ==
				"map[string]interface {}" {
				getobj = v.(map[string]interface{})
			}
		} else {
			*docin = false
		}
	}
}
func pushDocs(all bool, wherejson []map[string]interface{},
	newmap *[]map[string]interface{},
	fetchKeys []string) []map[string]interface{} {
	// if the select param is all return all data
	if all {
		return wherejson
	}
	// this is essentially a Primary Index
	// for each doc in filtered query
	for _, doc := range wherejson {
		docin := true
		// allocate memory for the obj
		newobj := make(map[string]interface{})
		getobj := doc
		// for each key in the token fetch keys
		for _, fetchkey := range fetchKeys {
			// if the key exists in the data
			if strings.Contains(fetchkey, ".") {
				//genius
				nestparams := strings.Split(fetchkey, ".")
				checkNest(nestparams, getobj, &docin)
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
		// append this data to the new map to return
		if docin {
			*newmap = append(*newmap, newobj)
		}
	}
	return nil
}

func checkParamHolds(ok bool, paramsplit []string,
	getobj map[string]interface{}, mathOp MathOp, holds *int) {
	if !ok {
		*holds = *holds & 0
		return
	}
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
			*holds = *holds & 0
		}
		// objval := getobj[nestparams[len(nestparams)-1]]
		wherevalue = regexp.MustCompile("(?i)[a-z ]+|[0-9]+").FindString(fmt.Sprintf("%v", wherevalue))
		if v, exists := getobj[nestparams[len(nestparams)-1]]; exists {
			checkIfDocHolds(mathOp, v, wherevalue, holds)
		}
	}
	if v, exists := getobj[wherekey]; exists {
		checkIfDocHolds(mathOp, v, wherevalue, holds)
	}
	*holds = *holds & 1
}

// executeQuery is the query chain
func executeQuery(queryType QueryType, whereParams []string,
	fetchKeys []string, jsondata interface{},
	all bool) []map[string]interface{} {

	var newjsondata = jsondata.(map[string]interface{})
	// jsondata = jsondata.(map[string]interface{})
	var newmap []map[string]interface{}

	switch queryType {
	case SELECT:
		var wherejson []map[string]interface{}
		// if there are whereParams
		//filter
		if len(whereParams) > 0 {
			// for each document
			for _, doc := range newjsondata {
				// this is seeing if the doc matches the query
				holds := 1
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
					checkParamHolds(ok, paramsplit, getobj,
						mathOp, &holds)
				}
				// if the doc matches the query
				if holds == 1 {
					wherejson = append(wherejson, getobj)
				}
			}
		} else {
			for _, doc := range newjsondata {
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
