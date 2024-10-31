package main

import (
	"context"
	"fmt"
)

func withAuthentication(handler authenticatedHandler) handler {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUserByName(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("error fetching current user %s: %w", s.cfg.CurrentUserName, err)
		}
		return handler(s, cmd, user)
	}
}
