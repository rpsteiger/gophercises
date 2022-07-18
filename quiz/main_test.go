package main

import (
	"fmt"
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

func ExampleReadProblems() {
	r := strings.NewReader("5+5,10\n1+1,2\n8+3,11")
	fmt.Println(readProblems(r))
	// Output:
	// [{5+5 10} {1+1 2} {8+3 11}]
}

func TestReadProblems_TabsAtStartOfLine(t *testing.T) {
	r := strings.NewReader("\t5+5,\t10\n\t1+1,\t2\n\t8+3,\t11")
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

func TestReadProblems_TabsAtEndOfline(t *testing.T) {
	r := strings.NewReader("5+5\t,10\t\n1+1\t,2\t\n8+3\t,11\t")
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
