package job

import (
	"context"
	"dqueue/store"
	"dqueue/utils"
	"encoding/json"
	"time"
)

type Job struct {
	ID         string   `json:"id"`
	Task       string   `json:"task"`
	Args       []string `json:"args"`
	Delay      int      `json:"delay"`
	TTL        int      `json:"ttl"`
	CreatedAt  int      `json:"created_at"`
	StartAt    *int     `json:"start_at"`
	FinishedAt *int     `json:"finished_at"`
}

func NewJobWithId(taskName string, delay, ttl int, args ...string) *Job {
	return &Job{
		ID:        utils.NewUUID(),
		Task:      taskName,
		Delay:     delay,
		TTL:       ttl,
		Args:      args,
		CreatedAt: int(time.Now().Unix()),
	}
}

func NewJob(jobId, taskName string, delay, ttl int, args ...string) *Job {
	return &Job{
		ID:        jobId,
		Task:      taskName,
		Delay:     delay,
		TTL:       ttl,
		Args:      args,
		CreatedAt: int(time.Now().Unix()),
	}
}

func JobKey(jobId string) string {
	return store.JobInfoKeyPrefix + jobId
}

func GetJob(jobId string) (*Job, error) {
	job := &Job{}
	result, err := store.RedisCli.Get(context.Background(), JobKey(jobId)).Result()
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal([]byte(result), job); err != nil {
		return nil, err
	}
	return job, nil
}

func DeleteJob(jobId string) error {
	_, err := store.RedisCli.Del(context.Background(), JobKey(jobId)).Result()
	return err
}

func AddJob(job *Job) error {
	jobByte, err := json.Marshal(job)
	if err != nil {
		return err
	}
	_, err = store.RedisCli.Set(context.Background(), JobKey(job.ID), jobByte, 0).Result()
	return err
}
