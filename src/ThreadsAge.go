package main

import (
	"errors"
	"strings"
)

type ThreadsAge int

var (
	lastDay     ThreadsAge = 1
	last2Days   ThreadsAge = 2
	last5Days   ThreadsAge = 5
	last10Days  ThreadsAge = 10
	last20Days  ThreadsAge = 20
	last30Days  ThreadsAge = 30
	last45Days  ThreadsAge = 45
	last60Days  ThreadsAge = 60
	last75Days  ThreadsAge = 75
	last100Days ThreadsAge = 100
	lastYear    ThreadsAge = 365
	all         ThreadsAge = -1
)

func (hc ThreadsAge) String() (string, error) {
	switch hc {
	case lastDay:
		return "lastDay", nil
	case last2Days:
		return "last2Days", nil
	case last5Days:
		return "last2Days", nil
	case last10Days:
		return "last10days", nil
	case last20Days:
		return "last20days", nil
	case last30Days:
		return "last30days", nil
	case last45Days:
		return "last45days", nil
	case last60Days:
		return "last60days", nil
	case last75Days:
		return "last75days", nil
	case last100Days:
		return "last100days", nil
	case lastYear:
		return "lastYear", nil
	case all:
		return "all", nil
	default:
		return "", errors.New("invalid ThreadsAge value")
	}
}

func ThreadsAgeFromString(input string) ThreadsAge {
	input = strings.ToLower(input)
	switch input {
	case "lastday":
		return lastDay
	case "last2days":
		return last2Days
	case "last5days":
		return last5Days
	case "last10days":
		return last10Days
	case "last20days":
		return last20Days
	case "last30days":
		return last30Days
	case "last45days":
		return last45Days
	case "last60days":
		return last60Days
	case "last75days":
		return last75Days
	case "last100days":
		return last100Days
	case "lastyear":
		return lastYear
	case "all":
		return all
	default:
		return last30Days
	}
}
