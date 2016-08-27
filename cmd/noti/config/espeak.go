// +build !darwin
// +build !windows

package config

import "github.com/variadico/noti/espeak"

// MergeESpeak combines espeak notification structs. Each subsequent struct
// field overwrites the previous one, if it was set.
func MergeESpeak(ns ...*espeak.Notification) *espeak.Notification {
	out := new(espeak.Notification)

	for _, n := range ns {
		if n == nil {
			continue
		}

		if n.Text != "" {
			out.Text = n.Text
		}
		if n.VoiceName != "" {
			out.VoiceName = n.VoiceName
		}
		if n.Rate != 0 {
			out.Rate = n.Rate
		}
	}

	return out
}
