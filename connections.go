package connections

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strconv"
)

var Configs map[string]string

func init() {
	loadConf()
}

func loadConf() {

	// ================ Read 'Configs.json' file =========
	_, filename, _, _ := runtime.Caller(1)
	configFile := path.Join(path.Dir(filename), "/configs.json")
	jsonFile, err := os.Open(configFile)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &Configs)

	// ================ Read 'Configs.json' file End =====

	// ======== Define logs info =========================

	f, err := os.OpenFile(Configs["LOG_FILE_PATH"], os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

	if err != nil {
		fmt.Println(err.Error() + " : " + Configs["LOG_FILE_PATH"])
	} else {
		log.SetOutput(f)
	}

	// ======== Define Logs End ==========================
}

func ConnectRedis() *redis.Client {

	c := Configs["REDIS_HOST"] + ":" + Configs["REDIS_PORT"]

	rd := redis.NewClient(&redis.Options{
		Addr:     c,
		Password: Configs["REDIS_PASS"], // no password set
		DB:       0,                     // use default DB
	})

	res, err := rd.Ping().Result()

	if err != nil {
		log.WithFields(log.Fields{
			"REDIS_HOST": Configs["REDIS_HOST"],
			"REDIS_PORT": Configs["REDIS_PORT"],
		}).Fatal(err.Error())
	}

	if res != "PONG" {
		log.WithFields(log.Fields{
			"REDIS_HOST": Configs["REDIS_HOST"],
			"REDIS_PORT": Configs["REDIS_PORT"],
		}).Fatal("Invalid Redis Password")
	}

	return rd
}

func ConnectMysql() *sql.DB {

	conn := Configs["MYSQL_USER"] + ":" + Configs["MYSQL_PASS"] + "@tcp(" + Configs["MYSQL_HOST"] + ":" + Configs["MYSQL_PORT"] + ")/" + Configs["MYSQL_DB"]
	db, err := sql.Open("mysql", conn)
	err = db.Ping()

	if err != nil {
		log.WithFields(log.Fields{
			"MYSQL_HOST":Configs["MYSQL_HOST"],
			"MYSQL_USER":Configs["MYSQL_USER"],
			"MYSQL_DB":Configs["MYSQL_DB"],
		}).Fatal(err.Error())
	}

	no_conn, err := strconv.Atoi(Configs["MYSQL_NOCONN"])

	if err != nil {
		log.WithFields(log.Fields{
			"MYSQL_HOST":Configs["MYSQL_HOST"],
			"MYSQL_USER":Configs["MYSQL_USER"],
			"MYSQL_DB":Configs["MYSQL_DB"],
		}).Fatal(err.Error())
	}

	db.SetMaxIdleConns(no_conn)
	db.SetMaxOpenConns(no_conn)
	db.SetConnMaxLifetime(0)

	return db
}
