package connections

import (
	"fmt"
	"database/sql"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"encoding/json"
	"io/ioutil"
	"strconv"
	log "github.com/sirupsen/logrus"
	"runtime"
	"path"
)

var configs map[string]string

func init(){
	loadConf()
}

func loadConf(){

	// ================ Read 'configs.json' file =========
	_, filename, _, _ := runtime.Caller(1)
	configFile := path.Join(path.Dir(filename), "/configs.json")
	jsonFile, err := os.Open(configFile);
    if err != nil {
        fmt.Println(err)
    }
	defer jsonFile.Close()
    byteValue, _ := ioutil.ReadAll(jsonFile)
    json.Unmarshal([]byte(byteValue), &configs)

	// ================ Read 'configs.json' file End =====

	// ======== Define logs info =========================

    f, err := os.OpenFile(configs["LOG_FILE_PATH"], os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    log.SetReportCaller(true)
    log.SetFormatter(&log.JSONFormatter{})
    log.SetLevel(log.InfoLevel)

    if err != nil {
        fmt.Println(err.Error() + " : " + configs["LOG_FILE_PATH"])
    }else{
        log.SetOutput(f)
    }

	// ======== Define Logs End ==========================
}


func ConnectRedis() *redis.Client {

	c := configs["REDIS_HOST"] + ":" + configs["REDIS_PORT"]

	rd := redis.NewClient(&redis.Options{
		Addr:     c,
		Password: configs["REDIS_PASS"], // no password set
		DB:       0,         // use default DB
	})

	res, err := rd.Ping().Result()

	if err != nil {
		log.WithFields(log.Fields{
			"REDIS_HOST": configs["REDIS_HOST"],
			"REDIS_PORT": configs["REDIS_PORT"],
		}).Fatal(err.Error())
	}

	if res != "PONG" {
		log.WithFields(log.Fields{
			"REDIS_HOST": configs["REDIS_HOST"],
			"REDIS_PORT": configs["REDIS_PORT"],
		}).Fatal("Invalid Redis Password")
	}

	return rd
}

func ConnectMysql() *sql.DB {

	conn := configs["MYSQL_USER"] + ":" + configs["MYSQL_PASS"] + "@tcp(" + configs["MYSQL_HOST"] + ":" + configs["MYSQL_PORT"] + ")/" + configs["MYSQL_DB"]
	db, err := sql.Open("mysql", conn)
	err = db.Ping()

	if err != nil {
		log.Fatal(err.Error())
	}

	no_conn, err := strconv.Atoi(configs["MYSQL_NOCONN"])

	if err != nil {
		log.Fatal(err.Error())
	}

	db.SetMaxIdleConns(no_conn)
	db.SetMaxOpenConns(no_conn)
	db.SetConnMaxLifetime(0)

	return db
}
