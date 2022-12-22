package main

import {
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
}

func main() {
	http.HandleFunc("/hello", func())
}