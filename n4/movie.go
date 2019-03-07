package n4

import "github.com/neo4j/neo4j-go-driver/neo4j"

type Movie struct {
	Title     string  `json:"t"`
	Year      int64   `json:"y"`
	Ratings   float64 `json:"r"`
	NmbrVotes int64   `json:"nv"`
	TitleID   string  `json:"id"`
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
