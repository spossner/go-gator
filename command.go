package main

import (
	"context"
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
	name := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("error fetching user %v: %w", name, err)
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("error logging in: %w", err)
	}
	fmt.Printf("User %v logged in successfully\n", s.cfg.CurrentUserName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing username")
	}
	user, err := s.db.CreateUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("error logging in %v after registration: %w", user.Name, err)
	}

	fmt.Printf("User %s generated and successfully logged in\n", user.Name)
	fmt.Println(user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.Reset(context.Background()); err != nil {
		return fmt.Errorf("error wiping users table: %w", err)
	}
	return nil
}
