
##Gin Framework Web Application

### Prerequisites

- Go (version 1.13 or later)
- MySQL database

This is a web application built using the Gin framework and MySQL database. The application provides various routes for user authentication, admission forms, enquery forms, and student data management.

##Features

- User authentication and authorization
- Admission form with data validation
- Enquery form with data validation
- Student data viewing and management
- Fee payment management
- Enquery student data viewing

##Getting Started

- Clone the repository and run go build to build the application.
- Run go run main.go to start the application.
- Open a web browser and navigate to http://localhost:90 to access the application.

##Database Setup

- Create a MySQL database and update the mysql.DataBaseURL variable in the init function with your database credentials.
- Run the database migration script to create the necessary tables.
