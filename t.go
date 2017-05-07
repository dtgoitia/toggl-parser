// go run t.go "inputFile.csv"

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
	s "strings"

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
	summmary := GetDataSummary(data)

	// Create test xlsx file
	WriteListToXlxs()

	fmt.Println(summmary)
}

// DataSummaryRecord : struct to store every record value
type DataSummaryRecord struct {
	project, activity string
	duration          float64
}

// GetDataSummary : return a table with duplicated records united
func GetDataSummary(data [][]string) []DataSummaryRecord {
	// Declare an empty slice to collect parsed records
	var dataSummary []DataSummaryRecord

	for i := range data {
		recordStruct := DataSummaryRecord{}
		recordStruct.project = data[i][3]
		recordStruct.activity = data[i][5]
		recordStruct.duration = ParseDuration(data[i][11])

		// Append parsed record to the slice
		dataSummary = append(dataSummary, recordStruct)
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
func WriteListToXlxs() {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1") // TODO: sheet name = csv file name
	if err != nil {
		log.Fatal(err.Error())
	}
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "I am a cell!"
	err = file.Save("test.xlsx") // TODO: file name = csv file name
	if err != nil {
		log.Fatal(err.Error())
	}
}