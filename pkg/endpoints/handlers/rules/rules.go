package rules

import (
	"blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/global"
	global_types "blazem/pkg/domain/global"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// We want to add a rule to blazem
func AddRuleHandler(r *endpoint.Respond) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var rule global_types.Rule
		var taskForRule = make([]global.Task, 0)

		var err error

		body, _ := ioutil.ReadAll(req.Body)
		json.Unmarshal(body, &rule)

		if len(rule.Tasks) == 0 {
			return
		}
		for _, task := range rule.Tasks {
			taskForRule = append(taskForRule, global.Task{
				Data:    task.Data,
				Require: task.Require,
				Type:    task.Type,
			})
		}
		if rule.Time != "" {
			if err != nil {
				fmt.Println("Failed to add rule")
				json.NewEncoder(w).Encode("fail")
			}
		}
		var ruleId = "rule" + strconv.Itoa(len(r.Node.Rules))
		r.Node.Rules[ruleId] = global.Rule{}
		json.NewEncoder(w).Encode("done")
	}
}

// We want to remove a rule from Blazem
func RemoveRuleHandler(r *endpoint.Respond) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		var ruleId = req.URL.Query().Get("ruleId")
		_, ok := r.Node.Rules[ruleId]
		if !ok {
			json.NewEncoder(w).Encode("fail")
			return
		}
		delete(r.Node.Rules, ruleId)
		json.NewEncoder(w).Encode("done")
	}
}

// We want to fetch all of the rules currently available in Blazem
func GetRulesHandler(r *endpoint.Respond) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		var jsonRules = make([]map[string]interface{}, 0)
		for _, rule := range r.Node.Rules {
			jsonTasks := make([]global_types.JSONTask, 0)
			for _, task := range rule.Tasks {
				jsonTasks = append(jsonTasks, global_types.JSONTask{
					Type:    task.Type,
					Data:    task.Data,
					Require: task.Require,
				})
			}
			sendRule := map[string]interface{}{}
			jsonRules = append(jsonRules, sendRule)
		}
		json.NewEncoder(w).Encode(jsonRules)
	}
}
