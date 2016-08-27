// +build darwin

package config

import "github.com/variadico/noti/nsuser"

// MergeNSUser combines banner notification structs. Each subsequent struct
// field overwrites the previous one, if it was set.
func MergeNSUser(ns ...*nsuser.Notification) *nsuser.Notification {
	out := new(nsuser.Notification)

	for _, n := range ns {
		if n == nil {
			continue
		}

		if n.Title != "" {
			out.Title = n.Title
		}
		if n.Subtitle != "" {
			out.Subtitle = n.Subtitle
		}
		if n.InformativeText != "" {
			out.InformativeText = n.InformativeText
		}
		if n.ContentImage != "" {
			out.ContentImage = n.ContentImage
		}
		if n.SoundName != "" {
			out.SoundName = n.SoundName
		}
	}

	return out
}
