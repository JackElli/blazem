package endpoints

import (
	handlers "blazem/endpoints/handlers"
	"blazem/global"
	"net/http"
)

type EndpointType string

const (
	ASYNC EndpointType = "async"
	SYNC  EndpointType = "sync"
)

type Endpoint struct {
	Node        *global.Node
	Route       string
	Handler     func(w http.ResponseWriter, req *http.Request)
	Description string
	Type        EndpointType
}

func SetupEndpoints(node *global.Node) {

	var endpoints = []Endpoint{
		{
			Route:       "/nodemap",
			Handler:     handlers.NodeMapHandler,
			Description: "We return the status of the nodemap; a list of nodes",
			Type:        SYNC,
			Node:        node,
		},
		{
			Route:       "/connect",
			Handler:     handlers.ConnectHandler((*handlers.Node)(node)),
			Description: "We want to connect a node to the cluster",
			Type:        SYNC,
		},
		{
			Route:       "/getDoc",
			Handler:     handlers.GetDocHandler((*handlers.Node)(node)),
			Description: "We want to fetch a doc from blazem and send it back",
			Type:        SYNC,
		},
		{
			Route:       "/folders",
			Handler:     handlers.FolderHandler((*handlers.Node)(node)),
			Description: "We want to fetch all of the root folders currently stored in blazem",
			Type:        SYNC,
		},
		{
			Route:       "/removeNode",
			Handler:     handlers.RemoveNodeHandler((*handlers.Node)(node)),
			Description: "We want to remove a node from the blazem cluster",
			Type:        SYNC,
		},
		{
			Route:       "/stats",
			Handler:     handlers.StatsHandler((*handlers.Node)(node)),
			Description: "We want to fetch the metrics of the current blazem OS",
			Type:        SYNC,
		},
		{
			Route:       "/getQuery",
			Handler:     handlers.QueryHandler((*handlers.Node)(node)),
			Description: "We want to execute a query on blazem",
			Type:        SYNC,
		},
		{
			Route:       "/getRecentQueries",
			Handler:     handlers.GetRecentQueriesHandler((*handlers.Node)(node)),
			Description: "We want to fetch all of the recent queries",
			Type:        SYNC,
		},
		{
			Route:       "/addRule",
			Handler:     handlers.AddRuleHandler((*handlers.Node)(node)),
			Description: "We want to add a rule to blazem",
			Type:        SYNC,
		},
		{
			Route:       "/removeRule",
			Handler:     handlers.RemoveNodeHandler((*handlers.Node)(node)),
			Description: "We want to remove a rule from blazem",
			Type:        SYNC,
		},
		{
			Route:       "/runRule",
			Handler:     nil,
			Description: "We want to run a rule now",
			Type:        SYNC,
		},
		{
			Route:       "/getRules",
			Handler:     handlers.GetRulesHandler((*handlers.Node)(node)),
			Description: "We want to get all of the rules currently in blazem",
			Type:        SYNC,
		},
		{
			Route:       "/getDataInFolder",
			Handler:     handlers.GetDataInFolder((*handlers.Node)(node)),
			Description: "We want to fetch all of the data currently in the specified folder",
			Type:        ASYNC,
		},
		{
			Route:       "/addDoc",
			Handler:     handlers.AddDocHandler((*handlers.Node)(node)),
			Description: "We want to add a document to blazem",
			Type:        ASYNC,
		},
		{
			Route:       "/deleteDoc",
			Handler:     handlers.DeleteDocHandler((*handlers.Node)(node)),
			Description: "We want to delete a document from blazem",
			Type:        ASYNC,
		},
		{
			Route:       "/ping",
			Handler:     handlers.PingHandler((*handlers.Node)(node)),
			Description: "We want to send a message to all followers with any data changes",
			Type:        ASYNC,
		},
	}

	for _, endpoint := range endpoints {
		if endpoint.Handler == nil {
			continue
		}
		if endpoint.Type == SYNC {
			http.HandleFunc(endpoint.Route, endpoint.Handler)
		} else {
			go http.HandleFunc(endpoint.Route, endpoint.Handler)
		}
	}
}

// DEPRECATED
// "addFolder":       node.addFolderHandler,
// "replicateFolder": node.replicateFolderHandler,
