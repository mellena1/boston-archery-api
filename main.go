package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mellena1/boston-archery-api/db"
)

func main() {
	ctx := context.Background()
	dynamoClient, err := db.CreateLocalClient(ctx)
	if err != nil {
		panic(err)
	}
	database := db.NewDB("ArcheryTag", "EntityTypeGSI", dynamoClient)
	err = database.AddSeason(ctx, db.SeasonInput{
		Name:      "Fall 2023",
		StartDate: time.Date(2023, 9, 11, 0, 0, 0, 0, time.Local),
		EndDate:   time.Date(2023, 11, 20, 0, 0, 0, 0, time.Local),
		ByeWeeks: []time.Time{
			time.Date(2023, 10, 9, 0, 0, 0, 0, time.Local),
		},
	})
	if err != nil {
		panic(err)
	}
	seasons, err := database.GetAllSeasons(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(seasons)
}
