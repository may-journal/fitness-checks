package acceptance

import "testing"

// TestSlotsHoldTheMachineToTheirCount locks every slot in a fresh folder and
// expects the next try to find none free.
func TestSlotsHoldTheMachineToTheirCount(t *testing.T) {
	dir := t.TempDir()
	for i := range 2 {
		if _, ok := trySlot(dir, 2); !ok {
			t.Fatalf("slot %d of 2 was not free in an empty folder", i)
		}
	}
	if _, ok := trySlot(dir, 2); ok {
		t.Fatal("a third binary got a slot while both of 2 were held")
	}
}

// TestFreedSlotIsTakenAgain closes a held slot and expects the next try to
// take it.
func TestFreedSlotIsTakenAgain(t *testing.T) {
	dir := t.TempDir()
	f, ok := trySlot(dir, 1)
	if !ok {
		t.Fatal("the only slot was not free in an empty folder")
	}
	held = held[:len(held)-1]
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	if _, ok := trySlot(dir, 1); !ok {
		t.Fatal("a freed slot was not taken again")
	}
}
