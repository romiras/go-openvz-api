package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/romiras/go-openvz-api/api"
	"github.com/romiras/go-openvz-api/models"
	openvzcmd "github.com/romiras/go-openvz-cmd"
)

const (
	AddContainerType = "add-container"
	ContainerType    = "container"
)

type (
	AddContainerJob struct {
		Name       string `json:"name"`
		OSTemplate string `json:"ostemplate"`
	}

	JobService struct {
		DB        DBConnection
		Commander *openvzcmd.POCCommanderStub
	}
)

func NewJobService(db DBConnection, cmd *openvzcmd.POCCommanderStub) *JobService {
	return &JobService{
		DB:        db,
		Commander: cmd,
	}
}

func (j *JobService) ConsumeJobs(jobInterval time.Duration) {
	var err error
	for {
		err = j.consumeJob()
		if err != nil {
			log.Println(err.Error()) // just log...
		}
		time.Sleep(jobInterval)
	}
}

func (j *JobService) consumeJob() error {
	var err error

	job := j.pickJob()
	if job == nil {
		log.Printf("No jobs.")
		return nil
	}

	err = j.lockJob(job.ID)
	if err != nil {
		return err
	}

	req := j.parseAddContainerRequest(job)
	if req == nil {
		return errors.New("CANNOT be parsed - skipped!")
	}

	id := uuid.New().String() // UUID of container

	err = j.Commander.CreateContainer(req.Name, req.OSTemplate, nil)

	err = j.updateJobStatus(job.ID, id, err)
	if err != nil {
		return err
	}

	_, err = j.DB.Exec("INSERT INTO containers (id, name, os_template) VALUES (?, ?, ?)", id, req.Name, req.OSTemplate)

	return err
}

func (j *JobService) pickJob() *models.Job {
	var job models.Job

	err := j.DB.Get(&job, "SELECT id, payload, type FROM jobs WHERE status=? AND locked_at IS NULL ORDER BY created_at LIMIT 1", models.PENDING)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		log.Fatal(err)
	}

	return &job
}

func (j *JobService) parseAddContainerRequest(job *models.Job) *api.AddContainerRequest {
	var err error
	var req *api.AddContainerRequest

	switch job.Type {
	case AddContainerType:
		err = json.Unmarshal(job.Payload, &req)
		if err != nil {
			log.Fatal(err)
		}
	}

	return req
}

func (j *JobService) lockJob(jobID string) error {
	_, err := j.DB.Exec("UPDATE jobs SET locked_at=? WHERE id=?", time.Now().UTC(), jobID)
	return err
}

func (j *JobService) updateJobStatus(jobID, id string, err error) error {
	if err != nil {
		_, err = j.DB.Exec("UPDATE jobs SET status=?, error_descr=?, locked_at=NULL WHERE id=?", models.FAILED, err.Error(), jobID)
		return err
	}
	_, err = j.DB.Exec("UPDATE jobs SET status=?, locked_at=NULL, entity_type=?, entity_id=? WHERE id=?", models.DONE, ContainerType, id, jobID)

	return err
}
