package handlers

import (
	"blazem/global"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func AddRuleHandler(node *Node) func(w http.ResponseWriter, req *http.Request) {
	return node.addRuleHandler
}

func RemoveRuleHandler(node *Node) func(w http.ResponseWriter, req *http.Request) {
	return node.removeRuleHandler
}

func GetRulesHandler(node *Node) func(w http.ResponseWriter, req *http.Request) {
	return node.getRulesHandler
}

func (node *Node) addRuleHandler(w http.ResponseWriter, req *http.Request) {
	// We want to add a rule to blazem
	WriteHeaders(w, []string{})
	var rule Rule
	var taskForRule = make([]global.Task, 0)
	var executeTime *time.Time
	var t time.Time
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
		t, err = time.Parse("2006-01-02 15:04:05", rule.Time)
		executeTime = &t
		if err != nil {
			fmt.Println("Failed to add rule")
			json.NewEncoder(w).Encode("fail")
		}
	} else {
		executeTime = nil
	}
	var ruleId = "rule" + strconv.Itoa(len(node.Rules))
	node.Rules[ruleId] = global.Rule{
		Id:          ruleId,
		Tasks:       taskForRule,
		ExecuteTime: executeTime,
	}
	json.NewEncoder(w).Encode("done")
}

func (node *Node) removeRuleHandler(w http.ResponseWriter, req *http.Request) {
	// We want to remove a rule from Blazem
	WriteHeaders(w, []string{"ruleId"})

	var ruleId = req.URL.Query().Get("ruleId")
	_, ok := node.Rules[ruleId]
	if !ok {
		json.NewEncoder(w).Encode("fail")
		return
	}
	delete(node.Rules, ruleId)
	json.NewEncoder(w).Encode("done")
}

func (node *Node) getRulesHandler(w http.ResponseWriter, req *http.Request) {
	// We want to fetch all of the rules currently available in Blazem
	WriteHeaders(w, []string{})

	var jsonRules = make([]map[string]interface{}, 0)
	for _, rule := range node.Rules {
		jsonTasks := make([]JSONTask, 0)
		for _, task := range rule.Tasks {
			jsonTasks = append(jsonTasks, JSONTask{
				Type:    task.Type,
				Data:    task.Data,
				Require: task.Require,
			})
		}
		sendRule := map[string]interface{}{
			"tasks":    jsonTasks,
			"execTime": rule.ExecuteTime.Format("2006-01-02 15:04:05"),
			"id":       rule.Id,
		}
		jsonRules = append(jsonRules, sendRule)
	}
	json.NewEncoder(w).Encode(jsonRules)
}
