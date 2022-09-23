package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/hamzaanis/graphql-test/graph/dal"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/hamzaanis/graphql-test/graph"
	"github.com/hamzaanis/graphql-test/graph/generated"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	db, err := dal.Connect()
	if err != nil {
		log.Fatal(err)
	}
	initDB(db)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(graph.NewRootResolvers(db)))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func initDB(db *sql.DB) {
	dal.MustExec(db, "DROP TABLE IF EXISTS reviews")
	dal.MustExec(db, "DROP TABLE IF EXISTS screenshots")
	dal.MustExec(db, "DROP TABLE IF EXISTS videos")
	dal.MustExec(db, "DROP TABLE IF EXISTS users")
	dal.MustExec(db, "CREATE TABLE public.users (id SERIAL PRIMARY KEY, name varchar(255), email varchar(255))")
	dal.MustExec(db, "CREATE TABLE public.videos (id SERIAL PRIMARY KEY, name varchar(255), description varchar(255), url text,created_at TIMESTAMP, user_id int, FOREIGN KEY (user_id) REFERENCES users (id))")
	dal.MustExec(db, "CREATE TABLE public.screenshots (id SERIAL PRIMARY KEY, video_id int, url text, FOREIGN KEY (video_id) REFERENCES videos (id))")
	dal.MustExec(db, "CREATE TABLE public.reviews (id SERIAL PRIMARY KEY, video_id int,user_id int, description varchar(255), rating varchar(255), created_at TIMESTAMP, FOREIGN KEY (user_id) REFERENCES users (id), FOREIGN KEY (video_id) REFERENCES videos (id))")
	dal.MustExec(db, "INSERT INTO users(name, email) VALUES('Ridham', 'contact@ridham.me')")
	dal.MustExec(db, "INSERT INTO users(name, email) VALUES('Tushar', 'tushar@ridham.me')")
	dal.MustExec(db, "INSERT INTO users(name, email) VALUES('Dipen', 'dipen@ridham.me')")
	dal.MustExec(db, "INSERT INTO users(name, email) VALUES('Harsh', 'harsh@ridham.me')")
	dal.MustExec(db, "INSERT INTO users(name, email) VALUES('Priyank', 'priyank@ridham.me')")
}
