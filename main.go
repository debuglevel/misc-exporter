package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

type randomCollector struct {
	doubleRandomMetric     *prometheus.Desc
	loggedInUsersMetric    *prometheus.Desc
	sshSessionsMetric      *prometheus.Desc
	ansibleProcessesMetric *prometheus.Desc
}

func newMiscCollector() *randomCollector {
	return &randomCollector{
		doubleRandomMetric: prometheus.NewDesc(
			"misc_random_double",
			"Shows whether a bar has occurred in our cluster",
			nil,
			nil,
		),
		loggedInUsersMetric: prometheus.NewDesc(
			"misc_logged_in_users_count",
			"How many users are logged in right now.",
			nil,
			nil,
		),
		sshSessionsMetric: prometheus.NewDesc(
			"misc_ssh_sessions_count",
			"How many SSH sessions there are now.",
			nil,
			nil,
		),
		ansibleProcessesMetric: prometheus.NewDesc(
			"misc_ansible_processes_count",
			"How many Ansible processes there are now.",
			nil,
			nil,
		),
	}
}

func (collector *randomCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.doubleRandomMetric
	ch <- collector.loggedInUsersMetric
}

// See https://stackoverflow.com/a/43246464/4764279
func Shellout(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func GetLoggedInUsers() int {
	out, err := exec.Command("who").Output()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	scanner.Split(bufio.ScanLines)

	count := 0
	for scanner.Scan() {
		count++
	}

	return count
}

func GetSshSessions() (int, error) {
	out, errout, err := Shellout("netstat -tnpa | grep 'ESTABLISHED.*sshd' | wc -l\n")
	if err != nil {
		log.Printf("GetSshSessions error: %v\n", err)
	}
	if errout != "" {
		fmt.Println("GetSshSessions ERROR: " + errout)
		return -1, errors.New("GetSshSessions error getting SSH sessions from netstat")
	}

	marks, err := strconv.Atoi(strings.Trim(out, "\n "))
	if err != nil {
		return -1, errors.New("GetSshSessions error during conversion")
	}

	return marks, nil
}

func GetAnsibleProcesses() (int, error) {
	out, errout, err := Shellout("ps -Af | grep ansible | grep -v grep | wc -l\n")
	log.Printf("GetAnsibleProcesses: %v\n", out)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	if errout != "" {
		fmt.Println("GetAnsibleProcesses ERROR: " + errout)
	}

	marks, err := strconv.Atoi(strings.Trim(out, "\n "))

	if err != nil {
		fmt.Println("GetAnsibleProcesses Error during conversion")
		return -1, errors.New("GetAnsibleProcesses Error during conversion")
	}

	return marks, nil
}

func (collector *randomCollector) Collect(ch chan<- prometheus.Metric) {
	sshSessions_, sshSessionsErr := GetSshSessions()
	ansibleProcesses_, ansibleProcessesErr := GetAnsibleProcesses()

	m2 := prometheus.MustNewConstMetric(collector.doubleRandomMetric, prometheus.GaugeValue, rand.Float64())
	loggedInUsers := prometheus.MustNewConstMetric(collector.loggedInUsersMetric, prometheus.GaugeValue, float64(GetLoggedInUsers()))
	sshSessions := prometheus.MustNewConstMetric(collector.sshSessionsMetric, prometheus.GaugeValue, float64(sshSessions_))
	ansibleProcesses := prometheus.MustNewConstMetric(collector.ansibleProcessesMetric, prometheus.GaugeValue, float64(ansibleProcesses_))

	ch <- m2
	ch <- loggedInUsers
	if sshSessionsErr == nil {
		ch <- sshSessions
	}
	if ansibleProcessesErr == nil {
		ch <- ansibleProcesses
	}
}

func main() {
	myMiscCollector := newMiscCollector()
	prometheus.MustRegister(myMiscCollector)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))
}
