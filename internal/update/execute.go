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

package deploy

import (
	"fmt"

	"github.com/eclipse-iofog/iofog-go-sdk/pkg/apps"
	"github.com/eclipse-iofog/iofogctl/internal/config"
	updateagent "github.com/eclipse-iofog/iofogctl/internal/update/agent"
	//updateagentconfig "github.com/eclipse-iofog/iofogctl/internal/update/agent-config"
	//updateapplication "github.com/eclipse-iofog/iofogctl/internal/update/application"
	//updateconnector "github.com/eclipse-iofog/iofogctl/internal/update/connector"
	//updatecontroller "github.com/eclipse-iofog/iofogctl/internal/update/controller"
	//updatecontrolplane "github.com/eclipse-iofog/iofogctl/internal/update/controlplane"
	//updatemicroservice "github.com/eclipse-iofog/iofogctl/internal/update/microservice"
	"github.com/eclipse-iofog/iofogctl/internal/execute"
)

var kindOrder = []apps.Kind{
	apps.ControlPlaneKind,
	apps.ControllerKind,
	apps.ConnectorKind,
	apps.AgentKind,
	config.AgentConfigKind,
	apps.ApplicationKind,
	apps.MicroserviceKind,
}

type Options struct {
	Namespace string
	InputFile string
}

//func updateApplication(namespace, name string, yaml []byte) (exe execute.Executor, err error) {
//	return updateapplication.NewExecutor(updateapplication.Options{Namespace: namespace, Yaml: yaml, Name: name})
//}
//
//func updateMicroservice(namespace, name string, yaml []byte) (exe execute.Executor, err error) {
//	return updatemicroservice.NewExecutor(updatemicroservice.Options{Namespace: namespace, Yaml: yaml, Name: name})
//}
//
//func updateControlPlane(namespace, name string, yaml []byte) (exe execute.Executor, err error) {
//	return updatecontrolplane.NewExecutor(updatecontrolplane.Options{Namespace: namespace, Yaml: yaml, Name: name})
//}

func updateAgent(namespace, name string, yaml []byte) (exe execute.Executor, err error) {
	return updateagent.NewExecutor(updateagent.Options{Namespace: namespace, Yaml: yaml, Name: name})
}

//func updateAgentConfig(namespace, name string, yaml []byte) (exe execute.Executor, err error) {
//	return updateagentconfig.NewExecutor(updateagentconfig.Options{Namespace: namespace, Yaml: yaml, Name: name})
//}
//
//func updateConnector(namespace, name string, yaml []byte) (exe execute.Executor, err error) {
//	return updateconnector.NewExecutor(updateconnector.Options{Namespace: namespace, Yaml: yaml, Name: name})
//}
//
//func updateController(namespace, name string, yaml []byte) (exe execute.Executor, err error) {
//	return updatecontroller.NewExecutor(updatecontroller.Options{Namespace: namespace, Yaml: yaml, Name: name})
//}

var kindHandlers = map[apps.Kind]func(string, string, []byte) (execute.Executor, error){
	//apps.ApplicationKind:   updateApplication,
	//apps.MicroserviceKind:  updateMicroservice,
	//apps.ControlPlaneKind:  updateControlPlane,
	apps.AgentKind: updateAgent,
	//config.AgentConfigKind: updateAgentConfig,
	//apps.ConnectorKind:     updateConnector,
	//apps.ControllerKind:    updateController,
}

// Execute deploy from yaml file
func Execute(opt Options) (err error) {
	executorsMap, err := execute.GetExecutorsFromYAML(opt.InputFile, opt.Namespace, kindHandlers)
	if err != nil {
		return err
	}

	// Run all kinds
	for idx := range kindOrder {
		if err = execute.RunExecutors(executorsMap[kindOrder[idx]], fmt.Sprintf("update %s", kindOrder[idx])); err != nil {
			return
		}
	}

	return nil
}
