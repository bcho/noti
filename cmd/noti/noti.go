package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/variadico/noti/cmd/noti/config"
	"github.com/variadico/noti/cmd/noti/run"
	"github.com/variadico/vbs"
)

type command struct {
	Version bool

	flags *flag.FlagSet
}

func (c *command) Run() error {
	vbs.Println("Executing command")
	stats := run.Exec(c.flags.Args()...)
	vbs.Println("Executed command")
	vbs.Printf("Run stats: %+v\n", stats)

	conf, err := config.File()
	if err != nil {
		vbs.Println(err)
	} else {
		vbs.Println("Found config file")
	}

	if len(conf.DefaultSet) == 0 {
		conf.DefaultSet = append(conf.DefaultSet, "banner")
	}

	for _, sub := range conf.DefaultSet {
		var err error

		if fn, found := notiCmds[sub]; found {
			subCmd := fn(c.flags.Args())
			err = subCmd.Notify(stats)
		} else {
			err = fmt.Errorf("unknown notification type: %s", sub)
		}

		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func newCommand(args []string) *command {
	cmd := new(command)

	cmd.flags = flag.NewFlagSet("noti", flag.ContinueOnError)
	cmd.flags.BoolVar(&vbs.Verbose, "verbose", false, "Enable verbose mode")
	cmd.flags.BoolVar(&cmd.Version, "version", false, "Print version and exit")

	cmd.flags.Parse(args)
	return cmd

}
