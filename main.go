package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/spossner/gator/internal/config"
	"github.com/spossner/gator/internal/database"
	"log"
	"os"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("error reading config file", err)
	}
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatal("error opening database", err)
	}
	s := &state{
		db:  database.New(db),
		cfg: cfg,
	}
	cmds := commands{
		list: make(map[string]handler),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("feeds", handlerFeeds)
	cmds.register("addfeed", withAuthentication(handlerAddFeed))
	cmds.register("follow", withAuthentication(handlerFollow))
	cmds.register("unfollow", withAuthentication(handlerUnfollow))
	cmds.register("following", withAuthentication(handlerFollowing))
	cmds.register("browse", withAuthentication(handlerBrowse))

	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("Usage...")
	}
	cmd := command{
		name: args[0],
		args: args[1:],
	}
	if err := cmds.run(s, cmd); err != nil {
		log.Fatalf("error executing %v: %v\n", cmd, err)
	}
}
