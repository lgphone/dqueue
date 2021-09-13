package client

import (
	"context"
	"dqueue/deley_queue"
	"dqueue/job"
	"dqueue/store"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"
)

type Client struct {
	Name string
	task sync.Map
}

func NewClient(name string) *Client {
	return &Client{
		Name: name,
	}
}

func (w *Client) RegisterTask(taskName string, taskFn interface{}) {
	w.task.LoadOrStore(taskName, taskFn)
}

func (w *Client) ListTask() []string {
	taskList := make([]string, 0)
	w.task.Range(func(key, val interface{}) bool {
		if keyString, ok := key.(string); ok {
			taskList = append(taskList, keyString)
		}
		return true
	})
	return taskList
}

func (w *Client) IsRegisteredTask(taskName string) bool {
	_, ok := w.task.Load(taskName)
	return ok
}

func (w *Client) processTask() error {
	result, err := store.RedisCli.BRPop(context.Background(), time.Second*30, store.ReadyKey).Result()
	if err != nil {
		return err
	}
	jobInfo, err := job.GetJob(result[1])
	if err != nil {
		return err
	}
	if err = deley_queue.PushJob(jobInfo); err != nil {
		return err
	}
	taskFn, ok := w.task.Load(jobInfo.Task)
	if !ok {
		log.Printf("not registered task: %s!", jobInfo.Task)
		return errors.New(fmt.Sprintf("not registered task: %s!", jobInfo.Task))
	}
	func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("process task err: %s!", err)
			}
		}()
		refTaskFn := reflect.ValueOf(taskFn)
		refTaskFnParams := make([]reflect.Value, 0)
		for _, item := range jobInfo.Args {
			refTaskFnParams = append(refTaskFnParams, reflect.ValueOf(item))
		}
		refTaskFn.Call(refTaskFnParams)
	}()
	return nil
}

func (w *Client) LoopProcessTask() {
	for {
		if err := w.processTask(); err != nil {
			log.Printf("process task err: %s!", err)
		}
	}
}

func (w *Client) PushTask(job *job.Job) error {
	if !w.IsRegisteredTask(job.Task) {
		return errors.New(fmt.Sprintf("not registered task: %s!", job.Task))
	}
	return deley_queue.PushJob(job)
}
