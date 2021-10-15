package slotlimiter

import (
	"sync"
	"time"
)

type Slot chan int

type Slotlimiter struct {
	slots map[string]Slot
	mutex sync.Mutex
}

func New() *Slotlimiter {
	return &Slotlimiter{slots: make(map[string]Slot), mutex: sync.Mutex{}}
}

func (this *Slotlimiter) GetSlot(key string, concurrencyLimit int) *Slot {
	return this.getChannel(key, concurrencyLimit)
}

func (this *Slotlimiter) GetSlots() []string {
  keys := make([]string, 0, len(this.slots))
  for key := range this.slots {
  	keys = append(keys, key)
  }
	return keys
}

func (this *Slotlimiter) SlotsUsed(key string) int {
	return len(this.slots[key])
}

/* Wait for a slot forever */
func (this *Slot) Occupy() {
	if this != nil {
		/* Send to the channel to take a slot. It will block if not enough capacity */
		*this <- 0
	}
}

/* Wait for a slot with timeout. Returns false if we timed out */
func (this *Slot) OccupyWithTimeout(timeout time.Duration) bool {
	select {
	case *this <- 0:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (this *Slot) SlotsUsed() int {
	return len(*this)
}

func (this *Slot) Leave() {
	if this != nil {
		select {
		case <-*this: /* take something out of the channel to leave the slot. */
		default: /* ignore if there was nothing to take. */
		}
	}
}

func (this *Slotlimiter) getChannel(key string, concurrencyLimit int) *Slot {
	/* Get the channel associated with this slot */
	this.mutex.Lock()
	slot, exists := this.slots[key]
	if !exists || cap(slot) != concurrencyLimit {
		slot = make(Slot, concurrencyLimit)
		this.slots[key] = slot
	}
	this.mutex.Unlock()
	return &slot
}
