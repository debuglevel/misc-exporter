package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/shirou/gopsutil/cpu"
	"go.uber.org/zap"
)

func GetCpuIdentifier() (string, error) {
	zap.L().Debug("Getting CPU identifier...")

	info, err := cpu.Info()
	if err != nil {
		zap.L().Error(err.Error())
		return "", fmt.Errorf("getting cpu info failed")
	}

	var cpuIdentifier = ""

	// Print the CPU manufacturer and model
	for _, ci := range info {
		//fmt.Printf("CPU manufacturer: %s\n", ci.VendorID)
		//fmt.Printf("CPU model: %s\n", ci.ModelName)
		zap.S().Debugf("CPU info: %v\n", ci)
		cpuIdentifier = ci.ModelName
	}

	zap.S().Debugf("Got CPU identifier: %v\n", cpuIdentifier)
	return cpuIdentifier, nil
}

func GetBenchmarkPage(cpuIdentifier string) (string, error) {
	zap.L().Debug("Getting Passmark CPU benchmark page...")

	url_ := "https://www.cpubenchmark.net/cpu.php?cpu=" + url.QueryEscape(cpuIdentifier)

	response, err := http.Get(url_)
	if err != nil {
		zap.L().Debug(err.Error())
		return "", fmt.Errorf("getting URL %s failed", url_)
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		zap.L().Debug(err.Error())
		return "", fmt.Errorf("reading response body failed")
	}

	// Print the response body
	//zap.L().Debug(string(body))

	zap.L().Debug("Got Passmark CPU benchmark page")
	return string(body), nil
}

func ExtractSingleThreadedRating(html string) (int, error) {
	zap.L().Debug("Extracting single-threaded rating...")

	re := regexp.MustCompile(`<strong>Single Thread Rating:</strong> (.*?)<br>`)
	match := re.FindStringSubmatch(html)
	if match == nil {
		return -1, fmt.Errorf("regexp did not match anything")
	}
	result := match[1]

	singleThreadedRating, _ := strconv.Atoi(result)

	zap.S().Debugf("Extracted single-threaded rating: %v\n", singleThreadedRating)
	return singleThreadedRating, nil
}

func GetSingleThreadedRating() (int, error) {
	cpuIdentifier, err := GetCpuIdentifier()
	if err != nil {
		zap.L().Debug(err.Error())
		return -1, errors.New("getting cpu identifier failed")
	}

	html, _ := GetBenchmarkPage(cpuIdentifier)
	if err != nil {
		zap.L().Debug(err.Error())
		return -1, errors.New("getting PassMark webpage failed")
	}

	singleThreadedRating, _ := ExtractSingleThreadedRating(html)
	if err != nil {
		zap.L().Debug(err.Error())
		return -1, errors.New("extracting single-threaded rating failed")
	}

	return singleThreadedRating, nil
}
