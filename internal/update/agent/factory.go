/*
 *  *******************************************************************************
 *  * Copyright (c) 2019 Edgeworx, Inc.
 *  *
 *  * This program and the accompanying materials are made available under the
 *  * terms of the Eclipse Public License v. 2.0 which is available at
 *  * http://www.eclipse.org/legal/epl-2.0
 *  *
 *  * SPDX-License-Identifier: EPL-2.0
 *  *******************************************************************************
 *
 */

package updateagent

import (
	"github.com/eclipse-iofog/iofogctl/internal/config"
	"github.com/eclipse-iofog/iofogctl/internal/execute"
	"github.com/eclipse-iofog/iofogctl/pkg/util"
	"gopkg.in/yaml.v2"
)

type Options struct {
	Namespace string
	Name      string
	Yaml      []byte
}

type executor struct {
	namespace string
	name      string
	currAgent config.Agent
	newAgent  config.Agent
}

func (exe executor) GetName() string {
	return exe.name
}

func (exe executor) Execute() error {
	// Point for brevity
	newAgent := &exe.newAgent
	currAgent := &exe.currAgent

	if newAgent.Host != "" {
		util.PrintNotify("Updating Agent Host field is not supported.")
	}

	if newAgent.KeyFile != "" {
		currAgent.KeyFile = newAgent.KeyFile
	}

	if newAgent.Port != 0 {
		currAgent.Port = newAgent.Port
	}

	if newAgent.User != "" {
		currAgent.User = newAgent.User
	}

	// Write to config the current Agent as updated with new Agent details
	config.UpdateAgent(exe.namespace, exe.currAgent)

	return config.Flush()
}

func NewExecutor(opt Options) (exe execute.Executor, err error) {
	// Check the namespace exists
	_, err = config.GetNamespace(opt.Namespace)
	if err != nil {
		return
	}

	// Check the agent exists
	currAgent, err := config.GetAgent(opt.Namespace, opt.Name)
	if err != nil {
		return
	}

	// Unmarshal file
	var newAgent config.Agent
	if err = yaml.Unmarshal(opt.Yaml, &newAgent); err != nil {
		err = util.NewInputError("Could not unmarshall\n" + err.Error())
		return
	}

	// Return executor
	exe = newExecutor(opt.Namespace, opt.Name, currAgent, newAgent)

	return
}

func newExecutor(namespace, name string, currAgent, newAgent config.Agent) (exe execute.Executor) {
	return executor{
		namespace: namespace,
		name:      name,
		currAgent: currAgent,
		newAgent:  newAgent,
	}
}
