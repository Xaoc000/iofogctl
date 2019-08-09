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

package deployagent

import (
	"fmt"
	"sync"

	"github.com/eclipse-iofog/iofogctl/internal/config"
	"github.com/eclipse-iofog/iofogctl/pkg/iofog/install"
	"github.com/eclipse-iofog/iofogctl/pkg/util"
)

type executor interface {
	execute() error
}

type Options struct {
	Namespace string
	InputFile string
}

type jobResult struct {
	name string
	err  error
}

func Deploy(opt Options) error {
	// Check the namespace exists
	_, err := config.GetNamespace(opt.Namespace)
	if err != nil {
		return err
	}

	// Read the input file
	agents, err := UnmarshallYAML(opt.InputFile)
	if err != nil {
		return err
	}

	// Instantiate wait group for parallel tasks
	var wg sync.WaitGroup
	localAgentCount := 0
	errChan := make(chan jobResult, len(agents))
	for _, agent := range agents {

		// Check local deploys
		if util.IsLocalHost(agent.Host) {
			localAgentCount++
			if localAgentCount > 1 {
				fmt.Printf("Agent [%v] not deployed, you can only run one local agent.\n", agent.Name)
				continue
			}
		}

		var exe executor
		exe, err := newExecutor(opt.Namespace, agent)
		if err != nil {
			return err
		}

		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			err := exe.execute()
			errChan <- jobResult{
				err:  err,
				name: name,
			}
		}(agent.Name)
	}
	wg.Wait()
	close(errChan)

	// Output any errors
	failed := false
	for result := range errChan {
		if result.err != nil {
			failed = true
			util.PrintNotify("Failed to deploy " + result.name + ". " + result.err.Error())
		}

		if failed {
			return util.NewError("Failed to deploy one or more resources")
		}
	}

	return nil
}

func newExecutor(namespace string, agent config.Agent) (executor, error) {
	// Check the namespace exists
	ns, err := config.GetNamespace(namespace)
	if err != nil {
		return nil, err
	}

	// Check Controller exists
	nbControllers := len(ns.ControlPlane.Controllers)
	if nbControllers != 1 {
		errMessage := fmt.Sprintf("This namespace contains %d Controller(s), you must have one, and only one.", nbControllers)
		return nil, util.NewInputError(errMessage)
	}

	// Local executor
	if util.IsLocalHost(agent.Host) {
		cli, err := install.NewLocalContainerClient()
		if err != nil {
			return nil, err
		}
		exe, err := newLocalExecutor(namespace, agent, cli)
		if err != nil {
			return nil, err
		}
		return exe, nil
	}

	// Default executor
	if agent.Host == "" || agent.KeyFile == "" || agent.User == "" {
		return nil, util.NewInputError("Must specify user, host, and key file flags for remote deployment")
	}
	return newRemoteExecutor(namespace, agent), nil
}
