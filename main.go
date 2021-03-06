package main

import (
	"database/sql"
	"github.com/foolin/gin-template"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"

	_ "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	Id   int
	Name string
	City string
}

func dbConn() (db *sql.DB) {
	dbDriver := os.Getenv("DB_ENGINE")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func New(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "new", nil)
}

func Delete(ctx *gin.Context) {
	db := dbConn()
	emp := ctx.Request.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM employee WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("DELETE")
	defer db.Close()
	ctx.Redirect(http.StatusMovedPermanently, "/home")
}

func Show(ctx *gin.Context) {
	db := dbConn()
	nId := ctx.Request.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM employee WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id int
		var name, city string
		err = selDB.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
	}
	ctx.HTML(http.StatusOK, "show", gin.H{
		"user": emp,
	})
	defer db.Close()
}

func Edit(ctx *gin.Context) {
	db := dbConn()
	nId := ctx.Request.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM employee WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id int
		var name, city string
		err = selDB.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
	}
	ctx.HTML(http.StatusOK, "edit", gin.H{
		"user": emp,
	})
	defer db.Close()
}

func Insert(ctx *gin.Context) {
	db := dbConn()
	name := ctx.PostForm("name")
	city := ctx.PostForm("city")
	insForm, err := db.Prepare("INSERT INTO employee(name, city) VALUES(?,?)")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(name, city)
	log.Println("INSERT: Name: " + name + " | City: " + city)

	defer db.Close()
	ctx.Redirect(http.StatusMovedPermanently, "/home")
}

func Update(ctx *gin.Context) {
	db := dbConn()

	name := ctx.PostForm("name")
	city := ctx.PostForm("city")
	id := ctx.PostForm("uid")
	insForm, err := db.Prepare("UPDATE employee SET name=?, city=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(name, city, id)
	log.Println("UPDATE: Name: " + name + " | City: " + city)

	defer db.Close()
	ctx.Redirect(http.StatusMovedPermanently, "/home")
}

func Index(ctx *gin.Context) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM employee ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	res := []Employee{}
	for selDB.Next() {
		var id int
		var name, city string
		err = selDB.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
		res = append(res, emp)
	}

	session := sessions.Default(ctx)
	user := session.Get("user")
	if user != nil {
		ctx.HTML(http.StatusOK, "index", gin.H{
			"res":  res,
			"user": user,
		})
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/login")
	}
	defer db.Close()
}

func signUpView(ctx *gin.Context) {
	session := sessions.Default(ctx)
	flashes := session.Flashes()
	session.Save()
	ctx.HTML(http.StatusOK, "signup", gin.H{
		"flashes": flashes,
	})
}
func loginView(ctx *gin.Context) {
	session := sessions.Default(ctx)
	flashes := session.Flashes()
	session.Save()
	ctx.HTML(http.StatusOK, "login", gin.H{
		"flashes": flashes,
	})

}
func signUp(ctx *gin.Context) {

	db := dbConn()
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	var user string

	err := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)

	switch {
	case username == "":
		setFlashMessage("Username can not be empty", ctx)
		ctx.Redirect(http.StatusMovedPermanently, "/signup")
		return
	case password == "":
		setFlashMessage("Password can not be empty", ctx)
		ctx.Redirect(http.StatusMovedPermanently, "/signup")
		return
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			setFlashMessage("Server error while hashing provided password value.", ctx)
			ctx.Redirect(http.StatusMovedPermanently, "/signup")
			return
		}

		_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
		if err != nil {
			setFlashMessage("Server error while register in DB provided data.", ctx)
			ctx.Redirect(http.StatusMovedPermanently, "/signup")
			return
		}
		ctx.Redirect(http.StatusMovedPermanently, "/home")

	case err != nil:
		setFlashMessage("Server error, unable to create your account.", ctx)
		ctx.Redirect(http.StatusMovedPermanently, "/signup")
		return
	default:
		setFlashMessage("Server error.Such username already exist.", ctx)
		ctx.Redirect(http.StatusMovedPermanently, "/signup")
	}
}

func loginPage(ctx *gin.Context) {

	db := dbConn()
	session := sessions.Default(ctx)
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	var databaseUsername string
	var databasePassword string

	err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)

	if err != nil {
		setFlashMessage("Provided Username or Password are invalid.", ctx)
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		setFlashMessage("Provided password is invalid.", ctx)
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	session.Set("user", username)
	session.Save()
	ctx.Redirect(http.StatusMovedPermanently, "/home")

}

func Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Delete("user")
	session.Save()

	ctx.Redirect(http.StatusMovedPermanently, "/login")
}
func authPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login_index", nil)
}

func setFlashMessage(message string, ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.AddFlash(message)
	session.Save()
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//-- Enable debug color in console
	gin.ForceConsoleColor()
	router := gin.Default()

	//-- Initialize Session based on Cookies
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	//new template engine
	router.HTMLRender = gintemplate.Default()

	router.GET("/", authPage)
	router.GET("/signup", signUpView)
	router.POST("/signup", signUp)
	router.GET("/login", loginView)
	router.GET("/logout", Logout)
	router.POST("/login", loginPage)

	router.GET("/home", Index)
	router.GET("/edit", Edit)
	router.POST("/update", Update)
	router.GET("/show", Show)
	router.GET("/delete", Delete)
	router.GET("/new", New)
	router.POST("/insert", Insert)

	router.Run(":9090")
}




