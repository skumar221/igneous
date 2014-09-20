package main

import (
	//"io/ioutil"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"strings"
)

const SP_QPREFIX = "kinds"
const STATIC_PATH string = "./static/"
const STATIC_HTML_PATH string = STATIC_PATH + "html/"
const STATIC_IMAGE_PATH string = STATIC_PATH + "images/"
const STATIC_ICON_PATH string = STATIC_IMAGE_PATH + "icons/"



type AppUrlFragments struct {
    Display string
    Select string
    Data string
    QueryPrefix string
}


type Graph struct {
    WeeklyData [][]int
    HourlyData [][]int
    Query string
    Label string
    Color string
    IconSrc string
}



func create2dslice(dimensionX, dimensionY int) [][]int { 
	_2d := make([][]int, dimensionX) 
	j := 0 
	for i := 0; i < dimensionX; i++ { 
		_2d[i] = make([]int, dimensionY)
		_2d[i][0] = j
		_2d[i][1] = j * 2
		j++
		
	} 
	return _2d
}


type Response map[string]interface{}

func graphDataHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// We'll default to JSON for now...
	w.Header().Set("Content-Type", "application/json")

	// Make the graphs slice to send over as JSON
	graphs := make([]Graph, 0, len(Graphs))

	// If there's a query string, return on the the selected graphs,
	// otherwise return all
	//fmt.Println(r.URL.RawQuery)
	//fmt.Println(strings.Contains(r.URL.RawQuery, SP_QPREFIX + "="))
	if strings.Contains(r.URL.RawQuery, SP_QPREFIX + "=") {
		graphTypes := strings.Split(r.URL.RawQuery[len(SP_QPREFIX + "="):], "&")
		fmt.Println(graphTypes)
		for i := 0; i < len(Graphs); i++ {
			for j := 0; j<len(graphTypes); j++ {
				if (Graphs[i].Query == graphTypes[j]){
					graphs = append(graphs, Graphs[i])
				}
			}
		}
	} else {
		for i := 0; i < len(Graphs); i++ {
			graphs = append(graphs, Graphs[i])
		}
	}


	// Marshal the json
	json, _ := json.Marshal(graphs)

	// Send the json over!
	w.Write(json)
}


func graphDisplayHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	template.Must(template.ParseFiles(
		STATIC_HTML_PATH + "graph-display.html")).Execute(w, nil);

	/*
	selectorParse, a := template.ParseFiles(STATIC_HTML_PATH + "graph.html");
	for i := 0; i < len(Graphs); i++ {
		for j := 0; j<len(graphTypes); j++ {
			if (Graphs[i].Query == graphTypes[j]){
				template.Must(selectorParse, a).Execute(w, Graphs[i])
			}
		}
	}
*/

}



func graphSelectHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Graph Select")
	selectorParse, a := template.ParseFiles(
		STATIC_HTML_PATH + "graph-selection.html");
	template.Must(template.ParseFiles(
		STATIC_HTML_PATH + "graph-select.html")).Execute(w, AppUrlFrags);
	for i := 0; i < len(Graphs); i++ {
		template.Must(selectorParse, a).Execute(w, Graphs[i])
	}
}



var twoD = create2dslice(20, 2)
var networkGraph = &Graph{ 
	WeeklyData: twoD, 
	HourlyData: twoD, 
	Query: "network", 
	Label: "Network Traffic", 
	IconSrc: STATIC_ICON_PATH[1:] + "network.png", 
	Color: "rgb(240,240,0)"}
var memoryGraph = &Graph{ 
	WeeklyData: twoD, 
	HourlyData: twoD, 
	Query: "memory",
	Label: "Memory", 
	IconSrc: STATIC_ICON_PATH[1:] + "memory.png", 
	Color: "rgb(150,240,200)"}
var energyGraph = &Graph{ 
	WeeklyData: twoD,
	HourlyData: twoD,  
	Query: "energy", 
	Label: "Energy Consumption", 
	IconSrc: STATIC_ICON_PATH[1:] + "energy.png", 
	Color: "rgb(240,240,0)"}
var cpuGraph = &Graph{ 
	WeeklyData: twoD, 
	HourlyData: twoD, 
	Query: "cpu", 
	Label: "CPU", 
	IconSrc: STATIC_ICON_PATH[1:] + "cpu.png", 
	Color: "rgb(240,240,0)"}
var diskCapacityGraph = &Graph{ 
	WeeklyData: twoD, 
	Query: "disk", 
	Label: "Disk Capacity", 
	IconSrc: STATIC_ICON_PATH[1:] + "diskCapacity.png", 
	Color: "rgb(240,240,0)"}
var temperaturesGraph = &Graph{ 
	WeeklyData: twoD, 
	HourlyData: twoD, 
	Query: "temps", 
	Label: "Temperatures", 
	IconSrc: STATIC_ICON_PATH[1:] + "temperatures.png", 
	Color: "rgb(240,240,0)"}
var Graphs = make([]Graph, 6)



var AppUrlFrags = &AppUrlFragments{
	Display: "/system-performance-display/",
	Select: "/system-performance-select/",
	Data: "/system-performance-data/",
        QueryPrefix: SP_QPREFIX}


func main() {	
	Graphs[0] = * networkGraph
	Graphs[1] = * memoryGraph
	Graphs[2] = * energyGraph
	Graphs[3] = * cpuGraph
	Graphs[4] = * diskCapacityGraph
	Graphs[5] = * temperaturesGraph

	http.HandleFunc(AppUrlFrags.Select, graphSelectHandler)
	http.HandleFunc(AppUrlFrags.Display, graphDisplayHandler)
	http.HandleFunc(AppUrlFrags.Data, graphDataHandler)

	//
	// Include the static path
	//
	http.Handle("/static/", http.StripPrefix("/static/", 
		http.FileServer(http.Dir("static"))));
	http.ListenAndServe (":8080", nil);
}
