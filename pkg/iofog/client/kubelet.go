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

package client

import "encoding/json"

func (clt *Client) GetKubeletToken() (accessToken string, err error) {
	// Prepare request
	body, err := clt.doRequest("GET", "/k8s/vk-token", nil)
	if err != nil {
		return
	}

	var response GetKubeletTokenResponse
	// Return body
	if err = json.Unmarshal(body, &response); err != nil {
		return
	}

	accessToken = response.Token
	return
}
