package rooms

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

type CarState string

const (
	CarFree     CarState = "FREE"
	CarReserved CarState = "RESERVED"
	CarBusy     CarState = "BUSY"
)

type Car struct {
	ID         string    `json:"id"`
	State      CarState  `json:"state"`
	AssignedTo string    `json:"assignedTo,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type Registry struct {
	mu   sync.RWMutex
	cars map[string]*Car
}

func NewRegistry(n int) *Registry {
	r := &Registry{
		cars: make(map[string]*Car, n),
	}
	for i := 1; i <= n; i++ {
		id := "car-" + strconv.Itoa(i)
		r.cars[id] = &Car{
			ID:        id,
			State:     CarFree,
			UpdatedAt: time.Now(),
		}
	}
	return r
}

func (r *Registry) List() []*Car {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*Car, 0, len(r.cars))
	for _, c := range r.cars {
		copy := *c
		out = append(out, &copy)
	}
	return out
}

func (r *Registry) Claim(userID string) (*Car, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, c := range r.cars {
		if c.State == CarFree {
			c.State = CarReserved
			c.AssignedTo = userID
			c.UpdatedAt = time.Now()
			cp := *c
			return &cp, nil
		}
	}
	return nil, errors.New("no free cars")
}

func (r *Registry) Release(userID, carID string) (*Car, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.cars[carID]
	if !ok {
		return nil, errors.New("car not found")
	}
	if c.AssignedTo != userID {
		return nil, errors.New("car not owned by user")
	}
	c.State = CarFree
	c.AssignedTo = ""
	c.UpdatedAt = time.Now()
	cp := *c
	return &cp, nil
}

func (r *Registry) MarkBusy(carID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.cars[carID]
	if !ok {
		return errors.New("car not found")
	}
	c.State = CarBusy
	c.UpdatedAt = time.Now()
	return nil
}

func (r *Registry) MarkFree(carID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.cars[carID]
	if !ok {
		return errors.New("car not found")
	}
	c.State = CarFree
	c.AssignedTo = ""
	c.UpdatedAt = time.Now()
	return nil
}
