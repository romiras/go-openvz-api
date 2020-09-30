package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/romiras/go-openvz-api/api"
	openvzcmd "github.com/romiras/go-openvz-cmd"
)

const (
	JOB_CHECK_PERIOD = 3 * time.Second
	AddContainerType = "add-container"
)

type (
	Job struct {
		ID      string          `json:"id"`
		Type    string          `json:"type"`
		Payload json.RawMessage `json:"payload"`
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
	for {
		j.consumeJob()
		time.Sleep(JOB_CHECK_PERIOD)
	}
}

func (j *JobService) consumeJob() {
	var err error

	job := j.pickJob()
	if job == nil {
		log.Printf("No jobs.")
		return
	}

	err = j.lockJob(job.ID)
	if err != nil {
		log.Println(err.Error())
		return
	}

	req := j.parseAddContainerRequest(job)
	if req == nil {
		log.Println("CANNOT be parsed - skipped!")
		return
	}

	err = j.Commander.CreateContainer(req.Name, req.OSTemplate, nil)

	err = j.updateJobStatus(job.ID, err)
	if err != nil {
		log.Println(err.Error())
	}

	id := uuid.New().String() // UUID of container

	_, err = j.DB.Exec("INSERT INTO containers (id, name, os_template) VALUES (?, ?, ?)", id, req.Name, req.OSTemplate)
	if err != nil {
		log.Println(err.Error())
	}
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

func (j *JobService) updateJobStatus(jobID string, err error) error {
	if err != nil {
		_, err = j.DB.Exec("UPDATE jobs SET status=?, error_descr=?, locked_at=NULL WHERE id=?", FAILED, err.Error(), jobID)
		return err
	}
	_, err = j.DB.Exec("UPDATE jobs SET status=?, locked_at=NULL WHERE id=?", DONE, jobID)

	return err
}
