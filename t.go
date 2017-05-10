// go run t.go "inputFile.csv"

package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"os/user"
	"strconv"
	s "strings"

	"fmt"

	"github.com/tealeg/xlsx"
)

func main() {
	// Store argument values passed on the command-line
	inputPath := os.Args[1]
	// outputPath := os.Args[2]

	// Load the file
	file, err := os.Open(inputPath)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new reader
	reader := csv.NewReader(bufio.NewReader(file))

	// Get all data in the reader
	dataWithHeader, err := reader.ReadAll()
	data := dataWithHeader[1:]

	// Get data summarized
	summary := GetDataSummary(data)

	// TODO get week date to add it to excel sheet

	// Create test xlsx file
	WriteListToXlxs("Week date goes here", summary, "test.xlsx")
}

// DataSummaryRecord : struct to store every record value
type DataSummaryRecord struct {
	project, activity string
	duration          float64
}

// GetDataSummary : return a table with duplicated records united
func GetDataSummary(data [][]string) []DataSummaryRecord {
	// Declare an empty slice to collect parsed records
	var dataSummary, dataTrimmed []DataSummaryRecord

	fmt.Println("\n\nTRIMMED")
	for i := range data {
		recordStruct := DataSummaryRecord{}
		recordStruct.project = data[i][3]
		recordStruct.activity = data[i][5]
		recordStruct.duration = ParseDuration(data[i][11])

		// Append parsed record to the slice
		dataTrimmed = append(dataTrimmed, recordStruct)
		fmt.Println(recordStruct)
	}

	for i := range dataTrimmed {
		recordStruct := DataSummaryRecord{}
		activityFound := false
		// Check if this activity is been already added to dataSummary
		for ii := range dataSummary {
			if dataSummary[ii].project == dataTrimmed[i].project && dataSummary[ii].activity == dataTrimmed[i].activity {
				dataSummary[ii].duration = dataSummary[ii].duration + dataTrimmed[i].duration
				activityFound = true
				break
			}
			//activityFound = false
		}
		if activityFound == false {
			recordStruct.project = dataTrimmed[i].project
			recordStruct.activity = dataTrimmed[i].activity
			recordStruct.duration = dataTrimmed[i].duration

			dataSummary = append(dataSummary, recordStruct)
		}

	}
	return dataSummary
}

// ParseDuration : return a float64 with the number of hours rounded to the nearest 0.25h
func ParseDuration(timeString string) float64 {
	if timeString == "" {
		log.Fatal("timeString = \"\"")
	}
	timeArray := s.Split(timeString, ":")
	h, _ := strconv.ParseFloat(timeArray[0], 64)
	m, _ := strconv.ParseFloat(timeArray[1], 64)
	s, _ := strconv.ParseFloat(timeArray[2], 64)
	h = h + (m / 60) + (s / 3600)

	roundh := RoundNea(h, 0.25)
	return roundh
}

// RoundNea : return "val" rounded to the nearest "nea"
func RoundNea(val, nea float64) float64 {
	a := val / nea
	b := float64(Round(a))
	return nea * b
}

// Round : float64 to int
func Round(f float64) int {
	if f < -0.5 {
		return int(f - 0.5)
	}
	if f > 0.5 {
		return int(f + 0.5)
	}
	return 0
}

// GetUserPath : Return user path if found
func GetUserPath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

// WriteListToXlxs : create a xlxs file from a list
func WriteListToXlxs(sheetName string, sheetData []DataSummaryRecord, outputPath string) {

	// Create file object
	file := xlsx.NewFile()

	// Add sheet
	sheet, err := file.AddSheet(sheetName) // TODO: sheet name = csv file name
	if err != nil {
		log.Fatal(err.Error())
	}

	// Add header to sheet
	sheet.Cell(0, 0).Value = "PROYECT"
	sheet.Cell(0, 1).Value = "ACTIVITY"
	sheet.Cell(0, 2).Value = "DURATION"
	// Populate sheet
	for i := range sheetData {
		sheet.Cell((1 + i), 0).Value = sheetData[i].project
		sheet.Cell((1 + i), 1).Value = sheetData[i].activity
		durationString := strconv.FormatFloat(sheetData[i].duration, 'f', -1, 64)
		sheet.Cell((1 + i), 2).Value = durationString
	}

	// Create excel file
	err = file.Save(outputPath) // TODO: file name = csv file name
	if err != nil {
		log.Fatal(err.Error())
	}
}
