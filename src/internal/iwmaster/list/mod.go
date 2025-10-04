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
	mu    sync.Mutex
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
	srv := &ServerList{
		lists: map[string][]*Server{},
		mu:    sync.Mutex{},
	}
	return srv, nil
}

func (s *ServerList) Add(game string, server *Server) {
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
	servers, exists := s.lists[game]
	if !exists {
		return fmt.Errorf("list: game is not registered")
	}

	idx := slices.IndexFunc(servers, func(se *Server) bool {
		return se.Challenge == challenge
	})

	if idx == -1 {
		return fmt.Errorf("list: index of slice was invalid")
	}

	return fn(servers[idx])
}

func (s *ServerList) Iterate(fn func(game string, s *Server)) {
	for game, servers := range s.lists {
		for _, server := range servers {
			fn(game, server)
		}
	}
}

func (s *ServerList) Lock(fn func()) {
	s.mu.Lock()
	fn()
	s.mu.Unlock()
}
