package slotlimiter

import (
	"testing"
	"time"
)

func TestSlotlimiter_Occupy(t *testing.T) {
	pl := New()
	slot := pl.GetSlot("test1", 1)
	startedAt := time.Now()

	if slot == nil {
		t.Fatal("didn't expect GetSlot to return nil")
	}

	if cap(*slot) != 1 {
		t.Fatal("expected GetSlot to return a channel of capacity 1")
	}

	if slot.SlotsUsed("test1") != 0 {
		t.Fatal("expected new slot to have 0 users")
	}

	slot.Occupy()
	go func() {
		time.Sleep(100 * time.Millisecond)
		slot.Leave()
	}()

	if slot.SlotsUsed("test1") != 1 {
		t.Fatal("expected occupied slot to have 1 users")
	}

	slot.Occupy()
	if time.Now().Sub(startedAt).Seconds() < 0.1 {
		t.Error("Occupy should have taken at least 0.1s")
	}
	slot.Leave()

	if slot.SlotsUsed("test1") != 0 {
		t.Fatal("expected left slot to have 0 users")
	}
}

func TestSlotlimiter_OccupyWithTimeout(t *testing.T) {
	pl := New()
	slot := pl.GetSlot("test2", 1)
	startedAt := time.Now()

	slot.Occupy()
	go func() {
		time.Sleep(100 * time.Millisecond)
		slot.Leave()
	}()

	if slot.OccupyWithTimeout(0) != false {
		t.Fatal("OccupyWithTimeout(0) should have given up immediately")
	}
	if slot.OccupyWithTimeout(2*time.Second) != true {
		t.Fatal("OccupyWithTimeout should have gotten a slot")
	}
	if time.Now().Sub(startedAt).Seconds() < 0.1 {
		t.Error("OccupyWithTimeout should have taken at least 0.1s")
	}
}
