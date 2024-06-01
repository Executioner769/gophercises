package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var (
	csvFilename *string
	timeLimit   *int
	shuffle     *bool
)

func main() {

	parseFlags()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the csv file: %s", *csvFilename))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintln("Failed to parse the provided CSV file."))
	}

	problems := parsLines(lines)

	startQuiz(problems, time.Duration(*timeLimit), *shuffle)
}

func startQuiz(problems []*Problem, limit time.Duration, shuffle bool) int32 {

	var score int32
	count := len(problems)

	if shuffle {
		rand.Shuffle(count, func(i, j int) {
			problems[i], problems[j] = problems[j], problems[i]
		})
	}

	fmt.Printf("Total Problems: %d \n", count)
	fmt.Printf("Time Limit set to %d sec\n", limit)
	fmt.Println("Press Enter to begin!")
	fmt.Scanln()

	timer := time.NewTimer(limit * time.Second)

problemloop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.question)

		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println()
			break problemloop
		case answer := <-answerCh:
			if answer == p.answer {
				score += 1
			}
		}
	}

	fmt.Printf("You scored %d out of %d \n", score, count)

	return score
}

func parseFlags() {
	csvFilename = flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit = flag.Int("limit", 5, "The time limit for the quiz")
	shuffle = flag.Bool("shuffle", true, "Shuffle the quiz order")
	flag.Parse()
}

func parsLines(lines [][]string) []*Problem {
	res := make([]*Problem, len(lines))
	for i, line := range lines {
		res[i] = NewProblem(line[0], line[1])
	}
	return res
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
