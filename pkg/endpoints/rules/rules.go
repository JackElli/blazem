package rules

import (
	"blazem/pkg/domain/global"
	global_types "blazem/pkg/domain/global"
	"blazem/pkg/domain/logger"
	"blazem/pkg/domain/node"
	"blazem/pkg/domain/responder"
	"encoding/json"
	"net/http"
	"strconv"
)

type RulesMgr struct {
	Node      *node.Node
	Responder responder.Responder
}

func NewRulesMgr(node *node.Node, responder responder.Responder) *RulesMgr {
	return &RulesMgr{
		Node:      node,
		Responder: responder,
	}
}

// We want to add a rule to blazem
func (e *RulesMgr) AddRuleHandler() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var taskForRule = make([]global.Task, 0)
		var rule global_types.Rule
		err := json.NewDecoder(req.Body).Decode(&rule)
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
				logger.Logger.Error("Failed to add rule")
				json.NewEncoder(w).Encode("fail")
			}
		}

		var ruleId = "rule" + strconv.Itoa(len(e.Node.Rules))
		e.Node.Rules[ruleId] = global.Rule{}
		json.NewEncoder(w).Encode("done")
	}
}
