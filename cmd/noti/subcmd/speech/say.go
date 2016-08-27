// +build darwin

package speech

import (
	"flag"

	"github.com/variadico/noti/cmd/noti/config"
	"github.com/variadico/noti/cmd/noti/run"
	"github.com/variadico/noti/cmd/noti/subcmd"
	"github.com/variadico/noti/say"
	"github.com/variadico/vbs"
)

var cmdDefault = &say.Notification{
	Voice: "Alex",
	Text:  "{{.Cmd}} done!",
	Rate:  200,
}

type template struct {
	*say.Notification
}

func (t template) TmplFields() []*string {
	if t.Notification == nil {
		return nil
	}

	return []*string{
		&t.Text,
	}
}

type Command struct {
	n     *say.Notification
	flags *flag.FlagSet
}

func (c *Command) Notify(stats run.Stats) error {
	conf, err := config.File()
	if err != nil {
		vbs.Println(err)
	} else {
		vbs.Println("Found config file")
	}

	fromFlags := new(say.Notification)
	if config.WasSet(c.flags, "voice") {
		fromFlags.Voice = c.n.Voice
	}
	if config.WasSet(c.flags, "message") || config.WasSet(c.flags, "m") {
		fromFlags.Text = c.n.Text
	}
	if config.WasSet(c.flags, "rate") {
		fromFlags.Rate = c.n.Rate
	}

	vbs.Println("Evaluating")
	vbs.Printf("Default: %+v\n", cmdDefault)
	vbs.Printf("Config: %+v\n", conf.Speech)
	vbs.Printf("Flags: %+v\n", fromFlags)

	config.EvalTmplFields(template{cmdDefault}, stats)
	config.EvalTmplFields(template{conf.Speech}, stats)
	config.EvalTmplFields(template{fromFlags}, stats)

	vbs.Println("Merging")
	n := config.MergeSay(cmdDefault, conf.Speech, fromFlags)
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
		n:     new(say.Notification),
		flags: flag.NewFlagSet("speech", flag.ContinueOnError),
	}

	cmd.flags.StringVar(&cmd.n.Voice, "voice", cmdDefault.Voice, "Voice")
	cmd.flags.StringVar(&cmd.n.Text, "message", cmdDefault.Text, "Message")
	cmd.flags.StringVar(&cmd.n.Text, "m", cmdDefault.Text, "Message")
	cmd.flags.IntVar(&cmd.n.Rate, "rate", cmdDefault.Rate, "rate")

	cmd.flags.Parse(args)
	return cmd
}
