package main

import (
	"flag"
	"os"
	"log"
)
var commands = make(map[string]*Command)

func main() {
	currentpath, _ := os.Getwd()

	flag.Usage = func(){
		log.Println("Usage")
	}

	flag.Parse()
	log.SetFlags(0)
	args := flag.Args()

	if len(args) < 1 {
		log.Println("Usage")
		os.Exit(2)
		return
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}

	if cmd, ok := commands[args[0]]; ok{
		if cmd.IsRunning{
			os.Exit(2)
		}

		cmd.Flag.Usage = func() { cmd.Usage() }
		if cmd.Interceptor != nil {
			cmd.Interceptor(args)
		}
		os.Exit(callRun(cmd, args))
		return
	}
}


func help(name []string) string{
	return "Help " + name[0];
}