package rules

import (
	"blazem/pkg/domain/endpoint_manager"
	"blazem/pkg/domain/global"
	global_types "blazem/pkg/domain/global"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// We want to add a rule to blazem
func AddRuleHandler(e *endpoint_manager.EndpointManager) func(w http.ResponseWriter, req *http.Request) {
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
				global.Logger.Error("Failed to add rule")
				json.NewEncoder(w).Encode("fail")
			}
		}
		var ruleId = "rule" + strconv.Itoa(len(e.Node.Rules))
		e.Node.Rules[ruleId] = global.Rule{}
		json.NewEncoder(w).Encode("done")
	}
}
