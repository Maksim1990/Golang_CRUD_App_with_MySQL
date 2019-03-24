package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "github.com/dgrijalva/jwt-go"
    _ "github.com/gin-gonic/gin"
    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "github.com/mitchellh/mapstructure"
    "net/http"
    "strings"
)

type Employee struct {
    Id   int `json:"id"`
    Name string `json:"name"`
    City string `json:"city"`
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

type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
}
type Data map[string][]Employee

type JwtToken struct {
    Token string `json:"token"`
}

type Exception struct {
    Message string `json:"message"`
}

func CreateTokenEndpoint(w http.ResponseWriter, req *http.Request) {
    var user User
    user.Username="maksim2"
    user.Password="maksimqwerty2"
    _ = json.NewDecoder(req.Body).Decode(&user)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": user.Username,
        "password": user.Password,
    })
    tokenString, error := token.SignedString([]byte("secret"))
    if error != nil {
        fmt.Println(error)
    }
    json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
}

func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        authorizationHeader := req.Header.Get("authorization")
        if authorizationHeader != "" {
            bearerToken := strings.Split(authorizationHeader, " ")
            if len(bearerToken) == 2 {
                token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
                    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                        return nil, fmt.Errorf("There was an error")
                    }
                    return []byte("secret"), nil
                })
                if error != nil {
                    json.NewEncoder(w).Encode(Exception{Message: error.Error()})
                    return
                }
                if token.Valid {
                    context.Set(req, "decoded", token.Claims)
                    next(w, req)
                } else {
                    json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
                }
            }
        } else {
            json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
        }
    })
}
func TestEndpoint(w http.ResponseWriter, req *http.Request) {
    decoded := context.Get(req, "decoded")
    var user User
    mapstructure.Decode(decoded.(jwt.MapClaims), &user)
    json.NewEncoder(w).Encode(user)
}
func UserEndpoint(w http.ResponseWriter, req *http.Request) {
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

    m := make(map[string][]Employee)
    m["data"] = res

    json.NewEncoder(w).Encode(m)

    defer db.Close()
}

func ProtectedEndpoint(w http.ResponseWriter, req *http.Request) {
    params := req.URL.Query()
    token, _ := jwt.Parse(params["token"][0], func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("There was an error")
        }
        return []byte("secret"), nil
    })
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        var user User
        mapstructure.Decode(claims, &user)
        json.NewEncoder(w).Encode(user)
    } else {
        json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
    }
}

func main() {
    router := mux.NewRouter()
    router.Use(commonMiddleware)
    fmt.Println("Starting the application...")
    router.HandleFunc("/authenticate", CreateTokenEndpoint).Methods("POST")
    router.HandleFunc("/protected", ProtectedEndpoint).Methods("GET")
    router.HandleFunc("/test", ValidateMiddleware(TestEndpoint)).Methods("GET")
    router.HandleFunc("/users", ValidateMiddleware(UserEndpoint)).Methods("GET")
    http.ListenAndServe(":9090", router)
}

//-- Set middleware for requests
func commonMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}