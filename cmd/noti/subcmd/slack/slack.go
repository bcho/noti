package slack

import (
	"flag"
	"net/http"
	"time"

	"github.com/variadico/noti/cmd/noti/config"
	"github.com/variadico/noti/cmd/noti/run"
	"github.com/variadico/noti/cmd/noti/subcmd"
	"github.com/variadico/noti/slack"
	"github.com/variadico/vbs"
)

var cmdDefault = &slack.Notification{
	Text:        "{{.Cmd}} done!",
	Parse:       slack.ParseFull,
	LinkNames:   slack.LinkNamesOn,
	UnfurlLinks: true,
	UnfurlMedia: true,
	Username:    "Noti",
}

type template struct {
	*slack.Notification
}

func (t template) TmplFields() []*string {
	if t.Notification == nil {
		return nil
	}

	return []*string{
		&t.Text,
		&t.Username,
	}
}

type Command struct {
	n     *slack.Notification
	flags *flag.FlagSet
}

func (c *Command) Notify(stats run.Stats) error {
	conf, err := config.File()
	if err != nil {
		vbs.Println(err)
	} else {
		vbs.Println("Found config file")
	}

	fromFlags := new(slack.Notification)
	if config.WasSet(c.flags, "message") || config.WasSet(c.flags, "m") {
		fromFlags.Text = c.n.Text
	}
	if config.WasSet(c.flags, "token") {
		fromFlags.Token = c.n.Token
	}
	if config.WasSet(c.flags, "channel") {
		fromFlags.Channel = c.n.Channel
	}
	if config.WasSet(c.flags, "parse") {
		fromFlags.Parse = c.n.Parse
	}
	if config.WasSet(c.flags, "link-names") {
		fromFlags.LinkNames = c.n.LinkNames
	}
	// Attachments
	if config.WasSet(c.flags, "unfurl-links") {
		fromFlags.UnfurlLinks = c.n.UnfurlLinks
	}
	if config.WasSet(c.flags, "unfurl-media") {
		fromFlags.UnfurlMedia = c.n.UnfurlMedia
	}
	if config.WasSet(c.flags, "username") {
		fromFlags.Username = c.n.Username
	}
	if config.WasSet(c.flags, "as-user") {
		fromFlags.AsUser = c.n.AsUser
	}
	if config.WasSet(c.flags, "icon-url") {
		fromFlags.IconURL = c.n.IconURL
	}
	if config.WasSet(c.flags, "icon-emoji") {
		fromFlags.IconEmoji = c.n.IconEmoji
	}

	vbs.Println("Evaluating")
	vbs.Printf("Default: %+v\n", cmdDefault)
	vbs.Printf("Config: %+v\n", conf.Banner)
	vbs.Printf("Flags: %+v\n", fromFlags)

	config.EvalTmplFields(template{cmdDefault}, stats)
	config.EvalTmplFields(template{conf.Slack}, stats)
	config.EvalTmplFields(template{fromFlags}, stats)

	vbs.Println("Merging")
	n := config.MergeSlack(cmdDefault, conf.Slack, fromFlags)
	vbs.Printf("Merge result: %+v\n", n)

	n.SetClient(&http.Client{Timeout: 5 * time.Second})

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
		n:     new(slack.Notification),
		flags: flag.NewFlagSet("slack", flag.ContinueOnError),
	}

	cmd.flags.StringVar(&cmd.n.Token, "token", cmdDefault.Token, "Token (Required)")
	cmd.flags.StringVar(&cmd.n.Channel, "channel", cmdDefault.Channel, "Channel (Required)")
	cmd.flags.StringVar(&cmd.n.Text, "message", cmdDefault.Text, "Message")
	cmd.flags.StringVar(&cmd.n.Text, "m", cmdDefault.Text, "Message")
	cmd.flags.StringVar(&cmd.n.Parse, "parse", cmdDefault.Parse, "Parse")
	cmd.flags.IntVar(&cmd.n.LinkNames, "link-names", cmdDefault.LinkNames, "LinkNames")
	// Attachments -- should go here.
	cmd.flags.BoolVar(&cmd.n.UnfurlLinks, "unfurl-links", cmdDefault.UnfurlLinks, "UnfurlLinks")
	cmd.flags.BoolVar(&cmd.n.UnfurlMedia, "unfurl-media", cmdDefault.UnfurlMedia, "UnfurlMedia")
	cmd.flags.StringVar(&cmd.n.Username, "username", cmdDefault.Username, "Username")
	cmd.flags.BoolVar(&cmd.n.AsUser, "as-user", cmdDefault.AsUser, "AsUser")
	cmd.flags.StringVar(&cmd.n.IconURL, "icon-url", cmdDefault.IconURL, "Username")
	cmd.flags.StringVar(&cmd.n.IconEmoji, "icon-emoji", cmdDefault.IconEmoji, "IconEmoji")

	cmd.flags.Parse(args)
	return cmd
}
