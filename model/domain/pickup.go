package domain

import "time"

type Pickup struct {
	PickupId int
	Book     Book
	Schedule time.Time
}