package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")

	if err != nil {
		t.Fatalf("opening test db error: %v", err)
	}
	err = Migrate(db, "../../migrations")
	if err != nil {
		t.Fatalf("migrating test db error: %v", err)
	}
	_, err = db.Exec(`TRUNCATE TABLE workouts, workout_entries CASCADE`)
	if err != nil {
		t.Fatalf("truncating test db error: %v", err)
	}
	return db
}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresWorkoutStore(db)

	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Title:           "Test Workout",
				Description:     "This is a test workout",
				DurationMinutes: 60,
				CaloriesBurned:  500,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Push Up",
						Sets:         3,
						Reps:         IntPtr(10),
						Weight:       FloatPtr(123.6),
						Notes:        "Good form",
						OrderIndex:   1,
					},
					{
						ExerciseName: "Squat",
						Sets:         3,
						Reps:         IntPtr(10),
						Weight:       FloatPtr(125.8),
						Notes:        "Good form",
						OrderIndex:   2,
					},
					{
						ExerciseName: "Deadlift",
						Sets:         3,
						Reps:         IntPtr(8),
						Weight:       FloatPtr(150.0),
						Notes:        "Focus on form",
						OrderIndex:   3,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "workout with invalid entries",
			workout: &Workout{
				Title:           "Full body workout",
				Description:     "This is a test workout",
				DurationMinutes: 60,
				CaloriesBurned:  500,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Plank",
						Sets:         3,
						Reps:         IntPtr(16),
						Notes:        "Good form",
						OrderIndex:   1,
					},
					{
						ExerciseName: "Squat",
						Sets:         3,
						Reps:         IntPtr(10),
						Notes:        "Good form",
						OrderIndex:   2,
					},
					{
						ExerciseName:    "Deadlift",
						Sets:            3,
						Reps:            IntPtr(8),
						DurationSeconds: IntPtr(60),
						Weight:          FloatPtr(150.0),
						Notes:           "Focus on form",
						OrderIndex:      3,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, errt := store.CreateWorkout(tt.workout)
			if tt.wantErr {
				assert.Error(t, errt)
				return

			}
			assert.NoError(t, errt)
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.DurationMinutes, createdWorkout.DurationMinutes)
			assert.Equal(t, tt.workout.CaloriesBurned, createdWorkout.CaloriesBurned)

			retrieved, err := store.GetWorkoutByID(int64(createdWorkout.ID))
			require.NoError(t, err)
			assert.Equal(t, tt.workout.ID, retrieved.ID)
			assert.Equal(t, len(tt.workout.Entries), len(retrieved.Entries))

			for i := range retrieved.Entries {
				assert.Equal(t, tt.workout.Entries[i].ExerciseName, retrieved.Entries[i].ExerciseName)
				assert.Equal(t, tt.workout.Entries[i].Sets, retrieved.Entries[i].Sets)
				assert.Equal(t, tt.workout.Entries[i].OrderIndex, retrieved.Entries[i].OrderIndex)
			}
		})
	}
}

func IntPtr(i int) *int {
	return &i
}

func FloatPtr(f float64) *float64 {
	return &f
}
