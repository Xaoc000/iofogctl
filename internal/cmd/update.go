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
	update "github.com/eclipse-iofog/iofogctl/internal/update"
	"github.com/eclipse-iofog/iofogctl/pkg/util"
	"github.com/spf13/cobra"
)

func newUpdateCommand() *cobra.Command {
	// Instantiate options
	var opt update.Options

	// Instantiate command
	cmd := &cobra.Command{
		Use:     "update",
		Example: `update -f resource.yaml`,
		Short:   "Update existing ioFog resources",
		Long:    `Update ioFog resources that have already been created.`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			// Get namespace
			opt.Namespace, err = cmd.Flags().GetString("namespace")
			util.Check(err)

			// Execute command
			err = update.Execute(opt)
			util.Check(err)

			util.PrintSuccess("Successfully updated resources")
		},
	}

	// Register flags
	cmd.Flags().StringVarP(&opt.InputFile, "file", "f", "", "YAML file containing one or more ioFog resource specifications")

	return cmd
}
