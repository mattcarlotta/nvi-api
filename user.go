package main

import (
	"fmt"
	"net/http"
)

func allUsers(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, "All users endpoint\n")
}
