package main

import (
	"fmt"
	"github.com/docker/go-units"
	"github.com/montanaflynn/stats"
	"regexp"
	"strconv"
)

const (
	nanoSecond  = 1e-9
	microSecond = 1e-6
	milliSecond = 1e-3
	second      = float64(1.0)
	minute      = 60.0 * second
	hour        = 60.0 * minute
	day         = 24.0 * hour
	week        = 7.0 * day
)

type unitMap map[string]float64

var (
	durationMap = unitMap{"ns": nanoSecond, "us": microSecond, "ms": milliSecond,
		"s": second, "m": minute, "h": hour, "d": day, "w": week}
	durationRegex = regexp.MustCompile("ns|us|ms|s|m|h|d|w")
)

func Sum(data stats.Float64Data) float64 {
	sum, err := stats.Sum(data)
	if err != nil {
		fmt.Println("Error ", err)
	}
	return sum
}

func Median(data stats.Float64Data) float64 {
	median, err := stats.Median(data)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	return median
}

func fromHumanDuration(data string) (float64, error) {
	unit := durationRegex.FindString(data)
	durationVal, err := strconv.ParseFloat(data[0:len(data)-1], 64)
	return durationVal * durationMap[unit], err
}

func fromHumanSize(data string) (float64, error) {
	machineReadableSize, err := units.FromHumanSize(data)
	return float64(machineReadableSize), err
}

func boolToInt(data bool) int64 {
	if data {
		return 1
	} else {
		return 0
	}
}
