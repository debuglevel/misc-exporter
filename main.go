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
	performanceMetric      *prometheus.Desc
	passmarkMetric         *prometheus.Desc
}

func newMiscCollector() *miscCollector {
	log.Println("Creating collector...")

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
		performanceMetric: prometheus.NewDesc(
			"misc_performance",
			"How often the performance evaluation function could be called within a certain duration",
			nil,
			nil,
		),
		passmarkMetric: prometheus.NewDesc(
			"misc_passmark_singlethreadedrating",
			"Single-Threaded Rating of the CPU on PassMark",
			nil,
			nil,
		),
	}
}

func (collector *miscCollector) Describe(ch chan<- *prometheus.Desc) {
	log.Println("Describing...")

	ch <- collector.loggedInUsersMetric
	ch <- collector.sshSessionsMetric
	ch <- collector.ansibleProcessesMetric
	ch <- collector.performanceMetric
	ch <- collector.passmarkMetric
}

func (collector *miscCollector) Collect(ch chan<- prometheus.Metric) {
	log.Println("Collecting...")

	loggedInUsers := prometheus.MustNewConstMetric(collector.loggedInUsersMetric, prometheus.GaugeValue, float64(GetLoggedInUsers()))
	ch <- loggedInUsers

	rating, err := GetSingleThreadedRating()
	if err == nil {
		passmark := prometheus.MustNewConstMetric(collector.passmarkMetric, prometheus.GaugeValue, float64(rating))
		ch <- passmark
	}

	ansibleProcesses_, err := GetAnsibleProcesses()
	if err == nil {
		ansibleProcesses := prometheus.MustNewConstMetric(collector.ansibleProcessesMetric, prometheus.GaugeValue, float64(ansibleProcesses_))
		ch <- ansibleProcesses
	}

	sshSessions_, err := GetSshSessions()
	if err == nil {
		sshSessions := prometheus.MustNewConstMetric(collector.sshSessionsMetric, prometheus.GaugeValue, float64(sshSessions_))
		ch <- sshSessions
	}

	if IsItAGoodTimeToEvaluatePerformance() {
		performance := prometheus.MustNewConstMetric(collector.performanceMetric, prometheus.GaugeValue, GetPerformance())
		ch <- performance
	}
}

func main() {
	log.Println("Starting...")

	myMiscCollector := newMiscCollector()
	prometheus.MustRegister(myMiscCollector)

	log.Println("Serving...")
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe("127.0.0.1:9886", nil))
}
