package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

//Simple hint query graph Neo4j movie database service
//Application is conected with database through official Neo4j bolt driver for Golang.
//One driver interface is created than new session for each query.
//Database was populated from IMDB shared data.
//Query is based on Lucene fulltext indexes, only for titles, not for persons.

type movieResult struct {
	Title     string  `json:"t"`
	Year      int64   `json:"y"`
	Ratings   float64 `json:"r"`
	NmbrVotes int64   `json:"nv"`
	TitleID   string  `json:"id"`
}

var (
	driver  neo4j.Driver
	session neo4j.Session
	result  neo4j.Result
	err     error
)

func hint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	session, err = driver.Session(neo4j.AccessModeRead)
	if err != nil {
		log.Println("error 1 quering database:", err)
		w.WriteHeader(500)
		w.Write([]byte("error 2 quering database:"))
	}
	defer session.Close()

	letters := mux.Vars(r)
	mapa := map[string]interface{}{
		"prTit": letters["lett"],
	}

	result, err = session.Run(`CALL db.index.fulltext.queryNodes("titles2", $prTit) YIELD node, score
	RETURN node.primaryTitle, node.startYear, node.ratings, node.numberOfVotes, node.titleID, score LIMIT 7
									`, mapa)
	if err != nil {
		log.Println("error 3 quering database:", err)
		w.WriteHeader(500)
		w.Write([]byte("error 4 quering database:"))
	}

	res := []movieResult{}
	tmp := movieResult{}
	for result.Next() {
		tmp.Title, _ = result.Record().GetByIndex(0).(string)
		tmp.Year, _ = result.Record().GetByIndex(1).(int64)
		tmp.Ratings, _ = result.Record().GetByIndex(2).(float64)
		tmp.NmbrVotes, _ = result.Record().GetByIndex(3).(int64)
		tmp.TitleID, _ = result.Record().GetByIndex(4).(string)

		res = append(res, tmp)

	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Println("error writing search response:", err)
		w.WriteHeader(500)
		w.Write([]byte("An error occurred writing response"))
	}
}

func main() {
	boltAdress := flag.String("bolt", "bolt://localhost:7687", "bolt adress of database (default: bolt://localhost:7687)")
	port := flag.Int("port", 9090, "port number (default 9090)")
	flag.Parse()

	driver, err = neo4j.NewDriver(*boltAdress, neo4j.BasicAuth(string("neo4j"), string("trustno1"), ""))
	if err != nil {
		log.Fatalln("error connecting to database:", err)
	}
	defer driver.Close()

	r := mux.NewRouter()
	r.HandleFunc("/hint/{lett}", hint).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), r))
}
