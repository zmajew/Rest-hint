# Rest-hint
Simple rest Go service for hints from Neo4j graph IMDB based movie database.

Service is conected with database through official Neo4j bolt driver for Golang.
One driver interface is created than new session for each query.
Database was populated from IMDB shared data.
Query is based on Lucene fulltext indexes, only for titles, not for persons.


