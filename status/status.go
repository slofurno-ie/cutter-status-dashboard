package status

import (
	"time"

	"github.com/IdeaEvolver/cutter-pkg/clog"
	"github.com/IdeaEvolver/cutter-pkg/cuterr"
	"github.com/garyburd/redigo/redis"
)

type StatusStore struct {
	pool *redis.Pool
}

type AllStatuses struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

type Status struct {
	Status string `json:"status"`
}

var services = []string{"platform", "platform-ui", "fulfillment", "crm", "study", "study-ui"}

func InitRedis(url string) *redis.Pool {
	// Establish a connection to the Redis server listening on port
	// 6379 of the local machine. 6379 is the default port, so unless
	// you've already changed the Redis configuration file this should
	// work.
	// conn, err := redis.Dial("tcp", cfg.REDIS_URL)
	// if err != nil {
	// 	clog.Fatalf("error dialing redis", err)
	// }

	pool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", url)
		},
	}

	return pool
}

func New(pool *redis.Pool) *StatusStore {
	return &StatusStore{
		pool: pool,
	}
}

func (s *StatusStore) InitStatuses() {
	// Use the connection pool's Get() method to fetch a single Redis
	// connection from the pool.
	conn := s.pool.Get()

	// Importantly, use defer and the connection's Close() method to
	// ensure that the connection is always returned to the pool before
	// FindAlbum() exits.
	defer conn.Close()

	// Send our command across the connection. The first parameter to
	// Do() is always the name of the Redis command (in this example
	// HMSET), optionally followed by any necessary arguments (in this
	// example the key, followed by the various hash fields and values).

	for _, s := range services {
		_, err := conn.Do("HSET", "service:"+s, "status", "200")
		if err != nil {
			clog.Fatalf("error updating redis", err)
		}
	}
}

func (s *StatusStore) GetAllStatuses() ([]*AllStatuses, error) {
	statusArr := []*AllStatuses{}

	conn := s.pool.Get()
	defer conn.Close()

	for _, s := range services {
		status := &AllStatuses{}
		code, err := redis.String(conn.Do("HGET", "service:"+s, "status"))
		if err != nil {
			clog.Fatalf("error getting artist", err)
			return nil, cuterr.New(cuterr.Internal, "could not get status", err)
		}
		status.Service = s
		status.Status = code
		statusArr = append(statusArr, status)
	}

	return statusArr, nil
}

func (s *StatusStore) UpdateStatus(service, status string) error {
	conn := s.pool.Get()
	defer conn.Close()

	_, err := conn.Do("HSET", "service:"+service, "status", status)
	if err != nil {
		return cuterr.New(cuterr.Internal, "could not update status", err)
	}

	return nil
}
