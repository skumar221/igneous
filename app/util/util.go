/*
Package regexp implements a simple library for regular expressions.
*/

package util

import (
	"os"
	"bytes"
	"io"
	"fmt"
	"encoding/csv"
	"strconv"
	"strings"
)


func SplitQuery(prefix string, queryStr string)[]string{
	values := make([]string, 0);
	if strings.Contains(queryStr, prefix + "=") {
		values = strings.Split(queryStr[len(prefix + "="):], "&")	
	} 
	return values
}



//
// Reads a file WARNING: Does not close the file!
//
func ReadFile(filename string)*os.File {
	// Open the csv file
	file, error := os.Open(filename)
	if error != nil {
		fmt.Println("Error:", error)
		return nil
	}
	return file
}


//
// Reads a file, provided by file name
// returns the 
//
func GetFileContents(filename string)[]byte {
	buf := bytes.NewBuffer(nil)
	f := ReadFile(filename) 
	if f == nil {
		fmt.Println("Error reading file:", filename)
		return nil;
	}
	_, err := io.Copy(buf, f)  
	if err != nil {
		fmt.Println("Error reading file:", filename)
		return nil;
	}
	f.Close()	
	return buf.Bytes()
}




func CsvToTwoD(filename string) [][]float64 {

	// The twod array
	_2d := make([][]float64, 0) 

	// Read the csv, defer close
	file := ReadFile(filename)
	defer file.Close()

	// parse the csv line by line...
	reader := csv.NewReader(file)
	reader.Comma = ','
	lineCount := 0
	for {
		// read just one record, but we could ReadAll() as well
		record, error := reader.Read()
		// end-of-file is fitted into error
		if error == io.EOF {
			break
		} else if error != nil {
			fmt.Println("Error:", error)
			return nil
		}
		// Populate array, skip csv header 	
		if lineCount > 0 {
			pt := make([]float64, 2)
			val1, err1:= strconv.ParseFloat(record[0], 64)
			val2, err2 := strconv.ParseFloat(record[1], 64)
			if err1 == nil && err2 == nil {
				pt[0] = val1;
				pt[1] = val2;
			} else {
				fmt.Println("Errors found converting to 2d array:", err1, err2)
				return nil;
			}
			_2d = append(_2d, pt);
		}
		lineCount += 1
	}

	return _2d;
}
