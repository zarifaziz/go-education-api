package server

import (
	"context"
	"education-api/api/models"
	"education-api/pkg/worker"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	client             *mongo.Client
	coursesCollection  *mongo.Collection
	studentsCollection *mongo.Collection
	jobQueue           = make(chan worker.Job, 100)
	resultChan         = make(chan worker.Result, 100)
)

func Start() {
	client = connectDB()
	defer client.Disconnect(context.Background())

	db := client.Database("education")
	coursesCollection = db.Collection("courses")
	studentsCollection = db.Collection("students")

	// Initialize workers
	for i := 1; i <= 3; i++ {
		w := worker.NewWorker(i, jobQueue, resultChan, db)
		w.Start(context.Background())
	}

	// Start result handler
	go handleResults()

	// Initialize Gin router
	r := gin.Default()

	// Routes
	r.POST("/courses", createCourse)
	r.GET("/courses", getCourses)
	r.POST("/students", createStudent)
	r.GET("/students/:id", getStudent)
	r.POST("/students/:id/enroll/:courseId", enrollStudent)

	// Start server
	r.Run(":8080")
}

func handleResults() {
	for result := range resultChan {
		if result.Error != nil {
			log.Printf("Job %s failed: %v", result.JobType, result.Error)
			continue
		}
		log.Printf("Job %s completed successfully", result.JobType)
	}
}

func createCourse(c *gin.Context) {
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course.ID = primitive.NewObjectID()
	course.CreatedAt = time.Now()

	jobQueue <- worker.Job{
		Type:    "course_creation",
		Payload: course,
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Course creation in progress",
		"id":      course.ID,
	})
}

func getCourses(c *gin.Context) {
	var courses []models.Course
	cursor, err := coursesCollection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching courses"})
		return
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &courses); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing courses"})
		return
	}

	c.JSON(http.StatusOK, courses)
}

func createStudent(c *gin.Context) {
	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student.ID = primitive.NewObjectID()
	student.CreatedAt = time.Now()
	student.EnrolledCourses = []primitive.ObjectID{} // Initialize empty enrolled courses

	_, err := studentsCollection.InsertOne(context.Background(), student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating student"})
		return
	}

	c.JSON(http.StatusCreated, student)
}

func getStudent(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID format"})
		return
	}

	var student models.Student
	err = studentsCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&student)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching student"})
		return
	}

	c.JSON(http.StatusOK, student)
}

func enrollStudent(c *gin.Context) {
	studentID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID format"})
		return
	}

	courseID, err := primitive.ObjectIDFromHex(c.Param("courseId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID format"})
		return
	}

	jobQueue <- worker.Job{
		Type: "student_enrollment",
		Payload: struct {
			StudentID primitive.ObjectID
			CourseID  primitive.ObjectID
		}{
			StudentID: studentID,
			CourseID:  courseID,
		},
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":   "Enrollment in progress",
		"studentId": studentID,
		"courseId":  courseID,
	})
}
