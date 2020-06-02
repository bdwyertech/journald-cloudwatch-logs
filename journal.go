package main

import (
	"github.com/coreos/go-systemd/v22/sdjournal"
	"strconv"
)

func AddLogFilters(journal *sdjournal.Journal, config *Config) {

	// Add Priority Filters
	if config.LogPriority < DEBUG {
		//lint:ignore S1005 Not sure staticcheck is correct
		for p, _ := range PriorityJSON {
			if p <= config.LogPriority {
				journal.AddMatch("PRIORITY=" + strconv.Itoa(int(p)))
			}
		}
		journal.AddDisjunction()
	}
}
