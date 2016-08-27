package subcmd

import "github.com/variadico/noti/cmd/noti/run"

type Cmd interface {
	Run() error
	Notify(stats run.Stats) error
}
