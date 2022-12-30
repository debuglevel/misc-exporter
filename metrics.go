package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// GetLoggedInUsers gets the current count of logged-in users (i.e. have a shell associated).
func GetLoggedInUsers() int {
	log.Println("Getting logged-in users...")
	stdout, err := exec.Command("who").Output()
	if err != nil {
		log.Println(err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))
	scanner.Split(bufio.ScanLines)

	count := 0
	for scanner.Scan() {
		count++
	}

	log.Printf("Got logged-in users: %v\n", count)
	return count
}

// GetSshSessions gets the current SSH sessions.
func GetSshSessions() (int, error) {
	log.Println("Getting SSH sessions...")

	stdout, stderr, err := ShellExecute("netstat -tnpa | grep 'ESTABLISHED.*sshd' | wc -l\n")
	if err != nil {
		log.Printf("Error getting SSH sessions from: %v\n", err)
		return -1, errors.New("command for getting SSH sessions returned error return code")
	}
	if stderr != "" {
		fmt.Println("Command stderr was not empty: " + stderr)
		return -1, errors.New("command stderr was not empty: " + stderr)
	}

	sshSessions, err := strconv.Atoi(strings.Trim(stdout, "\n "))
	if err != nil {
		fmt.Println("Converting output to integer failed")
		return -1, errors.New("converting output to integer failed")
	}

	log.Printf("Got SSH sessions: %v\n", sshSessions)
	return sshSessions, nil
}

// GetAnsibleProcesses gets the count of current ansible processes
func GetAnsibleProcesses() (int, error) {
	log.Println("Getting Ansible processes...")

	stdout, stderr, err := ShellExecute("ps -Af | grep ansible | grep -v grep | wc -l\n")
	if err != nil {
		log.Printf("Error getting Ansible processes: %v\n", err)
		return -1, errors.New("command for getting Ansible processes returned error return code")
	}
	if stderr != "" {
		fmt.Println("Command stderr was not empty: " + stderr)
		return -1, errors.New("command stderr was not empty: " + stderr)
	}

	ansibleProcesses, err := strconv.Atoi(strings.Trim(stdout, "\n "))
	if err != nil {
		fmt.Println("Converting output to integer failed")
		return -1, errors.New("converting output to integer failed")
	}

	log.Printf("Got Ansible processes: %v\n", ansibleProcesses)
	return ansibleProcesses, nil
}

func IsItAGoodTimeToEvaluatePerformance() bool {
	now := time.Now()
	hour := now.Hour()
	minute := now.Minute()

	if hour == 6 && minute >= 00 || hour == 6 && minute <= 10 {
		return true
	} else {
		return false
	}
}

// GetPerformance calculates a single-threaded performance indicator
func GetPerformance() float64 {
	log.Println("Getting performance...")

	maximumPrime := 1000
	seconds := 500 * time.Millisecond

	start := time.Now()
	iterations := 0

	for ; time.Since(start) < seconds; iterations++ {
		getPrimes(maximumPrime)
	}
	elapsed := time.Since(start)
	secondsPerIteration := elapsed.Seconds() / float64(iterations)
	iterationsPerSecond := float64(iterations) / elapsed.Seconds()
	fmt.Printf("%v iterations took %s; %vs per iteration; %v iterations per second\n", iterations, elapsed, secondsPerIteration, iterationsPerSecond)

	performance := iterationsPerSecond

	log.Printf("Got performance: %v\n", performance)
	return performance
}
