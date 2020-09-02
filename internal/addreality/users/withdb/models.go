package withdb

import (
	"github.com/paulbellamy/ratecounter"
)

type User struct {
	ID      string
	Counter *ratecounter.RateCounter
}
