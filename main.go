package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

type miscCollector struct {
	loggedInUsersMetric    *prometheus.Desc
	sshSessionsMetric      *prometheus.Desc
	ansibleProcessesMetric *prometheus.Desc
}

func newMiscCollector() *miscCollector {
	return &miscCollector{
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

func (collector *miscCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.loggedInUsersMetric
	ch <- collector.sshSessionsMetric
	ch <- collector.ansibleProcessesMetric
}

func (collector *miscCollector) Collect(ch chan<- prometheus.Metric) {
	sshSessions_, sshSessionsErr := GetSshSessions()
	ansibleProcesses_, ansibleProcessesErr := GetAnsibleProcesses()

	loggedInUsers := prometheus.MustNewConstMetric(collector.loggedInUsersMetric, prometheus.GaugeValue, float64(GetLoggedInUsers()))
	sshSessions := prometheus.MustNewConstMetric(collector.sshSessionsMetric, prometheus.GaugeValue, float64(sshSessions_))
	ansibleProcesses := prometheus.MustNewConstMetric(collector.ansibleProcessesMetric, prometheus.GaugeValue, float64(ansibleProcesses_))

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
