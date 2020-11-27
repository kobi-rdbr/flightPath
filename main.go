package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// FlightsList is an list of flights
type FlightsList struct {
	Flights [][]string
}

// TotalFlightPath is used for JSON package to encode data to json
type TotalFlightPath struct {
	FlightPath []string
}

// subtruct the same nodes counter from
func subtructMaps(firstMap map[string]int, secondMap map[string]int) map[string]int {
	result := map[string]int{}
	for k, v := range firstMap {
		if val, ok := secondMap[k]; ok {
			result[k] = v - val
		} else {
			result[k] = v
		}
	}
	return result
}

// the function counts the occurance for each uniq entry in the list
// genereates two maps/counters for start and end airport respectivly
// map's key is an airport name, value is the counter, the frequency it occured in the list
func convertFlightsListToMaps(f FlightsList) (map[string]int, map[string]int) {
	startMap := map[string]int{}
	endMap := map[string]int{}

	for _, v := range f.Flights {
		start := strings.ToUpper(v[0])
		end := strings.ToUpper(v[1])

		if _, ok := startMap[start]; ok {
			startMap[start]++
		} else {
			startMap[start] = 1
		}

		if _, ok := endMap[end]; ok {
			endMap[end]++
		} else {
			endMap[end] = 1
		}
	}

	return startMap, endMap
}

func getPosition(mapDelta map[string]int) []string {
	positions := []string{}
	for k, v := range mapDelta {
		if v > 0 {
			positions = append(positions, k)
		}
	}
	return positions
}

func findTotalFlightPath(f FlightsList) []string {
	startMap, endMap := convertFlightsListToMaps(f)
	startMapDelta := subtructMaps(startMap, endMap)
	endMapDelta := subtructMaps(endMap, startMap)

	startPositions := getPosition(startMapDelta)
	endPositions := getPosition(endMapDelta)

	// if start and end positions are not equal to 1, there was either loop or a fork
	if len(startPositions) == 1 && len(endPositions) == 1 {
		flightPath := []string{}
		flightPath = append(flightPath, startPositions[0])
		flightPath = append(flightPath, endPositions[0])
		return flightPath
	}

	return nil
}

func validateInput(r *http.Request) (FlightsList, error) {
	var f FlightsList
	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		return f, errors.New("invalid JSON data")
	}

	if len(f.Flights) == 0 {
		return f, errors.New("flight list is empty")
	}

	for _, v := range f.Flights {
		if len(v) != 2 {
			return f, errors.New("incorrect number of airport names in one of the elements")
		}
	}

	return f, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	// validate input and get flight list from json input
	flightList, err := validateInput(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	flightPath := findTotalFlightPath(flightList)
	if flightPath != nil {
		data, err := json.Marshal(TotalFlightPath{flightPath})
		if err != nil {
			log.Fatalf("JSON marshaling failed: %s", err)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(data))
		return
	}

	http.Error(w, "corrupted flight list contains loop or fork", http.StatusBadRequest)
	return
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
