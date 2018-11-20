# Golang_CRUD_App_with_MySQL
Golcang CRUD and Authentification application with MySQL database connection

### How To Run

1) Create a new database with a users table

```sql
CREATE TABLE users(
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50),
    password VARCHAR(120)
);
```

2) Go get both required packages listed below

```
go get golang.org/x/crypto/bcrypt

go get github.com/go-sql-driver/mysql
```

3) In  **signup.go** and **main.go** set correct DB connection

4) For CRUD example run
```
go run main.go
```
and navigate to [http://localhost:8090/](http://localhost:8090/)

4) For Authentification example run
```
go run signup.go
```
and navigate to [http://localhost:8080/](http://localhost:8080/)
