# Access Patterns

| Access Pattern              | Table/GSI/LSI | Key Condition                       | Filter Expression |
| --------------------------- | ------------- | ----------------------------------- | ----------------- |
| Get all seasons             | EntityTypeGSI | PK=EntityType                       |                   |
| Get all games for a season  | GSI1          | PK=SeasonID and SK begins_with "G#" |                   |
| Get all teams for a season  | GSI1          | PK=SeasonID and SK begins_with "T#" |                   |
| Get team for a given teamID | Table         | PK=TeamID                           |                   |
