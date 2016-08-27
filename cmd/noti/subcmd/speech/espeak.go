// +build !darwin
// +build !windows

package speech

import (
	"flag"

	"github.com/variadico/noti/cmd/noti/config"
	"github.com/variadico/noti/cmd/noti/run"
	"github.com/variadico/noti/cmd/noti/subcmd"
	"github.com/variadico/noti/espeak"
	"github.com/variadico/vbs"
)

var cmdDefault = &espeak.Notification{
	Text:      "{{.Cmd}} done!",
	VoiceName: "english-us",
}

type template struct {
	*espeak.Notification
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
	n     *espeak.Notification
	flags *flag.FlagSet
}

func (c *Command) Notify(stats run.Stats) error {
	conf, err := config.File()
	if err != nil {
		vbs.Println(err)
	} else {
		vbs.Println("Found config file")
	}

	fromFlags := new(espeak.Notification)
	if config.WasSet(c.flags, "message") || config.WasSet(c.flags, "m") {
		fromFlags.Text = c.n.Text
	}
	if config.WasSet(c.flags, "word-gap") {
		fromFlags.WordGap = c.n.WordGap
	}
	if config.WasSet(c.flags, "pitch") {
		fromFlags.PitchAdjustment = c.n.PitchAdjustment
	}
	if config.WasSet(c.flags, "rate") {
		fromFlags.Rate = c.n.Rate
	}
	if config.WasSet(c.flags, "voice-name") {
		fromFlags.VoiceName = c.n.VoiceName
	}

	vbs.Println("Evaluating")
	vbs.Printf("Default: %+v\n", cmdDefault)
	vbs.Printf("Config: %+v\n", conf.Speech)
	vbs.Printf("Flags: %+v\n", fromFlags)

	config.EvalTmplFields(template{cmdDefault}, stats)
	config.EvalTmplFields(template{conf.Speech}, stats)
	config.EvalTmplFields(template{fromFlags}, stats)

	vbs.Println("Merging")
	n := config.MergeESpeak(cmdDefault, conf.Speech, fromFlags)
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
		n:     new(espeak.Notification),
		flags: flag.NewFlagSet("speech", flag.ContinueOnError),
	}

	cmd.flags.StringVar(&cmd.n.VoiceName, "voice", cmdDefault.VoiceName, "Voice")
	cmd.flags.StringVar(&cmd.n.Text, "message", cmdDefault.Text, "Message")
	cmd.flags.StringVar(&cmd.n.Text, "m", cmdDefault.Text, "Message")
	cmd.flags.IntVar(&cmd.n.Rate, "rate", cmdDefault.Rate, "Rate")
	cmd.flags.IntVar(&cmd.n.PitchAdjustment, "pitch", cmdDefault.PitchAdjustment, "PitchAdjustment")
	cmd.flags.IntVar(&cmd.n.WordGap, "word-gap", cmdDefault.WordGap, "WordGap")

	cmd.flags.Parse(args)
	return cmd
}
