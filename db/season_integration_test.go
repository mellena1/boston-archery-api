//go:build integration

package db

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mellena1/boston-archery-api/model"
	"github.com/mellena1/boston-archery-api/ptr"
	"github.com/stretchr/testify/require"
)

func Test_CRUDSeason(t *testing.T) {
	ctx := context.Background()

	defer resetTable(ctx)

	_, err := db.GetSeason(ctx, uuid.New())
	require.ErrorIs(t, err, ErrItemNotFound, "getting non-found season should fail")

	id := uuid.New()
	season := model.Season{
		ID:        id,
		Name:      "Fall 2024",
		StartDate: stringToDate("2024-08-01"),
		EndDate:   stringToDate("2024-10-30"),
		ByeWeeks:  []time.Time{},
	}
	_, err = db.AddSeason(ctx, season)
	require.Nil(t, err, "add season should not fail")

	fetchedSeason, err := db.GetSeason(ctx, id)
	require.Nil(t, err, "getting season by id should not fail")

	require.Equal(t, season, *fetchedSeason)

	_, err = db.AddSeason(ctx, season)
	require.ErrorIs(t, err, ErrItemAlreadyExists, "trying to add season with existing UUID should fail")

	_, err = db.UpdateSeason(ctx, id, UpdateSeasonInput{
		Name:      ptr.Ptr("Winter 2024"),
		StartDate: ptr.Ptr(stringToDate("2024-11-01")),
		EndDate:   ptr.Ptr(stringToDate("2024-12-31")),
		ByeWeeks:  ptr.Ptr([]time.Time{stringToDate("2024-11-07")}),
	})
	require.Nil(t, err, "update existing season should not fail")

	fetchedSeason, err = db.GetSeason(ctx, id)
	require.Nil(t, err, "getting season by id should not fail")

	require.Equal(t, model.Season{
		ID:        id,
		Name:      "Winter 2024",
		StartDate: stringToDate("2024-11-01"),
		EndDate:   stringToDate("2024-12-31"),
		ByeWeeks: []time.Time{
			stringToDate("2024-11-07"),
		},
	}, *fetchedSeason)

	_, err = db.UpdateSeason(ctx, uuid.New(), UpdateSeasonInput{
		Name: ptr.Ptr("Winter 2024"),
	})
	require.ErrorIs(t, err, ErrItemNotFound, "trying to update non-existing season should fail")
}

func Test_GetAllSeasons(t *testing.T) {
	ctx := context.Background()

	defer resetTable(ctx)

	seasons := []model.Season{
		{
			ID:        uuid.New(),
			Name:      "Fall 2024",
			StartDate: stringToDate("2024-08-01"),
			EndDate:   stringToDate("2024-10-30"),
		},
		{
			ID:        uuid.New(),
			Name:      "Winter 2024",
			StartDate: stringToDate("2024-11-01"),
			EndDate:   stringToDate("2024-12-31"),
		},
	}
	sort.Slice(seasons, func(i, j int) bool { return seasons[i].StartDate.Before(seasons[j].StartDate) })

	for _, s := range seasons {
		_, err := db.AddSeason(ctx, s)
		require.Nil(t, err, "adding season should not fail")
	}

	fetchedSeasons, err := db.GetAllSeasons(ctx)
	require.Nil(t, err, "should not fail getting all seasons")

	require.Equal(t, seasons, fetchedSeasons)
}
