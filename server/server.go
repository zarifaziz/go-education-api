package server

import (
	"context"
	"education-api/api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client
var coursesCollection *mongo.Collection
var studentsCollection *mongo.Collection

func Start() {
	client = connectDB()
	defer client.Disconnect(context.Background())

	coursesCollection = client.Database("education").Collection("courses")
	studentsCollection = client.Database("education").Collection("students")

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

func createCourse(c *gin.Context) {
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course.ID = primitive.NewObjectID()
	course.CreatedAt = time.Now()

	_, err := coursesCollection.InsertOne(context.Background(), course)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating course"})
		return
	}

	c.JSON(http.StatusCreated, course)
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

	// Verify course exists
	var course models.Course
	err = coursesCollection.FindOne(context.Background(), bson.M{"_id": courseID}).Decode(&course)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verifying course"})
		return
	}

	// Update student's enrolled courses
	update := bson.M{
		"$addToSet": bson.M{
			"enrolled_courses": courseID,
		},
	}

	result, err := studentsCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": studentID},
		update,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error enrolling student"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student enrolled successfully"})
}
