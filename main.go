package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

const moviebuff = "http://data.moviebuff.com/"

var (
	seen    map[string]bool
	degrees int
	trace   map[string]traceData
	now     time.Time
)

//General error panic
func ErrHandle(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func notSeen(in string) bool {
	if !seen[in] {
		seen[in] = true
		return true
	} else {
		return false
	}
}

func loopMovies(argument, parent, destination string) []string {
	var retList []string
	url := moviebuff + argument
	json, err := getData(url)
	defer ErrHandle(err)
	for _, movie := range json.Movies {
		if notSeen(movie.Url) {
			retList = loopActors(movie.Url, movie.Name, argument,
				destination, retList)
		}
	}
	return retList
}

func procD(movie, parent, destination, Url, Name, Role string) {
	var t traceData
	t.addTrace(movie, parent, Name, Role)
	trace[Url] = t
	if Url == destination {
		fmt.Println("Degree of Separation: ", degrees)
		tracer(Url, parent)
		fmt.Println("Time taken: ", time.Since(now))
		os.Exit(1)
	}
}

func loopActors(argument, movie, parent, destination string,
	retList []string) []string {
	url := moviebuff + argument
	json, err := getData(url)
	defer ErrHandle(err)

	for _, cast := range json.Cast {
		if notSeen(cast.Url) {
			retList = append(retList, cast.Url)
			procD(movie, parent, destination, cast.Url, cast.Name, cast.Role)
		}
	}
	for _, crew := range json.Crew {
		if notSeen(crew.Url) {
			retList = append(retList, crew.Url)
			procD(movie, parent, destination, crew.Url, crew.Name, crew.Role)
		}
	}
	return retList
}

func master(q queue, retList map[string][]string) {
	rateLimit := time.Tick(time.Millisecond / 10)
	for len(q.value) != 0 {
		degrees++ //global
		for _, k := range q.value {
			q.dequeue()
			for _, v := range retList[k] {
				<-rateLimit
				fmt.Println(v)
				retList[v] = loopMovies(v, v, os.Args[2])
				q.enqueue(v)
			}
		}
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Print("Usage Example : degrees vn-mayekar magie-mathur")
	}
	now = time.Now()
	seen = make(map[string]bool)
	trace = make(map[string]traceData)
	retList := make(map[string][]string)
	runtime.GOMAXPROCS(runtime.NumCPU())

	var q queue

	degrees++
	retList[os.Args[1]] = loopMovies(os.Args[1], os.Args[1], os.Args[2])

	//Queue to employ BFS
	for k := range retList {
		q.enqueue(k)
	}
	master(q, retList)
}
