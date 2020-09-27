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

package services

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/romiras/go-openvz-api/api"
	"github.com/romiras/go-openvz-api/models"
)

type (
	DBConnection = *sqlx.DB

	ContainerAPIService struct {
		DB DBConnection
	}
)

func InitializeDB() DBConnection {
	db := sqlx.MustConnect("sqlite3", ":memory:")
	if err := db.Ping(); err != nil {
		log.Fatal(err.Error())
	}

	return db
}

func NewContainerAPIService(db DBConnection) *ContainerAPIService {
	return &ContainerAPIService{
		DB: db,
	}
}

func (srv *ContainerAPIService) Create(req *api.AddContainerRequest) (*api.AddContainerResponse, error) {
	container := stubContainer()

	return &api.AddContainerResponse{
		ID: container.ID,
	}, nil
}

func (srv *ContainerAPIService) Update(req *api.UpdateContainerRequest) (*api.ApiResponse, error) {
	return &api.ApiResponse{
		Code:    0,
		Message: "success",
	}, nil
}

func (srv *ContainerAPIService) GetById(id string) (*api.GetContainerByIdResponse, error) {
	container := stubContainer()

	return &api.GetContainerByIdResponse{
		Container: container,
	}, nil
}

func (srv *ContainerAPIService) List() (*api.ListContainersResponse, error) {
	containers := make([]models.Container, 0)

	containers = append(containers, stubContainer())

	return &api.ListContainersResponse{
		ApiResponse: api.ApiResponse{
			Code:    0,
			Message: "success",
		},
		Containers: containers,
	}, nil
}

func stubContainer() models.Container {
	return models.Container{
		ID:         "100",
		Name:       "name1",
		OSTemplate: "template1",
	}
}
