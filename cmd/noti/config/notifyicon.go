// +build windows

package config

import "github.com/variadico/noti/notifyicon"

// MergeNotifyIcon combines banner notification structs. Each subsequent struct
// field overwrites the previous one, if it was set.
func MergeNotifyIcon(ns ...*notifyicon.Notification) *notifyicon.Notification {
	out := new(notifyicon.Notification)

	for _, n := range ns {
		if n == nil {
			continue
		}

		if n.BalloonTipIcon != "" {
			out.BalloonTipIcon = n.BalloonTipIcon
		}
		if n.BalloonTipText != "" {
			out.BalloonTipText = n.BalloonTipText
		}
		if n.BalloonTipTitle != "" {
			out.BalloonTipTitle = n.BalloonTipTitle
		}
		if n.Icon != "" {
			out.Icon = n.Icon
		}
		if n.Text != "" {
			out.Text = n.Text
		}
		if n.Duration != 0 {
			out.Duration = n.Duration
		}
	}

	return out
}
