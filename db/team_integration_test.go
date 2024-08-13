//go:build integration

package db

import (
	"context"
	"sort"
	"testing"

	"github.com/google/uuid"
	"github.com/mellena1/boston-archery-api/model"
	"github.com/mellena1/boston-archery-api/ptr"
	"github.com/stretchr/testify/require"
)

func Test_CRUDTeam(t *testing.T) {
	ctx := context.Background()

	defer resetTable(ctx)

	id := uuid.New()
	team := model.Team{
		ID:         id,
		Name:       "Hold The Lime",
		TeamColors: []string{"#32CD32", "#000000"},
	}
	_, err := db.AddTeam(ctx, team)
	require.Nil(t, err, "add team should not fail")

	fetchedTeam, err := db.GetTeam(ctx, id)
	require.Nil(t, err, "getting team by id should not fail")

	require.Equal(t, team, *fetchedTeam)

	_, err = db.AddTeam(ctx, team)
	require.Error(t, err, "trying to add team with existing UUID should fail")

	_, err = db.UpdateTeam(ctx, id, UpdateTeamInput{
		Name:       ptr.Ptr("No Strings Attached"),
		TeamColors: ptr.Ptr([]string{"#FFC0CB", "#87CEEB"}),
	})
	require.Nil(t, err, "update existing team should not fail")

	fetchedTeam, err = db.GetTeam(ctx, id)
	require.Nil(t, err, "getting team by id should not fail")

	require.Equal(t, model.Team{
		ID:         id,
		Name:       "No Strings Attached",
		TeamColors: []string{"#FFC0CB", "#87CEEB"},
	}, *fetchedTeam)

	_, err = db.UpdateTeam(ctx, uuid.New(), UpdateTeamInput{
		Name: ptr.Ptr("Bowstars"),
	})
	require.Error(t, err, "trying to update non-existing team should fail")
}

func Test_GetAllTeams(t *testing.T) {
	ctx := context.Background()

	defer resetTable(ctx)

	teams := []model.Team{
		{
			ID:         uuid.New(),
			Name:       "Hold The Lime",
			TeamColors: []string{"#32CD32", "#000000"},
		},
		{
			ID:         uuid.New(),
			Name:       "No Strings Attached",
			TeamColors: []string{"#FFC0CB", "#87CEEB"},
		},
	}
	sort.Slice(teams, func(i, j int) bool { return teams[i].ID.String() < teams[j].ID.String() })

	for _, team := range teams {
		_, err := db.AddTeam(ctx, team)
		require.Nil(t, err, "adding team should not fail")
	}

	fetchedTeams, err := db.GetAllTeams(ctx)
	require.Nil(t, err, "should not fail getting all seasons")

	require.Equal(t, teams, fetchedTeams)
}
