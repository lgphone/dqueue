package deley_queue

import (
	"context"
	"dqueue/job"
	"dqueue/store"
	"errors"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

const (
	pushExpiredJobLua = `
local zset_key = KEYS[1]
local job_key_prefix = KEYS[2]
local ready_key = KEYS[3]
local now = ARGV[1]
local limit = ARGV[2]

local expiredMembers = redis.call("ZRANGEBYSCORE", zset_key, 0, now, "LIMIT", 0, limit)
if #expiredMembers == 0 then
	return {}
end
for _,member in ipairs(expiredMembers) do
    local job_key = table.concat({job_key_prefix, member})
    if redis.call("GET", job_key) > 0 then
        redis.call("LPUSH", ready_key, member)
    end
end
redis.call("ZREM", zset_key, unpack(expiredMembers))
return expiredMembers
`
)

func PushJob(newJob *job.Job) error {
	if err := job.AddJob(newJob); err != nil {
		return err
	}
	ctx := context.Background()
	if newJob.Delay > 0 {
		_, err := store.RedisCli.ZAdd(ctx, store.QueueKey, &redis.Z{Member: newJob.ID, Score: float64(newJob.Delay + int(time.Now().Unix()))}).Result()
		return err
	}
	_, err := store.RedisCli.LPush(ctx, newJob.ID).Result()
	return err
}

func pushExpiredJobToReadyQueue(limit int) ([]string, error) {
	script := redis.NewScript(pushExpiredJobLua)
	now := time.Now().Unix()
	result, err := script.Run(context.Background(), store.RedisCli, []string{store.QueueKey, store.JobInfoKeyPrefix, store.ReadyKey}, now, limit).Result()
	if err != nil {
		return nil, err
	}
	fields, ok := result.([]interface{})
	if !ok {
		return nil, errors.New("lua return value should be elements array")
	}
	jobIds := make([]string, 0)
	for _, item := range fields {
		jobId, ok := item.(string)
		if !ok {
			return nil, errors.New("invalid lua value type")
		}
		jobIds = append(jobIds, jobId)
	}
	return jobIds, err
}

func LoopPushExpiredJobReadyQueue(limit int) {
	for {
		func() {
			defer func() {
				if err := recover(); err != nil {
					log.Println("ERROR: err to loop push expired job to readyQueue: ", err)
				}
			}()
			jobIds, err := pushExpiredJobToReadyQueue(limit)
			if err != nil {
				log.Println(err.Error())
			}
			if len(jobIds) > 0 {
				for _, item := range jobIds {
					log.Printf("INFO: %v job expired and send to readyQueue.", item)
				}
			}
			time.Sleep(time.Millisecond * 100)
		}()
	}
}
