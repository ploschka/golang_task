package main

import (
	"net/http"

	log "github.com/ploschka/golang_task/internal/logger"
	sw "github.com/ploschka/golang_task/internal/swagger"
)

func main() {
	log.Info("Server started")

	router := sw.NewRouter()

	log.Error(http.ListenAndServe(":88", router))
}
