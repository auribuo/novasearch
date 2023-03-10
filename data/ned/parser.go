package ned

import (
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

func parseResponse(response string) ([]Response, error) {
	lines := strings.Split(response, "\n")

	i := 0

	for ; i < len(lines) && !strings.HasPrefix(lines[i], "UGC"); i++ {
	}

	responses := make([]Response, 0)

	for ; i < len(lines) && strings.HasPrefix(lines[i], "UGC"); i++ {
		responses = append(responses, parseLine(lines[i]))
	}

	return responses, nil
}

func parseLine(line string) Response {
	columns := strings.Split(line, "|")

	for i := range columns {
		columns[i] = strings.TrimSpace(columns[i])
	}

	//fmt.Printf("parsing: %s\n", strings.Join(columns, "|"))

	ugcNumber, err := strconv.Atoi(columns[0][3:])
	cobra.CheckErr(err)
	preferredName := columns[1]
	redshiftString := columns[2]
	if redshiftString == "" {
		redshiftString = "-1"
	}
	redshift, err := strconv.ParseFloat(redshiftString, 64)
	cobra.CheckErr(err)
	magnitude := parseMagnitude(columns[3])
	hubbleType := columns[4]

	return Response{
		UgcNumber:     ugcNumber,
		PreferredName: preferredName,
		Redshift:      redshift,
		Magnitude:     magnitude,
		HubbleType:    hubbleType,
	}
}

func parseMagnitude(magnitudeString string) float64 {
	if magnitudeString == "" {
		return -1
	}

	str := strings.Builder{}

	for _, char := range magnitudeString {
		if (char > '0' && char < '9') || char == '.' {
			str.WriteRune(char)
		}
	}

	magnitude, err := strconv.ParseFloat(str.String(), 64)
	if err != nil {
		return -1
	}
	cobra.CheckErr(err)
	return magnitude
}
