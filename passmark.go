package main

import (
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

func GetCpuIdentifier() (string, error) {
	log.Println("Getting CPU identifier...")

	info, err := cpu.Info()
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("getting cpu info failed")
	}

	var cpuIdentifier = ""

	// Print the CPU manufacturer and model
	for _, ci := range info {
		//fmt.Printf("CPU manufacturer: %s\n", ci.VendorID)
		//fmt.Printf("CPU model: %s\n", ci.ModelName)

		cpuIdentifier = ci.ModelName
	}

	log.Printf("Got CPU identifier: %v\n", cpuIdentifier)
	return cpuIdentifier, nil
}

func GetBenchmarkPage(cpuIdentifier string) (string, error) {
	log.Println("Getting Passmark CPU benchmark page...")

	url_ := "https://www.cpubenchmark.net/cpu.php?cpu=" + url.QueryEscape(cpuIdentifier)

	response, err := http.Get(url_)
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("getting URL %s failed", url_)
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("reading response body failed")
	}

	// Print the response body
	//fmt.Println(string(body))

	log.Println("Got Passmark CPU benchmark page")
	return string(body), nil
}

func ExtractSingleThreadedRating(html string) (int, error) {
	log.Println("Extracting single-threaded rating...")

	re := regexp.MustCompile(`<strong>Single Thread Rating:</strong> (.*?)<br>`)
	match := re.FindStringSubmatch(html)
	if match == nil {
	  return nil, fmt.Errorf("regexp did not match anything")
	}
	result := match[1]

	singleThreadedRating, _ := strconv.Atoi(result)

	log.Printf("Extracted single-threaded rating: %v\n", singleThreadedRating)
	return singleThreadedRating, nil
}

func GetSingleThreadedRating() (int, error) {
	cpuIdentifier, err := GetCpuIdentifier()
	if err != nil {
		fmt.Println(err)
		return -1, errors.New("getting cpu identifier failed")
	}

	html, _ := GetBenchmarkPage(cpuIdentifier)
	if err != nil {
		fmt.Println(err)
		return -1, errors.New("getting PassMark webpage failed")
	}

	singleThreadedRating, _ := ExtractSingleThreadedRating(html)
	if err != nil {
		fmt.Println(err)
		return -1, errors.New("extracting single-threaded rating failed")
	}

	return singleThreadedRating, nil
}
