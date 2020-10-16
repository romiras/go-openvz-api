/*
 * Copyright 2020 Roman Miro
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api

import (
	"errors"

	openvzcmd "github.com/romiras/go-openvz-cmd"
)

const MissingParamError = " is missing or empty"

type (
	AddContainerRequest struct {
		Name       string `json:"name"`
		OSTemplate string `json:"ostemplate"`
	}

	// ContainerParameters map[string]string

	UpdateContainerRequest struct {
		ID         string
		Parameters openvzcmd.Options `json:"parameters"`
	}
)

func missingParam(param string) error {
	return errors.New(param + MissingParamError)
}

func InvalidRequest(err error) *ApiResponse {
	return &ApiResponse{
		Code:    100,
		Message: err.Error(),
		// Type: "",
	}
}

func ValidateAddContainerRequest(req *AddContainerRequest) error {
	if req.Name == "" {
		return missingParam("name")
	}
	if req.OSTemplate == "" {
		return missingParam("ostemplate")
	}

	return nil
}

func ValidateGetContainerByIdRequest(id string) error {
	if id == "" {
		return missingParam("id")
	}

	return nil
}

func ValidateGetJobByIdRequest(id string) error {
	if id == "" {
		return missingParam("id")
	}

	return nil
}

func ValidateUpdateContainerRequest(req *UpdateContainerRequest) error {
	return nil
}
