package main

import (
	"fmt"
	"os"
)

const moviebuff = "http://data.moviebuff.com/"

var (
	seen    map[string]bool
	degrees int
	trace   map[string]traceData
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

func loopMovies(argument,
	parent,
	destination string) []string {
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

func loopActors(argument,
	movie,
	parent,
	destination string,
	retList []string) []string {
	url := moviebuff + argument
	json, err := getData(url)
	defer ErrHandle(err)

	for _, cast := range json.Cast {
		if notSeen(cast.Url) {
			retList = append(retList, cast.Url)
			var t traceData
			t.addTrace( movie, parent, cast.Name, cast.Role)
			trace[cast.Url] = t
			if cast.Url == destination {
				fmt.Println("Degree of Separation: ", degrees)
				tracer(cast.Url, parent)
				os.Exit(1)
			}
		}
	}
	for _, crew := range json.Crew {
		if notSeen(crew.Url) {
			retList = append(retList, crew.Url)
			var t traceData
			t.addTrace(movie, parent, crew.Name, crew.Role)
			trace[crew.Url] = t
			if crew.Url == destination {
				fmt.Println("Degree of Separation: ", degrees)
				tracer(crew.Url, parent)
				os.Exit(1)
			}
		}
	}
	return retList
}

func main() {
	if len(os.Args) != 3 {
		fmt.Print("Usage Example : degrees vn-mayekar magie-mathur")
	}
	seen = make(map[string]bool)
	trace = make(map[string]traceData)
	retList := make(map[string][]string)
	var q queue

	degrees++
	retList[os.Args[1]] = loopMovies(os.Args[1], os.Args[1], os.Args[2])

	/*Queue to employ BFS*/
	for k := range retList {
		q.enqueue(k)
	}
	for len(q.value) != 0 {
		degrees++
		for _, k := range q.value {
			q.dequeue()
			for _, v := range retList[k] {
				retList[v] = loopMovies(v, v, os.Args[2])
				q.enqueue(v)
			}
		}
	}
}