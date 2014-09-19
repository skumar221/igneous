package main

import (
	//"io/ioutil"
	"fmt"
	"net/http"
	"text/template"
	"strings"
)


const STATIC_PATH string = "./static/"
const STATIC_HTML_PATH string = STATIC_PATH + "html/"
const STATIC_IMAGE_PATH string = STATIC_PATH + "images/"
const STATIC_ICON_PATH string = STATIC_IMAGE_PATH + "/icons/"



type AppUrlFragments struct {
    GraphDisplay string
    GraphSelect string
}


type Graph struct {
    Data [][]int
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
		j++
		
	} 
	return _2d
}



func graphDisplayHandler(w http.ResponseWriter, r *http.Request) {
	graphTypes := strings.Split(r.URL.RawQuery[len("types="):], "&")
	fmt.Println(graphTypes)
	template.Must(template.ParseFiles(STATIC_HTML_PATH + "graph-display.html")).Execute(w, nil)

	//selectorParse, a := template.ParseFiles(STATIC_HTML_PATH + "graph-display.html");

	/*
	template.Must(template.ParseFiles(STATIC_HTML_PATH + "graph-select.html")).Execute(w, nil)
	for i := 0; i < len(Graphs); i++ {
		template.Must(selectorParse, a).Execute(w, Graphs[i])
	}
*/
}



func graphSelectHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Graph Select")
	selectorParse, a := template.ParseFiles(STATIC_HTML_PATH + "graph-selection.html");
	template.Must(template.ParseFiles(STATIC_HTML_PATH + "graph-select.html")).Execute(w, AppUrlFrags)
	for i := 0; i < len(Graphs); i++ {
		template.Must(selectorParse, a).Execute(w, Graphs[i])
	}
}



var twoD = create2dslice(20, 2)
var networkGraph = &Graph{ 
	Data: twoD, 
	Query: "network", 
	Label: "Network Traffic", 
	IconSrc: STATIC_ICON_PATH[1:] + "network.png", 
	Color: "rgb(240,240,0)"}
var memoryGraph = &Graph{ 
	Data: twoD, 
	Query: "memory", 
	Label: "Memory", 
	IconSrc: STATIC_ICON_PATH[1:] + "memory.png", 
	Color: "rgb(240,240,0)"}
var energyGraph = &Graph{ 
	Data: twoD, 
	Query: "energy", 
	Label: "Energy Consumption", 
	IconSrc: STATIC_ICON_PATH[1:] + "energy.png", 
	Color: "rgb(240,240,0)"}
var cpuGraph = &Graph{ 
	Data: twoD, 
	Query: "cpu", 
	Label: "CPU", 
	IconSrc: STATIC_ICON_PATH[1:] + "cpu.png", 
	Color: "rgb(240,240,0)"}
var diskCapacityGraph = &Graph{ 
	Data: twoD, 
	Query: "disk", 
	Label: "Disk Capacity", 
	IconSrc: STATIC_ICON_PATH[1:] + "diskCapacity.png", 
	Color: "rgb(240,240,0)"}
var temperaturesGraph = &Graph{ 
	Data: twoD, 
	Query: "temps", 
	Label: "Temperatures", 
	IconSrc: STATIC_ICON_PATH[1:] + "temperatures.png", 
	Color: "rgb(240,240,0)"}
var Graphs = make([]Graph, 6)



var AppUrlFrags = &AppUrlFragments{
	GraphDisplay: "/graph-display/",
	GraphSelect: "/graph-select/"}


func main() {

	Graphs[0] = * networkGraph
	Graphs[1] = * memoryGraph
	Graphs[2] = * energyGraph
	Graphs[3] = * cpuGraph
	Graphs[4] = * diskCapacityGraph
	Graphs[5] = * temperaturesGraph


	http.HandleFunc(AppUrlFrags.GraphSelect, graphSelectHandler)
	http.HandleFunc(AppUrlFrags.GraphDisplay, graphDisplayHandler)

	http.Handle("/static/", http.StripPrefix("/static/", 
		http.FileServer(http.Dir("static"))));
	http.ListenAndServe (":8080", nil);
}
