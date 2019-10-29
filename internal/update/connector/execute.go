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

package updateconnector

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
	namespace     string
	name          string
	currConnector config.Connector
	newConnector  config.Connector
}

func (exe executor) GetName() string {
	return exe.name
}

func (exe executor) Execute() error {
	// Point for brevity
	newConnector := &exe.newConnector
	currConnector := &exe.currConnector

	if newConnector.Host != "" {
		util.PrintNotify("Updating Connector Host field is not supported.")
	}

	// Disallow editing vanilla fields for k8s Connector
	if currConnector.KubeConfig != "" && (newConnector.Host != "" || newConnector.Port != 0 || newConnector.KeyFile != "") {
		return util.NewInputError("Connector " + exe.name + " is deployed on Kubernetes. You cannot add SSH details to this Connector")
	}

	// Disallow editing k8s fields for vanilla Connector
	if currConnector.Host != "" && currConnector.KeyFile != "" && newConnector.KubeConfig != "" {
		return util.NewInputError("Connector " + exe.name + " is not deployed on Kubernetes. You cannot add Kube Config details to this Connector")
	}

	if newConnector.KeyFile != "" {
		currConnector.KeyFile = newConnector.KeyFile
	}

	if newConnector.Port != 0 {
		currConnector.Port = newConnector.Port
	}

	if newConnector.User != "" {
		currConnector.User = newConnector.User
	}

	if newConnector.KubeConfig != "" {
		currConnector.KubeConfig = newConnector.KubeConfig
	}

	// Write to config the current Connector as updated with new Connector details
	config.UpdateConnector(exe.namespace, exe.currConnector)

	return config.Flush()
}

func NewExecutor(opt Options) (exe execute.Executor, err error) {
	// Check the namespace exists
	_, err = config.GetNamespace(opt.Namespace)
	if err != nil {
		return
	}

	// Check the agent exists
	currConnector, err := config.GetConnector(opt.Namespace, opt.Name)
	if err != nil {
		return
	}

	// Unmarshal file
	var newConnector config.Connector
	if err = yaml.Unmarshal(opt.Yaml, &newConnector); err != nil {
		err = util.NewInputError("Could not unmarshall\n" + err.Error())
		return
	}

	// Return executor
	exe = newExecutor(opt.Namespace, opt.Name, currConnector, newConnector)

	return
}

func newExecutor(namespace, name string, currConnector, newConnector config.Connector) (exe execute.Executor) {
	return executor{
		namespace:     namespace,
		name:          name,
		currConnector: currConnector,
		newConnector:  newConnector,
	}
}
