package main

import (
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/rs/zerolog"

	configDb "github.com/AgieAja/go-config-db/database/mysql"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	// users
	usersHandler "golang-backend-test/users/usersdelivery"
	usersRepo "golang-backend-test/users/usersrepository"
	usersUC "golang-backend-test/users/usersusecase"
)

var rxURL = regexp.MustCompile(`^/regexp\d*`)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Error().Msg("Failed read configuration database")
		return
	}

	dbHOST := os.Getenv("DB_HOST")
	dbPORT := os.Getenv("DB_PORT")
	dbUSER := os.Getenv("DB_USER")
	dbPASS := os.Getenv("DB_PASS")
	dbNAME := os.Getenv("DB_NAME")
	maxIdle := os.Getenv("MAX_IDLE")
	maxConn := os.Getenv("MAX_CONN")

	myMaxIdle, errMaxidle := strconv.Atoi(maxIdle)

	if errMaxidle != nil {
		log.Error().Msg("Failed convert errMaxIdle = " + errMaxidle.Error())
		return
	}

	myMaxConn, errMaxConn := strconv.Atoi(maxConn)

	if errMaxConn != nil {
		log.Error().Msg("Failed convert errMaxConn = " + errMaxConn.Error())
		return
	}

	configDb.InitConnMySQLORM(dbHOST, dbPORT, dbUSER, dbPASS, dbNAME, myMaxIdle, myMaxConn)

}
func main() {
	port := os.Getenv("PORT")
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"POST", "DELETE", "OPTIONS", "PUT", "GET"},
		AllowHeaders:     []string{"Origin", "Content-Type", " Authorization", "userid"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           720 * time.Hour,
	}))

	// Add a logger middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	r.Use(logger.SetLogger())

	// custom logger
	subLog := zerolog.New(os.Stdout).With().Logger()

	r.Use(logger.SetLogger(logger.Config{
		Logger:         &subLog,
		UTC:            true,
		SkipPath:       []string{"/skip"},
		SkipPathRegexp: rxURL,
	}))

	r.Use(gin.Recovery())

	if port == "" {
		log.Error().Msg("Port Can't Empty")
		return
	}

	db, err := configDb.GetMySQLORM()
	if err != nil {
		log.Error().Msg("Failed Connect to Database = " + err.Error())
		return
	}

	defer db.Close()

	// module users
	usersRepository := usersRepo.NewUsersRepository(db)
	usersUsecase := usersUC.NewUsersUsecase(usersRepository)
	usersHandler.NewUsersHttpHandler(r, usersUsecase)

	log.Info().Msg("Last Update : " + time.Now().Format("2006-01-02 15:04:05"))
	log.Info().Msg("Service Running version 0.0.1 at port : " + port)

	if errHTTP := http.ListenAndServe(":"+port, r); errHTTP != nil {
		log.Error().Msg(errHTTP.Error())
	}
}
