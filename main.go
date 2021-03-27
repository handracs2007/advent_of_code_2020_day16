// Advent of code 2020 day 16: https://adventofcode.com/2020/day/16
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// YourTicket defines the indicator telling that the content after this is your Ticket details.
const YourTicket = "your ticket"

// NearbyTickets defines the indicator telling that the content after this is the nearby tickets details.
const NearbyTickets = "nearby tickets"

// ValidRange stores the valid range (minimum and maximum Values). Both inclusive.
type ValidRange struct {
	Min int
	Max int
}

// Ticket stores the Ticket details.
type Ticket struct {
	Values []int
}

// Configuration stores the Ticket Configuration.
type Configuration struct {
	Field  string
	Ranges []ValidRange
}

// parseConfiguration parses the Configuration string. It returns the Configuration object.
// We assume that the config string is always valid.
func parseConfiguration(config string) Configuration {
	// Format is <Field>: <range> [or <range]...
	// Get the Field first.
	colonIdx := strings.Index(config, ":")
	field := config[:colonIdx]

	// Get the range string, by removing everything before the range indicator.
	config = config[colonIdx+2:]

	// Separate the range by " or " separator.
	ranges := strings.Split(config, " or ")

	// Build the ValidRange for each Ranges.
	validRanges := make([]ValidRange, len(ranges))

	for idx, rng := range ranges {
		// Separate the range by "-" to get the minimum and maximum value.
		minMax := strings.Split(rng, "-")
		min, _ := strconv.Atoi(minMax[0]) // We assume that the value is always a valid integer value,hence ignoring the error handling. Don't do this in production.
		max, _ := strconv.Atoi(minMax[1]) // We assume that the value is always a valid integer value,hence ignoring the error handling. Don't do this in production.

		validRanges[idx] = ValidRange{Min: min, Max: max}
	}

	return Configuration{
		Field:  field,
		Ranges: validRanges,
	}
}

// parseTicket parses the Ticket string. It returns a Ticket object that contains
// all the Values found inside the Ticket.
// We assume that the Ticket data is always valid.
func parseTicket(ticketData string) Ticket {
	data := strings.Split(ticketData, ",")
	values := make([]int, len(data))

	for idx, datum := range data {
		values[idx], _ = strconv.Atoi(datum) // We assume that the value is always a valid integer value, hence ignoring the error handling. Don't do this in production.
	}

	return Ticket{Values: values}
}

// isValidTicket checks whether a given Ticket is valid or not based on the set of configurations.
// It returns 2 Values. First value is either true (if the Ticket is valid) or false if the Ticket is
// invalid. When the Ticket is invalid, the second return value should contain list of invalid Values,
// otherwise, it will be nil.
func isValidTicket(ticket Ticket, configs []Configuration) (bool, []int) {
	invalidValues := make([]int, 0)

	for _, value := range ticket.Values {
		foundValid := false

		for _, config := range configs {
			// Now we have the value and a config, let's check against it.
			for _, minMax := range config.Ranges {
				if value >= minMax.Min && value <= minMax.Max {
					// The value is valid
					foundValid = true
				}
			}
		}

		if !foundValid {
			invalidValues = append(invalidValues, value)
		}
	}

	if len(invalidValues) > 0 {
		return false, invalidValues
	}

	return true, nil
}

// getOrdering gets the ordering of the fields in the ticket.
func getOrdering(tickets []Ticket, configs []Configuration) []string {
	fieldSize := len(tickets[0].Values)
	orderedFields := make([]string, fieldSize)

	for len(configs) > 0 {
		// We process from the first position, second position, and so on.
		for fieldPos := 0; fieldPos < fieldSize; fieldPos++ {
			// Get all the values of the given position in all tickets.
			values := make([]int, 0)
			for _, ticket := range tickets {
				values = append(values, ticket.Values[fieldPos])
			}

			// We now check the validity of all those values against the configurations.
			// All must pass to be considered that the values belong to a field.
			validConfigCount := 0
			validConfig := Configuration{}
			validConfigIdx := -1
			for idx, config := range configs {
				isValidConfig := true

				for _, value := range values {
					// Now we have the value and a config, let's check against it.
					isValidRange := false
					for _, allowedRange := range config.Ranges {
						if value >= allowedRange.Min && value <= allowedRange.Max {
							// The value is valid
							isValidRange = true
							break
						}
					}

					if !isValidRange {
						isValidConfig = false
						break
					}
				}

				if isValidConfig {
					validConfigCount++

					validConfig = config
					validConfigIdx = idx

					if validConfigCount > 1 {
						// If the value has more than one valid config, let's skip first.
						// We must find configuration that is really unique.
						break
					}
				}
			}

			if validConfigCount == 1 {
				// The values fulfill a specific configuration.
				orderedFields[fieldPos] = validConfig.Field

				// We remove this config from the list of configurations so that
				// its not being checked in further iteration.
				configs = append(configs[:validConfigIdx], configs[validConfigIdx+1:]...)
			}
		}
	}

	return orderedFields
}

func main() {
	// Let's open the file
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("Unable to open input file. %s.", err)
	}
	defer file.Close() // Close the file

	readConfiguration := true // First reading will be the configuration.
	readYourTicket := false   // We are not reading "your ticket" details until told to.
	readNearbyTicket := false // We are not reading "nearby tickets" details until told to.

	configs := make([]Configuration, 0)
	myTicket := Ticket{}
	validNearbyTickets := make([]Ticket, 0)
	invalidValues := make([]int, 0)

	// Create a reader to read line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Ignore empty line
		if len(line) == 0 {
			continue
		}

		// Check if we are reading "your ticket" or "nearby tickets". Other than that, we are reading
		// configuration.
		if strings.HasPrefix(line, YourTicket) {
			// Reading our own ticket later. Set the flag.
			readConfiguration = false
			readYourTicket = true
			readNearbyTicket = false
		} else if strings.HasPrefix(line, NearbyTickets) {
			// Reading nearby tickets. Set the flag
			readConfiguration = false
			readYourTicket = false
			readNearbyTicket = true
		} else {
			// Reading the data and process based on the flag.
			if readConfiguration {
				// Process the configuration
				newConfig := parseConfiguration(line)
				configs = append(configs, newConfig)
			} else if readYourTicket {
				// Process our own ticket. Our own ticket is assumed to be always valid.
				myTicket = parseTicket(line)
				validNearbyTickets = append(validNearbyTickets, myTicket)
			} else if readNearbyTicket {
				// Process the nearby ticket
				nearbyTicket := parseTicket(line)

				valid, invalids := isValidTicket(nearbyTicket, configs)
				if !valid {
					invalidValues = append(invalidValues, invalids...)
				} else {
					validNearbyTickets = append(validNearbyTickets, nearbyTicket)
				}
			}
		}
	}

	sum := 0
	for _, value := range invalidValues {
		sum += value
	}
	fmt.Println(sum)

	// Part 2, determine the fields ordering.
	mul := 1
	orderedFields := getOrdering(validNearbyTickets, configs)
	for idx, field := range orderedFields {
		if strings.HasPrefix(field, "departure ") {
			mul *= myTicket.Values[idx]
		}
	}
	fmt.Println(mul)
}
