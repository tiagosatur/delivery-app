package route

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
)

type Position struct {
	Lat float64 `json:lat`
	Long float64 `json:long`
}

// These template
type PartialRoutePosition struct {
	ID string `json:routeId`
	ClientID string `json:clientId`
	Position []float64 `position`
	Finished bool `finished`
}

type Route struct {
	ID string `json:routeId`
	ClientID string `json:clientid`
	Positions 	[]Position `json:positions`
}

// If there is an error, returns it
func(r *Route) LoadPositions() error {
	if r.ID == "" {
		return errors.New("route id required")
	}
	
	file, err := os.Open("destinations/" + r.ID + ".txt")

	if err != nil { 
		return err // return error if there is one
	}
	
	defer file.Close() // Close file after everything

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), ",")
		lat, err := strconv.ParseFloat(data[0], 64)
		if err != nil {
			return err
		}

		long, err := strconv.ParseFloat(data[1], 64)
		if err != nil {
			return err
		}

		// Add lat and long inside our positions list
		r.Positions = append(r.Positions, Position{
			Lat: lat,
			Long: long,
		})
	}

	return nil
}

// We gonna send to kafka the vehicle position after certain intervals
// ExportJsonPositinos returns a string list or error
func (r *Route) ExportJsonPositinos() ([]string, error) {
	var route PartialRoutePosition
	var result []string
	totalPositions := len(r.Positions)

	for k, v := range r.Positions {
		route.ID = r.ID
		route.ClientID = r.ClientID
		route.Position = []float64{v.Lat, v.Long}
		route.Finished = false

		if totalPositions - 1 == k {
			route.Finished = true
		}

		// byte slice
		jsonRoute, err := json.Marshal(route)

		if err != nil {
			return nil, err // return a list of empty strings and the error
		}

		// converting byte slice to strings
		result = append(result, string(jsonRoute))
	}

	// return the result and the empty error
	return result, nil 

}