package main

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"
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

func TestReadProblems_EmptyReader(t *testing.T) {
	r := strings.NewReader("")
	actual := readProblems(r)
	if actual == nil || len(actual) != 0 {
		t.Error("did not return empty []question from empty io.Reader!")
	}
}

func TestReadProblems_InvalidFormat_SingleColumn(t *testing.T) {
	r := strings.NewReader(`
	5+5
	1+1,
	8+3`)
	defer func() {
		if r := recover(); r == nil {
			t.Error("invalid format not detected in readProblems. Two columns are necessary!", r)
		}
	}()
	readProblems(r)
}

func TestReadProblems_InvalidFormat_TrippleColumn(t *testing.T) {
	r := strings.NewReader(`
	5+5,10,really???
	1+1,2,zzzz
	8+3,11,blubbi`)
	defer func() {
		if r := recover(); r == nil {
			t.Error("invalid format not detected in readProblems. Not more than two columns are allowed!", r)
		}
	}()
	readProblems(r)
}

func TestReadProblems_InvalidFormat_LastLineInvalid(t *testing.T) {
	r := strings.NewReader(`
	5+5,10
	1+1,2
	blubbi!!!!!!blub :O`)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Format of last line is not correct!", r)
		}
	}()
	readProblems(r)
}

func TestPlayGame_Normal(t *testing.T) {
	r := strings.NewReader("5+5,10\n1+1,2\n8+3,11")
	ps := readProblems(r)
	obuf := bytes.NewBufferString("10\n")
	playGame(ps, obuf, time.Second*time.Duration(5))

	timer := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)

	go func() {
		var i = 0
		for {
			select {
			case <-done:
				return
			case <-timer.C:
				s := strings.TrimSuffix(ps[i].answer, "\n") + "\n"
				fmt.Println("!", s)
				io.WriteString(obuf, s)
				i++
			}
		}
	}()

	time.Sleep(1600 * time.Millisecond)
	timer.Stop()
	done <- true
}
