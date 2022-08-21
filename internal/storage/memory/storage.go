package memorystorage

import (
	"context"
	"net"
	"sync"

	"github.com/sn5ake6/otus-final-project/internal/storage"
)

type Storage struct {
	mu        sync.RWMutex
	blacklist map[string]string
	whitelist map[string]string
}

func New() *Storage {
	return &Storage{
		blacklist: make(map[string]string),
		whitelist: make(map[string]string),
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	return nil
}

func (s *Storage) AddToBlacklist(subnet string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.blacklist[subnet]; ok {
		return storage.ErrSubnetAlreadyExists
	}

	s.blacklist[subnet] = subnet

	return nil
}

func (s *Storage) DeleteFromBlacklist(subnet string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.blacklist[subnet]; !ok {
		return storage.ErrSubnetNotExists
	}

	delete(s.blacklist, subnet)

	return nil
}

func (s *Storage) FindIPInBlacklist(ip string) (bool, error) {
	ipForSearch := net.ParseIP(ip)

	for subnet := range s.blacklist {
		_, subnetNet, err := net.ParseCIDR(subnet)
		if err != nil {
			return false, err
		}

		if subnetNet.Contains(ipForSearch) {
			return true, nil
		}
	}

	return false, nil
}

func (s *Storage) AddToWhitelist(subnet string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.whitelist[subnet]; ok {
		return storage.ErrSubnetAlreadyExists
	}

	s.whitelist[subnet] = subnet

	return nil
}

func (s *Storage) DeleteFromWhitelist(subnet string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.whitelist[subnet]; !ok {
		return storage.ErrSubnetNotExists
	}

	delete(s.whitelist, subnet)

	return nil
}

func (s *Storage) FindIPInWhitelist(ip string) (bool, error) {
	ipForSearch := net.ParseIP(ip)

	for subnet := range s.whitelist {
		_, subnetNet, err := net.ParseCIDR(subnet)
		if err != nil {
			return false, err
		}

		if subnetNet.Contains(ipForSearch) {
			return true, nil
		}
	}

	return false, nil
}
