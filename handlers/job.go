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

package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romiras/go-openvz-api/api"
	"github.com/romiras/go-openvz-api/registries"
)

// GetJobById - Find container by ID
func GetJobById(c *gin.Context, registry *registries.Registry) {
	id, err := handleFindByID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.InvalidRequest(err))
	}

	resp, err := registry.JobAPIService.GetById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, api.InvalidRequest(errors.New("no such job")))
			return
		}
		c.JSON(http.StatusInternalServerError, api.InvalidRequest(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// unused
func handleFindJobByID(c *gin.Context) (string, error) {
	id := c.Param("id")

	err := api.ValidateGetJobByIdRequest(id)
	if err != nil {
		return id, err
	}

	return id, nil
}
