package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"go.uber.org/dig"
)

//The same script with Dependency Injection approach

type Movie struct {
	Title     string  `json:"t"`
	Year      int64   `json:"y"`
	Ratings   float64 `json:"r"`
	NmbrVotes int64   `json:"nv"`
	TitleID   string  `json:"id"`
}

type Config struct {
	Enabled      bool
	DatabasePath string
	Port         string
	UserName     string
	Password     string
}

func NewConfig() *Config {
	databasePath := flag.String("bolt", "bolt://localhost:7687", "bolt adress of database (default: bolt://localhost:7687)")
	port := flag.Int("port", 9090, "port number (default 9090)")
	userName := flag.String("User", "neo4j", "")
	password := flag.String("Password", "", "")
	flag.Parse()
	return &Config{
		Enabled:      true,
		DatabasePath: *databasePath,
		Port:         strconv.Itoa(*port),
		UserName:     *userName,
		Password:     *password,
	}
}

func ConnectDatabase(config *Config) (neo4j.Driver, error) {
	return neo4j.NewDriver(config.DatabasePath, neo4j.BasicAuth(string(config.UserName), string(config.Password), ""))
}

type MovieRep struct {
	driver neo4j.Driver
}

func (rep *MovieRep) GiveHint(letters map[string]string) ([]Movie, error) {
	session, err := rep.driver.Session(neo4j.AccessModeRead)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	mapa := map[string]interface{}{
		"prTit": letters["lett"],
	}

	result, err := session.Run(`CALL db.index.fulltext.queryNodes("titles2", $prTit) YIELD node, score
	RETURN node.primaryTitle, node.startYear, node.ratings, node.numberOfVotes, node.titleID, score LIMIT 7
									`, mapa)
	if err != nil {
		return nil, err
	}

	res := []Movie{}
	tmp := Movie{}
	for result.Next() {
		tmp.Title, _ = result.Record().GetByIndex(0).(string)
		tmp.Year, _ = result.Record().GetByIndex(1).(int64)
		tmp.Ratings, _ = result.Record().GetByIndex(2).(float64)
		tmp.NmbrVotes, _ = result.Record().GetByIndex(3).(int64)
		tmp.TitleID, _ = result.Record().GetByIndex(4).(string)

		res = append(res, tmp)
	}
	return res, nil
}

func NewMovieRep(driver neo4j.Driver) *MovieRep {
	return &MovieRep{driver: driver}
}

type MovieService struct {
	config *Config
	rep    *MovieRep
}

func (service *MovieService) GiveHint(letters map[string]string) ([]Movie, error) {
	if service.config.Enabled {
		return service.rep.GiveHint(letters)
	}
	return []Movie{}, nil
}

func NewMovieService(config *Config, rep *MovieRep) *MovieService {
	return &MovieService{config: config, rep: rep}
}

type Server struct {
	config       *Config
	movieService *MovieService
}

func (s *Server) Handler() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hint/{lett}", s.movies).Methods("GET")
	return r
}

func (s *Server) Run() {
	httpServer := &http.Server{
		Addr:    ":" + s.config.Port,
		Handler: s.Handler(),
	}

	httpServer.ListenAndServe()
}

func (s *Server) movies(w http.ResponseWriter, r *http.Request) {
	letters := mux.Vars(r)

	movies, err := s.movieService.GiveHint(letters)
	if err != nil {
		log.Println("error 3 quering database:", err)
		w.WriteHeader(500)
		w.Write([]byte("error 4 quering database:"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(movies)
	if err != nil {
		log.Println("error writing search response:", err)
		w.WriteHeader(500)
		w.Write([]byte("An error occurred writing response"))
	}
}

func NewServer(config *Config, service *MovieService) *Server {
	return &Server{
		config:       config,
		movieService: service,
	}
}

func BuildContainer() *dig.Container {
	container := dig.New()

	container.Provide(NewConfig)
	container.Provide(ConnectDatabase)
	container.Provide(NewMovieRep)
	container.Provide(NewMovieService)
	container.Provide(NewServer)

	return container
}

func main() {
	container := BuildContainer()

	err := container.Invoke(func(server *Server) {
		server.Run()
	})

	if err != nil {
		panic(err)
	}
}
