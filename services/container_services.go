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
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

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

func InitializeDB() DBConnection {
	db := sqlx.MustConnect("sqlite3", ":memory:")
	if err := db.Ping(); err != nil {
		log.Fatal(err.Error())
	}
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS containers (id CHAR(36), name VARCHAR(255) NOT NULL, os_template VARCHAR(255) NOT NULL, parameters TEXT, created_at datetime default current_timestamp, CONSTRAINT rid_pkey PRIMARY KEY (id))")
	_, err := statement.Exec()
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	return db
}

func NewContainerAPIService(db DBConnection) *ContainerAPIService {
	cmd, err := openvzcmd.NewPOCCommanderStub("vz_commands.yml")
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	return &ContainerAPIService{
		DB:        db,
		Commander: cmd,
	}
}

func (srv *ContainerAPIService) Create(req *api.AddContainerRequest) (*api.AddContainerResponse, error) {
	id := uuid.New().String()

	_, err := srv.DB.Exec("INSERT INTO containers (id, name, os_template) VALUES (?, ?, ?)", id, req.Name, req.OSTemplate)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	err = srv.Commander.CreateContainer(req.Name, req.OSTemplate, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &api.AddContainerResponse{
		ApiResponse: api.ApiResponse{
			Code:    0,
			Message: "success",
		},
		ID: id,
	}, nil
}

func (srv *ContainerAPIService) Update(req *api.UpdateContainerRequest) (*api.ApiResponse, error) {
	jsonData, err := json.Marshal(req.Parameters)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	err = srv.Commander.SetContainerParameters(req.Parameters)
	if err != nil {
		log.Fatal(err.Error())
	}

	statement, err := srv.DB.Prepare("UPDATE containers SET parameters=? WHERE id=?")
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	_, err = statement.Exec(jsonData, req.ID)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	return &api.ApiResponse{
		Code:    0,
		Message: "success",
	}, nil
}

func (srv *ContainerAPIService) findContainer(id string) (*models.Container, error) {
	var container *models.Container

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

	return container, nil
}

func (srv *ContainerAPIService) GetById(id string) (*api.GetContainerByIdResponse, error) {
	container, err := srv.findContainer(id)
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
	container, err := srv.findContainer(id)
	if err != nil {
		return nil, err
	}

	err = srv.Commander.DeleteContainer(container.Name)
	if err != nil {
		log.Fatal(err.Error())
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
