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

package movemicroservice

import (
	"fmt"

	"github.com/eclipse-iofog/iofog-go-sdk/pkg/client"
	"github.com/eclipse-iofog/iofogctl/internal"
	"github.com/eclipse-iofog/iofogctl/pkg/util"
)

func Execute(namespace, name, agent string) error {
	util.SpinStart(fmt.Sprintf("Moving microservice %s", name))

	// Init remote resources
	clt, err := internal.NewControllerClient(namespace)
	if err != nil {
		return err
	}

	msvc, err := clt.GetMicroserviceByName(name)
	if err != nil {
		return err
	}

	destAgent, err := clt.GetAgentByName(agent)
	if err != nil {
		return err
	}

	_, err = clt.UpdateMicroservice(client.MicroserviceUpdateRequest{
		UUID:      msvc.UUID,
		AgentUUID: &destAgent.UUID,
		// Bug in Controller, fails if empty because images should be an array
		Images: msvc.Images,
		// Ports and Routes get automatically updated by the SDK, to avoid deletion of port mapping or route, those fields are mandatory
		Ports:  msvc.Ports,
		Routes: msvc.Routes,
	})

	return err
}
