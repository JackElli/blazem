package global

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func IsDocInIndex(key string) bool {
	var data map[string]interface{}
	backup, _ := ioutil.ReadFile("backup/primary.json")
	json.Unmarshal(backup, &data)
	_, ok := data[key]
	return ok
}

func (node *Node) AppendDataJson(key string, value JsonData) {
	readBackup, _ := ioutil.ReadFile("backup/primary.json")
	writeBackup, _ := os.OpenFile("backup/primary.json",
		os.O_WRONLY, 0644)

	data := make(map[string]JsonData)
	data[key] = value
	writeDataJson, _ := json.Marshal(data)

	writeData := ",\n" + string(writeDataJson)[1:]
	writeBackup.WriteAt([]byte(writeData), (int64)(len(readBackup)-1))

}

// This needs a lot of work
func (node *Node) ReplaceDataJson(key string, value JsonData) {
	readBackup, _ := os.Open("backup/primary.json")
	writeBackup, _ := os.OpenFile("backup/primary.json",
		os.O_WRONLY, 0644)

	fileScanner := bufio.NewScanner(readBackup)
	fileScanner.Split(bufio.ScanLines)

	data := make(map[string]JsonData)
	data[key] = value
	writeDataJson, _ := json.Marshal(data)

	indToWrite := 1
	lineNum := 0
	var fileBytes []byte
	readBackup.Read(fileBytes)

	numOfLines := len(strings.Split(string(fileBytes), "\n"))

	for fileScanner.Scan() {
		lineSize := len(fileScanner.Bytes())
		lineText := fileScanner.Text()
		checkKey := "\"" + key + "\""

		if strings.Contains(lineText, checkKey) {
			// doc is on this line
			writeData := ""

			writeData += string(writeDataJson)[:len(writeDataJson)-1]
			writeData = writeData[1:]
			if lineNum != numOfLines {
				writeData += ",\n"

			} else {
				//end of lines
				writeData += "}"

			}
			// will get some weird formatting with
			// this but it will still work
			nothingData := strings.Repeat(" ", lineSize)
			//clear line
			writeBackup.WriteAt([]byte(nothingData), (int64)(indToWrite))
			//write data
			writeBackup.WriteAt([]byte(writeData), (int64)(indToWrite))
		}
		indToWrite += lineSize
		lineNum++
	}
}

func (node *Node) SaveDataJson() {
	if node.Rank != MASTER {
		return
	}

	backup, err := os.OpenFile("backup/primary.json",
		os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		err = os.MkdirAll("backup/", os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
		backup, err = os.Create("backup/primary.json")
	}

	writeDataJson, err := json.Marshal(node.Data)
	writeData := strings.ReplaceAll(string(writeDataJson), "},\"", "},\n\"")

	backup.Write([]byte(writeData))
}
