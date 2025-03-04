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

package namespace

import (
	"fmt"
	"github.com/eclipse-iofog/iofogctl/internal/config"
	"github.com/eclipse-iofog/iofogctl/pkg/util"
)

func Execute(name, newName string) error {
	if name == "" || name == "default" {
		return util.NewError("Cannot rename default or nonexistant namespaces")
	}

	util.SpinStart(fmt.Sprintf("Renaming Namespace %s", name))

	if err := config.RenameNamespace(name, newName); err != nil {
		return err
	}
	return nil

}
