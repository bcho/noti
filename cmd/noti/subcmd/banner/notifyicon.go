// +build windows

package banner

import (
	"flag"

	"github.com/variadico/noti/cmd/noti/config"
	"github.com/variadico/noti/cmd/noti/run"
	"github.com/variadico/noti/cmd/noti/subcmd"
	"github.com/variadico/noti/notifyicon"
	"github.com/variadico/vbs"
)

var cmdDefault = &notifyicon.Notification{
	BalloonTipTitle: "{{.Cmd}}",
	BalloonTipText:  "Done!",
	BalloonTipIcon:  notifyicon.BalloonTipIconInfo,
	Icon:            notifyicon.DefaultIcon,
	Duration:        1000,
}

type template struct {
	*notifyicon.Notification
}

func (t template) TmplFields() []*string {
	if t.Notification == nil {
		return nil
	}

	return []*string{
		&t.BalloonTipTitle,
		&t.BalloonTipText,
		&t.Text,
	}
}

type Command struct {
	n     *notifyicon.Notification
	flags *flag.FlagSet
}

func (c *Command) Notify(stats run.Stats) error {
	conf, err := config.File()
	if err != nil {
		vbs.Println(err)
	} else {
		vbs.Println("Found config file")
	}

	fromFlags := new(notifyicon.Notification)
	if config.WasSet(c.flags, "title") || config.WasSet(c.flags, "t") {
		fromFlags.BalloonTipTitle = c.n.BalloonTipTitle
	}
	if config.WasSet(c.flags, "message") || config.WasSet(c.flags, "m") {
		fromFlags.BalloonTipText = c.n.BalloonTipText
	}
	if config.WasSet(c.flags, "balloon-icon") {
		fromFlags.BalloonTipIcon = c.n.BalloonTipIcon
	}
	if config.WasSet(c.flags, "app-icon") {
		fromFlags.Icon = c.n.Icon
	}
	if config.WasSet(c.flags, "hover-text") {
		fromFlags.Text = c.n.Text
	}
	if config.WasSet(c.flags, "duration") {
		fromFlags.Duration = c.n.Duration
	}

	vbs.Println("Evaluating")
	vbs.Printf("Default: %+v\n", cmdDefault)
	vbs.Printf("Config: %+v\n", conf.Banner)
	vbs.Printf("Flags: %+v\n", fromFlags)

	config.EvalTmplFields(template{cmdDefault}, stats)
	config.EvalTmplFields(template{conf.Banner}, stats)
	config.EvalTmplFields(template{fromFlags}, stats)

	vbs.Println("Merging")
	n := config.MergeNotifyIcon(cmdDefault, conf.Banner, fromFlags)
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
		n:     new(notifyicon.Notification),
		flags: flag.NewFlagSet("banner", flag.ContinueOnError),
	}

	cmd.flags.StringVar(&cmd.n.BalloonTipTitle, "title", cmdDefault.BalloonTipTitle, "Title")
	cmd.flags.StringVar(&cmd.n.BalloonTipTitle, "t", cmdDefault.BalloonTipTitle, "Title")
	cmd.flags.StringVar(&cmd.n.BalloonTipText, "message", cmdDefault.BalloonTipText, "Message")
	cmd.flags.StringVar(&cmd.n.BalloonTipText, "m", cmdDefault.BalloonTipText, "Message")
	cmd.flags.StringVar(&cmd.n.BalloonTipIcon, "balloon-icon", cmdDefault.BalloonTipIcon, "BalloonTipIcon")
	cmd.flags.StringVar(&cmd.n.Icon, "app-icon", cmdDefault.Icon, "App icon")
	cmd.flags.StringVar(&cmd.n.Text, "hover-text", cmdDefault.Text, "Hover text")
	cmd.flags.IntVar(&cmd.n.Duration, "duration", cmdDefault.Duration, "Duration")

	cmd.flags.Parse(args)
	return cmd
}
