package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
)

// HANDLERS ARE HERE

func mainHandler(w http.ResponseWriter, r *http.Request) {
	// if config.firstStart {
	// 	http.Redirect(w, r, "/setup", http.StatusFound)
	// 	return
	// }
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	search := r.URL.Query().Get("d") + ""
	usetime := r.URL.Query().Get("usetime") + ""
	tfrom := r.URL.Query().Get("timefrom") + ""
	tto := r.URL.Query().Get("timeto") + ""
	//cxt := r.URL.Query().Get("context") + ""
	cxt := r.URL.Query().Get("cxt") + ""
	if cxt != "" {
		pageData.CurrentContext = cxt
	}

	// parse numbers to int
	t_from, _ := strconv.Atoi(tfrom)
	t_to, _ := strconv.Atoi(tto)

	//fmt.Println("Context: ", pageData.CurrentContext)

	tmpl := "mainView"

	if r.Method == "POST" {
		refreshList(search)
	} else if search != "" {
		if usetime == "1" {
			pageData.MovieList = searchList(search, cxt, t_from, t_to)
		} else {
			pageData.MovieList = searchList(search, cxt, 0, 0)
		}
	}
	renderTemplate(nil, w, tmpl)
}

func sameTimeHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: do search
	cat := r.URL.Query().Get("cxt")
	dt := r.URL.Query().Get("dt")
	tm := r.URL.Query().Get("tm")
	template := "similarTimeView"

	tmp, _ := strconv.Atoi(tm)
	tmp = int(math.Floor(float64(tmp) / 100.0))
	// if tmp > 999 {
	// 	tmp = tmp * 100
	// } else {
	// 	tmp = tmp * 1000
	// }

	s := similarMovies{}
	s.MovieList = searchList(dt, cat, tmp-15, tmp+15)
	fmt.Println(cat, dt, tmp, tmp-15, tmp+15)
	fmt.Println(s.MovieList)

	renderTemplate(s, w, template)
}

func videoHandler(w http.ResponseWriter, r *http.Request) {
	// vdo := r.URL.Query().Get("v")
	// cat := r.URL.Query().Get("c")
	// tm := r.URL.Query().Get("t")
	template := "playView"
	k := r.URL.Query().Get("k")
	key, _ := strconv.ParseInt(k, 10, 64)
	m := Movie{}
	if key > 0 {
		m = allMovies[key]
	}

	// pageData.CurrentFilm = vdo
	// pageData.CurrentFilmContext = cat
	// pageData.CurrentFilmTime, _ = strconv.Atoi(tm)

	renderTemplate(m, w, template)
	//http.Redirect(w, r, "static/html/play.html", http.StatusMovedPermanently)
	//http.Redirect(w, r, "file://///tmp/camplayer.html", http.StatusMovedPermanently)
}

func camsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(nil, w, "camView")
}
func menuHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(nil, w, "menu")
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(nil, w, "indexView")
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(GetStats(), w, "sysStats")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(nil, w, "aboutView")
}

// func setupHandler(w http.ResponseWriter, r *http.Request) {
// 	tmpl := "setup.html"
// 	err := r.ParseForm()
// 	if err != nil {
// 		panic(err)
// 	}
// 	if r.Method == "POST" {
// 		config.firstStart = false
// 		if _, ok := r.Form["addFilePath"]; ok {
// 			alreadyExists := false
// 			newPath := r.Form["filepath"][0]
// 			for _, path := range config.FilePathList {
// 				if path == newPath {
// 					alreadyExists = true
// 					break
// 				}
// 			}
// 			if !alreadyExists {
// 				config.FilePathList = append(config.FilePathList, r.Form["filepath"][0])
// 			}
// 			refreshCheck(&tmpl)
// 		} else if _, ok := r.Form["deleteFilePath"]; ok {
// 			if i, err := strconv.Atoi(r.Form["deleteFilePath"][0]); err == nil {
// 				config.FilePathList = append(config.FilePathList[:i], config.FilePathList[i+1:]...)
// 				if len(config.FilePathList) == 0 {
// 					config.firstStart = true
// 				}
// 				refreshCheck(&tmpl)
// 			} else {
// 				panic(err)
// 			}
// 		} else if _, ok := r.Form["saveSetup"]; ok {
// 			// check if there is something to save
// 			if len(config.FilePathList) > 0 {
// 				// Open file using READ & WRITE permission.
// 				var file, err = os.OpenFile("config", os.O_RDWR, 0644)
// 				if err != nil {
// 					panic(err)
// 				}
// 				defer file.Close()

// 				for _, path := range config.FilePathList {
// 					_, err := file.WriteString(path + "\n")
// 					if err != nil {
// 						panic(err)
// 					}
// 				}
// 				err = file.Sync()
// 				if err != nil {
// 					panic(err)
// 				}
// 				tmpl = "index.html"
// 			}
// 		}
// 	}
// 	renderTemplate(config, w, tmpl)
// }

// func playmodeHandler(w http.ResponseWriter, r *http.Request) {
// 	tmpl := "playmode.html"
// 	err := r.ParseForm()
// 	if err != nil {
// 		panic(err)
// 	}
// 	if r.Method == "POST" {
// 		if _, ok := r.Form["pmode"]; ok {
// 			pm := r.Form["pmode"][0]
// 			if pm == "" {
// 				panic(fmt.Errorf("Invalid playmode : %s", pm))
// 			}
// 			config.Playmode = pm
// 			tmpl = "index.html"
// 		}
// 	}
// 	renderTemplate(config, w, tmpl)
// }

// func refreshCheck(tmpl *string) {
// 	if err := refreshList(""); err != nil {
// 		config.firstStart = true
// 		*tmpl = "nothingfound.html"
// 	}
// }
