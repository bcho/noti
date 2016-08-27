package main

import (
	"log"
	"os"

	"github.com/variadico/noti/cmd/noti/subcmd"
	"github.com/variadico/noti/cmd/noti/subcmd/banner"
	"github.com/variadico/noti/cmd/noti/subcmd/slack"
	"github.com/variadico/noti/cmd/noti/subcmd/speech"
)

type commandFunc func([]string) subcmd.Cmd

var notiCmds = map[string]commandFunc{
	"banner": banner.NewCommand,
	"speech": speech.NewCommand,
	"slack":  slack.NewCommand,
}

func main() {
	log.SetFlags(0)

	notiCmd := newCommand(os.Args[1:])

	if len(notiCmd.flags.Args()) == 0 {
		// Use Notify instead.
		if err := notiCmd.Run(); err != nil {
			log.Fatalln(err)
		}
		return
	}

	var err error
	sub := notiCmd.flags.Args()[0]

	if fn, found := notiCmds[sub]; found {
		subCmd := fn(notiCmd.flags.Args()[1:])
		err = subCmd.Run()
	} else {
		err = notiCmd.Run()
	}
	if err != nil {
		log.Fatalln("Error:", err)
	}
}
