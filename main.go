package main

import (
	"net/http"

	"flag"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type miscCollector struct {
	loggedInUsersMetric    *prometheus.Desc
	sshSessionsMetric      *prometheus.Desc
	ansibleProcessesMetric *prometheus.Desc
	performanceMetric      *prometheus.Desc
	passmarkMetric         *prometheus.Desc
}

func newMiscCollector() *miscCollector {
	zap.L().Debug("Creating collector...")

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
	zap.L().Debug("Describing...")

	ch <- collector.loggedInUsersMetric
	ch <- collector.sshSessionsMetric
	ch <- collector.ansibleProcessesMetric
	ch <- collector.performanceMetric
	ch <- collector.passmarkMetric
}

func (collector *miscCollector) Collect(ch chan<- prometheus.Metric) {
	zap.L().Debug("Collecting...")

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
	loglevelPtr := flag.String("loglevel", "info", "log level (debug, info, warn, error, fatal)")
	flag.Parse()

	var level zap.AtomicLevel
	if *loglevelPtr == "debug" {
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else if *loglevelPtr == "info" {
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else if *loglevelPtr == "warn" {
		level = zap.NewAtomicLevelAt(zap.WarnLevel)
	} else if *loglevelPtr == "error" {
		level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	} else {
		level = zap.NewAtomicLevelAt(zap.FatalLevel)
	}

	loggerConfig := zap.Config{
		Level:            level,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, _ := loggerConfig.Build()
	defer logger.Sync() // nolint:errcheck // not sure how to errcheck a deferred call like this
	zap.ReplaceGlobals(logger)

	zap.L().Info("Starting...")

	myMiscCollector := newMiscCollector()
	prometheus.MustRegister(myMiscCollector)

	zap.L().Info("Serving...")
	http.Handle("/metrics", promhttp.Handler())
	zap.L().Fatal(http.ListenAndServe("127.0.0.1:9886", nil).Error())
}
