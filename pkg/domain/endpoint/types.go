package endpoint

import "blazem/pkg/domain/global"

type WebNodeMap struct {
	Ip     string      `json:"ip"`
	Active bool        `json:"active"`
	Rank   global.Rank `json:"rank"`
}

type Stats struct {
	Cpu float64 `json:"cpu"`
	Ram float64 `json:"ram"`
}

type SendData struct {
	Key  string                 `json:"key"`
	Data map[string]interface{} `json:"data"`
}

type SendQueryData struct {
	Docs      []SendData `json:"docs"`
	TimeTaken int64      `json:"timeTaken"`
}

type Folder struct {
	Folder      string `json:"folder,omitempty"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	DocCount    int    `json:"docCount"`
	BackedUp    bool   `json:"backedUp"`
	CreatedBy   string `json:"createdBy"`
	Global      bool   `json:"global"`
	DateCreated string `json:"date"`
	Type        string `json:"type,omitempty"`
}

type DataInFolder struct {
	FolderName string     `json:"folderName"`
	Data       []SendData `json:"data"`
}
