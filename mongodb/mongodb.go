package mongodb

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mcp struct {
	mutex *sync.Mutex
	count uint
	last  uint
	murl  string
	con   []*mongo.Client
}

// createConnection internal func for raw MongoDB connection
func createConnection(murl string) (*mongo.Client, error) {
	client, err := mongo.
		Connect(context.Background(), options.
			Client().
			ApplyURI(murl))
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, err
}

// Create MongoDB connections pool with 'n' connections
func Create(murl string, n int) (mcp, error) {
	var wg sync.WaitGroup

	cp := mcp{murl: murl, count: 0, last: 0}
	cp.mutex = &sync.Mutex{}

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			nc, err := createConnection(cp.murl)
			if err == nil {
				cp.mutex.Lock()
				cp.con = append(cp.con, nc)
				cp.count++
				cp.mutex.Unlock()
			}
		}()
	}
	wg.Wait()

	return cp, nil
}

// Length return size of MongoDB connections pool
func (cp *mcp) Length() uint {
	return cp.count
}

// Close single connection from pool
func (cp *mcp) Close(n uint) error {
	if n >= 0 && n < cp.count {
		err := cp.con[n].Ping(context.Background(), nil)
		if err != nil {
			return err
		}
		return cp.con[n].Disconnect(context.Background())
	}
	return errors.New("try to close not exists client")
}

// Destroy (close) all conections in pool
func (cp *mcp) Destroy() {
	var i uint
	for i = 0; i < cp.count; i++ {
		cp.Close(i)
	}
}

// Get single connection from connections pool
func (cp *mcp) Get(n uint) *mongo.Client {
	if n < 0 || n > cp.count {
		return nil
	}

	err := cp.con[n].Ping(context.Background(), nil)
	if err != nil {
		c, err := createConnection(cp.murl)
		if err != nil {
			return nil
		}
		cp.con[n] = c
	}

	return cp.con[n]
}

// GetRandom single connection from connections pool
func (cp *mcp) GetRandom() (*mongo.Client, uint32) {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(int(cp.count))
	return cp.con[i], uint32(i)
}

// GetRoundRobin gets single connection from connections pool by round robin
func (cp *mcp) GetRoundRobin() (*mongo.Client, uint32) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	cp.last++
	if cp.last == cp.count {
		cp.last = 0
	}
	n := cp.last

	return cp.con[n], uint32(n)
}
