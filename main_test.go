package main

import "testing"

func TestImageServerStatsEmpty(t *testing.T) {
	is := ImageServer{}
	num, avgW, avgH := is.GetStats()
	if num != 0 {
		t.Error("Number of images should be 0")
	}
	if avgW != 0 {
		t.Error("Avg weight should be 0")
	}
	if avgH != 0 {
		t.Error("Avg height should be 0")
	}
}

func TestImageServerStatsAvg(t *testing.T) {
	is := ImageServer{}
	is.AddStats(100, 100)
	is.AddStats(200, 20)
	num, avgW, avgH := is.GetStats()
	if num != 2 {
		t.Error("Number of images should be 2")
	}
	if avgW != 150 {
		t.Error("Avg weight should be 150")
	}
	if avgH != 60 {
		t.Error("Avg height should be 60")
	}
}
