package main

import (
	"os"
	"io"
	"io/ioutil"
	"strconv"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"strings"
)


const MAX_DATA_POINTS = 1000
const SP_QPREFIX = "type"
const STATIC_PATH string = "./static/"
const STATIC_HTML_PATH string = STATIC_PATH + "html/"
var Graphs map[string]Graph


type AppUrlFragments struct {
    Display string
    Select string
    Data string
    QueryPrefix string
}

var FRAGS = &AppUrlFragments{
	Display: "/system-performance-display/",
	Select: "/system-performance-select/",
	Data: "/system-performance-data/",
        QueryPrefix: SP_QPREFIX}


type Graph struct {
    WeeklyData [][]float64
    HourlyData [][]float64
    Query string
    Label string
    Color string
    IconSrc string
    Unit string
}



func graphDataHandler(w http.ResponseWriter, r *http.Request) {
	// parse the request
	r.ParseForm()

	// We'll default to JSON for now...
	w.Header().Set("Content-Type", "application/json")

	// Get the graphs
	graphs := graphsFromQuery(r.URL.RawQuery, true);

	// Marshal the json
	json, err := json.Marshal(graphs)
	if err != nil {
		fmt.Println(err)
		return
	} 
	// Send the json over if good!
	w.Write(json)
}




func populateDataPoints(graphs []Graph){
	for k, _ := range graphs {
		pref := "./static/data/" + graphs[k].Query
		graphs[k].WeeklyData = csvToTwoD(pref + "-week.csv");
		graphs[k].HourlyData = csvToTwoD(pref + "-hour.csv");
	}	
}


func graphsFromQuery(queryStr string, populateData bool)[]Graph{
	// Construct graphs slice
	graphs := make([]Graph, 0, len(Graphs))
	// Break apart query if it exists
	if strings.Contains(queryStr, SP_QPREFIX + "=") {
		graphTypes := strings.Split(queryStr[len(SP_QPREFIX + "="):], "&")
		for j := 0; j<len(graphTypes); j++ {
			if _,ok := Graphs[graphTypes[j]]; ok {
				graphs = append(graphs, Graphs[graphTypes[j]])
			}
		}
	// Otherwise use all graphs	
	} else {
		for _, v := range Graphs {
			graphs = append(graphs, v)
		}
	}
	// Populate the data, if specified
	if populateData == true {
		populateDataPoints(graphs)
	}
	return graphs
}



func graphDisplayHandler(w http.ResponseWriter, r *http.Request) {
	// Get the graphs specified, populate data
	graphs := graphsFromQuery(r.URL.RawQuery, true);
	// Add displayer page
	template.Must(template.ParseFiles(
		STATIC_HTML_PATH + "graph-display.html")).Execute(w, graphs);
	// Add header
	template.Must(template.ParseFiles(
		STATIC_HTML_PATH + "header.html")).
			Execute(w, map[string]string{"Header": "_system performance"});
}



func graphSelectHandler(w http.ResponseWriter, r *http.Request) {
	// Create the selector parser
	selectorParse, a := template.ParseFiles(
		STATIC_HTML_PATH + "graph-selection.html");
	// Add graph selections
	template.Must(template.ParseFiles(
		STATIC_HTML_PATH + "graph-select.html")).Execute(w, FRAGS);

	// Add header
	template.Must(template.ParseFiles(
		STATIC_HTML_PATH + "header.html")).
			Execute(w, map[string]string{"Header": "_performance select"});
	// Add graph selectors
	for _,v := range Graphs {
		template.Must(selectorParse, a).Execute(w, v)
	}
}



func readFile(filename string)*os.File {
	// Open the csv file
	file, error := os.Open(filename)
	if error != nil {
		fmt.Println("Error:", error)
		return nil
	}
	return file
}


func csvToTwoD(filename string) [][]float64 {

	// The twod array
	_2d := make([][]float64, 0, MAX_DATA_POINTS) 

	// Read the csv, defer close
	file := readFile(filename)
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
			}
			_2d = append(_2d, pt);
		}
		lineCount += 1
	}

	return _2d;
}



func getGraphTypesFromJson(jsonFilename string)(map[string]Graph) {
	// empty graph map
	var Graphs = make(map[string]Graph)
	// Read the json file
	content, err := ioutil.ReadFile(jsonFilename)
	if err!=nil{
		fmt.Print("Error:",err)
	}
	// Begin the unmarshaling...
	var data map[string]interface{}
	json.Unmarshal(content, &data)
	for _, v := range data {
		// Continue unwrapping into individual graphs
		strs := v.(map[string]interface{})
		aG := Graph{} 
		for k, a := range strs {
			switch k {
			case "Query": aG.Query = a.(string);
			case "IconSrc": aG.IconSrc = a.(string);
			case "Label": aG.Label = a.(string);
			case "Color": aG.Color = a.(string);
			case "Unit": aG.Unit = a.(string);	
			}
		}
		Graphs[aG.Query] = aG;
	}
	return Graphs
}




func main() {
	// Get the graph types
	Graphs = getGraphTypesFromJson("./static/graphtypes/graphtypes.json")
	// Handlers
	http.HandleFunc(FRAGS.Select, graphSelectHandler)
	http.HandleFunc(FRAGS.Display, graphDisplayHandler)
	http.HandleFunc(FRAGS.Data, graphDataHandler)
	// Include the static path
	http.Handle("/static/", http.StripPrefix("/static/", 
		http.FileServer(http.Dir("static"))));
	http.ListenAndServe (":8080", nil);
}
