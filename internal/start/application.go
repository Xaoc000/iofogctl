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

package startapplication

import (
	"github.com/eclipse-iofog/iofogctl/internal"
	"github.com/eclipse-iofog/iofogctl/internal/execute"
)

type Options struct {
	Namespace string
	Name      string
}

type executor struct {
	namespace string
	name      string
}

func NewExecutor(opt Options) (exe execute.Executor) {
	return &executor{
		name:      opt.Name,
		namespace: opt.Namespace,
	}
}

func (exe *executor) GetName() string {
	return exe.name
}

func (exe *executor) Execute() (err error) {
	clt, err := internal.NewControllerClient(exe.namespace)
	if err != nil {
		return err
	}

	flow, err := clt.GetFlowByName(exe.name)
	if err != nil {
		return err
	}

	_, err = clt.StartFlow(flow.ID)

	return
}
