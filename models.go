package main

// DATA STRUCTURES
type AppConfig struct {
	DefaultCam           string
	StatsDiskPath        string
	StatsRefreshInterval int
}

type ConfigDetails struct {
	firstStart        bool
	templateDirectory string
	FilePathList      []string
	//templateFileList  []string
	Playmode string
	basePath string
}

type Movie struct {
	FullFilePath string
	FileName     string
	Folder       string
	Category     string
	Dt           int    // date as int e.g. 20210914 parsed from folder name
	Tm           int    // time of file parsed from filename e.g. 13-05-01.mp4 will be 130501
	Key          string // key is foldername + filename e.g. 2021-09-14-13-05-01
	EventCam     string // if it is motion event, then name of cam from where event arises
}

type PageData struct {
	MovieList          []Movie
	CurrentFilm        string
	CurrentContext     string
	CurrentFilmContext string
	CurrentFilmTime    int
}

type similarMovies struct {
	MovieList []Movie
}
