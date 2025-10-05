package list

import (
	"fmt"
	"slices"
	"sync"
	"time"

	"buffersnow.com/spiritonline/internal/iwmaster/protocol"
)

type ServerState int

const (
	ServerState_Idle ServerState = iota
	ServerState_Refreshing
	ServerState_Looking
)

type ServerList struct {
	lists map[string][]*Server
	mu    sync.RWMutex
}

type Server struct {
	State      ServerState
	Protocol   int
	Challenge  string
	Name       string
	Players    int
	MaxPlayers int
	Context    *protocol.IWContext
	LastPing   time.Time
}

func New() (*ServerList, error) {
	return &ServerList{
		lists: make(map[string][]*Server),
	}, nil
}

func (s *ServerList) Add(game string, server *Server) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lists[game] = append(s.lists[game], server)
}

func (s *ServerList) Remove(game string, server *Server) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lists[game] = slices.DeleteFunc(s.lists[game], func(se *Server) bool {
		return se.Challenge == server.Challenge
	})
}

func (s *ServerList) Access(game, challenge string, fn func(s *Server) error) error {
	s.mu.RLock()
	servers, exists := s.lists[game]
	s.mu.RUnlock()
	if !exists {
		return fmt.Errorf("list: game is not registered")
	}

	idx := slices.IndexFunc(servers, func(se *Server) bool {
		return se.Challenge == challenge
	})
	if idx == -1 {
		return fmt.Errorf("list: index of slice was invalid")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	return fn(servers[idx])
}

func (s *ServerList) Lock(fn func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn()
}

func (s *ServerList) IterateRead(fn func(game string, s *Server)) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for game, servers := range s.lists {
		for _, server := range servers {
			fn(game, server)
		}
	}
}

func (s *ServerList) IterateMutable(fn func(game string, s *Server)) {
	s.mu.RLock()
	copy := make(map[string][]*Server, len(s.lists))
	for k, v := range s.lists {
		copy[k] = append([]*Server(nil), v...)
	}
	s.mu.RUnlock()

	for game, servers := range copy {
		for _, server := range servers {
			fn(game, server)
		}
	}
}
