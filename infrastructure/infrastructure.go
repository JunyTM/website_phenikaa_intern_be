package infrastructure

import (
	"crypto/rsa"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/go-chi/jwtauth"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	APPPORT    = "APP_PORT"
	DBHOST     = "DB_HOST"
	DBPORT     = "DB_PORT"
	DBUSER     = "DB_USER"
	DBPASSWORD = "DB_PASSWORD"
	DBNAME     = "DB_NAME"

	ROOTPATH    = "ROOT_PATH"
	HTTPURL     = "HTTP_URL"
	HTTPSWAGGER = "HTTP_SWAGGER"

	PRIVATEPASSWORD = "PRIVATE_PASSWORD"
	PRIVATEPATH     = "PRIVATE_PATH"
	PUBLICPATH      = "PUBLIC_PATH"

	REDISURL = "REDIS_URL"

	EXTENDHOUR         = "EXTEND_ACCESS_HOUR"
	EXTENDACCESSMINUTE = "EXTEND_ACCESS_MINUTE"
	EXTENDREFRESHHOUR  = "EXTEND_REFRESH_HOUR"

	MAILSERVER  = "MAIL_SERVER"
	MAILPORT    = "MAIL_PORT"
	MAILACCOUNT = "MAIL_ACCOUNT"
	MAILPASS    = "MAIL_PASS"
	ENV         = "ENV"
)

var (
	env string

	appPort    string
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string

	httpURL     string
	httpSwagger string
	rootPath    string
	staticPath  string

	InfoLog *log.Logger
	ErrLog  *log.Logger

	ZapSugar  *zap.SugaredLogger
	ZapLogger *zap.Logger

	db         *gorm.DB
	encodeAuth *jwtauth.JWTAuth
	decodeAuth *jwtauth.JWTAuth
	privateKey *rsa.PrivateKey
	publicKey  interface{}

	redisURL    string
	redisClient *redis.Client

	privatePassword    string
	privatePath        string
	extendAccessMinute int
	extendRefreshHour  int

	publicPath string

	extendHour int

	NameRefreshTokenInCookie string
	NameAccessTokenInCookie  string

	storagePath       string
	storagePublicPath string

	mailServer   string
	mailPort     string
	mailAccount  string
	mailPassword string
)

func getStringEnvParameter(envParam string, defaultValue string) string {
	if value, ok := os.LookupEnv(envParam); ok {
		return value
	}
	return defaultValue
}

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func loadEnvParameters(version int, dbNameArg string, dbPwdArg string) {
	root, _ := os.Getwd()
	env = getStringEnvParameter(ENV, goDotEnvVariable(ENV))
	appPort = getStringEnvParameter(APPPORT, goDotEnvVariable(APPPORT))
	dbPort = getStringEnvParameter(DBPORT, goDotEnvVariable(DBPORT))

	switch version {
	case 0:
		dbHost = getStringEnvParameter(DBHOST, "localhost")
		dbUser = getStringEnvParameter(DBUSER, "postgres")
		dbPassword = getStringEnvParameter(DBPASSWORD, dbPwdArg)
		dbName = getStringEnvParameter(DBNAME, dbNameArg)
		InfoLog.Println("Environment: LOCALHOST")

	default:
		dbHost = getStringEnvParameter(DBHOST, goDotEnvVariable(DBHOST))
		dbUser = getStringEnvParameter(DBUSER, goDotEnvVariable(DBUSER))
		dbPassword = getStringEnvParameter(DBPASSWORD, goDotEnvVariable(DBPASSWORD))
		dbName = getStringEnvParameter(DBNAME, goDotEnvVariable(DBNAME))

		InfoLog.Println("Environment: Development Default")
	}

	privatePath = getStringEnvParameter(PRIVATEPATH, root+"/infrastructure/private.pem")
	publicPath = getStringEnvParameter(PUBLICPATH, root+"/infrastructure/public.pem")

	extendHour, _ = strconv.Atoi(getStringEnvParameter(EXTENDHOUR, "720"))
	extendRefreshHour, _ = strconv.Atoi(getStringEnvParameter(EXTENDREFRESHHOUR, "1440"))
	// extendAccessMinute, _ = strconv.Atoi(getStringEnvParameter(EXTENDACCESSMINUTE, goDotEnvVariable(EXTENDACCESSMINUTE)))

	redisURL = getStringEnvParameter(REDISURL, goDotEnvVariable("REDIS_URL"))

	httpURL = getStringEnvParameter(HTTPURL, goDotEnvVariable(HTTPURL))
	httpSwagger = getStringEnvParameter(HTTPSWAGGER, goDotEnvVariable(HTTPSWAGGER))

	rootPath = getStringEnvParameter(ROOTPATH, root)

	staticPath = rootPath + "/static"
	storagePath = "pnk_intern_storage"

	storagePublicPath = "public"

	NameRefreshTokenInCookie = "RefreshToken"
	NameAccessTokenInCookie = "AccessToken"

	mailServer = getStringEnvParameter(MAILSERVER, goDotEnvVariable(MAILSERVER))
	mailPort = getStringEnvParameter(MAILPORT, goDotEnvVariable(MAILPORT))
	mailAccount = getStringEnvParameter(MAILACCOUNT, goDotEnvVariable(MAILACCOUNT))
	mailPassword = getStringEnvParameter(MAILPASS, goDotEnvVariable(MAILPASS))
}

