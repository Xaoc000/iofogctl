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

package updatecontroller

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
	namespace      string
	name           string
	currController config.Controller
	newController  config.Controller
}

func (exe executor) GetName() string {
	return exe.name
}

func (exe executor) Execute() error {
	// Point for brevity
	newController := &exe.newController
	currController := &exe.currController

	if newController.Host != "" {
		util.PrintNotify("Updating Controller Host field is not supported.")
	}

	// Disallow editing vanilla fields for k8s Controller
	if currController.KubeConfig != "" && (newController.Host != "" || newController.Port != 0 || newController.KeyFile != "") {
		return util.NewInputError("Controller " + exe.name + " is deployed on Kubernetes. You cannot add SSH details to this Controller")
	}

	// Disallow editing k8s fields for vanilla Controller
	if currController.Host != "" && currController.KeyFile != "" && newController.KubeConfig != "" {
		return util.NewInputError("Controller " + exe.name + " is not deployed on Kubernetes. You cannot add Kube Config details to this Controller")
	}

	if newController.KeyFile != "" {
		currController.KeyFile = newController.KeyFile
	}

	if newController.Port != 0 {
		currController.Port = newController.Port
	}

	if newController.User != "" {
		currController.User = newController.User
	}

	if newController.KubeConfig != "" {
		currController.KubeConfig = newController.KubeConfig
	}

	// Write to config the current Controller as updated with new Controller details
	config.UpdateController(exe.namespace, exe.currController)

	return config.Flush()
}

func NewExecutor(opt Options) (exe execute.Executor, err error) {
	// Check the namespace exists
	_, err = config.GetNamespace(opt.Namespace)
	if err != nil {
		return
	}

	// Check the agent exists
	currController, err := config.GetController(opt.Namespace, opt.Name)
	if err != nil {
		return
	}

	// Unmarshal file
	var newController config.Controller
	if err = yaml.Unmarshal(opt.Yaml, &newController); err != nil {
		err = util.NewInputError("Could not unmarshall\n" + err.Error())
		return
	}

	// Return executor
	exe = newExecutor(opt.Namespace, opt.Name, currController, newController)

	return
}

func newExecutor(namespace, name string, currController, newController config.Controller) (exe execute.Executor) {
	return executor{
		namespace:      namespace,
		name:           name,
		currController: currController,
		newController:  newController,
	}
}
