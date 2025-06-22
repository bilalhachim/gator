package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/bilalhachim/gator/internal/config"
	"github.com/bilalhachim/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error %v", err)
	}
	dbQueries := database.New(db)

	programState := &state{
		cfg: &cfg,
		db:  dbQueries,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerusers)
	cmds.register("agg", handleragg)
	cmds.register("addfeed", middlewareLoggedIn(handlerFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow",middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse",middlewareLoggedIn(handlerBrowse))
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		fmt.Println("Exit code 1")
		os.Exit(1)
	}
	fmt.Println("Exit code 0")
	os.Exit(0)
}
