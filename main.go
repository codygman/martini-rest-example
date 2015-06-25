package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
	"log"
	"time"
	"github.com/robfig/config"
	"fmt"
)

type Log struct {
	Id int64 `db:"id" json:"id"`
	Logtime time.Time `db:"logtime" json:"logtime"`
	Latitude float64 `db:"latitude" json:"latitude"`
	Longitude float64 `db:"longitude" json:"longitude"`
	VenueID string `db:"venue_id" json:"venue_id"`
}


func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}


func main() {
	r := gin.Default()

	v1 := r.Group("api/v1")
	{
		v1.GET("/logs", GetLogs)
		v1.GET("/logs/:id", GetLog)
		v1.POST("/logs", PostLog)
		// v1.PUT("/logs/:id", UpdateLog)
		// v1.DELETE("/logs/:id", DeleteLog)
	}
	r.Run(":8080")	
}

var dbmap = initDb()

func initDb() *gorp.DbMap {
	c, err := config.ReadDefault("config.cfg")
	checkErr(err, "reading config failed")
	dbName, err := c.String("DEFAULT", "db_name")
	checkErr(err, "reading config failed")
	dbUserName, err := c.String("DEFAULT", "db_username")
	checkErr(err, "reading config failed")
	dbHost, err := c.String("DEFAULT", "db_host")
	checkErr(err, "reading config failed")
	dbPassword, err := c.String("DEFAULT", "db_password")
	checkErr(err, "reading config failed")
	
	connStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", dbUserName, dbPassword, dbHost, dbName)
	db, err := sql.Open("mysql", connStr)
	checkErr(err, "sql.Open failed")
	
	dbmap := &gorp.DbMap{Db: db, Dialect:           gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(Log{}, "log").SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create table failed")
	return dbmap
}


func GetLogs(c *gin.Context) {
	var log []Log
	_, err := dbmap.Select(&log, "SELECT * FROM log")

	if err == nil {
		c.JSON(200, log)
	} else {
		c.JSON(404, gin.H{"error": "no log(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/logs
}

func GetLog(c *gin.Context) {
	id := c.Params.ByName("id")
	var log Log
	err := dbmap.SelectOne(&log, "SELECT * FROM log WHERE id=?", id)

	if err == nil {
		log_id, _ := strconv.ParseInt(id, 0, 64)

		content := &Log{
			Id: log_id,
			Logtime: log.Logtime,
			Latitude: log.Latitude,
			Longitude: log.Longitude,
			VenueID: log.VenueID,
		}
		c.JSON(200, content)
	} else {
		c.JSON(404, gin.H{"error": "log not found"})
	}

	// curl -i http://localhost:8080/api/v1/logs/1
}

func PostLog(c *gin.Context) {
	var log Log
	c.Bind(&log)
	insert, _ := dbmap.Exec(`INSERT INTO log (logtime,latitude,longitude,venue_id) VALUES (?, ?, ?, ?)`, log.Logtime, log.Latitude, log.Longitude, log.VenueID);
	if insert != nil {
		log_id, err := insert.LastInsertId()
		if err == nil {
			content := &Log{
				Id: log_id,
				Logtime: log.Logtime,
				Latitude: log.Latitude,
				Longitude: log.Longitude,
				VenueID: log.VenueID,
			}
			c.JSON(201, content)
		} else {
			checkErr(err, "Insert failed")
		}
	}  else {
		c.JSON(422, gin.H{"error": "fields are empty"})
	}

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/logs
}

// func UpdateLog(c *gin.Context) {
// 	id := c.Params.ByName("id")
// 	var log Log
// 	err := dbmap.SelectOne(&log, "SELECT * FROM log WHERE id=?", id)

// 	if err == nil {
// 		var json Log
// 		c.Bind(&json)

// 		log_id, _ := strconv.ParseInt(id, 0, 64)

// 		log := Log{
// 			Id: log_id,
// 			Firstname: json.Firstname,
// 			Lastname: json.Lastname,
// 		}

// 		if log.Firstname != "" && log.Lastname != ""{
// 			_, err = dbmap.Update(&log)

// 			if err == nil {
// 				c.JSON(200, log)
// 			} else {
// 				checkErr(err, "Updated failed")
// 			}

// 		} else {
// 			c.JSON(422, gin.H{"error": "fields are empty"})
// 		}

// 	} else {
// 		c.JSON(404, gin.H{"error": "log not found"})
// 	}

// 	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/logs/1
// }


// func DeleteLog(c *gin.Context) {
// 	id := c.Params.ByName("id")

// 	var log Log
// 	err := dbmap.SelectOne(&log, "SELECT id FROM log WHERE id=?", id)

// 	if err == nil {
// 		_, err = dbmap.Delete(&log)

// 		if err == nil {
// 			c.JSON(200, gin.H{"id #" + id: " deleted"})
// 		} else {
// 			checkErr(err, "Delete failed")
// 		}

// 	} else {
// 		c.JSON(404, gin.H{"error": "log not found"})
// 	}

// 	// curl -i -X DELETE http://localhost:8080/api/v1/logs/1
// }
