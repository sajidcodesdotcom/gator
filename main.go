package main

import (
	"fmt"
	"log"

	"github.com/sajidcodess/gator/internal/config"
)

func main() {
  cfg, err := config.Read()
  if err != nil {
    log.Fatalf("Error reading the config: %v", err)
  }

  fmt.Printf("REading conigi: %+v\n", cfg)
}
