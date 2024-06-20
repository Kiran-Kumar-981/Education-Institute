//gin-framework or gin package give us the used to reduce lines of code and helps in scaling the application and helps in easy routing of application
package main

import (
	"database/sql"
	"fmt"
	"institute/mysql"  //hot code is secured in different package to increase security and reduce boiler plate
	"net/http"

	"github.com/gin-gonic/gin" //third-party package 
	_ "github.com/go-sql-driver/mysql" //mysql driver
)

var (
	dataBase *sql.DB
	err      error
)

type UserData struct {
	ID            int    `form:"id"`
	Name          string `form:"Name" binding:"required"`
	FatherName    string `form:"FatherName" binding:"required"`
	Qualification string `form:"Qualification" binding:"required"`
	Email         string `form:"Email" binding:"required,email"`
	PhNumber      string `form:"PhNumber" binding:"required,len=10"`
	Course        string `form:"Course" binding:"required"`
	Address       string `form:"Address" binding:"required"`
	Duration      string `form:"Duration" binding:"required"`
	Fee           int    `form:"Fee" binding:"required"`
	BatchTiming   string `form:"BatchTiming" binding:"required"`
	FeePaid       int    `form:"FeePaid" binding:"required"`
}

type EnqueryData struct {
	ID            int    `form:"id"`
	Name          string `form:"Name" binding:"required"`
	Qualification string `form:"Qualification" binding:"required"`
	Email         string `form:"Email" binding:"required,email"`
	PhNumber      string `form:"PhNumber" binding:"required,len=10"`
	Course        string `form:"Course" binding:"required"`
}
//initialising the db is good practice as it need's to store and retrive the information from database as per the requirment
func init() {
	dataBase, err = sql.Open("mysql", mysql.DataBaseURL)
	if err != nil {
		fmt.Println("Database connection error:", err)
		return
	}

	err = dataBase.Ping()
	if err != nil {
		fmt.Println("Database ping error:", err)
		return
	}
}

func main() {
	defer dataBase.Close()
	router := gin.Default()
	//static files are given in static folder which are  css files
	router.Static("/static", "./static")
	//HTML files are stored in templet folder for the readability and modularity
	router.LoadHTMLGlob("templete/*.html")
	//localhost will show the login page as index page for this
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	//grouping and authorizing will help in use friendly and secure from outsiders
	authorized := router.Group("/")
	authorized.Use(authRequired())
	{
		authorized.GET("/home", home)
		authorized.GET("/admissionForm", admissionForm)
		authorized.GET("/admissionFormSubmitting", func(ctx *gin.Context) {
			go admissions(ctx)
			ctx.Redirect(http.StatusPermanentRedirect, "/DetailsSubmited")
		})
		authorized.GET("/viewstudent", viewStudent)
		authorized.GET("/DetailsSubmited", detailsSubmited)
		authorized.GET("/PayFee", payFee)
		authorized.GET("/enqueryForm", enqueryForm)
		authorized.POST("/enqueryformsubmiting", func(ctx *gin.Context) {
			go enquery(ctx)
			ctx.Redirect(http.StatusPermanentRedirect, "/DetailsEnqSubmited")
		})
		authorized.GET("/enqueryviewstudent", enqueryViewStudent)
		authorized.GET("/DetailsEnqSubmited", detailsEnqSubmited)
		authorized.GET("/logout", logout)
	}

	router.POST("/login", login)
	http.ListenAndServe(":90", router)
}

