// go run t.go "inputFile.csv" "test.xlsx"
// go run t.go "inputFile.csv" "test.xlsx" -v

package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"os/user"
	"strconv"
	s "strings"
	"text/tabwriter"

	"fmt"
	"time"

	"github.com/tealeg/xlsx"
)

func main() {
	// Load the file
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// Create a new reader
	fmt.Println("Openning CSV file...")
	reader := csv.NewReader(bufio.NewReader(file))

	// Extract CSV data
	fmt.Println("Extracting CSV file data...")
	dataWithHeader, err := reader.ReadAll()
	data := dataWithHeader[1:]

	// Summarize data
	fmt.Println("Summarising data...")
	summaryDataSummaryRecord := GetDataSummary(data)

	// Joining duplicated activity duration
	fmt.Println("Joining duplicated activity duration...")
	summaryDataToExport := GetDataToExport(summaryDataSummaryRecord)
	PrintDataToExport(summaryDataToExport)

	// Create test xlsx file
	fmt.Println("Writing data to", os.Args[2], "...")
	WriteDataToExportsToXlxs("Week date goes here", summaryDataToExport, os.Args[2])

}

// DataSummaryRecord : struct to store every record value
type DataSummaryRecord struct {
	project  string
	activity string
	date     time.Time
	duration float64
}

// DataToExport : struct to store week activity durations
type DataToExport struct {
	project, activity                                                                                                      string
	mondayDuration, tuesdayDuration, wednesdayDuration, thursdayDuration, fridayDuration, saturdayDuration, sundayDuration float64
}

// GetDataSummary : return a table with duplicated records united
func GetDataSummary(data [][]string) []DataSummaryRecord {
	// Declare an empty slice to collect parsed records
	var dataSummary, dataTrimmed []DataSummaryRecord

	fmt.Println("Trimming data...")
	for i := range data {
		record := DataSummaryRecord{}
		record.project = data[i][3]
		record.activity = data[i][5]
		t, _ := time.Parse("2006-01-02", data[i][7])
		record.date = t
		record.duration = ParseDuration(data[i][11])

		// Append parsed record to the slice
		dataTrimmed = append(dataTrimmed, record)
	}

	fmt.Println("Summarising data... ")
	for i := range dataTrimmed {
		record := DataSummaryRecord{}
		activityFound := false
		// Check if this activity is been already added to dataSummary
		for ii := range dataSummary {
			if dataSummary[ii].project == dataTrimmed[i].project && dataSummary[ii].activity == dataTrimmed[i].activity && dataSummary[ii].date == dataTrimmed[i].date {
				dataSummary[ii].duration = dataSummary[ii].duration + dataTrimmed[i].duration
				activityFound = true
				break
			}
		}
		if activityFound == false {
			record.project = dataTrimmed[i].project
			record.activity = dataTrimmed[i].activity
			record.duration = dataTrimmed[i].duration
			record.date = dataTrimmed[i].date

			dataSummary = append(dataSummary, record)
		}

	}

	return dataSummary
}

