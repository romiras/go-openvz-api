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
	"log"

	"github.com/romiras/go-openvz-api/api"
	"github.com/romiras/go-openvz-api/models"
	openvzcmd "github.com/romiras/go-openvz-cmd"
)

type JobAPIService struct {
	DB        DBConnection
	Commander *openvzcmd.POCCommanderStub
}

func NewJobAPIService(db DBConnection, cmd *openvzcmd.POCCommanderStub) *JobAPIService {
	return &JobAPIService{
		DB:        db,
		Commander: cmd,
	}
}

func (srv *JobAPIService) GetById(id string) (*api.GetJobByIdResponse, error) {
	job, err := srv.findJobByID(id)
	if err != nil {
		return nil, err
	}

	var entityType, entityID *string
	if job.EntityType.Valid {
		entityType = &job.EntityType.String
	}
	if job.EntityID.Valid {
		entityID = &job.EntityID.String
	}

	return &api.GetJobByIdResponse{
		ApiResponse: api.ApiResponse{
			Code:    0,
			Message: "success",
		},
		Status:     srv.getJobStatus(job.Status),
		EntityType: entityType,
		EntityID:   entityID,
	}, nil
}

func (srv *JobAPIService) findJobByID(id string) (*models.Job, error) {
	var job models.Job

	err := srv.DB.Get(&job, "SELECT id, type, status, entity_type, entity_id FROM jobs WHERE id=? LIMIT 1", id)
	switch {
	case err == sql.ErrNoRows:
		return nil, err
	case err != nil:
		log.Fatal(err)
	}

	return &job, nil
}

func (srv *JobAPIService) getJobStatus(status models.JobStatus) string {
	switch status {
	case models.PENDING:
		return "pending"
	case models.DONE:
		return "done"
	case models.FAILED:
		return "failed"
	default:
		return ""
	}
}
