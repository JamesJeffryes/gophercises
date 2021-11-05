package db

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"strconv"
	"time"
)

type Task struct {
	Name string
	Date time.Time
}

var db *bolt.DB

func Init(dbPath string) *bolt.DB {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	handleError(err)
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Tasks"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("CompletedTasks"))
		if err != nil {
			return err
		}
		return nil
	})
	return db
}

func AddTask(s string, duedate time.Time) {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))
		t := Task{s, duedate}
		buf, err := json.Marshal(t)
		if err != nil {
			return err
		}
		id, _ := b.NextSequence()
		err = b.Put([]byte(strconv.Itoa(int(id))), buf)
		return err
	})
	handleError(err)
	fmt.Printf("Added '%s' to db", s)
}

func GetTaskIDs() ([]string, error) {
	var ids []string
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))
		ids = make([]string, b.Stats().KeyN)
		i := 0
		err := b.ForEach(func(k, v []byte) error {
			ids[i] = string(k)
			i++
			return nil
		})
		return err
	})
	return ids, err
}

func ListTasks(complete bool) {
	err := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		var b *bolt.Bucket
		if !complete {
			b = tx.Bucket([]byte("Tasks"))
		} else {
			b = tx.Bucket([]byte("CompletedTasks"))
		}
		err := b.ForEach(func(k, v []byte) error {
			var task = Task{}
			err := json.Unmarshal(v, &task)
			if err != nil {
				return err
			}
			if task.Date == time.Unix(0, 0) {
				fmt.Printf("Task %s: %s\n", k, task.Name)
			} else {
				fmt.Printf("Task %s: %s %s\n", k, task.Name, task.Date.Format("Jan 2, 2006"))
			}
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	handleError(err)
}

func CompleteTask(s string) {
	db.Update(func(tx *bolt.Tx) error {
		bTodo := tx.Bucket([]byte("Tasks"))
		taskBytes := bTodo.Get([]byte(s))
		if taskBytes == nil {
			fmt.Printf("%s is not a valid taskBytes ID\n", s)
			ListTasks(false)
			os.Exit(1)
		}
		var task = Task{}
		err := json.Unmarshal(taskBytes, &task)
		if err != nil {
			return err
		}
		task.Date = time.Now()
		taskBytes, err = json.Marshal(task)
		if err != nil {
			return err
		}
		bComplete := tx.Bucket([]byte("CompletedTasks"))
		err = bComplete.Put([]byte(s), taskBytes)
		if err != nil {
			return err
		}
		err = bTodo.Delete([]byte(s))
		if err != nil {
			return err
		}
		fmt.Printf("%s has been marked complete\n", s)
		return nil
	})
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
