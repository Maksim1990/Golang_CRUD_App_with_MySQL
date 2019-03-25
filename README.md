# Golang_CRUD_App_with_MySQL
Golcang CRUD and Authentification application with MySQL database connection.

- *Based on main functionality of [Gin](https://gin-gonic.com/) framework*
- *Implemented routing functionality*
- *Implemented DotEnv functionality*
- *Added Session-based simple authentication*
- *Added support of reusable templates*
- *Implemented support of sessions functionality based on Redis*

### How To Run

1) Create a new database tables:
- **users** table

```sql
CREATE TABLE users(
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50),
    password VARCHAR(120)
);
```

- **employee** table

```sql
CREATE TABLE employee(
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50),
    city VARCHAR(120)
);
```

2) Go get both required packages listed below

```
go get github.com/foolin/gin-template

go get github.com/gin-contrib/sessions

go get github.com/gin-contrib/sessions

go get github.com/gin-gonic/gin

go get github.com/joho/godotenv

go get golang.org/x/crypto/bcrypt

go get github.com/go-sql-driver/mysql
```

3) Run following command to prepare **.env** file

```
cp .env.dist .env
```
4) Fill in correct DB settings in **.env** file

4) Run following comand to start Gin framework server
```
gin -i -all run main.go

```
and navigate to [http://localhost:9090/](http://localhost:9090/)


### Demo papp preview

1) Example of application **Login** page
![Mockup for feature A](https://github.com/Maksim1990/Golang_CRUD_App_with_MySQL/blob/master/github/1.PNG?raw=true)

1) Example of application **Home** page
![Mockup for feature A](https://github.com/Maksim1990/Golang_CRUD_App_with_MySQL/blob/master/github/2.PNG?raw=true)

1) Example of application **Create new** page
![Mockup for feature A](https://github.com/Maksim1990/Golang_CRUD_App_with_MySQL/blob/master/github/3.PNG?raw=true)
