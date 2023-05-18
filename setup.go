package main

import (
	"blazem/pkg/domain/logger"
	"blazem/pkg/domain/node"
	blazem_node "blazem/pkg/domain/node"
	"blazem/pkg/domain/user"
	"blazem/pkg/domain/users"
	"blazem/pkg/endpoints"
	"blazem/pkg/query"
	"fmt"
)

type SetupManager struct {
	Steps []blazem_node.SetupStep
	Node  *blazem_node.Node
}

// Returns a setupmgr with the steps to complete and the node
func CreateSetupMgr(node *blazem_node.Node, steps []node.SetupStep) SetupManager {
	return SetupManager{
		Steps: steps,
		Node:  node,
	}
}

// Runs all the steps in order
func (mgr *SetupManager) RunSteps() {
	logger.Logger.Info("Setting up Blazem.")
	for _, step := range mgr.Steps {
		if err := step.Fn(); err != nil {
			logger.Logger.Error("Found error in " + step.Description + " " + err.Error())
			return
		}
		logger.Logger.Info("Completed step.")
	}
	logger.Logger.Info("All steps completed successfully :)")
}

// Run the setup process by creating a setup mgr and running each
// step
func RunSetup(node *blazem_node.Node) {
	var masterip string = ""
	var localip = node.GetLocalIp()
	blazem_node.GlobalNode = node

	node.SetupLogger()

	mgr := CreateSetupMgr(node, []blazem_node.SetupStep{
		{
			Description: "Picks port for blazem to start on",
			Fn: func() error {
				go node.PickPort(localip)
				return nil
			},
		},
		{
			Description: "If this node is the master, set master attrs",
			Fn: func() error {
				if masterip == node.Ip {
					node.SetNodeMasterAttrs()
				}
				return nil
			},
		},
		{
			Description: "Loads users or creates an admin user",
			Fn: func() error {
				node.UserStore = users.NewUserStore()
				numOfUsers, err := node.UserStore.LoadUsers()
				if err != nil {
					return err
				}
				fmt.Println(numOfUsers)
				if numOfUsers == 0 {
					err = node.UserStore.Insert("user:1", &user.User{
						Id:       "user:1",
						Name:     "Jack Ellis",
						Username: "JackTest",
						Password: "test123",
						Role:     "admin",
					})
					if err != nil {
						return err
					}
				}
				return nil
			},
		},
		{
			Description: "Sets up blazem endpoints",
			Fn: func() error {
				if err := endpoints.SetupEndpoints(node); err != nil {
					return err
				}
				return nil
			},
		},
		{
			Description: "Adds this node to the nodemap",
			Fn: func() error {
				node.NodeMap = append(node.NodeMap, node)
				return nil
			},
		},

		{
			Description: "Read from local storage",
			Fn: func() error {
				node.ReadFromLocal()
				return nil
			},
		},
		{
			Description: "First ping and ping either the master or followers",
			Fn: func() error {
				go node.Ping()
				return nil
			},
		},
		{
			Description: "Load all query data into memory",
			Fn: func() error {
				query.LoadIntoMemory(node)
				return nil
			},
		},
	})
	mgr.RunSteps()
}
