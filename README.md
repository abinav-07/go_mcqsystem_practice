# GO Practice: MCQ TEST System 
This project was done in order to practice GO lang. This project helps to understand the infranstucture and basics of GO lang. The system include:
- Basic Authentication and Authorization with `JWT`.
- Database migrations and queries using `GORM`.
 - To add `Migrations`, make change in create command of MakeFile by changing the migration name and run `make create` in your power shell.
 **NOTE**: Need to do this RN because @read of MakeFile is not interpreted. 
- Handles `HTTP` requests.
- Utilizes `Gin` framework. 
- Uses `minio` as bucket. 

## Starting the System
- Use `air` to start the system.

## Features Implemented in this system
- Register and Login Users
- JWT Authorization for Users and Admin
- Create Test For Users as Admin
- Create Question and Answers as Admin
- Get Test Details as Admin and User
- Get all the available tests for User
- Submit test and Generate report for User
- Update test availability as Admin

**TODO**: Delete User API for users and use transaction to understand how to use it

- To USE FireStore set current ENV in current PowerShell to $env:GOOGLE_APPLICATION_CREDENTIALS=`serviceAccountPath`


