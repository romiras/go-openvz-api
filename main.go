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

package main

import (
	"flag"

	"github.com/romiras/go-openvz-api/registries"
	"github.com/romiras/go-openvz-api/routes"
)

func main() {
	dsn := flag.String("dsn", ":memory:", "Data source name.")
	flag.Parse()

	registry := registries.NewRegistry(dsn)
	defer registry.DB.Close()

	// Run a job service in background.
	go registry.JobService.ConsumeJobs()

	// Our server will live in the routes package
	routes.Run(registry)
}
