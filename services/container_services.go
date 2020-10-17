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
	"database/sql"
	"encoding/json"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/romiras/go-openvz-api/api"
	"github.com/romiras/go-openvz-api/models"

	openvzcmd "github.com/romiras/go-openvz-cmd"
)

type (
	DBConnection = *sqlx.DB

	ContainerAPIService struct {
		DB        DBConnection
		Commander *openvzcmd.POCCommanderStub
	}
)

func NewContainerAPIService(db DBConnection, cmd *openvzcmd.POCCommanderStub) *ContainerAPIService {
	return &ContainerAPIService{
		DB:        db,
		Commander: cmd,
	}
}

func (srv *ContainerAPIService) Create(req *api.AddContainerRequest) (*api.AddContainerResponse, error) {
	if srv.hasContainerWithName(req.Name) {
		return nil, errors.New("duplicate-name")
	}

	payload, err := json.Marshal(AddContainerJob{
		Name:       req.Name,
		OSTemplate: req.OSTemplate,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	jobID := uuid.New().String()

	_, err = srv.DB.Exec("INSERT INTO jobs (id, status, payload, type) VALUES (?, ?, ?, ?)", jobID, models.PENDING, payload, AddContainerType)
	if err != nil {
		return nil, err
	}

	return &api.AddContainerResponse{
		ApiResponse: api.ApiResponse{
			Code:    0,
			Message: "success",
		},
		JobID: jobID,
	}, nil
}

func (srv *ContainerAPIService) Update(req *api.UpdateContainerRequest) (*api.ApiResponse, error) {
	jsonData, err := json.Marshal(req.Parameters)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	container, err := srv.findContainerByID(req.ID)
	if err != nil {
		return nil, err
	}

	err = srv.Commander.SetContainerParameters(container.Name, req.Parameters)
	if err != nil {
		return nil, err
	}

	statement, err := srv.DB.Prepare("UPDATE containers SET parameters=? WHERE id=?")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = statement.Exec(jsonData, req.ID)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &api.ApiResponse{
		Code:    0,
		Message: "success",
	}, nil
}

func (srv *ContainerAPIService) findContainerByID(id string) (*models.Container, error) {
	var container models.Container

	err := srv.DB.Get(&container, "SELECT * FROM containers WHERE id=? LIMIT 1", id)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No rows.")
		return nil, err
	case err != nil:
		log.Fatal(err)
	}
	err = container.UnmarshalParametersDB()
	if err != nil {
		log.Fatal(err)
	}

	return &container, nil
}

func (srv *ContainerAPIService) hasContainerWithName(name string) bool {
	var i int
	err := srv.DB.DB.QueryRow("SELECT 1 FROM containers WHERE name=? LIMIT 1", name).Scan(&i)
	switch {
	case err == sql.ErrNoRows:
		return false
	case err != nil:
		log.Fatal(err)
	}

	return true
}

func (srv *ContainerAPIService) GetById(id string) (*api.GetContainerByIdResponse, error) {
	container, err := srv.findContainerByID(id)
	if err != nil {
		return nil, err
	}

	return &api.GetContainerByIdResponse{
		ApiResponse: api.ApiResponse{
			Code:    0,
			Message: "success",
		},
		Container: container,
	}, nil
}

func (srv *ContainerAPIService) List() (*api.ListContainersResponse, error) {
	containers := make([]*models.Container, 0)

	err := srv.DB.Select(&containers, "SELECT * FROM containers LIMIT 1000")
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No rows.")
	case err != nil:
		log.Fatal(err)
	}

	for _, container := range containers {
		err = container.UnmarshalParametersDB()
		if err != nil {
			log.Fatal(err)
		}
	}

	return &api.ListContainersResponse{
		ApiResponse: api.ApiResponse{
			Code:    0,
			Message: "success",
		},
		Containers: containers,
	}, nil
}

func (srv *ContainerAPIService) Delete(id string) (*api.ApiResponse, error) {
	container, err := srv.findContainerByID(id)
	if err != nil {
		return nil, err
	}

	err = srv.Commander.DeleteContainer(container.Name)
	if err != nil {
		return nil, err
	}

	_, err = srv.DB.Exec("DELETE FROM containers WHERE id=?", id)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &api.ApiResponse{
		Code:    0,
		Message: "success",
	}, nil
}
