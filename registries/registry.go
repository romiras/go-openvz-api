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

package registries

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/romiras/go-openvz-api/services"
	openvzcmd "github.com/romiras/go-openvz-cmd"
)

const (
	SQL_CREATE_CONTAINERS = "CREATE TABLE IF NOT EXISTS containers (id CHAR(36) NOT NULL, name VARCHAR(255) NOT NULL, os_template VARCHAR(255) NOT NULL, parameters TEXT, created_at datetime default current_timestamp, CONSTRAINT rid_pkey PRIMARY KEY (id))"
	SQL_CREATE_JOBS       = "CREATE TABLE jobs (id uuid NOT NULL, type VARCHAR(255) NOT NULL, payload text NOT NULL, status integer NOT NULL, created_at timestamp NOT NULL default current_timestamp, locked_at timestamp, error_descr varchar(255), CONSTRAINT rid_pkey PRIMARY KEY (id))"
	SQL_CREATE_JOBS_INDEX = "CREATE INDEX jobs_status_locked_at_created_at_index ON jobs (status, locked_at, created_at)"
)

type Registry struct {
	ContainerAPIService *services.ContainerAPIService
	JobService          *services.JobService
	DB                  services.DBConnection
	Commander           *openvzcmd.POCCommanderStub
}

func NewRegistry() *Registry {
	db := InitializeDB()

	cmd, err := openvzcmd.NewPOCCommanderStub("vz_commands.yml")
	if err != nil {
		log.Fatal(err.Error())
	}

	return &Registry{
		ContainerAPIService: services.NewContainerAPIService(db, cmd),
		JobService:          services.NewJobService(db, cmd),
		DB:                  db,
		Commander:           cmd,
	}
}

func createTables(db services.DBConnection) error {
	_, err := db.Exec(SQL_CREATE_CONTAINERS)
	if err != nil {
		return err
	}

	_, err = db.Exec(SQL_CREATE_JOBS)
	if err != nil {
		return err
	}

	_, err = db.Exec(SQL_CREATE_JOBS_INDEX)
	if err != nil {
		return err
	}

	return nil
}

func InitializeDB() services.DBConnection {
	db := sqlx.MustConnect("sqlite3", ":memory:")
	if err := db.Ping(); err != nil {
		log.Fatal(err.Error())
	}

	err := createTables(db)
	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}
