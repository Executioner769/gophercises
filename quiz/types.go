package main

import "strings"

type Problem struct {
	question string
	answer   string
}

func NewProblem(question string, answer string) *Problem {
	return &Problem{
		question: question,
		answer:   strings.TrimSpace(answer),
	}
}
