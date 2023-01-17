package endpoints

import (
	"blazem/global"
	"blazem/query"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/couchbase/gocb/v2"
)

type JSONTask struct {
	Type    string
	Data    string
	Require int
}

var taskFncDecoder = map[string]func(interface{}, interface{}) (interface{}, error){
	"query": func(queryVal interface{}, requirePass interface{}) (interface{}, error) {
		queryResult, _, _, _ := query.Execute(queryVal.(string), "")
		return queryResult, nil
	},
	"export": func(hostName interface{}, requirePass interface{}) (interface{}, error) {
		getHost, ok := hostName.(string)
		if !ok {
			return "", fmt.Errorf("not a string host")
		}
		getDocs, ok := requirePass.([]map[string]interface{})
		if !ok {
			return "", fmt.Errorf("cannot find docs")
		}

		if strings.Contains(getHost, "couchbase") {
			err := addToCouchbase(getHost, getDocs)
			if err != nil {
				return "", fmt.Errorf("cannot connect to couchbase")
			}
		}

		return "", nil
	},
}

func (node *Node) addRuleHandler(w http.ResponseWriter, req *http.Request) {
	writeHeaders(w, []string{})

	var tasks []JSONTask

	body, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(body, &tasks)

	if len(tasks) == 0 {
		return
	}

	// will do this is a strange order
	var taskForRule = make([]global.Task, 0)

	for _, task := range tasks {
		taskForRule = append(taskForRule, global.Task{
			Data:    task.Data,
			Require: task.Require,
			Type:    task.Type,
		})
	}

	ruleId := "rule" + strconv.Itoa(len(node.Rules))
	node.Rules[ruleId] = global.Rule{
		Id:    ruleId,
		Tasks: taskForRule,
		ExecuteTime: time.Date(
			2023, 01, 14, 22, 34, 58, 651387237, time.UTC),
	}

	json.NewEncoder(w).Encode("done")

}

func (node *Node) runRuleHandler(w http.ResponseWriter, req *http.Request) {
	writeHeaders(w, []string{"ruleId"})

	ruleId := req.URL.Query().Get("ruleId")
	ruleTasks, ok := node.Rules[ruleId]
	if !ok {
		json.NewEncoder(w).Encode("fail")
		return
	}

	var taskOutput = make([]interface{}, 0)

	// getting input and
	// running output
	for _, task := range ruleTasks.Tasks {
		runTask := taskFncDecoder[task.Type]
		data := task.Data
		if task.Require == -1 {
			out, err := runTask(data, "")
			if err != nil {
				json.NewEncoder(w).Encode("fail")
				return
			}
			taskOutput = append(taskOutput, out)
			continue
		}
		taskOutput = append(taskOutput, "")
		passData := taskOutput[task.Require]
		_, err := runTask(data, passData)
		if err != nil {
			json.NewEncoder(w).Encode("fail")
			return
		}
	}

	json.NewEncoder(w).Encode("done")

}
func (node *Node) getRulesHandler(w http.ResponseWriter, req *http.Request) {
	writeHeaders(w, []string{})

	jsonRules := make([]map[string]interface{}, 0)
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
			"execTime": rule.ExecuteTime,
			"id":       rule.Id,
		}
		jsonRules = append(jsonRules, sendRule)
	}

	json.NewEncoder(w).Encode(jsonRules)
}

func addToCouchbase(connectionString string, docs []map[string]interface{}) error {
	// Update this to your cluster details

	username := "Administrator"
	password := "password"

	initConnectionStringSplit := strings.Split(connectionString, "couchbase://")
	findTerms := strings.Split(initConnectionStringSplit[1], "/")
	var scope, collection string = "", ""
	hostName := findTerms[0]
	bucketName := findTerms[1]
	if len(findTerms) > 2 {
		scope = findTerms[2]
		collection = findTerms[3]
	}

	connectionString = "couchbase://" + hostName

	// For a secure cluster connection, use `couchbases://<your-cluster-ip>` instead.
	cluster, err := gocb.Connect(connectionString, gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
	})

	if err != nil {
		log.Fatal(err)
		return err
	}

	bucket := cluster.Bucket(bucketName)

	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		return err
	}
	var col *gocb.Collection
	if scope == "" {
		col = bucket.DefaultCollection()
	} else {
		col = bucket.Scope(scope).Collection(collection)
	}

	var getDocJSON map[string]interface{}
	for _, doc := range docs {
		docJSON, _ := json.Marshal(doc)
		json.Unmarshal(docJSON, &getDocJSON)

		key := getDocJSON["key"].(string)
		_, err = col.Upsert(key, doc, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
