//go:build integration

package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/mellena1/boston-archery-api/model"
	"github.com/mellena1/boston-archery-api/ptr"
	"github.com/stretchr/testify/require"
)

func Test_CRUDPlayer(t *testing.T) {
	ctx := context.Background()

	defer resetTable(ctx)

	id := uuid.New()
	player := model.Player{
		ID:        id,
		FirstName: "Joe",
		LastName:  "Burrow",
	}
	_, err := db.AddPlayer(ctx, player)
	require.Nil(t, err, "add player should not fail")

	fetchedPlayer, err := db.GetPlayer(ctx, id)
	require.Nil(t, err, "getting player by id should not fail")

	require.Equal(t, player, *fetchedPlayer)

	_, err = db.AddPlayer(ctx, player)
	require.Error(t, err, "trying to add player with existing UUID should fail")

	_, err = db.UpdatePlayer(ctx, id, UpdatePlayerInput{
		FirstName: ptr.Ptr("Ja'Marr"),
		LastName:  ptr.Ptr("Chase"),
	})
	require.Nil(t, err, "updating existing player should not fail")

	fetchedPlayer, err = db.GetPlayer(ctx, id)
	require.Nil(t, err, "getting player by id should not fail")

	require.Equal(t, model.Player{
		ID:        id,
		FirstName: "Ja'Marr",
		LastName:  "Chase",
	}, *fetchedPlayer)

	_, err = db.UpdatePlayer(ctx, uuid.New(), UpdatePlayerInput{
		FirstName: ptr.Ptr("John"),
	})
	require.Error(t, err, "trying to update non-existing player should fail")
}

func Test_GetAllPlayers(t *testing.T) {
	ctx := context.Background()

	defer resetTable(ctx)

	players := []model.Player{
		{
			ID:        uuid.New(),
			FirstName: "Joe",
			LastName:  "Burrow",
		},
		{
			ID:        uuid.New(),
			FirstName: "Ja'Marr",
			LastName:  "Chase",
		},
	}

	for _, p := range players {
		_, err := db.AddPlayer(ctx, p)
		require.Nil(t, err, "adding player should not fail")
	}

	fetchedPlayers, err := db.GetAllPlayers(ctx)
	require.Nil(t, err, "should not fail getting all players")

	require.Equal(t, players, fetchedPlayers)
}
