package main

import (
	"fmt"
	"time"

	"github.com/osamingo/indigo"
)

var g *indigo.Generator

func init() {
	t := time.Unix(1257894000, 0) // 2009-11-10 23:00:00 UTC
	g = indigo.New(nil, indigo.StartTime(t))
	_, err := g.NextID()
	if err != nil {
		fmt.Println("error")
	}
}

func Next() (*string, error) {
	id, err := g.NextID()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		return &id, nil
	}

	return nil, err
}
