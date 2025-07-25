package model

import (
	"encoding/json"
	"goexamples/gateway/openapi/proto"
	"os"
	"sync"
	"sync/atomic"

	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserModel struct {
	maxId   atomic.Int64
	Records []*proto.User
	mu      sync.RWMutex
}

func (m *UserModel) Total() int {
	return len(m.Records)
}

func (m *UserModel) MustLoad(file string) *UserModel {
	if file == "" {
		panic("file is empty")
	}

	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	users := []*proto.User{}
	if err := json.Unmarshal(data, &users); err != nil {
		panic(err)
	}
	for i := range users {
		now := timestamppb.Now()
		users[i].CreateAt = now
		users[i].UpdateAt = now
	}

	m.Records = append(m.Records, users...)
	m.maxId.Add(int64(len(users)))
	return m
}

func (m *UserModel) Create(record *proto.User, _ *fieldmaskpb.FieldMask) *proto.User {
	m.mu.Lock()
	defer m.mu.Unlock()

	record.Id = m.maxId.Add(1)
	now := timestamppb.Now()

	record.CreateAt = now
	record.UpdateAt = now
	m.Records = append(m.Records, record)
	return record
}

func (m *UserModel) Delete(id int64) (*proto.User, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, user := range m.Records {
		if user.Id == id {
			m.Records = append(m.Records[:i], m.Records[i+1:]...)
			return user, true
		}
	}
	return nil, false
}

func (m *UserModel) Update(record *proto.User, mask *fieldmaskpb.FieldMask) (*proto.User, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, user := range m.Records {
		if user.Id == record.Id {
			if mask == nil {
				return user, false
			}
			for _, field := range mask.Paths {
				switch field {
				case "name":
					user.Name = record.Name
				case "email":
					user.Email = record.Email
				}
			}
			user.UpdateAt = timestamppb.Now()
			return user, true
		}
	}
	return nil, false
}

func (m *UserModel) Get(id int64) (*proto.User, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, user := range m.Records {
		if user.Id == id {
			return user, true
		}
	}
	return nil, false
}

func (m *UserModel) List() []*proto.User {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.Records
}

func NewUserModel() *UserModel {
	return &UserModel{maxId: atomic.Int64{}}
}
