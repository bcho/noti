// +build !darwin
// +build !windows

package config

import "github.com/variadico/noti/freedesktop"

// MergeFreedesktop combines banner notification structs. Each subsequent struct
// field overwrites the previous one, if it was set.
func MergeFreedesktop(ns ...*freedesktop.Notification) *freedesktop.Notification {
	out := new(freedesktop.Notification)

	for _, n := range ns {
		if n == nil {
			continue
		}

		if n.AppName != "" {
			out.AppName = n.AppName
		}
		if n.ReplacesID != 0 {
			out.ReplacesID = n.ReplacesID
		}
		if n.AppIcon != "" {
			out.AppIcon = n.AppIcon
		}
		if n.Summary != "" {
			out.Summary = n.Summary
		}
		if n.Body != "" {
			out.Body = n.Body
		}
		if n.Actions != nil {
			out.Actions = n.Actions
		}
		// if n.Hints != "" {
		// 	out.Hints = n.Hints
		// }
		if n.ExpireTimeout != 0 {
			out.ExpireTimeout = n.ExpireTimeout
		}
	}

	return out
}
