package connection

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
)

var configs map[string]string
var configFile string =  "config.json"

func init(){
	loadConf()
}

func loadConf(){

	jsonFile, err := os.Open(configFile);
    if err != nil {
        fmt.Println(err)
    }
	defer jsonFile.Close()
    byteValue, _ := ioutil.ReadAll(jsonFile)
    json.Unmarshal([]byte(byteValue), &configs)

//  Define logging info

    f, _ := os.OpenFile(configs["LOG_FILE_PATH"], os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    log.SetOutput(f)
    log.SetReportCaller(true)
    log.SetFormatter(&log.JSONFormatter{})
    log.SetLevel(log.InfoLevel)


}


func ConnectRedis() *redis.Client {

	c := configs["REDIS_HOST"] + ":" + configs["REDIS_PORT"]

	rd := redis.NewClient(&redis.Options{
		Addr:     c,
		Password: configs["REDIS_AUTH"], // no password set
		DB:       0,         // use default DB
	})

	res, err := rd.Ping().Result()

	if err != nil {
		log.Fatal(err.Error())
	}

	if res != "PONG" {
		log.Fatal("Invalid Redis Password")
	}

	return rd
}

func ConnectMysql() *sql.DB {

	conn := configs["MYSQL_USER"] + ":" + configs["MYSQL_PASS"] + "@tcp(" + configs["MYSQL_HOST"] + ":" + configs["MYSQL_PORT"] + ")/" + configs["MYSQL_DB"]
	db, err := sql.Open("mysql", conn)
	err = db.Ping()

	if err != nil {
		log.Error("Exit. Cannot make connection with host: '" + configs["MYSQL_HOST"] + "' user: '" + configs["MYSQL_USER"] + "' pass: '" + configs["MYSQL_PASS"] + "' database: '" + configs["MYSQL_DB"] + "'")
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
