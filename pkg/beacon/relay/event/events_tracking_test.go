package event

import (
	"sync"
	"testing"
)

func TestGroupSelectionTrack_Add(t *testing.T) {
	type groupSelected struct {
		Data  map[string]bool
		Mutex *sync.Mutex
	}
	tests := map[string]struct {
		groupSelected groupSelected
		arg           string
		expected      bool
	}{
		"adding unique entry key": {
			groupSelected: groupSelected{Data: make(map[string]bool), Mutex: &sync.Mutex{}},
			arg:           "0x123456",
			expected:      true,
		},
		"adding unique entry key to non-empty map": {
			groupSelected: groupSelected{Data: map[string]bool{"0x123456": true}, Mutex: &sync.Mutex{}},
			arg:           "0x7891011",
			expected:      true,
		},
		"adding existing entry key": {
			groupSelected: groupSelected{Data: map[string]bool{"0x123456": true}, Mutex: &sync.Mutex{}},
			arg:           "0x123456",
			expected:      false,
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			gst := &GroupSelectionTrack{
				Data:  test.groupSelected.Data,
				Mutex: test.groupSelected.Mutex,
			}
			if actual := gst.Add(test.arg); actual != test.expected {
				t.Errorf("GroupSelectionTrack.Add() = %v, expected %v", actual, test.expected)
			}
		})
	}
}

func TestGroupSelectionTrack_Remove(t *testing.T) {
	type groupSelected struct {
		Data  map[string]bool
		Mutex *sync.Mutex
	}
	tests := map[string]struct {
		groupSelected groupSelected
		arg           string
		expected      bool
	}{
		"adding same entry key after removing it from the map": {
			groupSelected: groupSelected{Data: map[string]bool{"0x123456": true}, Mutex: &sync.Mutex{}},
			arg:           "0x123456",
			expected:      true,
		},
		"removing from the empty map": {
			groupSelected: groupSelected{Data: make(map[string]bool), Mutex: &sync.Mutex{}},
			arg:           "0x123456",
			expected:      true,
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			gst := &GroupSelectionTrack{
				Data:  test.groupSelected.Data,
				Mutex: test.groupSelected.Mutex,
			}
			gst.Remove(test.arg)
			if actual := gst.Add(test.arg); actual != test.expected {
				t.Errorf("GroupSelectionTrack.Remove() should have removed entry %v", test.arg)
			}
		})
	}
}
