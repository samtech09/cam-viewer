package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

// GLOBALS DECLARED HERE

var extensionList = [][]byte{
	{'.', 'm', 'k', 'v'},
	{'.', 'm', 'p', 'g'},
	{'.', 'a', 'v', 'i'},
	{'.', 'm', '4', 'v'},
	{'.', 'm', 'p', '4'},
	{'.', 'd', 'a', 't'}}

var config = ConfigDetails{}
var appConfig = AppConfig{}

//var templates map[string]*template.Template
var templates *template.Template
var pageData = PageData{}
var currentPath string
var currentBase string
var allMovies map[int64]Movie
var contexts map[string]struct{}

// THE VIEW CODE IS HERE

func generateTemplates() {
	// templates = make(map[string]*template.Template)
	// modulus := template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }}
	// for _, tmpl := range config.templateFileList {
	// 	t := template.New("base.html").Funcs(modulus)
	// 	templates[tmpl] = template.Must(t.ParseFiles(config.templateDirectory+"base.html", config.templateDirectory+tmpl))
	// }

	// templates["play.html"] = template.Must(template.ParseFiles(config.templateDirectory + "play.html"))
	// templates["index.html"] = template.Must(template.ParseFiles(config.templateDirectory + "index.html"))
	// templates["about.html"] = template.Must(template.ParseFiles(config.templateDirectory + "about.html"))
	// templates["cams.html"] = template.Must(template.ParseFiles(config.templateDirectory + "cams.html"))
	// templates["menu.html"] = template.Must(template.ParseFiles(config.templateDirectory + "menu.html"))
	// templates["main.html"] = template.Must(template.ParseFiles(config.templateDirectory + "main.html"))

	// compile all templates and cache them
	templates = template.Must(template.ParseGlob(config.templateDirectory + "*"))
}

func renderTemplate(pageStruct interface{}, w http.ResponseWriter, tmpl string) {
	var err error
	if pageStruct == nil {
		//err = templates[tmpl].Execute(w, pageData)
		err = templates.ExecuteTemplate(w, tmpl, pageData)
	} else {
		//err = templates[tmpl].Execute(w, pageStruct)
		err = templates.ExecuteTemplate(w, tmpl, pageStruct)
	}
	if err != nil {
		log.Printf("The follwing error occurred: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// IT ALL STARTS HERE

func initConfigDetails() {
	// config.templateFileList = append(config.templateFileList,
	// 	"index.html", "about.html", "movie.html", "alreadyplaying.html", "setup.html", "nothingfound.html")

	conffile := "config.json"
	file, err := os.Open(conffile)
	if err != nil {
		log.Fatal("Missing config file.\n", err)
	}

	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&appConfig)
	if err != nil {
		log.Fatal("Cannot parse config.\n", err)
	}

	config.firstStart = true
	config.Playmode = "keyboard"
	config.templateDirectory = "./templates/"
	config.basePath = "/videos/"
	pageData.CurrentContext = appConfig.DefaultCam // "porch" // initially display recording for this cam

}

func initPaths() {
	config.FilePathList = append(config.FilePathList, "static/videos/")

	if len(config.FilePathList) > 0 {
		config.firstStart = false
	}
}

func main() {
	initConfigDetails()
	generateTemplates()
	initPaths()
	refreshList("d")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/menu", menuHandler)
	http.HandleFunc("/main", mainHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/play", videoHandler)
	http.HandleFunc("/cams", camsHandler)
	http.HandleFunc("/similar-time", sameTimeHandler)
	http.HandleFunc("/sys-stats", statsHandler)
	//http.HandleFunc("/setup", setupHandler)
	//http.HandleFunc("/playmode", playmodeHandler)
	//http.HandleFunc("/movie", movieHandler)

	// getCPU()
	// getMEM()
	// getSWAP()
	// getDisk()

	fmt.Println("Server started, browse http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

//Initconfig read config file and initialize AppConfig struct
func initConfig() {

}
