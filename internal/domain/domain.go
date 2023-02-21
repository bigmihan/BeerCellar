package domain

import "time"

type CellarRemains struct {
	Sum_beer int
	Batch    time.Time
	Beer_id  int
}