// GetDataToExport : return a table with information sorted ready to export
func GetDataToExport(data []DataSummaryRecord) []DataToExport {
	var dataExport []DataToExport
	for i := range data {
		x := DataToExport{}
		activityFound := false
		// Check if this activity has been already added to dataExport
		for ii := range dataExport {
			if dataExport[ii].project == data[i].project && dataExport[ii].activity == data[i].activity {
				// edit existing tasks to add values
				switch weekday := data[i].date.Weekday().String(); weekday {
				case "Monday":
					if dataExport[ii].mondayDuration != 0 {
						dataExport[ii].mondayDuration = dataExport[ii].mondayDuration + data[i].duration
					} else {
						dataExport[ii].mondayDuration = data[i].duration
					}
				case "Tuesday":
					if dataExport[ii].tuesdayDuration != 0 {
						dataExport[ii].tuesdayDuration = dataExport[ii].tuesdayDuration + data[i].duration
					} else {
						dataExport[ii].tuesdayDuration = data[i].duration
					}
				case "Wednesday":
					if dataExport[ii].tuesdayDuration != 0 {
						dataExport[ii].wednesdayDuration = dataExport[ii].wednesdayDuration + data[i].duration
					} else {
						dataExport[ii].wednesdayDuration = data[i].duration
					}
				case "Thursday":
					if dataExport[ii].thursdayDuration != 0 {
						dataExport[ii].thursdayDuration = dataExport[ii].thursdayDuration + data[i].duration
					} else {
						dataExport[ii].thursdayDuration = data[i].duration
					}
				case "Friday":
					if dataExport[ii].fridayDuration != 0 {
						dataExport[ii].fridayDuration = dataExport[ii].fridayDuration + data[i].duration
					} else {
						dataExport[ii].fridayDuration = data[i].duration
					}
				case "Saturday":
					if dataExport[ii].saturdayDuration != 0 {
						dataExport[ii].saturdayDuration = dataExport[ii].saturdayDuration + data[i].duration
					} else {
						dataExport[ii].saturdayDuration = data[i].duration
					}
				case "Sunday":
					if dataExport[ii].sundayDuration != 0 {
						dataExport[ii].sundayDuration = dataExport[ii].sundayDuration + data[i].duration
					} else {
						dataExport[ii].sundayDuration = data[i].duration
					}
				}

				activityFound = true
				break
			}
		}
		if activityFound == false {
			// add new task
			x.project = data[i].project
			x.activity = data[i].activity
			switch weekday := data[i].date.Weekday().String(); weekday {
			case "Monday":
				x.mondayDuration = data[i].duration
			case "Tuesday":
				x.tuesdayDuration = data[i].duration
			case "Wednesday":
				x.wednesdayDuration = data[i].duration
			case "Thursday":
				x.thursdayDuration = data[i].duration
			case "Friday":
				x.fridayDuration = data[i].duration
			case "Saturday":
				x.saturdayDuration = data[i].duration
			case "Sunday":
				x.sundayDuration = data[i].duration
			}
			dataExport = append(dataExport, x)
		}
	}
	return dataExport
}

// PrintDataToExport : print data table
func PrintDataToExport(data []DataToExport) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	for i := range data {
		// Concatenata all data in a single string
		var s string
		s = data[i].project + "\t" + data[i].activity + "\t"
		s = s + strconv.FormatFloat(data[i].mondayDuration, 'f', -1, 64) + "\t"
		s = s + strconv.FormatFloat(data[i].tuesdayDuration, 'f', -1, 64) + "\t"
		s = s + strconv.FormatFloat(data[i].wednesdayDuration, 'f', -1, 64) + "\t"
		s = s + strconv.FormatFloat(data[i].thursdayDuration, 'f', -1, 64) + "\t"
		s = s + strconv.FormatFloat(data[i].fridayDuration, 'f', -1, 64) + "\t"
		s = s + strconv.FormatFloat(data[i].saturdayDuration, 'f', -1, 64) + "\t"
		s = s + strconv.FormatFloat(data[i].sundayDuration, 'f', -1, 64)

		// Print contenated string
		fmt.Fprintln(w, s)
	}
	w.Flush()
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

// WriteDataToExportsToXlxs : create a xlxs file from a list
func WriteDataToExportsToXlxs(sheetName string, sheetData []DataToExport, outputPath string) {

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
	sheet.Cell(0, 2).Value = "DATE"
	sheet.Cell(0, 3).Value = "MON"
	sheet.Cell(0, 4).Value = "TUE"
	sheet.Cell(0, 5).Value = "WED"
	sheet.Cell(0, 6).Value = "THU"
	sheet.Cell(0, 7).Value = "FRI"

	// Populate sheet
	for i := range sheetData {
		sheet.Cell((1 + i), 0).Value = sheetData[i].project
		sheet.Cell((1 + i), 1).Value = sheetData[i].activity
		sheet.Cell((1 + i), 2).Value = strconv.FormatFloat(sheetData[i].mondayDuration, 'f', -1, 64)
		sheet.Cell((1 + i), 3).Value = strconv.FormatFloat(sheetData[i].tuesdayDuration, 'f', -1, 64)
		sheet.Cell((1 + i), 4).Value = strconv.FormatFloat(sheetData[i].wednesdayDuration, 'f', -1, 64)
		sheet.Cell((1 + i), 5).Value = strconv.FormatFloat(sheetData[i].thursdayDuration, 'f', -1, 64)
		sheet.Cell((1 + i), 6).Value = strconv.FormatFloat(sheetData[i].fridayDuration, 'f', -1, 64)
		sheet.Cell((1 + i), 7).Value = strconv.FormatFloat(sheetData[i].saturdayDuration, 'f', -1, 64)
		//fmt.Println(sheetData[i])
	}

	// Create excel file
	err = file.Save(outputPath) // TODO: file name = csv file name
	if err != nil {
		log.Fatal(err.Error())
	}
}
