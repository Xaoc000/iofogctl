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

package application

import (
	"fmt"
	"github.com/eclipse-iofog/iofog-go-sdk/pkg/client"
	"github.com/eclipse-iofog/iofogctl/internal"
	"github.com/eclipse-iofog/iofogctl/internal/config"
	"github.com/eclipse-iofog/iofogctl/pkg/util"
)

func Execute(namespace, name, newName string) error {
	util.SpinStart(fmt.Sprintf("Renaming Application %s", name))

	// Init remote resources
	clt, err := internal.NewControllerClient(namespace)
	if err != nil {
		return err
	}

	flow, err := clt.GetFlowByName(name)
	if err != nil {
		return err
	}

	flow.Name = newName
	_, err = clt.UpdateFlow(&client.FlowUpdateRequest{
		ID:   flow.ID,
		Name: &newName,
	})
	if err != nil {
		return err
	}
	config.Flush()
	return nil
}
