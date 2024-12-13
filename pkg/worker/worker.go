package worker

import (
	"context"
	"education-api/api/models"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Job represents a task to be processed
type Job struct {
	Type    string      // "course_creation" or "student_enrollment"
	Payload interface{} // Course or EnrollmentData
}

// Result represents the outcome of a job
type Result struct {
	JobType string
	Status  string
	Error   error
}

// Worker represents our background processor
type Worker struct {
	ID         int
	JobQueue   chan Job
	ResultChan chan Result
	db         *mongo.Database
}

func NewWorker(id int, jobQueue chan Job, resultChan chan Result, db *mongo.Database) *Worker {
	return &Worker{
		ID:         id,
		JobQueue:   jobQueue,
		ResultChan: resultChan,
		db:         db,
	}
}

func (w *Worker) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case job := <-w.JobQueue:
				// Process the job
				result := w.processJob(job)
				w.ResultChan <- result
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (w *Worker) processJob(job Job) Result {
	log.Printf("Worker %d processing %s\n", w.ID, job.Type)
	ctx := context.Background()

	switch job.Type {
	case "course_creation":
		course := job.Payload.(models.Course)
		_, err := w.db.Collection("courses").InsertOne(ctx, course)
		if err != nil {
			return Result{
				JobType: job.Type,
				Status:  "failed",
				Error:   fmt.Errorf("failed to insert course: %v", err),
			}
		}
		log.Printf("Worker %d created course: %s\n", w.ID, course.ID.Hex())
		return Result{
			JobType: job.Type,
			Status:  "completed",
		}

	case "student_enrollment":
		data := job.Payload.(struct {
			StudentID primitive.ObjectID
			CourseID  primitive.ObjectID
		})

		// Verify course exists
		var course models.Course
		err := w.db.Collection("courses").FindOne(ctx, bson.M{"_id": data.CourseID}).Decode(&course)
		if err != nil {
			return Result{
				JobType: job.Type,
				Status:  "failed",
				Error:   fmt.Errorf("course not found: %v", err),
			}
		}

		// Update student's enrolled courses
		update := bson.M{
			"$addToSet": bson.M{
				"enrolled_courses": data.CourseID,
			},
		}

		result, err := w.db.Collection("students").UpdateOne(
			ctx,
			bson.M{"_id": data.StudentID},
			update,
		)

		if err != nil {
			return Result{
				JobType: job.Type,
				Status:  "failed",
				Error:   fmt.Errorf("failed to enroll student: %v", err),
			}
		}

		if result.MatchedCount == 0 {
			return Result{
				JobType: job.Type,
				Status:  "failed",
				Error:   fmt.Errorf("student not found"),
			}
		}

		log.Printf("Worker %d enrolled student %s in course %s\n",
			w.ID, data.StudentID.Hex(), data.CourseID.Hex())

		return Result{
			JobType: job.Type,
			Status:  "completed",
		}
	}

	return Result{
		JobType: job.Type,
		Status:  "failed",
		Error:   fmt.Errorf("unknown job type"),
	}
}
