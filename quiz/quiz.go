package main
import (
	"fmt"
	"os"
	"log"
	"bufio"
	"strings"
	"time"
	"flag"
	//"strconv"
)

// create a struct to bind question and answer together
type problem struct {
	question string
	answer string
}

func main() {
	// read command line flag (time limit)
	quiz_dir := flag.String("csv", "./problems.csv",
		"a csv file in the form of [question,answer]")
	var timelimit int
	flag.IntVar(&timelimit, "limit", 30, 
		"Input an int to set the time limit to the quiz in seconds")
	flag.Parse()

	// open quiz file
	quiz_file, err := os.Open(*quiz_dir)
	if err != nil {
		log.Fatal(err)
	}
	defer quiz_file.Close()

	// create a slice of problems
	// this creates an empty slice
	problems := make([]problem, 0)

	// read from quiz file
	// once a line is read in, seperate the question and the answer
	scanner := bufio.NewScanner(quiz_file)
	for scanner.Scan() {
		current_line := scanner.Text()
		question_answer := strings.Split(current_line, ",")
		
		problems = append(problems, problem{question_answer[0], question_answer[1]})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// define vars related to quiz
	num_questions := len(problems)
	correct_count := 0
	// done is a channel that sends bool
	// it marks user has answered all questions
	done := make(chan bool)

	// start the quiz
	go func() {
		fmt.Println("--- The Quiz has Started ---")

		for _, p := range problems {
			fmt.Printf("%s = ", p.question)
			var answer_try string
			fmt.Scanln(&answer_try)
			//fmt.Printf("%s\n", answer_try)
			if answer_try == p.answer {
				correct_count++
			}
		}
		done <-true
	}()

	select{
	case <-time.After(time.Duration(timelimit) * time.Second):
		fmt.Println("\nTime is up!")
		fmt.Printf("%d correct out of %d questions\n", correct_count, num_questions)
		fmt.Println("--- Quiz has Ended ---")
	case <-done:
		fmt.Printf("%d correct out of %d questions\n", correct_count, num_questions)
		fmt.Println("--- Quiz has Ended ---")
	}
}