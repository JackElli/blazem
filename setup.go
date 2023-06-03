package main

import (
	"blazem/pkg/domain/logger"
	blazem_node "blazem/pkg/domain/node"
	"blazem/pkg/endpoints"
	"fmt"
)

type SetupStep struct {
	Description string
	Fn          func() error
}

type SetupManager struct {
	Steps []SetupStep
	Node  *blazem_node.Node
}

// Returns a setupmgr with the steps to complete and the node
func CreateSetupMgr(node *blazem_node.Node, steps []SetupStep) SetupManager {
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
		logger.Logger.Info(step.Description)
	}
	logger.Logger.Info("All steps completed successfully :)")
}

// Run the setup process by creating a setup mgr and running each
// step
func RunSetup(node *blazem_node.Node) {
	localip := node.GetLocalIp()
	_, err := node.SetupLogger()
	if err != nil {
		fmt.Println("Cannot set up logger")
		return
	}

	setupSteps := []SetupStep{
		{
			Description: "Picks port for blazem to start on",
			Fn: func() error {
				go node.PickPort(localip)
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
	}

	mgr := CreateSetupMgr(node, setupSteps)
	mgr.RunSteps()
}
