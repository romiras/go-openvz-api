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
	"log"

	"../registries"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

// Run will start the server
func Run(reg *registries.Registry) {
	getRoutes(reg)
	log.Fatal(router.Run(":5000"))
}

// getRoutes will create our routes of our entire application
// this way every group of routes can be defined in their own file
// so this one won't be so messy
func getRoutes(reg *registries.Registry) {
	v1 := router.Group("/v0.1")
	addContainerRoutes(reg, v1)
}

func withRegistry(handler func(*gin.Context, *registries.Registry), registry *registries.Registry) func(*gin.Context) {
	return func(ctx *gin.Context) {
		handler(ctx, registry)
	}
}
