// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package db

import (
	"auto-test-go/pkg/entities"
	"auto-test-go/pkg/util"
	"container/list"
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"time"
)

const (
	DoneUserContextBucket     = "DoneUserContextBucket"
	DoneScenarioContextBucket = "DoneScenarioContextBucket"
	BufferSize                = 10000
	QueueCount                = 1000
	Expire                    = "-10m"
)

type BufferTaskId struct {
	TaskIds      []string
	LastTimeSpan time.Time
}

type MemoryDbManager struct {
	currentDb      *bolt.DB
	MemoryMap      map[string]interface{}
	BufferQueue    *list.List
	Statistics     TaskStatistics
	finishTaskChan chan string
}

var BoltDbManager = &MemoryDbManager{}

type TaskStatistics struct {
	DoUserCommandCount       int `json:"doUserCommandCount"`
	DoScenarioCommandCount   int `json:"doScenarioCommandCount"`
	DoneScenarioCommandCount int `json:"doneScenarioCommandCount"`
	DoneUserCommandCount     int `json:"doneUserCommandCount"`
}

func (t *TaskStatistics) AddDoUserCaseCount() {
	t.DoUserCommandCount++
}

func (t *TaskStatistics) AddDoneUserCaseCount() {
	t.DoUserCommandCount--
	t.DoneUserCommandCount++
}

func (t *TaskStatistics) AddDoScenarioCaseCount() {
	t.DoScenarioCommandCount++
}
func (t *TaskStatistics) AddDoneScenarioCaseCount() {
	t.DoScenarioCommandCount--
	t.DoneScenarioCommandCount++
}

func (m *MemoryDbManager) Init() error {
	m.MemoryMap = map[string]interface{}{}
	db, err := bolt.Open("database/auto-test-engine.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		defer func(currentDb *bolt.DB) {
			err := currentDb.Close()
			if err != nil {
				util.Logger.Error("[MemoryDbManager] memory db close error: " + err.Error())
			}
		}(db)
		return err
	}

	m.currentDb = db
	m.BufferQueue = list.New()
	bucket := []string{DoneUserContextBucket, DoneScenarioContextBucket}
	for _, v := range bucket {
		m.createBucket(v)
	}
	go m.gc()
	m.finishTaskChan = make(chan string, 10)
	m.Statistics = TaskStatistics{}
	go m.addBuffer()
	return nil
}

func (m *MemoryDbManager) gc() {
	// 两种策略最近1w条数据不清理，大于1w小于1000w数据每隔半小时清理一次,大于1000w数据每隔10秒清理一次
	myTimer := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-myTimer.C:
			item := m.BufferQueue.Back()
			if item != nil {
				buffer := item.Value.(BufferTaskId)
				if m.BufferQueue.Len() > QueueCount {
					m.BufferQueue.Remove(item)
					go m.removeBuffer(buffer.TaskIds)
				} else {
					mm, _ := time.ParseDuration(Expire)
					date := time.Now().Add(mm)
					if date.After(buffer.LastTimeSpan) {
						m.BufferQueue.Remove(item)
						go m.removeBuffer(buffer.TaskIds)
					}
				}
			}
			myTimer.Reset(time.Second * 10)
		}
	}
}

func (m *MemoryDbManager) addBuffer() {
	buffer := make([]string, BufferSize)
	i := 0
	for {
		taskId := <-m.finishTaskChan
		if i > len(buffer)-1 {
			bufferTask := BufferTaskId{
				TaskIds:      buffer,
				LastTimeSpan: time.Now(),
			}
			m.BufferQueue.PushFront(bufferTask)
			buffer = make([]string, BufferSize)
			i = 0
		}

		buffer[i] = taskId
		i++
	}
}

func (m *MemoryDbManager) removeBuffer(buf []string) {
	bucket := []string{DoneUserContextBucket, DoneScenarioContextBucket}
	for _, v := range buf {
		for _, name := range bucket {
			err := m.remove(name, v)
			if err != nil {
				util.Logger.Error("[MemoryDbManager] memory db close error: " + err.Error())
			}
		}
	}
}

func (m *MemoryDbManager) createBucket(name string) {
	if m.currentDb == nil {
		panic(errors.New("current db is error"))
	}

	err := m.currentDb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(name))
		if b == nil {
			_, err := tx.CreateBucket([]byte(name))
			if err != nil {
				panic(err)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func (m *MemoryDbManager) RefreshUserContext(ctx *entities.ExecContext, finish bool) error {
	if ctx.TaskId == "" {
		return errors.New("context task id can not be null")
	}

	if finish {
		str, _ := json.Marshal(ctx)
		delete(m.MemoryMap, ctx.TaskId)
		err := m.update(DoneUserContextBucket, ctx.TaskId, str)
		m.finishTaskChan <- ctx.TaskId
		m.Statistics.AddDoneUserCaseCount()
		return err
	} else {
		m.Statistics.AddDoUserCaseCount()
		m.MemoryMap[ctx.TaskId] = ctx
		return nil
	}
}

func (m *MemoryDbManager) RefreshScenarioContext(ctx *entities.ScenarioContext, finish bool) error {
	if ctx.Self.TaskId == "" {
		return errors.New("scenario context task id can not be null")
	}

	if finish {
		str, _ := json.Marshal(ctx)
		delete(m.MemoryMap, ctx.Self.TaskId)
		err := m.update(DoneScenarioContextBucket, ctx.Self.TaskId, str)
		m.finishTaskChan <- ctx.Self.TaskId
		m.Statistics.AddDoneScenarioCaseCount()
		return err
	} else {
		m.Statistics.AddDoScenarioCaseCount()
		m.MemoryMap[ctx.Self.TaskId] = ctx
		return nil
	}
}

func (m *MemoryDbManager) QueryUserTask(taskId string) (ctx *entities.ExecContext, err error) {
	obj := m.MemoryMap[taskId]
	if obj != nil {
		ctx = obj.(*entities.ExecContext)
		return ctx, nil
	}

	err = m.currentDb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DoneUserContextBucket))
		v := b.Get([]byte(taskId))
		if v != nil {
			json.Unmarshal(v, &ctx)
		}
		return nil
	})
	return ctx, err
}

func (m *MemoryDbManager) QueryScenarioTask(taskId string) (ctx *entities.ScenarioContext, err error) {
	obj := m.MemoryMap[taskId]
	if obj != nil {
		ctx = obj.(*entities.ScenarioContext)
		return ctx, nil
	}

	err = m.currentDb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DoneScenarioContextBucket))
		v := b.Get([]byte(taskId))
		if v != nil {
			json.Unmarshal(v, &ctx)
		}
		return nil
	})
	return ctx, err
}

func (m *MemoryDbManager) remove(bucket string, key string) error {
	return m.currentDb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			err := b.Delete([]byte(key))
			if err != nil {
				util.Logger.Error("[MemoryDbManager] Current Context delete error: " + err.Error())
			}
		}
		return nil
	})
}

func (m *MemoryDbManager) update(bucket string, key string, value []byte) error {
	return m.currentDb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			err := b.Put([]byte(key), value)
			if err != nil {
				util.Logger.Error("[MemoryDbManager] Current Context update error: " + err.Error())
			}
		}
		return nil
	})
}
