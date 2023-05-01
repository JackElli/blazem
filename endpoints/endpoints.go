package endpoints

import (
	types "blazem/domain/endpoint"
	"blazem/domain/global"
	"blazem/endpoints/handlers/connect"
	"blazem/endpoints/handlers/datainfolder"
	"blazem/endpoints/handlers/doc"
	"blazem/endpoints/handlers/folder"
	"blazem/endpoints/handlers/nodemap"
	"blazem/endpoints/handlers/parent"
	"blazem/endpoints/handlers/ping"
	"blazem/endpoints/handlers/query"
	"blazem/endpoints/handlers/recentquery"
	"blazem/endpoints/handlers/removenode"
	"blazem/endpoints/handlers/rules"
	"blazem/endpoints/handlers/stats"
	"net/http"
)

// Create all of the endpoints for Blazem
func SetupEndpoints(node *global.Node) error {
	var endpoints = []types.Endpoint{
		{
			Route:       "/nodemap",
			Handler:     nodemap.NewNodeMapHandler(nil),
			Description: "We return the status of the nodemap; a list of nodes",
			Type:        types.SYNC,
			Node:        node,
		},
		{
			Route:       "/connect",
			Handler:     connect.NewConnectHandler(nil),
			Description: "We want to connect a node to the cluster",
			Type:        types.SYNC,
			Node:        node,
		},
		{
			Route:       "/getDoc",
			Handler:     doc.NewGetDocHandler(nil),
			Description: "We want to fetch a doc from blazem and send it back",
			Type:        types.SYNC,
			Node:        node,
		},
		{
			Route:       "/folders",
			Handler:     folder.NewFolderHandler(nil),
			Description: "We want to fetch all of the root folders currently stored in blazem",
			Type:        types.SYNC,
			Node:        node,
		},
		{
			Route:       "/removeNode",
			Handler:     removenode.NewRemoveNodeHandler(nil),
			Description: "We want to remove a node from the blazem cluster",
			Type:        types.SYNC,
			Node:        node,
		},
		{
			Route:       "/stats",
			Handler:     stats.NewStatsHandler(nil),
			Description: "We want to fetch the metrics of the current blazem OS",
			Type:        types.SYNC,
			Node:        node,
		},
		{
			Route:       "/getQuery",
			Handler:     query.NewQueryHandler(nil),
			Description: "We want to execute a query on blazem",
			Type:        types.SYNC,
			Node:        node,
		},
		{
			Route:       "/getRecentQueries",
			Handler:     recentquery.NewRecentQueryHandler(nil),
			Description: "We want to fetch all of the recent queries",
			Type:        types.SYNC,
			Node:        node,
		},
		{
			Route:       "/addRule",
			Handler:     rules.NewAddRuleHandler(nil),
			Description: "We want to add a rule to blazem",
			Type:        types.SYNC,
			Node:        node,
		},
		{
			Route:       "/removeRule",
			Handler:     rules.NewRemoveRuleHandler(nil),
			Description: "We want to remove a rule from blazem",
			Type:        types.SYNC,
			Node:        node,
		},
		{
			Route:       "/runRule",
			Handler:     nil,
			Description: "We want to run a rule now",
			Type:        types.SYNC,
			Node:        node,
		},
		{
			Route:       "/getRules",
			Handler:     rules.NewGetRulesHandler(nil),
			Description: "We want to get all of the rules currently in blazem",
			Type:        types.SYNC,
			Node:        node,
		},
		{
			Route:       "/getDataInFolder",
			Handler:     datainfolder.NewGetDataFolderHandler(nil),
			Description: "We want to fetch all of the data currently in the specified folder",
			Type:        types.ASYNC,
			Node:        node,
		},
		{
			Route:       "/deleteDoc",
			Handler:     doc.NewDeleteDocHandler(nil),
			Description: "We want to delete a document from blazem",
			Type:        types.ASYNC,
			Node:        node,
		},
		{
			Route:       "/ping",
			Handler:     ping.NewPingHandler(nil),
			Description: "We want to send a message to all followers with any data changes",
			Type:        types.ASYNC,
			Node:        node,
		},
		{
			Route:       "/addDoc",
			Handler:     doc.NewAddDocHandler(nil),
			Description: "We want to add a document to blazem",
			Type:        types.ASYNC,
			Node:        node,
		},
		{
			Route:       "/parents",
			Handler:     parent.NewParentHandler(nil),
			Description: "We want get the parent folders of that folder",
			Type:        types.ASYNC,
			Node:        node,
		},
	}

	for _, e := range endpoints {
		if e.Handler == nil {
			continue
		}
		if e.Type == types.SYNC {
			http.HandleFunc(e.Route, e.Handler(&e))
		} else {
			go http.HandleFunc(e.Route, e.Handler(&e))
		}
	}
	return nil
}
