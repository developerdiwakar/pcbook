package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/developerdiwakar/pcbook/pb"
	"github.com/jinzhu/copier"
)

// ErrEmailAlreadyExists is retunred when a record with the same ID already exists in the store
var ErrAlreadyExists = errors.New("record already exists")

type LaptopStore interface {
	// Save saves the laptop to the store
	Save(laptop *pb.Laptop) error
	// Find finds the laptop to the store
	Find(id string) (*pb.Laptop, error)
	// Search searches for laptops with filter, returns one by one via the found function
	Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error
}

// InMemoryLaptopStore stores laptop in memory
type InMemoryLaptopStore struct {
	mutex sync.RWMutex
	data  map[string]*pb.Laptop
}

// NewInMemoryLaptopStore returns a new InMemoryLaptopStore
func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*pb.Laptop),
	}
}

// Save saves the laptop to the store
func (store *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}

	// Deep copy
	other, err := deepCopy(laptop)
	if err != nil {
		return fmt.Errorf("cannot copy laptop data: %w", err)
	}

	store.data[other.Id] = other
	return nil
}

// Find finds the laptop to the store
func (store *InMemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	laptop := store.data[id]

	if laptop == nil {
		return nil, nil
	}

	// Deep copy
	return deepCopy(laptop)
}

// Search searches for laptops with filter, returns one by one via the found function
func (store *InMemoryLaptopStore) Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, laptop := range store.data {
		// simulate heavy processing
		time.Sleep(time.Second)
		log.Print("checking laptop id: ", laptop.GetId())
		if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
			log.Println("request is canceled or exceed deadline")
			return errors.New("context is camceled")
		}
		if isQualified(filter, laptop) {
			other, err := deepCopy(laptop)
			if err != nil {
				return err
			}
			err = found(other)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isQualified(filter *pb.Filter, laptop *pb.Laptop) bool {
	if laptop.GetPriceUsd() > filter.MaxPriceUsd {
		return false
	}
	if laptop.GetCpu().NumberCores < filter.MinCpuCores {
		return false
	}
	if laptop.GetCpu().GetMinGhz() < filter.MinCpuGhz {
		return false
	}

	if toBit(laptop.GetRam()) < toBit(filter.GetMinRam()) {
		return false
	}

	return true
}

func toBit(memory *pb.Memory) uint64 {
	value := memory.GetValue()

	switch memory.GetUnit() {
	case pb.Memory_BIT:
		return value // 1 bit
	case pb.Memory_BYTE:
		return value << 3 // 8 = 2^3 bits
	case pb.Memory_KILOBYTE:
		return value << 13 // 1024*8 = 2^10 + 2^3 = 2^13 bits
	case pb.Memory_MEGABITE:
		return value << 23 // 1024 * 1024 * 8 = 2^10 + 2^10 + 2^3 bits
	case pb.Memory_GIGABYTE:
		return value << 33 // 1024 * 1024 * 1024 * 8 bits
	case pb.Memory_TERABYTE:
		return value << 43 // 1024 * 1024 * 1024 * 1024 * 8 bits
	default:
		return 0
	}
}

func deepCopy(laptop *pb.Laptop) (*pb.Laptop, error) {
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data: %w", err)
	}
	return other, nil
}
