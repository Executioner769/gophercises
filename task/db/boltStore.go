package db

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

type BoltStore struct {
	tasksBucket []byte
	db          *bolt.DB
}

func NewBoltStore(dbPath, tasksBkt string) (*BoltStore, error) {
	var tasksBucket []byte
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	tasksBucket = []byte(tasksBkt)
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(tasksBucket)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &BoltStore{
		tasksBucket: tasksBucket,
		db:          db,
	}, nil
}

func (bs *BoltStore) CreateTask(task *Task) (*Task, error) {
	err := bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bs.tasksBucket)
		id64, _ := b.NextSequence()
		task.ID = int(id64)
		key := itob(task.ID)

		taskBytes, err := json.Marshal(task)
		if err != nil {
			return err
		}

		return b.Put(key, taskBytes)
	})

	if err != nil {
		return nil, err
	}
	return task, nil
}

func (bs *BoltStore) GetTasks() ([]*Task, error) {
	var tasks []*Task
	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bs.tasksBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			task := Task{}
			if err := json.Unmarshal(v, &task); err != nil {
				fmt.Println("Something went wrong: ", err)
			}
			tasks = append(tasks, &task)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (bs *BoltStore) DeleteTask(id int) (*Task, error) {
	var task *Task
	var err error

	task, err = bs.getTask(id)
	if err != nil {
		return nil, err
	}
	err = bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bs.tasksBucket)
		return b.Delete(itob(id))
	})
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (bs *BoltStore) CompleteTask(id int) (*Task, error) {
	var task *Task
	var err error

	task, err = bs.getTask(id)
	if err != nil {
		return nil, err
	}
	err = bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bs.tasksBucket)
		task.Done = true
		taskBytes, err := json.Marshal(task)
		if err != nil {
			return err
		}

		return b.Put(itob(id), taskBytes)
	})
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (bs *BoltStore) getTask(key int) (*Task, error) {
	var taskBytes []byte
	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bs.tasksBucket)
		taskBytes = b.Get(itob(key))
		if taskBytes == nil {
			return fmt.Errorf("oops no task found with key: %d", key)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	task := &Task{}
	err = json.Unmarshal(taskBytes, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// Integer to Byte Slice
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// Byte Slice to Integer
func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
