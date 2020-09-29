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

// ListContainers - List containers
func ListContainers(c *gin.Context, registry *registries.Registry) {
	containers, err := registry.ContainerAPIService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.InvalidRequest(err))
		return
	}

	c.JSON(http.StatusOK, containers)
}

// CreateContainer - Create a new container
func CreateContainer(c *gin.Context, registry *registries.Registry) {
	var req *api.AddContainerRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.InvalidRequest(err))
		return
	}

	err = api.ValidateAddContainerRequest(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.InvalidRequest(err))
		return
	}

	container, err := registry.ContainerAPIService.Create(req)
	if err != nil {
		if err.Error() == "duplicate-name" {
			c.JSON(http.StatusUnprocessableEntity, api.InvalidRequest(errors.New("A container with given name already exists")))
			return
		}
		c.JSON(http.StatusInternalServerError, api.InvalidRequest(err))
		return
	}

	c.JSON(http.StatusOK, container)
}

// DeleteContainer - Deletes a container
func DeleteContainer(c *gin.Context, registry *registries.Registry) {
	_, err := handleFindByID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.InvalidRequest(err))
	}

	c.JSON(http.StatusOK, api.ApiResponse{
		Code:    0,
		Message: "success",
	})
}

// GetContainerById - Find container by ID
func GetContainerById(c *gin.Context, registry *registries.Registry) {
	id, err := handleFindByID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.InvalidRequest(err))
	}

	container, err := registry.ContainerAPIService.GetById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, api.InvalidRequest(errors.New("no such container")))
			return
		}
		c.JSON(http.StatusInternalServerError, api.InvalidRequest(err))
		return
	}

	c.JSON(http.StatusOK, container)
}

// UpdateContainer - Updates a container
func UpdateContainer(c *gin.Context, registry *registries.Registry) {
	var req *api.UpdateContainerRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.InvalidRequest(err))
		return
	}

	req.ID, err = handleFindByID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.InvalidRequest(err))
	}

	resp, err := registry.ContainerAPIService.Update(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.InvalidRequest(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func handleFindByID(c *gin.Context) (string, error) {
	id := c.Param("id")

	err := api.ValidateGetContainerByIdRequest(id)
	if err != nil {
		return id, err
	}

	return id, nil
}
