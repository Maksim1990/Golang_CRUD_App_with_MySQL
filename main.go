package main

import (
    "database/sql"
    "github.com/foolin/gin-template"
    "github.com/gin-gonic/gin"
    "log"
    "net/http"

    _ "github.com/go-sql-driver/mysql"
    _ "github.com/gin-gonic/gin"
)

type Employee struct {
    Id    int
    Name  string
    City string
}

func dbConn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := "root"
    dbName := "goblog"
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
    ctx.Redirect(http.StatusMovedPermanently, "/")
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
    ctx.Redirect(http.StatusMovedPermanently, "/")
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
    ctx.Redirect(http.StatusMovedPermanently, "/")
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

    ctx.HTML(http.StatusOK, "index", gin.H{
        "res": res,
    })
    defer db.Close()
}

func main() {
    //-- Enable debug color in console
    gin.ForceConsoleColor()
    router := gin.Default()

    //new template engine
    router.HTMLRender = gintemplate.Default()

    router.GET("/test", func(ctx *gin.Context) {
        //render with master
        ctx.HTML(http.StatusOK, "test", gin.H{
            "title": "Index title!",
            "add": func(a int, b int) int {
                return a + b
            },
        })
    })

    router.GET("/", Index)
    router.GET("/edit", Edit)
    router.POST("/update", Update)
    router.GET("/show", Show)
    router.GET("/delete", Delete)
    router.GET("/new", New)
    router.POST("/insert", Insert)

    router.Run(":9090")
}