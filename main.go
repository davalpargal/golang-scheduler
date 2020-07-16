package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"time"
)

//TODO: Add better sleep mechanism to save connections
//TODO: Add more worker threads for executeJobInQueueMethod

var client *redis.Client
var ctx = context.Background()

func main() {
	fmt.Println("Program started")
	connectToRedis()
	startPrintingList()

	fmt.Println("Entry into SCHEDULE sorted set will be picked based on least score")
}

func connectToRedis() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := client.Ping(ctx).Result()
	fmt.Println("redis connect response: ", pong, "err: ", err)
}

func startPrintingList() {
	//exit := make(chan string)

	executeJobsInQueue()

	for {
		fmt.Println("Checking...")

		result := client.ZPopMin(ctx, "SCHEDULE")
		fmt.Println("Found this: ", result)

		if len(result.Val()) > 0 {
			fmt.Println("Adding this to backup...")
			client.ZAdd(ctx, "BACKUP", &result.Val()[0])

			fmt.Println("Adding this to job queue: ", result.Val()[0].Member)
			client.LPush(ctx, "JOBS", result.Val()[0].Member)
		}

		time.Sleep(time.Second * 10)
		continue
	}
}

func executeJobsInQueue() {
	go func() {
		for {
			result := client.RPop(ctx, "JOBS")
			if result.Val() != "" {
				if isJobDone() {
					fmt.Println("Job done for: ", result.Val())
					client.ZRem(ctx, "BACKUP", result.Val())
					fmt.Println("removed ", result.Val(), " from backup")
				} else {
					fmt.Println("Job not done for: ", result.Val(), " still in backup list")
				}
			}
			time.Sleep(time.Second * 5)
		}
	}()
}

func isJobDone() bool {
	return (rand.Int63()*time.Now().UnixNano() % 2) == 0
}
