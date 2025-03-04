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

package describe

import (
	apps "github.com/eclipse-iofog/iofog-go-sdk/pkg/apps"
	"github.com/eclipse-iofog/iofog-go-sdk/pkg/client"
	"github.com/eclipse-iofog/iofogctl/internal"
	"github.com/eclipse-iofog/iofogctl/internal/config"
	"github.com/eclipse-iofog/iofogctl/pkg/util"
)

type applicationExecutor struct {
	namespace string
	name      string
	filename  string
	flow      *client.FlowInfo
	client    *client.Client
	msvcs     []client.MicroserviceInfo
	msvcPerID map[string]*client.MicroserviceInfo
}

func newApplicationExecutor(namespace, name, filename string) *applicationExecutor {
	a := &applicationExecutor{}
	a.namespace = namespace
	a.name = name
	a.filename = filename
	return a
}

func (exe *applicationExecutor) init() (err error) {
	exe.client, err = internal.NewControllerClient(exe.namespace)
	if err != nil {
		return
	}
	exe.flow, err = exe.client.GetFlowByName(exe.name)
	if err != nil {
		return
	}
	msvcListResponse, err := exe.client.GetMicroservicesPerFlow(exe.flow.ID)
	if err != nil {
		return
	}

	// Filter system microservices
	for _, msvc := range msvcListResponse.Microservices {
		if util.IsSystemMsvc(msvc) {
			continue
		}
		exe.msvcs = append(exe.msvcs, msvc)
	}
	exe.msvcPerID = make(map[string]*client.MicroserviceInfo)
	for i := 0; i < len(exe.msvcs); i++ {
		exe.msvcPerID[exe.msvcs[i].UUID] = &exe.msvcs[i]
	}
	return
}

func (exe *applicationExecutor) GetName() string {
	return exe.name
}

func (exe *applicationExecutor) Execute() error {
	// Fetch data
	if err := exe.init(); err != nil {
		return err
	}

	yamlMsvcs := []apps.Microservice{}
	yamlRoutes := []apps.Route{}

	for _, msvc := range exe.msvcs {
		yamlMsvc, err := MapClientMicroserviceToDeployMicroservice(&msvc, exe.client)
		if err != nil {
			return err
		}
		for _, route := range msvc.Routes {
			yamlRoutes = append(yamlRoutes, apps.Route{
				From: yamlMsvc.Name,
				To:   exe.msvcPerID[route].Name,
			})
		}
		// Remove fields
		yamlMsvc.Routes = nil
		yamlMsvc.Flow = nil
		yamlMsvcs = append(yamlMsvcs, *yamlMsvc)
	}

	application := apps.Application{
		Name:          exe.flow.Name,
		Microservices: yamlMsvcs,
		Routes:        yamlRoutes,
		ID:            exe.flow.ID,
	}

	header := config.Header{
		APIVersion: internal.LatestAPIVersion,
		Kind:       apps.ApplicationKind,
		Metadata: config.HeaderMetadata{
			Namespace: exe.namespace,
			Name:      exe.name,
		},
		Spec: application,
	}

	if exe.filename == "" {
		if err := util.Print(header); err != nil {
			return err
		}
	} else {
		if err := util.FPrint(header, exe.filename); err != nil {
			return err
		}
	}
	return nil
}
