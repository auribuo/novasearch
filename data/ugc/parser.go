package ugc

import (
	"github.com/auribuo/novasearch/log"
	"github.com/auribuo/novasearch/sql"
	"github.com/auribuo/novasearch/types/coordinates"
	"github.com/auribuo/novasearch/util"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

func parseResponse(response string) ([]Response, error) {
	lines := strings.Split(response, "\n")

	i := 0

	for ; i < len(lines) && !strings.HasPrefix(lines[i], "----------"); i++ {
	}

	i++

	for ; i < len(lines) && !strings.HasPrefix(lines[i], "----------"); i++ {
	}

	header := lines[i]

	i++

	responses := make([]Response, 0)

	for ; i < len(lines) && !strings.HasPrefix(lines[i], "#"); i++ {
		line := lines[i]
		if line == "" {
			continue
		}
		responses = append(responses, parseLine(line, header))
	}

	return responses, nil
}

func parseLine(line string, header string) Response {
	log.Logger.Debugf("parsing: %s", line)
	columnDelimiters := strings.Split(header, " ")
	columnWidths := sql.Map(columnDelimiters, func(t string) int {
		return len(t)
	})
	columns := util.SliceBasedOn(line, columnWidths)
	for i := range columns {
		columns[i] = strings.TrimSpace(columns[i])
	}

	raHours, err := strconv.ParseFloat(strings.Split(columns[0], " ")[0], 64)
	cobra.CheckErr(err)
	raMinutes, err := strconv.ParseFloat(strings.Split(columns[0], " ")[1], 64)
	cobra.CheckErr(err)
	raSeconds, err := strconv.ParseFloat(strings.Split(columns[0], " ")[2], 64)
	cobra.CheckErr(err)
	decDegrees, err := strconv.ParseFloat(strings.Split(columns[1], " ")[0], 64)
	cobra.CheckErr(err)
	decMinutes, err := strconv.ParseFloat(strings.Split(columns[1], " ")[1], 64)
	cobra.CheckErr(err)
	decSeconds, err := strconv.ParseFloat(strings.Split(columns[1], " ")[2], 64)
	cobra.CheckErr(err)
	ugcNumber, err := strconv.Atoi(columns[2])
	cobra.CheckErr(err)
	semiMajorAxisString := columns[3]
	if semiMajorAxisString == "" {
		semiMajorAxisString = "-1"
	}
	semiMajorAxis, err := strconv.ParseFloat(semiMajorAxisString, 64)
	cobra.CheckErr(err)
	semiMinorAxisString := columns[4]
	if semiMinorAxisString == "" {
		semiMinorAxisString = "-1"
	}
	semiMinorAxis, err := strconv.ParseFloat(semiMinorAxisString, 64)
	cobra.CheckErr(err)
	hubbleType := columns[6]
	cobra.CheckErr(err)
	magnitudeString := columns[7]
	if magnitudeString == "" {
		magnitudeString = "-1"
	}
	magnitude, err := strconv.ParseFloat(magnitudeString, 64)
	cobra.CheckErr(err)

	return Response{
		RightAscension: coordinates.RightAscensionToDegrees(raHours, raMinutes, raSeconds),
		Declination:    coordinates.DeclinationToDegrees(decDegrees, decMinutes, decSeconds),
		UgcNumber:      ugcNumber,
		SemiMajorAxis:  semiMajorAxis,
		SemiMinorAxis:  semiMinorAxis,
		HubbleType:     hubbleType,
		Magnitude:      magnitude,
	}
}
