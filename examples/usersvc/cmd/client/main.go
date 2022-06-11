package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/protogodev/httptest/examples/usersvc"
)

func main() {
	baseURL := flag.String("url", "http://localhost:8080", "The base URL")
	flag.Parse()

	client, err := usersvc.NewHTTPClient(
		&http.Client{Timeout: 10 * time.Second},
		*baseURL,
	)
	if err != nil {
		log.Fatalf("NewHTTPClient err: %v\n", err)
	}

	u := &usersvc.User{
		Name:  "foo",
		Sex:   "male",
		Birth: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	if err := client.CreateUser(context.Background(), u); err != nil {
		log.Fatalf("CreateUser err: %v\n", err)
	}

	user, err := client.GetUser(context.Background(), "foo")
	if err != nil {
		log.Fatalf("GetUser err: %v\n", err)
	}
	log.Printf("GetUser ok: %+v\n", user)

	u = &usersvc.User{
		Sex:   "female",
		Birth: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	if err := client.UpdateUser(context.Background(), "foo", u); err != nil {
		log.Fatalf("UpdateUser err: %v\n", err)
	}

	users, err := client.ListUsers(context.Background())
	if err != nil {
		log.Fatalf("ListUsers err: %v\n", err)
	}
	if len(users) > 0 {
		log.Printf("ListUsers ok: [%+v]\n", users[0])
	}

	if err := client.DeleteUser(context.Background(), "foo"); err != nil {
		log.Fatalf("DeleteUser err: %v\n", err)
	}
}
