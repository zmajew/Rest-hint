package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

//Simple hint query graph Neo4j movie database API
//Application is conected with database thrue official Neo4j bolt driver for Golang.
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

	session, err = driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		log.Println("error quering database:", err)
		w.WriteHeader(500)
		w.Write([]byte("error quering database:"))
	}
	defer session.Close()

	letters := mux.Vars(r)
	mapa := map[string]interface{}{
		"prTit": letters["lett"],
	}

	result, err = session.Run(`CALL db.index.fulltext.queryNodes("titles", $prTit) YIELD node, score
	RETURN node.primaryTitle, node.startYear, node.ratings, node.numberOfVotes, node.titleID, score LIMIT 7
									`, mapa)
	if err != nil {
		log.Println("error quering database:", err)
		w.WriteHeader(500)
		w.Write([]byte("error quering database:"))
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
	driver, err = neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("user_name", "password", ""))
	if err != nil {
		log.Println("error connecting to database:", err)
	}
	defer driver.Close()

	r := mux.NewRouter()
	r.HandleFunc("/hint/{lett}", hint).Methods("GET")

	log.Fatal(http.ListenAndServe(":9090", r))
}
