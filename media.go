package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

var fileCounter int

func visit(path string, f os.FileInfo, err error) error {
	bpath := []byte(strings.ToLower(path))
	bpath = bpath[len(bpath)-4:]
	for i := 0; i < len(extensionList); i++ {
		if reflect.DeepEqual(bpath, extensionList[i]) {
			fileCounter++
			// get folder names
			parentFld, recfld := parseFolderName(path)
			camname := ""
			// get time from filename e.g. 13-10-21.mp4 = 131021
			//fname_without_ext := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
			parts := strings.Split(f.Name(), ".")
			if len(parts) > 2 {
				// probably it is motion event where filename is like 13-10-21.camname.mp4
				camname = parts[1]
			}

			file_without_ext := parts[0]
			file_digits := strings.Replace(file_without_ext, "-", "", -1)
			tm, _ := strconv.Atoi(file_digits)

			// getfile_digits date from folder name 2021-09-18
			folder_digits := strings.Replace(recfld, "-", "", -1)
			dt, _ := strconv.Atoi(folder_digits)

			key := folder_digits + file_digits + strconv.Itoa(fileCounter) // e.g. 20210918131021
			key_int64, _ := strconv.ParseInt(key, 10, 64)
			movie := Movie{path, file_without_ext, recfld, parentFld, dt, tm, key, camname}
			fmt.Printf("%v\n", movie)
			//allMovies = append(allMovies, movie)
			allMovies[key_int64] = movie
		}
	}
	return nil
}

// returns last-most folder name and first-most folder after basepath
// last int is folder date as 20210912
// e.g flawn/media/2021-09-12
//   will give   2021-02-12 as folder   and  flawn as cetegory
func parseFolderName(path string) (parentFolder string, recFolder string) {
	f := filepath.Dir(path)
	// if f == currentPath {
	// 	return currentBase
	// }
	// //s := strings.Replace(f, currentPath, currentBase, 1)
	// s := strings.Replace(f, currentPath, "", 1)
	// fmt.Printf("F: %s     CP: %s    CB: %s\n", f, currentPath, currentBase)
	// fmt.Printf("S: %s\n", s)
	// return s
	s := strings.Replace(f, currentPath, "", 1)
	tmp := strings.Split(s, "/")
	//d, _ := strconv.Atoi(strings.Replace(tmp[len(tmp)-1], "-", "", -1))
	parentFolder = tmp[0]
	recFolder = tmp[len(tmp)-1]
	return
}

// func getFolderNameByDepth(path string, depth int) string {
// 	f := filepath.Dir(path)
// 	sep := string(os.PathSeparator)
// 	tmp := strings.Split(f, sep)
// 	d := len(tmp)
// 	if d <= depth {
// 		return f
// 	}
// 	f = ""
// 	for i := depth; i > 0; i-- {
// 		f += tmp[d-i]
// 		if i > 1 {
// 			f += sep
// 		}
// 	}
// 	return f
// }

func generateMovies() error {
	allMovies = map[int64]Movie{}
	startingCounter := 0 //len(allMovies)
	if len(config.FilePathList) > 0 {
		for _, path := range config.FilePathList { //for index, path := range config.FilePathList {
			currentPath = filepath.Dir(path) + "/"
			currentBase = filepath.Base(currentPath)
			//fmt.Println("currentBase: ", currentBase)
			//fmt.Println("Traversing: ", path)
			err := filepath.Walk(path, visit)
			if err != nil || len(allMovies) == startingCounter {
				// // DO not remove path - Santosh 2019-12-17
				// config.FilePathList = append(config.FilePathList[:index], config.FilePathList[index+1:]...)
				// if err == nil {
				// 	err = fmt.Errorf("No files found.")
				// }
				return err
			}
			startingCounter = len(allMovies)
		}
	}

	// sort ascending
	//sort.Slice(allMovies, func(i, j int) bool { return allMovies[i].Folder > allMovies[j].Folder })

	fmt.Printf("file import complete: %d files imported\n", len(allMovies))
	return nil
}

func refreshList(search string) error {
	oldContext := pageData.CurrentContext
	// if pageData.MovieList != nil {
	// 	oldContext = pageData.CurrentContext
	// }
	pageData = PageData{}
	err := generateMovies()

	pageData.CurrentContext = oldContext
	if search == "" {
		pageData.MovieList = getAllMovies()
		//fmt.Printf("All Movie count: %d,  %d\n", len(pageData.MovieList), len(allMovies))
	} else {
		pageData.MovieList = searchList(search, pageData.CurrentContext, 0, 0)
		//fmt.Printf("Movie count: %d\n", len(pageData.MovieList))
	}

	return err
}

func getAllMovies() []Movie {
	l := len(allMovies)
	mvs := make([]Movie, l)
	for _, m := range allMovies {
		mvs = append(mvs, m)
	}
	sort.Slice(mvs, func(i, j int) bool { return mvs[i].Key > mvs[j].Key })
	return mvs
}

