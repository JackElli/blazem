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
	var primaryData map[string]interface{}
	primaryIndex, _ := ioutil.ReadFile("index/primary.json")
	json.Unmarshal(primaryIndex, &primaryData)
	_, ok := primaryData[key]
	return ok
}

func (node *Node) AppendDataJson(key string, value JsonData) {
	readPrimaryIndex, _ := ioutil.ReadFile("index/primary.json")
	writePrimaryIndex, _ := os.OpenFile("index/primary.json",
		os.O_WRONLY, 0644)

	data := make(map[string]JsonData)
	data[key] = value
	writeDataJson, _ := json.Marshal(data)

	writeData := ",\n" + string(writeDataJson)[1:]
	writePrimaryIndex.WriteAt([]byte(writeData), (int64)(len(readPrimaryIndex)-1))

}

// This needs a lot of work
func (node *Node) ReplaceDataJson(key string, value JsonData) {
	readPrimaryIndex, _ := os.Open("index/primary.json")
	writePrimaryIndex, _ := os.OpenFile("index/primary.json",
		os.O_WRONLY, 0644)

	fileScanner := bufio.NewScanner(readPrimaryIndex)
	fileScanner.Split(bufio.ScanLines)

	data := make(map[string]JsonData)
	data[key] = value
	writeDataJson, _ := json.Marshal(data)

	indToWrite := 1
	lineNum := 0
	var fileBytes []byte
	readPrimaryIndex.Read(fileBytes)

	numOfLines := len(strings.Split(string(fileBytes), "\n"))

	for fileScanner.Scan() {
		lineSize := len(fileScanner.Bytes())
		lineText := fileScanner.Text()
		checkKey := "\"" + key + "\""

		// diff := lineSize - len(writeDataJson)
		if strings.Contains(lineText, checkKey) {
			// doc is on this line
			// fmt.Println("writing at ", indToWrite)
			writeData := ""
			// if lineNum == 0 {
			// 	writeData += "{"
			// }
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
			writePrimaryIndex.WriteAt([]byte(nothingData), (int64)(indToWrite))
			//write data
			writePrimaryIndex.WriteAt([]byte(writeData), (int64)(indToWrite))
		}
		indToWrite += lineSize
		lineNum++

	}
}

func (node *Node) SaveDataJson() {
	if node.Rank != MASTER {
		return
	}

	primaryIndex, err := os.OpenFile("index/primary.json",
		os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		err = os.MkdirAll("index/", os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
		primaryIndex, err = os.Create("index/primary.json")
	}

	writeDataJson, err := json.Marshal(node.Data)
	writeData := strings.ReplaceAll(string(writeDataJson), "},\"", "},\n\"")

	primaryIndex.Write([]byte(writeData))
}
