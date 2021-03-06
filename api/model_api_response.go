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

import "github.com/romiras/go-openvz-api/models"

type (
	ApiResponse struct {
		Code    int32  `json:"code,omitempty"`
		Type    string `json:"type,omitempty"`
		Message string `json:"message,omitempty"`
	}

	AddContainerResponse struct {
		ApiResponse
		JobID string `json:"job_id,omitempty"`
	}

	GetContainerByIdResponse struct {
		ApiResponse
		Container *models.Container `json:"container"`
	}

	ListContainersResponse struct {
		ApiResponse
		Containers []*models.Container `json:"containers"`
	}

	GetJobByIdResponse struct {
		ApiResponse
		Status     string  `json:"status"`
		EntityType *string `json:"entity_type,omitempty"`
		EntityID   *string `json:"entity_id,omitempty"`
	}
)
