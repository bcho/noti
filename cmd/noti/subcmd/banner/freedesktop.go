// +build !darwin
// +build !windows

package banner

import (
	"flag"

	"github.com/variadico/noti/cmd/noti/config"
	"github.com/variadico/noti/cmd/noti/run"
	"github.com/variadico/noti/cmd/noti/subcmd"
	"github.com/variadico/noti/freedesktop"
	"github.com/variadico/vbs"
)

var cmdDefault = &freedesktop.Notification{
	Summary:       "{{.Cmd}}",
	Body:          "Done!",
	ExpireTimeout: 500,
}

type template struct {
	*freedesktop.Notification
}

func (t template) TmplFields() []*string {
	if t.Notification == nil {
		return nil
	}

	return []*string{
		&t.Body,
		&t.Summary,
	}
}

type Command struct {
	n     *freedesktop.Notification
	flags *flag.FlagSet
}

func (c *Command) Notify(stats run.Stats) error {
	conf, err := config.File()
	if err != nil {
		vbs.Println(err)
	} else {
		vbs.Println("Found config file")
	}

	fromFlags := new(freedesktop.Notification)
	if config.WasSet(c.flags, "title") || config.WasSet(c.flags, "t") {
		fromFlags.Summary = c.n.Summary
	}
	if config.WasSet(c.flags, "message") || config.WasSet(c.flags, "m") {
		fromFlags.Body = c.n.Body
	}
	if config.WasSet(c.flags, "app-name") {
		fromFlags.AppName = c.n.AppName
	}
	if config.WasSet(c.flags, "replaces-id") {
		fromFlags.ReplacesID = c.n.ReplacesID
	}
	if config.WasSet(c.flags, "icon") {
		fromFlags.AppIcon = c.n.AppIcon
	}
	if config.WasSet(c.flags, "timeout") {
		fromFlags.ExpireTimeout = c.n.ExpireTimeout
	}

	vbs.Println("Evaluating")
	vbs.Printf("Default: %+v\n", cmdDefault)
	vbs.Printf("Config: %+v\n", conf.Banner)
	vbs.Printf("Flags: %+v\n", fromFlags)

	config.EvalTmplFields(template{cmdDefault}, stats)
	config.EvalTmplFields(template{conf.Banner}, stats)
	config.EvalTmplFields(template{fromFlags}, stats)

	vbs.Println("Merging")
	n := config.MergeFreedesktop(cmdDefault, conf.Banner, fromFlags)
	vbs.Printf("Merge result: %+v\n", n)

	vbs.Println("Sending notification")
	err = n.Send()
	vbs.Println("Sent notification")
	return err
}

func (c *Command) Run() error {
	vbs.Println("Executing command")
	stats := run.Exec(c.flags.Args()...)
	vbs.Println("Executed command")
	vbs.Printf("Run stats: %+v\n", stats)

	return c.Notify(stats)
}

func NewCommand(args []string) subcmd.Cmd {
	cmd := &Command{
		n:     new(freedesktop.Notification),
		flags: flag.NewFlagSet("banner", flag.ContinueOnError),
	}

	cmd.flags.StringVar(&cmd.n.Summary, "title", cmdDefault.Summary, "Title")
	cmd.flags.StringVar(&cmd.n.Summary, "t", cmdDefault.Summary, "Title")
	cmd.flags.StringVar(&cmd.n.Body, "message", cmdDefault.Body, "Message")
	cmd.flags.StringVar(&cmd.n.Body, "m", cmdDefault.Body, "Message")
	cmd.flags.StringVar(&cmd.n.AppName, "app-name", cmdDefault.AppName, "AppName")
	cmd.flags.UintVar(&cmd.n.ReplacesID, "replaces-id", cmdDefault.ReplacesID, "ReplacesID")
	cmd.flags.StringVar(&cmd.n.AppIcon, "icon", cmdDefault.AppIcon, "AppIcon")
	cmd.flags.IntVar(&cmd.n.ExpireTimeout, "timeout", cmdDefault.ExpireTimeout, "ExpireTimeout")

	cmd.flags.Parse(args)
	return cmd
}
