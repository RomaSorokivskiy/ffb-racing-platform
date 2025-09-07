package rooms

import (
	"encoding/json"
	"errors"
	"log"
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
	TTL        int64     `json:"ttl,omitempty"` // seconds left (for RESERVED)
}

type Event struct {
	Type string `json:"type"` // "snapshot" | "update"
	Data any    `json:"data"`
}

type claimInfo struct {
	user   string
	expire time.Time
}

type Registry struct {
	mu       sync.RWMutex
	cars     map[string]*Car
	claims   map[string]*claimInfo // carID -> claim
	subs     map[chan Event]struct{}
	stopGC   chan struct{}
	lifetime time.Duration
}

func NewRegistry(n int) *Registry {
	r := &Registry{
		cars:     make(map[string]*Car, n),
		claims:   make(map[string]*claimInfo),
		subs:     make(map[chan Event]struct{}),
		stopGC:   make(chan struct{}),
		lifetime: 2 * time.Minute, // default claim TTL
	}
	for i := 1; i <= n; i++ {
		id := "car-" + strconv.Itoa(i)
		r.cars[id] = &Car{
			ID:        id,
			State:     CarFree,
			UpdatedAt: time.Now(),
		}
	}
	go r.gcLoop()
	return r
}

func (r *Registry) Close() { close(r.stopGC) }

func (r *Registry) gcLoop() {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	for {
		select {
		case <-r.stopGC:
			return
		case <-tick.C:
			now := time.Now()
			var toBroadcast []*Car

			r.mu.Lock()
			for id, ci := range r.claims {
				if now.After(ci.expire) {
					if c := r.cars[id]; c != nil && c.State == CarReserved && c.AssignedTo == ci.user {
						c.State = CarFree
						c.AssignedTo = ""
						c.UpdatedAt = now
						c.TTL = 0
						delete(r.claims, id)
						cp := *c
						toBroadcast = append(toBroadcast, &cp)
					} else {
						delete(r.claims, id)
					}
				} else {
					if c := r.cars[id]; c != nil && c.State == CarReserved {
						c.TTL = int64(ci.expire.Sub(now).Seconds())
					}
				}
			}
			r.mu.Unlock()

			for _, c := range toBroadcast {
				r.broadcast(Event{Type: "update", Data: c})
			}
		}
	}
}

func (r *Registry) List() []*Car {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*Car, 0, len(r.cars))
	for _, c := range r.cars {
		cp := *c
		if c.State != CarReserved {
			cp.TTL = 0
		}
		out = append(out, &cp)
	}
	return out
}

func (r *Registry) Claim(userID string, ttl time.Duration) (*Car, error) {
	if ttl <= 0 || ttl > 10*time.Minute {
		ttl = r.lifetime
	}
	var out *Car

	r.mu.Lock()
	for _, c := range r.cars {
		if c.State == CarFree {
			c.State = CarReserved
			c.AssignedTo = userID
			c.UpdatedAt = time.Now()
			r.claims[c.ID] = &claimInfo{user: userID, expire: c.UpdatedAt.Add(ttl)}
			c.TTL = int64(ttl.Seconds())
			cp := *c
			out = &cp
			break
		}
	}
	r.mu.Unlock()

	if out == nil {
		return nil, errors.New("no free cars")
	}
	r.broadcast(Event{Type: "update", Data: out})
	return out, nil
}

func (r *Registry) Release(userID, carID string) (*Car, error) {
	var out *Car

	r.mu.Lock()
	c, ok := r.cars[carID]
	if !ok {
		r.mu.Unlock()
		return nil, errors.New("car not found")
	}
	if c.AssignedTo != userID {
		r.mu.Unlock()
		return nil, errors.New("car not owned by user")
	}
	delete(r.claims, carID)
	c.State = CarFree
	c.AssignedTo = ""
	c.UpdatedAt = time.Now()
	c.TTL = 0
	cp := *c
	out = &cp
	r.mu.Unlock()

	r.broadcast(Event{Type: "update", Data: out})
	return out, nil
}

func (r *Registry) MarkBusy(carID string) error {
	var out *Car

	r.mu.Lock()
	c, ok := r.cars[carID]
	if !ok {
		r.mu.Unlock()
		return errors.New("car not found")
	}
	c.State = CarBusy
	c.UpdatedAt = time.Now()
	c.TTL = 0
	cp := *c
	out = &cp
	r.mu.Unlock()

	r.broadcast(Event{Type: "update", Data: out})
	return nil
}

func (r *Registry) MarkFree(carID string) error {
	var out *Car

	r.mu.Lock()
	c, ok := r.cars[carID]
	if !ok {
		r.mu.Unlock()
		return errors.New("car not found")
	}
	delete(r.claims, carID)
	c.State = CarFree
	c.AssignedTo = ""
	c.UpdatedAt = time.Now()
	c.TTL = 0
	cp := *c
	out = &cp
	r.mu.Unlock()

	r.broadcast(Event{Type: "update", Data: out})
	return nil
}

// ---- SSE ----

func (r *Registry) Subscribe() chan Event {
	ch := make(chan Event, 16)

	// реєструємо підписника
	r.mu.Lock()
	r.subs[ch] = struct{}{}
	r.mu.Unlock()

	// відправляємо snapshot поза м'ютексом (без дедлоку)
	go func() {
		snap := r.List()
		ch <- Event{Type: "snapshot", Data: snap}
	}()

	return ch
}

func (r *Registry) Unsubscribe(ch chan Event) {
	r.mu.Lock()
	if _, ok := r.subs[ch]; ok {
		delete(r.subs, ch)
		close(ch)
	}
	r.mu.Unlock()
}

func (r *Registry) broadcast(ev Event) {
	// копіюємо список підписників під RLock
	r.mu.RLock()
	subs := make([]chan Event, 0, len(r.subs))
	for ch := range r.subs {
		subs = append(subs, ch)
	}
	r.mu.RUnlock()

	// шлемо поза м'ютексом
	for _, ch := range subs {
		select {
		case ch <- ev:
		default:
			log.Println("rooms: drop event to slow subscriber")
		}
	}
}

func MarshalEvent(ev Event) []byte {
	b, _ := json.Marshal(ev)
	return b
}
