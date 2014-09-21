/*
Package regexp implements a simple library for regular expressions.
*/

package main

import (
	"./app"
	"./app/util"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"bytes"
)


func dataHandler(w http.ResponseWriter, r *http.Request, App app.App, Urls Urls) {
	// parse the request
	r.ParseForm()

	// We'll default to JSON for now...
	w.Header().Set("Content-Type", "application/json")

	// Get the graphs
	graphs := App.GetGraphsById(util.SplitQuery(Urls.QueryPrefix, r.URL.RawQuery));

	// Populate the data
	App.PopulateData(graphs)

	// Marshal the json
	json, err := json.Marshal(graphs)
	if err != nil {
		fmt.Println(err)
		return
	} 
	// Send the json over, if good!
	w.Write(json)
}


func displayHandler(w http.ResponseWriter, r *http.Request, App app.App, Urls Urls) {
	// Get the graphs
	graphs := App.GetGraphsById(util.SplitQuery(Urls.QueryPrefix, r.URL.RawQuery));

	// Populate the data
	App.PopulateData(graphs)
	
	// Anonymous struct
	data := struct {
		Header string
		Graphs []app.Graph
	} {
		string(util.GetFileContents(Urls.Html + "header.html")),
		graphs,
	}

	// Create page!
	template.Must(template.ParseFiles(
		Urls.Html + "graph-display.html")).Execute(w, data);
}


func selectHandler(w http.ResponseWriter, r *http.Request, App app.App, Urls Urls) {

	// Add graph selectors
	selTemplate, err := template.ParseFiles(Urls.Html + "graph-selection.html");
	if err != nil {
		fmt.Println(err);
		return;
	}

	// Put selector html in a string
	var selectors bytes.Buffer 
	for _, graph := range App.Graphs {
		err := selTemplate.Execute(&selectors, graph)
		if err != nil{
			fmt.Println(err);
			return;
		}
	}

	// Anonymous struct
	data := struct {
		Header string
		Selectors string
		QueryPrefix string
		DisplayUrl string
	} {
		string(util.GetFileContents(Urls.Html + "header.html")),
		selectors.String(),
		Urls.QueryPrefix,
		Urls.Display}

	// Add graph selections
	template.Must(template.ParseFiles(Urls.Html + "graph-select.html")).Execute(w, data);

}


func createApp()app.App {
	app := app.App{
	        GraphConfig: "./static/graphtypes/graphtypes.json",
		DataPath: "./data/"}
	app.Init()
	return app;
}


type Urls struct {
	Display string
	Select string
	Data string
	Html string
	QueryPrefix string
}


func main() {

	// Set urls
	urls := Urls{
		Display: "/system-performance-display/",
		Select: "/system-performance-select/",
		Data: "/system-performance-data/",
		Html: "./static/html/",
		QueryPrefix: "type"}

	// Create the app
	app := createApp();

	// Handlers
	http.HandleFunc(urls.Display, func (w http.ResponseWriter, r *http.Request){
		displayHandler(w, r, app, urls)
	})
	http.HandleFunc(urls.Data, func (w http.ResponseWriter, r *http.Request){
		dataHandler(w, r, app, urls)
	})
	http.HandleFunc(urls.Select, func (w http.ResponseWriter, r *http.Request){
		selectHandler(w, r, app, urls)
	})

	// Include the static path
	http.Handle("/static/", http.StripPrefix("/static/", 
		http.FileServer(http.Dir("static"))));

	// Lisetn!
	http.ListenAndServe (":8080", nil);
}