func searchList(search, searchContext string, timeFrom, timeTo int) []Movie {
	s_date := 0
	e_date := 0
	t_from := 0
	t_to := 0
	usetime := false
	if timeFrom > 0 || timeTo > 0 {
		switch {
		case timeFrom > 999:
			t_from = timeFrom * 100
		case timeFrom > 99:
			t_from = timeFrom * 1000
		default:
			t_from = timeFrom * 10000
		}

		switch {
		case timeTo > 999:
			t_to = timeTo * 100
		case timeTo > 99:
			t_to = timeTo * 1000
		default:
			t_to = timeTo * 10000
		}

		// t_from = timeFrom * 10000
		// t_to = timeTo * 10000
		usetime = true
	}
	var err error
	m := []Movie{}

	switch search {
	case "d":
		s_date, _ = strconv.Atoi(time.Now().Format("20060102"))
		e_date, _ = strconv.Atoi(time.Now().Format("20060102"))
	case "d1":
		currentTime := time.Now().Add(time.Hour * 24 * -1)
		s_date, _ = strconv.Atoi(currentTime.Format("20060102"))
		e_date, _ = strconv.Atoi(currentTime.Format("20060102"))
	case "d2":
		currentTime := time.Now().Add(time.Hour * 24 * -2)
		s_date, _ = strconv.Atoi(currentTime.Format("20060102"))
		e_date, _ = strconv.Atoi(currentTime.Format("20060102"))
	case "w1":
		currentTime := time.Now().Add(time.Hour * 24 * -7)
		s_date, _ = strconv.Atoi(currentTime.Format("20060102"))
		e_date, _ = strconv.Atoi(time.Now().Format("20060102"))
	case "w2":
		currentTime := time.Now().Add(time.Hour * 24 * -15)
		s_date, _ = strconv.Atoi(currentTime.Format("20060102"))
		e_date, _ = strconv.Atoi(time.Now().Format("20060102"))
	case "m":
		currentTime := time.Now().Add(time.Hour * 24 * -30)
		s_date, _ = strconv.Atoi(currentTime.Format("20060102"))
		e_date, _ = strconv.Atoi(time.Now().Format("20060102"))
	case "m2":
		currentTime := time.Now().Add(time.Hour * 24 * -60)
		s_date, _ = strconv.Atoi(currentTime.Format("20060102"))
		e_date, _ = strconv.Atoi(time.Now().Format("20060102"))
	case "m3":
		currentTime := time.Now().Add(time.Hour * 24 * -90)
		s_date, _ = strconv.Atoi(currentTime.Format("20060102"))
		e_date, _ = strconv.Atoi(time.Now().Format("20060102"))
	default:
		//try to parse from passed string e.g. 20210912
		s_date, err = strconv.Atoi(search)
		if err != nil {
			// not a valid date
			s_date, _ = strconv.Atoi(time.Now().Format("20060102"))
			e_date, _ = strconv.Atoi(time.Now().Format("20060102"))
		} else {
			e_date = s_date
		}
	}

	// fmt.Printf("Search: %s,  Start: %d,  End: %d,  Context: %s\n", search, s_date, e_date, searchContext)
	// fmt.Printf("Time: %d  - %d,  UseTime:%v \n", t_from, t_to, usetime)

	if searchContext == "-" {
		for _, s := range allMovies {
			if s.Dt >= s_date && s.Dt <= e_date {
				if usetime {
					if s.Tm >= t_from && s.Tm <= t_to {
						m = append(m, s)
					}
				} else {
					m = append(m, s)
				}
			}
		}
	} else {
		//create map from context
		//contexts = contextToMap()
		contextToMap(searchContext)

		for _, s := range allMovies {
			//if s.dt >= s_date && s.dt <= e_date && s.Category == pageData.CurrentContext {
			//fmt.Println(s.Dt, s.FileName, s.Category, isExistinContext(s.Category))
			if s.Dt >= s_date && s.Dt <= e_date && isExistinContext(s.Category) {
				if usetime {
					if s.Tm >= t_from && s.Tm <= t_to {
						m = append(m, s)
					}
				} else {
					m = append(m, s)
				}
			}
		}
	}

	//sort.Slice(m, func(i, j int) bool { return m[i].FileName > m[j].FileName })
	sort.Slice(m, func(i, j int) bool { return m[i].Key > m[j].Key })

	return m
}

func contextToMap(searchContext string) {
	contexts = make(map[string]struct{})
	tmp := strings.Split(searchContext, ",")
	for _, v := range tmp {
		contexts[v] = struct{}{}
	}
	//fmt.Println(contexts)
}

func isExistinContext(s string) bool {
	_, ispresent := contexts[s]
	//fmt.Printf("%s exist in map = %v\n", s, ispresent)
	return ispresent
}
