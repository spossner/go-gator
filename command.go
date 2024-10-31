package main

import (
	"errors"
	"fmt"
)

type command struct {
	name string
	args []string
}

type handler func(*state, command) error

type commands struct {
	list map[string]handler
}

func (c *commands) register(name string, fn func(*state, command) error) {
	c.list[name] = fn
}

func (c *commands) run(s *state, cmd command) error {
	fn, ok := c.list[cmd.name]
	if !ok {
		return fmt.Errorf("unknwon command %v", cmd.name)
	}
	return fn(s, cmd)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing username")
	}
	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("error logging in: %w", err)
	}
	fmt.Printf("User %v logged in successfully\n", s.cfg.CurrentUserName)
	return nil
}
