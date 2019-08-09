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

package cmd

import (
	"github.com/eclipse-iofog/iofogctl/internal/deploy/agent"
	"github.com/eclipse-iofog/iofogctl/pkg/util"
	"github.com/spf13/cobra"
)

func newDeployAgentCommand() *cobra.Command {
	// Instantiate options
	var opt deployagent.Options

	// Instantiate command
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Bootstrap and provision edge hosts",
		Long: `Bootstrap edge hosts with the ioFog Agent stack and provision them with a Controller in the namespace.

A Controller must first be deployed within the corresponding namespace in order to provision the Agent.`,
		Example: `iofogctl deploy agent -f agent.yaml`,
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			// Get agent name and namespace
			opt.Namespace, err = cmd.Flags().GetString("namespace")
			util.Check(err)

			// Execute the command
			err = deployagent.Deploy(opt)
			util.Check(err)

			util.PrintSuccess("Successfully deployed Agents to namespace " + opt.Namespace)
		},
	}

	// Register flags
	cmd.Flags().StringVarP(&opt.InputFile, "file", "f", "", "YAML file containing resource definitions for Agents")

	return cmd
}
