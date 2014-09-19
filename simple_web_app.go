package main

import (
	//"io/ioutil"
	"fmt"
	"net/http"
	"text/template"
)


const TOTAL_GRAPHS int = 6


type Graph struct {
    Data [][]int
    Label string
    Color string
    IconSrc string
}


//
// From: https://groups.google.com/forum/#!topic/golang-nuts/rXHT2fJ0hG8
//
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


func handler(w http.ResponseWriter, r *http.Request) {

	var twoD = create2dslice(20, 2)
	var graph = &Graph{ 
		Data: twoD, 
		Label: "Network", 
		IconSrc: "/static/icons/network_icon.png", 
		Color: "rgb(240,240,0)"}
	fmt.Println(graph)

	template.Must(template.ParseFiles("index.html")).Execute(w, graph)
	//t, _ := template.ParseFiles("graph.html")
	//graph, _ := ioutil.ReadFile("graph.html")
	//graphStr := string(graph)
	//fmt.Println(graphStr)
	//gParse, a := template.ParseFiles("graph.html")
	selectorParse, a := template.ParseFiles("graph_selector.html");



	for i := 0; i < TOTAL_GRAPHS; i++ {
		template.Must(selectorParse, a).Execute(w, graph)
		//template.Must(gParse, a).Execute(w, graph)
	}
	
	//fmt.Fprintf(w, graphStr)
	//fmt.Fprintf(w, "<div class=\"ig-checkbox\">%s</div>", "hello")
	//fmt.Fprintf(w, "<div class=\"ig-checkbox\">%s</div>", "hello")
	//fmt.Fprintf(w, "<div class=\"ig-checkbox\">%s</div>", "hello")
}

func main() {
	http.HandleFunc("/", handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	//http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.ListenAndServe (":8080", nil);
}
