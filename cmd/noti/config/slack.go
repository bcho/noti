package config

import "github.com/variadico/noti/slack"

// MergeSlack combines slack notification structs. Each subsequent struct
// field overwrites the previous one, if it was set.
func MergeSlack(ns ...*slack.Notification) *slack.Notification {
	out := new(slack.Notification)

	for _, n := range ns {
		if n == nil {
			continue
		}

		if n.Token != "" {
			out.Token = n.Token
		}
		if n.Channel != "" {
			out.Channel = n.Channel
		}
		if n.Text != "" {
			out.Text = n.Text
		}
		if n.Parse != "" {
			out.Parse = n.Parse
		}
		if n.LinkNames != 0 {
			out.LinkNames = n.LinkNames
		}
		if n.Attachments != nil {
			out.Attachments = n.Attachments
		}
		if n.UnfurlLinks != false {
			out.UnfurlLinks = n.UnfurlLinks
		}
		if n.UnfurlMedia != false {
			out.UnfurlMedia = n.UnfurlMedia
		}
		if n.Username != "" {
			out.Username = n.Username
		}
		if n.AsUser != false {
			out.AsUser = n.AsUser
		}
		if n.IconURL != "" {
			out.IconURL = n.IconURL
		}
		if n.IconEmoji != "" {
			out.IconEmoji = n.IconEmoji
		}

		if n.Client() != nil {
			out.SetClient(n.Client())
		}
	}

	return out
}
