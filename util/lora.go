package util

import (
	"math"
	"regexp"
	"strconv"
)

// CalculatePacketTime calculates the time it took to transmit the packet, given the payload size, the bandwidth and the spreading factor.
func CalculatePacketTime(phyPayloadSize int, dataRate string) (float64, error) {
	dataRateRegexp := regexp.MustCompile(`SF(\d+)BW(\d+)`)
	matches := dataRateRegexp.FindStringSubmatch(dataRate)

	spreadingFactor, err := strconv.ParseInt(matches[1], 10, 0)
	if err != nil {
		return 0, err
	}

	bandwidth, err := strconv.ParseInt(matches[2], 10, 0)

	if err != nil {
		return 0, err
	}

	tSymb := math.Pow(2, float64(spreadingFactor)) / float64(bandwidth)

	dataRateOptimization := 0
	if int(spreadingFactor) >= 11 {
		dataRateOptimization = 1
	}

	preamble := float64(8 + 4.25)
	payload := math.Max(math.Ceil(float64(8*phyPayloadSize-4*int(spreadingFactor)+24)/float64(4*int(spreadingFactor)-2*dataRateOptimization))*5, 0)
	crc := float64(8)

	return (preamble + payload + crc) * tSymb, nil
}
