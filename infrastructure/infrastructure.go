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

	ADMIN_ROLE      = "ADMIN_ROLE"
	ENTERPRISE_ROLE = "ENTERPRISE_ROLE"
	STUDENT_ROLE    = "STUDENT_ROLE"
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

	adminRole      uint
	enterpriseRole uint
	studentRole    uint
)

// This function returns the default value if no value is specified
func getStringEnvParameter(envParam string, defaultValue string) string {
	if value, ok := os.LookupEnv(envParam); ok {
		return value
	}
	return defaultValue
}

func goDotEnvVariable(key string, version int) string {
	// load .env file
	switch version {
	case 1:
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("Error loading.env file")
		}
	case 2:
		if err := godotenv.Load(".env.dev"); err != nil {
			log.Fatal("Error loading.env.test file")
		}
	default:
		InfoLog.Printf("Environment: %s not found!\n", key)
		os.Exit(1)
	}
	return os.Getenv(key)
}

func loadEnvParameters(version int) {
	root, _ := os.Getwd()
	env = getStringEnvParameter(ENV, goDotEnvVariable(ENV, version))
	appPort = getStringEnvParameter(APPPORT, goDotEnvVariable(APPPORT, version))
	dbPort = getStringEnvParameter(DBPORT, goDotEnvVariable(DBPORT, version))

	InfoLog.Printf("Environment: %s\n", env)
	dbHost = getStringEnvParameter(DBHOST, goDotEnvVariable(DBHOST, version))
	dbUser = getStringEnvParameter(DBUSER, goDotEnvVariable(DBUSER, version))
	dbPassword = getStringEnvParameter(DBPASSWORD, goDotEnvVariable(DBPASSWORD, version))
	dbName = getStringEnvParameter(DBNAME, goDotEnvVariable(DBNAME, version))

	privatePath = getStringEnvParameter(PRIVATEPATH, root+"/infrastructure/private.pem")
	publicPath = getStringEnvParameter(PUBLICPATH, root+"/infrastructure/public.pem")

	extendHour, _ = strconv.Atoi(getStringEnvParameter(EXTENDHOUR, goDotEnvVariable(EXTENDHOUR, version)))
	extendRefreshHour, _ = strconv.Atoi(getStringEnvParameter(EXTENDREFRESHHOUR, goDotEnvVariable(EXTENDREFRESHHOUR, version)))
	extendAccessMinute, _ = strconv.Atoi(getStringEnvParameter(EXTENDACCESSMINUTE, goDotEnvVariable(EXTENDACCESSMINUTE, version)))

	redisURL = getStringEnvParameter(REDISURL, goDotEnvVariable("REDIS_URL", version))

	httpURL = getStringEnvParameter(HTTPURL, goDotEnvVariable(HTTPURL, version))
	httpSwagger = getStringEnvParameter(HTTPSWAGGER, goDotEnvVariable(HTTPSWAGGER, version))

	rootPath = getStringEnvParameter(ROOTPATH, root)

	staticPath = rootPath + "/static"
	storagePath = "pnk_intern_storage"

	storagePublicPath = "public"

	NameRefreshTokenInCookie = "RefreshToken"
	NameAccessTokenInCookie = "AccessToken"

	mailServer = getStringEnvParameter(MAILSERVER, goDotEnvVariable(MAILSERVER, version))
	mailPort = getStringEnvParameter(MAILPORT, goDotEnvVariable(MAILPORT, version))
	mailAccount = getStringEnvParameter(MAILACCOUNT, goDotEnvVariable(MAILACCOUNT, version))
	mailPassword = getStringEnvParameter(MAILPASS, goDotEnvVariable(MAILPASS, version))

	adminRoleStr, _ := strconv.Atoi(getStringEnvParameter(ADMIN_ROLE, goDotEnvVariable(ADMIN_ROLE, version)))
	enterpriseRoleStr, _ := strconv.Atoi(getStringEnvParameter(ENTERPRISE_ROLE, goDotEnvVariable(ENTERPRISE_ROLE, version)))
	studentRoleStr, _ := strconv.Atoi(getStringEnvParameter(STUDENT_ROLE, goDotEnvVariable(STUDENT_ROLE, version)))
	adminRole = uint(adminRoleStr)
	enterpriseRole = uint(enterpriseRoleStr)
	studentRole = uint(studentRoleStr)
}

func init() {
	InfoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Llongfile)
	ErrLog = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	zapSgLogger, _ := zap.NewDevelopment()
	ZapLogger, _ = zap.NewDevelopment()
	ZapSugar = zapSgLogger.Sugar()

	// Get version ARGS
	var version int
	flag.IntVar(&version, "v", 1, "select version product - v1 or dev - v2")

	var initDB bool
	flag.BoolVar(&initDB, "db", false, "allow recreate model database in postgres")

	flag.Parse()
	InfoLog.Println("database version: ", version)

	loadEnvParameters(version)
	if err := loadAuthToken(); err != nil {
		ErrLog.Println("error load auth token: ", err)
	}

	if err := InitRedis(); err != nil {
		log.Fatal("error initialize redis: ", err)
	}

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

// Get User role in environment
func GetAdminRole() uint {
	return adminRole
}
func GetEnterpriseRole() uint {
	return enterpriseRole
}
func GetStudentRole() uint {
	return studentRole
}
