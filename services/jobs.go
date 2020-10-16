package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/romiras/go-openvz-api/api"
	openvzcmd "github.com/romiras/go-openvz-cmd"
)

type JobStatus int

const (
	PENDING = iota
	DONE
	FAILED
)

const (
	JOB_CHECK_PERIOD = 3 * time.Second
	AddContainerType = "add-container"
	ContainerType    = "container"
)

type (
	Job struct {
		ID         string          `json:"id" db:"id"`
		Type       string          `json:"type" db:"type"`
		Status     JobStatus       `json:"status,omitempty" db:"status"`
		Payload    json.RawMessage `json:"payload" db:"payload"`
		EntityType string          `json:"entity_type,omitempty" db:"entity_type"`
		EntityID   string          `json:"entity_id,omitempty" db:"entity_id"`
	}

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

func (j *JobService) ConsumeJobs() {
	var err error
	for {
		err = j.consumeJob()
		if err != nil {
			log.Println(err.Error()) // just log...
		}
		time.Sleep(JOB_CHECK_PERIOD)
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

func (j *JobService) pickJob() *Job {
	var job Job

	err := j.DB.Get(&job, "SELECT id, payload, type FROM jobs WHERE status=? AND locked_at IS NULL ORDER BY created_at LIMIT 1", PENDING)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		log.Fatal(err)
	}

	return &job
}

func (j *JobService) parseAddContainerRequest(job *Job) *api.AddContainerRequest {
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
		_, err = j.DB.Exec("UPDATE jobs SET status=?, error_descr=?, locked_at=NULL WHERE id=?", FAILED, err.Error(), jobID)
		return err
	}
	_, err = j.DB.Exec("UPDATE jobs SET status=?, locked_at=NULL, entity_type=?, entity_id=? WHERE id=?", DONE, ContainerType, id, jobID)

	return err
}