func init() {
	InfoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Llongfile)
	ErrLog = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	zapSgLogger, _ := zap.NewDevelopment()
	ZapLogger, _ = zap.NewDevelopment()
	ZapSugar = zapSgLogger.Sugar()

	// Get version ARGS
	var version int
	flag.IntVar(&version, "v", 1, "select version dev v1 or dev v2")

	var dbNameArg string
	flag.StringVar(&dbNameArg, "dbname", "postgres", "database name need to connect")

	var dbPwdArg string
	flag.StringVar(&dbPwdArg, "dbpwd", "147563", "password in database need to connect")

	var initDB bool
	flag.BoolVar(&initDB, "db", false, "allow recreate model database in postgres")

	flag.Parse()
	InfoLog.Println("database version: ", version)

	loadEnvParameters(version, dbNameArg, dbPwdArg)
	if err := loadAuthToken(); err != nil {
		ErrLog.Println("error load auth token: ", err)
	}

	// if err := InitRedis(); err != nil {
	// 	log.Fatal("error initialize redis: ", err)
	// }

	if err := InitDatabase(initDB); err != nil {
		ErrLog.Println("error initialize database: ", err)
	}
}

func GetDBName() string {
	return dbName
}

// GetDB export db
func GetDB() *gorm.DB {
	return db
}

// GetHTTPURL export http url
func GetHTTPURL() string {
	return httpURL
}

// GetHTTPSwagger export link swagger
func GetHTTPSwagger() string {
	return httpSwagger
}

// GetAppPort export app port
func GetAppPort() string {
	return appPort
}

func GetRootPath() string {
	return rootPath
}

// GetStoragePath get path of storage
func GetStoragePath() string {
	return storagePath
}

// GetStaticPath export static path
func GetStaticPath() string {
	return staticPath
}

// GetEncodeAuth get token auth
func GetEncodeAuth() *jwtauth.JWTAuth {
	return encodeAuth
}

// GetDecodeAuth export decode auth
func GetDecodeAuth() *jwtauth.JWTAuth {
	return decodeAuth
}

// GetExtendAccessMinute export access extend minute
func GetExtendAccessHour() int {
	return extendHour
}

// GetExtendAccessMinute export access extend minute
func GetExtendAccessMinute() int {
	return extendAccessMinute
}

// GetExtendRefreshHour export refresh extends hour
func GetExtendRefreshHour() int {
	return extendRefreshHour
}

// GetMailParam
func GetMailParam() (string, string, string, string) {
	return mailServer, mailPort, mailAccount, mailPassword
}

// GetRedisClient export redis client
func GetRedisClient() *redis.Client {
	return redisClient
}

// GetPublicKey get public key
func GetPublicKey() interface{} {
	return publicKey
}

func GetEnvironments() string {
	return env
}
