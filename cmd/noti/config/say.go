// +build darwin

package config

import "github.com/variadico/noti/say"

// MergeSay combines say notification structs. Each subsequent struct
// field overwrites the previous one, if it was set.
func MergeSay(ns ...*say.Notification) *say.Notification {
	out := new(say.Notification)

	for _, n := range ns {
		if n == nil {
			continue
		}

		if n.Text != "" {
			out.Text = n.Text
		}
		if n.Voice != "" {
			out.Voice = n.Voice
		}
		if n.Rate != 0 {
			out.Rate = n.Rate
		}
	}

	return out
}
