package main

import (
	"backend/src"
	"crypto/rand"
	"crypto/rsa"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
)

func main() {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
	// keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	// if err != nil {
	// 	panic(err)
	// }

	connConf, err := pgx.ParseConnectionString("postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		panic(err)
	}

	sqlDb := stdlib.OpenDB(connConf)
	if err := sqlDb.Ping(); err != nil {
		panic(err)
	}

	jwtUtil := src.NewJwtUtil(key)
	passwordUtil := src.NewPasswordUtil()
	db := src.NewDB(sqlDb)
	userService := src.NewUserService(src.UserServiceDeps{
		Jwt:      jwtUtil,
		Password: passwordUtil,
		Db:       db,
		Cache:    src.NewCacheInmem[string, src.UserDTO](60 * 60),
	})

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	userGroup := r.Group("/user")
	userGroup.POST("/create", src.NewUserCreateHandler(userService))
	userGroup.POST("/login", src.NewUserLoginHandler(userService))

	dummyGroup := r.Group("/dummy")
	dummyGroup.Use(src.NewAuthMiddleware(userService))
	dummyGroup.GET("/", src.NewDummyHandler())

	r.Run(":8080")
}
