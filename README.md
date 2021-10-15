### Limit concurrency based on scarcity of virtual slots

This Go package provides a generic way to limit the concurrency based on occupying a limited number of available slots.

It supports an optional timeout to give up waiting for a slot.

You can have lots of different groups of slots, each with its own (string) name.

It uses a buffered Go Channel internally, which usefully, blocks when full.


#### Install
```shell
$ go get github.com/andys/slotlimiter
```


#### Usage

```go
import (
  "time"
  "github.com/andys/slotlimiter"
)

// Create a long-lived handle
slotLimiter := slotlimiter.New()

// Select a slot group labelled bigtask that only lets 3 run at once
// Changing the concurrency makes it start fresh with all the slots cleared
slot := slotLimiter.GetSlot("bigtask", 3)

// Wait for up to 10 seconds to get one of the 3 slots
// returns true if we got the slot, or false if we timed out.
// Use Occupy() to wait forever
result := slot.OccupyWithTimeout(10 * time.Second)

// Make sure you leave the slot when done otherwise it will be locked forever.
// defer is a good way to make sure this is done for you.
defer slot.Leave()

// How many slots are taken
n := slot.SlotsUsed()
n := slotLimiter.SlotsUsed("bigtask")

// List of all slots
strings := slotLimiter.GetSlots()

```

