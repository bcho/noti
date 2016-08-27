// +build darwin

package banner

import (
	"flag"

	"github.com/variadico/noti/cmd/noti/config"
	"github.com/variadico/noti/cmd/noti/run"
	"github.com/variadico/noti/cmd/noti/subcmd"
	"github.com/variadico/noti/nsuser"
	"github.com/variadico/vbs"
)

var cmdDefault = &nsuser.Notification{
	Title:           "{{.Cmd}}",
	InformativeText: "Done!",
	SoundName:       "Ping",
}

type template struct {
	*nsuser.Notification
}

func (t template) TmplFields() []*string {
	if t.Notification == nil {
		return nil
	}

	return []*string{
		&t.Title,
		&t.Subtitle,
		&t.InformativeText,
	}
}

type Command struct {
	n     *nsuser.Notification
	flags *flag.FlagSet
}

func (c *Command) Notify(stats run.Stats) error {
	conf, err := config.File()
	if err != nil {
		vbs.Println(err)
	} else {
		vbs.Println("Found config file")
	}

	fromFlags := new(nsuser.Notification)
	if config.WasSet(c.flags, "title") || config.WasSet(c.flags, "t") {
		fromFlags.Title = c.n.Title
	}
	if config.WasSet(c.flags, "subtitle") {
		fromFlags.Subtitle = c.n.Subtitle
	}
	if config.WasSet(c.flags, "message") || config.WasSet(c.flags, "m") {
		fromFlags.InformativeText = c.n.InformativeText
	}
	if config.WasSet(c.flags, "icon") {
		fromFlags.ContentImage = c.n.ContentImage
	}
	if config.WasSet(c.flags, "sound") {
		fromFlags.SoundName = c.n.SoundName
	}

	vbs.Println("Evaluating")
	vbs.Printf("Default: %+v\n", cmdDefault)
	vbs.Printf("Config: %+v\n", conf.Banner)
	vbs.Printf("Flags: %+v\n", fromFlags)

	config.EvalTmplFields(template{cmdDefault}, stats)
	config.EvalTmplFields(template{conf.Banner}, stats)
	config.EvalTmplFields(template{fromFlags}, stats)

	vbs.Println("Merging")
	n := config.MergeNSUser(cmdDefault, conf.Banner, fromFlags)
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
		n:     new(nsuser.Notification),
		flags: flag.NewFlagSet("banner", flag.ContinueOnError),
	}

	cmd.flags.StringVar(&cmd.n.Title, "title", cmdDefault.Title, "Title")
	cmd.flags.StringVar(&cmd.n.Title, "t", cmdDefault.Title, "Title")
	cmd.flags.StringVar(&cmd.n.Subtitle, "subtitle", cmdDefault.Subtitle, "Subtitle")
	cmd.flags.StringVar(&cmd.n.InformativeText, "message", cmdDefault.InformativeText, "Message")
	cmd.flags.StringVar(&cmd.n.InformativeText, "m", cmdDefault.InformativeText, "Message")
	cmd.flags.StringVar(&cmd.n.ContentImage, "icon", cmdDefault.ContentImage, "Icon")
	cmd.flags.StringVar(&cmd.n.SoundName, "sound", cmdDefault.SoundName, "Sound")

	cmd.flags.Parse(args)
	return cmd
}
