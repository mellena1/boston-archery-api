# Database

## Intro

This project uses AWS's DynamoDB as its database, utilizing single table design.

DynamoDB was chosen primarily to keep prices down, as it offers pay per request billing rather than paying to keep a database running 24/7.

## ERD

```mermaid
---
title: Boston Archery
---
erDiagram
SEASON {
    uuid id
    string name
    datetime startDate
    datetime endDate
    datetime[] byeWeeks
}
TEAM {
    uuid id
    string name
    string[] teamColors
    uuid captain
}
SEASON }o--o{ TEAM : HAS
PLAYER {
    uuid id
    string firstName
    string lastName
    string discordUser
    string avatarUrl
}
ROSTERED_PLAYER {
    uuid seasonId
    uuid teamId
    uuid playerId
}
ROSTERED_PLAYER ||--o{ PLAYER : JOINS
ROSTERED_PLAYER ||--o{ TEAM : JOINS
ROSTERED_PLAYER ||--o{ SEASON : JOINS
GAME {
    uuid id
    uuid seasonId
    int week
    uuid homeTeam
    uuid awayTeam
    int homeScore
    int awayScore
    datetime startTime
    string youtubeUrl
}
SEASON ||--o{ GAME : HAS
GAME ||--|{ TEAM : HAS
STATLINE {
    uuid gameId
    uuid playerId
    int shots
    int kills
    int deaths
    int catches
    int catchesGiven
    int revives
    int gongHits
}
GAME ||--o{ STATLINE : HAS
PLAYER ||--o{ STATLINE : HAS
```

## Access patterns

| Access Pattern                                        | Table/GSI/LSI | Key Condition                          | Filter Expression | Notes                                              |
| ----------------------------------------------------- | ------------- | -------------------------------------- | ----------------- | -------------------------------------------------- |
| Create/Update Season by ID                            | Table         | PK=SeasonID and SK=SeasonID            |                   |                                                    |
| Get season by ID                                      | Table         | PK=SeasonID and SK=SeasonID            |                   |                                                    |
| Get list of all seasons ordered by date               | GSI1          | PK="SEASONS"                           |                   | Can reverse sort order in query if needed          |
| Get dog eat dog winner/stats by season                | Table         | PK=SeasonID and SK=SeasonID            |                   |                                                    |
| Get stat totals by season?                            | Table         | PK=SeasonID and SK=SeasonID            |                   |                                                    |
| Get all teams by season                               | Table         | PK=SeasonID and SK startsWith("TEAM#") |                   |                                                    |
| Get all team records and points by season (standings) | Table         | PK=SeasonID and SK startsWith("TEAM#") |                   |                                                    |
| Get all stat records by season                        | TODO          | TODO                                   |                   |                                                    |
| Create/Update Team by ID                              | Table         | PK=TeamID and SK=TeamID                |                   |                                                    |
| Get team by ID                                        | Table         | PK=TeamID and SK=TeamID                |                   |                                                    |
| Get all teams                                         | GSI1          | PK="TEAMS"                             |                   |                                                    |
| Get all seasons by team                               | GSI1          | PK=TeamID and SK startsWith("SEASON#") |                   |                                                    |
| Get all team stats by season                          | Table         | PK=SeasonID and SK=TeamID              |                   |                                                    |
| Get team stats all time                               | Table         | PK=TeamID and SK=TeamID                |                   |                                                    |
| Get team roster by season                             | TODO          | TODO                                   |                   |                                                    |
| Get record by team all time                           | Table         | PK=TeamID and SK=TeamID                |                   |                                                    |
| Get points by team all time                           | Table         | PK=TeamID and SK=TeamID                |                   |                                                    |
| Create/Update Player by ID                            | Table         | PK=PlayerID and SK=PlayerID            |                   |                                                    |
| Get player by ID                                      | Table         | PK=PlayerID and SK=PlayerID            |                   |                                                    |
| Get all players                                       | GSI1          | PK="PLAYERS"                           |                   |                                                    |
| Get players by season                                 | TODO          | TODO                                   |                   |                                                    |
| Get player stats by season                            | TODO          | TODO                                   |                   |                                                    |
| Get player stats by year                              | TODO          | TODO                                   |                   |                                                    |
| Get dog eat dog stats by player                       | TODO          | TODO                                   |                   |                                                    |
| Get record by player all time                         | TODO          | TODO                                   |                   |                                                    |
| Get record by player by season                        | TODO          | TODO                                   |                   |                                                    |
| Get player stat improvements by season                | TODO          | TODO                                   |                   |                                                    |
| Create/Update Game by ID and Season ID                | Table         | PK=SeasonID and SK=GameID              |                   |                                                    |
| Get game by ID?                                       | TODO          | TODO                                   |                   |                                                    |
| Get game by ID and Season ID                          | Table         | PK=SeasonID and SK=GameID              |                   |                                                    |
| Get all games by player                               | TODO          | TODO                                   |                   |                                                    |
| Get all games by season                               | Table         | PK=SeasonID and SK startsWith("GAME#") |                   |                                                    |
| Get next week games by season                         | GSI1          | PK=SeasonID and SK startsWith("GAME#") |                   | Reverse sort, get most recent 3                    |
| Get last week games by season                         | GSI1          | PK=SeasonID and SK startsWith("GAME#") |                   | Reverse sort, get most recent 6, take last 3       |
| Get all games by team                                 | GSI2 and GSI3 | PK=TeamID and SK startsWith("GAME#")   |                   | Will be ordered games furtherest in the past first |
| Get all player stats by game                          | TODO          | TODO                                   |                   |                                                    |
| Get team stats by game                                | TODO          | TODO                                   |                   |                                                    |
| Add dog eat dog winners by game                       | Table         | PK=SeasonID and SK=GameID              |                   |                                                    |
| Set rosters/lineup/subs per game                      | TODO          | TODO                                   |                   |                                                    |
| Get individual stat records by year                   | TODO          | TODO                                   |                   |                                                    |
| Get individual stat records all time                  | TODO          | TODO                                   |                   |                                                    |
