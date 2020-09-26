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

package routes

import (
	"../handlers"
	"../registries"
	"github.com/gin-gonic/gin"
)

func addContainerRoutes(reg *registries.Registry, grp *gin.RouterGroup) {
	containers := grp.Group("/containers")

	containers.GET("/", withRegistry(handlers.ListContainers, reg))
	containers.POST("/", withRegistry(handlers.CreateContainer, reg))
	containers.GET("/:id", withRegistry(handlers.GetContainerById, reg))
	containers.PATCH("/:id", withRegistry(handlers.UpdateContainer, reg))
	containers.DELETE("/:id", withRegistry(handlers.DeleteContainer, reg))
}
