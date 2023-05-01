package rules

import (
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/global"
	global_types "blazem/pkg/domain/global"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func NewAddRuleHandler(e *types.Endpoint) func(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	return AddRuleHandler
}

func NewGetRulesHandler(e *types.Endpoint) func(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	return GetRulesHandler
}

func NewRemoveRuleHandler(e *types.Endpoint) func(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	return RemoveRuleHandler
}

func AddRuleHandler(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	pe := &RuleEndpoint{
		Endpoint: *e,
	}
	return pe.addRuleHandler
}

func RemoveRuleHandler(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	pe := &RuleEndpoint{
		Endpoint: *e,
	}
	return pe.removeRuleHandler
}

func GetRulesHandler(e *types.Endpoint) func(w http.ResponseWriter, req *http.Request) {
	pe := &RuleEndpoint{
		Endpoint: *e,
	}
	return pe.getRulesHandler
}

// We want to add a rule to blazem
func (e *RuleEndpoint) addRuleHandler(w http.ResponseWriter, req *http.Request) {
	e.Endpoint.WriteHeaders(w, []string{})
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
	var ruleId = "rule" + strconv.Itoa(len(e.Endpoint.Node.Rules))
	e.Endpoint.Node.Rules[ruleId] = global.Rule{}
	json.NewEncoder(w).Encode("done")
}

// We want to remove a rule from Blazem
func (e *RuleEndpoint) removeRuleHandler(w http.ResponseWriter, req *http.Request) {
	e.Endpoint.WriteHeaders(w, []string{"ruleId"})

	var ruleId = req.URL.Query().Get("ruleId")
	_, ok := e.Endpoint.Node.Rules[ruleId]
	if !ok {
		json.NewEncoder(w).Encode("fail")
		return
	}
	delete(e.Endpoint.Node.Rules, ruleId)
	json.NewEncoder(w).Encode("done")
}

// We want to fetch all of the rules currently available in Blazem
func (e *RuleEndpoint) getRulesHandler(w http.ResponseWriter, req *http.Request) {
	e.Endpoint.WriteHeaders(w, []string{})

	var jsonRules = make([]map[string]interface{}, 0)
	for _, rule := range e.Endpoint.Node.Rules {
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
