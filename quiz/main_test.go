package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestPrintBanner_DoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("printBanner panicked:", r)
		}
	}()
	printBanner()
}

func TestPrintResult_DoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("printResultPanicked", r)
		}
	}()
	printResult(0, 0)
}

func TestOpenFile_BadFilename(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("readFile is supposed to panic when given a bogus filename!")
		}
	}()
	openFile("./ljkasjf908asfasfia.gopher")
}

func TestReadProblems_Normal(t *testing.T) {
	r := strings.NewReader("5+5,10\n1+1,2\n8+3,11")
	expected := []question{
		{"5+5", "10"},
		{"1+1", "2"},
		{"8+3", "11"},
	}
	actual := readProblems(r)
	if !reflect.DeepEqual(actual, expected) {
		t.Error("expected", expected, "got", actual)
	}
}
