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

package deletecatalogitem

import (
	"github.com/eclipse-iofog/iofogctl/internal"
	"github.com/eclipse-iofog/iofogctl/internal/execute"
	"github.com/eclipse-iofog/iofogctl/pkg/util"
)

type Executor struct {
	namespace string
	name      string
}

func NewExecutor(namespace, name string) (execute.Executor, error) {
	exe := &Executor{
		namespace: namespace,
		name:      name,
	}

	return exe, nil
}

// GetName returns application name
func (exe *Executor) GetName() string {
	return exe.name
}

// Execute deletes application by deleting its associated flow
func (exe *Executor) Execute() (err error) {
	util.SpinStart("Deleting Microservice")
	// Init remote resources
	clt, err := internal.NewControllerClient(exe.namespace)
	if err != nil {
		return err
	}

	item, err := clt.GetMicroserviceByName(exe.name)
	if err != nil {
		return err
	}

	if err = clt.DeleteMicroservice(item.UUID); err != nil {
		return err
	}

	return nil
}