func home(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "home.html", nil)
}
//opens the admission form folder to get the use input
func admissionForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admissionForm.html", nil)
}
//viewStudent function gets the details of students from the daytabase and displays in the table fomate
func viewStudent(ctx *gin.Context) {
	students, err := fetchStudentData()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.HTML(http.StatusOK, "viewStudents.html", gin.H{"students": students})
}
//after succesfully inserting the details of new student it displays succesfull insertion for admission
func detailsSubmited(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "inserted.html", nil)
}
//after succesfully inserting the details of new student it displays succesfull insertion for enquery
func detailsEnqSubmited(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "DetailsEnqSubmited.html", nil)
}
//student can pay fee directly in viewstudents portal onclick redirects to payment portal
func payFee(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "PayFee.html", nil)
}
//gives the faculty to store the details of the students for enquery about the courese present and they want
func enqueryForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "enqueryForm.html", nil)
}
//admissions gets the details from the students and inserts to the database
func admissions(ctx *gin.Context) {
	var userData UserData
	if err := ctx.ShouldBind(&userData); err != nil {
		ctx.HTML(http.StatusBadRequest, "insertionfailed.html", nil)
		return
	}

	statement, err := dataBase.Prepare("INSERT INTO students(Name, FatherName, Qualification, Email, PhNumber, Course, Address, Duration, Fee, BatchTiming, FeePaid) VALUES(?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Println("Prepare error:", err)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(userData.Name, userData.FatherName, userData.Qualification, userData.Email, userData.PhNumber, userData.Course, userData.Address, userData.Duration, userData.Fee, userData.BatchTiming, userData.FeePaid)
	if err != nil {
		fmt.Println("Exec error:", err)
		ctx.HTML(http.StatusInternalServerError, "insertionfailed.html", nil)
		return
	}

	ctx.HTML(http.StatusOK, "inserted.html", nil)
}
//this is the login form that is authentication and sets the cookie in browser for further purposes
func login(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if username == "username" && password == "password" {
		ctx.SetCookie("authenticated", "true", 600, "/", "", true, true)
		ctx.Redirect(http.StatusSeeOther, "/home")
	} else {
		ctx.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "Invalid credentials",
		})
	}
}
//authRequired handles the cookie to be set or not 
func authRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authenticated, err := ctx.Cookie("authenticated")
		if err != nil || authenticated != "true" {
			ctx.Redirect(http.StatusSeeOther, "/")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
//logout removes the cookie and brings to the login page
func logout(ctx *gin.Context) {
	ctx.SetCookie("authenticated", "", -1, "/", "", false, true)
	ctx.Redirect(http.StatusSeeOther, "/")
}
//enquery gets the detaies of the students that are intrested and willing to join the institute
func enquery(ctx *gin.Context) {
	var enqueryData EnqueryData
	if err := ctx.ShouldBind(&enqueryData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	statement, err := dataBase.Prepare("INSERT INTO enqueries(Name, Qualification, Email, PhNumber, Course) VALUES(?,?,?,?,?)")
	if err != nil {
		fmt.Println("Prepare error:", err)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(enqueryData.Name, enqueryData.Qualification, enqueryData.Email, enqueryData.PhNumber, enqueryData.Course)
	if err != nil {
		fmt.Println("Exec error:", err)
		ctx.HTML(http.StatusInternalServerError, "insertionfailed.html", nil)
		return
	}

	ctx.HTML(http.StatusOK, "inserted.html", nil)
}
//enqueryViewStudent give acceses to faculty to contact the students visited and still need to join 
func enqueryViewStudent(ctx *gin.Context) {
	enqStudents, err := fetchEnqStudentData()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.HTML(http.StatusOK, "enqueryviewstudents.html", gin.H{"students": enqStudents})
}
//fetchStudentData featches the raw data from database and send to the viewStudents table
func fetchStudentData() ([]UserData, error) {
	rows, err := dataBase.Query("SELECT * FROM students")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []UserData
	for rows.Next() {
		var student UserData
		if err := rows.Scan(&student.ID, &student.Name, &student.FatherName, &student.Qualification, &student.Email, &student.PhNumber, &student.Course, &student.Address, &student.Duration, &student.Fee, &student.BatchTiming, &student.FeePaid); err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, nil
}
//fetchEnqStudentData featches the raw data from database and send to the viewEnqStudents table
func fetchEnqStudentData() ([]EnqueryData, error) {
	rows, err := dataBase.Query("SELECT * FROM enqueries")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enqStudents []EnqueryData
	for rows.Next() {
		var enqStudent EnqueryData
		if err := rows.Scan(&enqStudent.ID, &enqStudent.Name, &enqStudent.Qualification, &enqStudent.Email, &enqStudent.PhNumber, &enqStudent.Course); err != nil {
			return nil, err
		}
		enqStudents = append(enqStudents, enqStudent)
	}
	return enqStudents, nil
}
