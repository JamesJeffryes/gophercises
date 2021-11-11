package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	questionFile := *flag.String("f", "problems.csv", "Path to CSV of problems")
	timeLimit := *flag.Int("t", 30, "Max quiz time in seconds")
	flag.Parse()
	problems := readProblems(questionFile)
	score := runQuiz(problems, timeLimit)
	fmt.Printf("You got %d of %d questions right!", score, len(problems))
}

type problemMap map[string]string

func readProblems(filename string) problemMap {
	file, err := os.Open(filename)
	handleError(err)
	defer file.Close()

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	handleError(err)
	problems := make(problemMap, len(lines))

	for _, record := range lines {
		problems[record[0]] = record[1]
	}
	fmt.Printf("Imported %d problems\n", len(problems))
	return problems
}

func runQuiz(problems problemMap, timeLimit int) int {
	correct := 0
	consoleReader := bufio.NewReader(os.Stdin)
	fmt.Printf("You will have %d seconds to complete this quiz. Hit enter to begin.", timeLimit)
	consoleReader.ReadString('\n')
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	quizDone := make(chan bool)

	go func() {
		for question, answer := range problems {
			fmt.Println(question)
			// also consider ScanF
			response, err := consoleReader.ReadString('\n')
			handleError(err)
			if strings.TrimSpace(response) == strings.TrimSpace(answer) {
				correct++
			}
		}
		quizDone <- true
	}()
	// select chooses the first channel with an available value
	// if done is available first, the user finished
	// if ticker is available first, the time limit has been reached
	select {
	case <-quizDone:
		fmt.Println("All Done!")
	case <-timer.C:
		fmt.Println("Time's up!")
	}
	return correct
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
