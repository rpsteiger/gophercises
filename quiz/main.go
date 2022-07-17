package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	ct "github.com/daviddengcn/go-colortext"
)

const (
	GameTime = time.Second * 30
)

type question struct {
	questionStr string
	answer      string
}

func printBanner() {
	fmt.Printf(`   _____ ______  _______  __    ______   ____  __  ___________
  / ___//  _/  |/  / __ \/ /   / ____/  / __ \/ / / /  _/__  /
  \__ \ / // /|_/ / /_/ / /   / __/    / / / / / / // /   / / 
 ___/ // // /  / / ____/ /___/ /___   / /_/ / /_/ // /   / /__
/____/___/_/  /_/_/   /_____/_____/   \___\_\____/___/  /____/`)
	fmt.Printf("\n")
	fmt.Printf("\n")
}

func printResult(correct int, wrong int) {
	ct.ResetColor()
	fmt.Println("--------------------------------------------------------------")
	fmt.Printf("You guessed %d out of %d correctly!\n", correct, (correct + wrong))
	fmt.Println("--------------------------------------------------------------")
}

func printStatsLine(correct int, wrong int) {
	ct.ResetColor()
	fmt.Printf("correct: ")
	ct.Foreground(ct.Green, false)
	fmt.Printf("%d", correct)
	ct.ResetColor()
	fmt.Printf(", wrong: ")
	ct.Foreground(ct.Red, false)
	fmt.Printf("%d", wrong)
	fmt.Printf("\n")
}

var questionsFilename string
var timeLimit int

func init() {
	const (
		defaultFilename      = "problems.csv"
		filenameExplanation  = `the text file containing the questions in csv format.`
		defaultTimeLimit     = 30
		timeLimitExplanation = `the time limite for the quiz in seconds.`
	)
	flag.StringVar(&questionsFilename, "filename", defaultFilename, filenameExplanation)
	flag.StringVar(&questionsFilename, "f", defaultFilename, filenameExplanation+" (shorthand)")
	flag.IntVar(&timeLimit, "limit", defaultTimeLimit, timeLimitExplanation)
	flag.IntVar(&timeLimit, "l", defaultTimeLimit, timeLimitExplanation+" (shorthand")
}

func main() {
	flag.Parse()

	questions := readFile(questionsFilename)
	timer := time.NewTimer(time.Second * time.Duration(timeLimit))

	correct := 0
	wrong := 0
	printBanner()

	for i, q := range questions {
		printStatsLine(correct, wrong)
		ct.Foreground(ct.Blue, false)
		fmt.Printf("#%d: %s ?\n", i+1, q.questionStr)

		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanln(&answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			printResult(correct, wrong)
			return
		case answer := <-answerCh:
			if answer == q.answer {
				correct++
			} else {
				wrong++
			}
			close(answerCh)
		}
	}
	printResult(correct, wrong)
}

func readFile(filename string) []question {
	f, err := os.Open(filename)
	if err != nil {
		log.Panicf("Error while opening file '%s': %v\n", filename, err)
	}
	defer f.Close()

	csvR := csv.NewReader(f)
	results := make([]question, 0, 20)

	for {
		record, err := csvR.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Panicf("Error while reading csv line: %v\n", err)
		}

		q := question{
			questionStr: record[0],
			answer:      strings.TrimSpace(record[1]),
		}
		results = append(results, q)
	}

	return results
}
