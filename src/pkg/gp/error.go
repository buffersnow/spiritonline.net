package gp

import "fmt"

type GameSpyError struct {
	ErrorCode int
	Message   string
	IsFatal   bool
}

// Error implements the error interface.
func (g *GameSpyError) Error() string {
	if g.IsFatal {
		return fmt.Sprintf("GameSpy Fatal Error: %d - %s", g.ErrorCode, g.Message)
	}

	return fmt.Sprintf("GameSpy Error: %d - %s", g.ErrorCode, g.Message)
}
