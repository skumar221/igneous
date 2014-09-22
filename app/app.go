/*
Package regexp implements a simple library for regular expressions.
*/

package app

import (
	"encoding/json"
	"./util"
	"fmt"
)

type App struct {
	MaxDataPoints int
	GraphConfig string
	DataPath string
	Graphs map[string]Graph
}

type Graph struct {
	WeeklyData [][]float64
	HourlyData [][]float64
	Query string
	Label string
	Color string
	IconSrc string
	Unit string
}




func (app * App) Init(){
	app.MaxDataPoints = 20000
	app.Graphs = app.graphsFromJson(app.GraphConfig);
}



//
// Gets graphs
//
func (app App) PopulateData(graphs []Graph){
	for k, _ := range graphs {
		pref := app.DataPath + graphs[k].Query
		graphs[k].WeeklyData = util.CsvToTwoD(pref + "-week.csv");
		graphs[k].HourlyData = util.CsvToTwoD(pref + "-hour.csv");
	}	
}



func (app App) getSortedGraphIds()([]string){
	// Sort the keys
	var keys []string
	for k, _ := range app.Graphs {
		keys = append(keys, k)
	}
	return util.SortStringArray(keys)
}



//
// Gets graphs
//
func (app App) GetGraphs()([]Graph) {
	// Sort the keys
	keys := app.getSortedGraphIds()

	// Return the graphs
	graphs := make([]Graph, 0, len(app.Graphs))
	for _, key := range keys {
		graphs = append(graphs, app.Graphs[key])
	}
	return graphs
}


//
// Gets graphs by provided ids
//
func (app App) GetGraphsById(keys []string)([]Graph) {
	graphs := make([]Graph, 0, len(app.Graphs))
	if len(keys) == 0 {
		graphs = app.GetGraphs()
	} else {	
		keys := util.SortStringArray(keys)
		for _, key := range keys{
			if _, ok := app.Graphs[key]; ok {
				graphs = append(graphs, app.Graphs[key])
			}
		}
	}
	return graphs	
}


//
// Generates the Graphs map from the config file 
//
func (app App) graphsFromJson(jsonFilename string)(map[string]Graph) {
	// Make empty graph map
	var Graphs = make(map[string]Graph)

	// Begin the unmarshaling...
	var data map[string]interface{}
	json.Unmarshal(util.GetFileContents(jsonFilename), &data)
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
	fmt.Println("")
	return Graphs
}
