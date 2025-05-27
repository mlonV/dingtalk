package search

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DatabaseConfig holds the configuration for a database connection.
type DatabaseConfig struct {
	Name string
	DSN  string // Data Source Name, e.g., "user:password@tcp(host:port)/dbname"
}
type Result struct {
	DBInstance string
	Schemas    string
	TableName  string
	Columns    string
}

// Initialize with your actual database configurations.
var dbConfigs = []DatabaseConfig{
	//  // 测试用的DSN
	// {Name: "prometheus", DSN: "root:!QAZ2wsx#EDC@tcp(172.31.1.243:3306)/information_schema"},
	// {Name: "localdockermysql", DSN: "root:123456@tcp(127.0.0.1:3306)/information_schema"},

	// // 实际用的DSN
	{Name: "mysql1.karawangroup.com", DSN: "admin:hw-1YFM=E8+YWAR%fJLoxLC9@tcp(mysql1.karawangroup.com:3306)/information_schema"},
	{Name: "mysql2.karawangroup.com", DSN: "admin:hw-2YFM=E8+YWAR%fJLoxLC9@tcp(mysql2.karawangroup.com:3306)/information_schema"},
	{Name: "mysql3.karawangroup.com", DSN: "admin:hw-3YFM=E8+YWAR%fJLoxLC9@tcp(mysql3.karawangroup.com:3306)/information_schema"},
	{Name: "mysql4.karawangroup.com", DSN: "admin:hw-4YFM=E8+YWAR%fJLoxLC9@tcp(mysql4.karawangroup.com:3306)/information_schema"},
	{Name: "mysql5.karawangroup.com", DSN: "admin:hw-5YFM=E8+YWAR%fJLoxLC9@tcp(mysql5.karawangroup.com:3306)/information_schema"},
	{Name: "mysql6.karawangroup.com", DSN: "admin:hw-6YFM=E8+YWAR%fJLoxLC9@tcp(mysql6.karawangroup.com:3306)/information_schema"},
	{Name: "mysql7.karawangroup.com", DSN: "admin:hw-7YFM=E8+YWAR%fJLoxLC9@tcp(mysql7.karawangroup.com:3306)/information_schema"},
	{Name: "mysql8.karawangroup.com", DSN: "admin:hw-8YFM=E8+YWAR%fJLoxLC9@tcp(mysql8.karawangroup.com:3306)/information_schema"},
	{Name: "mysql9.karawangroup.com", DSN: "admin:hw-9YFM=E8+YWAR%fJLoxLC9@tcp(mysql9.karawangroup.com:3306)/information_schema"},
	{Name: "bigdata1", DSN: "admin:bigdata-1YFM=E8+YWAR%fJLoxLC9@tcp(bigdata-1.cqczrfqimrrs.eu-central-1.rds.amazonaws.com:3306)/information_schema"},
	{Name: "bigdata2", DSN: "admin:bigdata-2YFM=E8+YWAR%fJLoxLC9@tcp(bigdata-2.cqczrfqimrrs.eu-central-1.rds.amazonaws.com:3306)/information_schema"},
}

type Tables struct {
	TABLE_SCHEMA string `gorm:"column:TABLE_SCHEMA"`
	TABLE_NAME   string `gorm:"column:TABLE_NAME"`
}

// 明确指定表名为 user
func (Tables) TableName() string {
	return "TABLES"
}

type Columns struct {
	TABLE_SCHEMA string `gorm:"column:TABLE_SCHEMA"`
	TABLE_NAME   string `gorm:"column:TABLE_NAME"`
	COLUMN_NAME  string `gorm:"column:COLUMN_NAME"`
}

func (Columns) TableName() string {
	return "COLUMNS"

}

func queryDatabase(q string, methed string) ([]Result, []string) {

	var (
		finalResultsMu sync.Mutex // 保护 finalResults 的并发写入
		finalResults   []Result
		allErrorsMu    sync.Mutex // 保护 allErrors 的并发写入
		allErrors      []string
		wg             sync.WaitGroup
	)
	for _, config := range dbConfigs {
		wg.Add(1)

		go func(cfg DatabaseConfig) {
			defer wg.Done()
			var query_tbs []Tables
			var query_cos []Columns
			db, err := gorm.Open(mysql.Open(config.DSN), &gorm.Config{})
			if err != nil {
				allErrorsMu.Lock()
				allErrors = append(allErrors, err.Error())
				allErrorsMu.Unlock()
				return
			}
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()
			switch methed {
			case "schema":
				var schema Tables
				if err := db.Where("TABLE_SCHEMA = ?", q).First(&schema).Error; err == nil {
					finalResultsMu.Lock()
					finalResults = append(finalResults, Result{
						DBInstance: cfg.Name,
						Schemas:    schema.TABLE_SCHEMA,
					})
					finalResultsMu.Unlock()
				} else {
					allErrorsMu.Lock()
					allErrors = append(allErrors, fmt.Sprintf("%s: %v", cfg.Name, err.Error()))
					allErrorsMu.Unlock()
				}
			case "table":
				if err := db.Where("TABLE_NAME = ?", q).Find(&query_tbs).Error; err == nil {
					finalResultsMu.Lock()
					for _, tb := range query_tbs {
						finalResults = append(finalResults, Result{
							DBInstance: cfg.Name,
							Schemas:    tb.TABLE_SCHEMA,
							TableName:  tb.TABLE_NAME,
						})
					}
					finalResultsMu.Unlock()
				} else {
					allErrorsMu.Lock()
					allErrors = append(allErrors, fmt.Sprintf("%s: %v", cfg.Name, err.Error()))
					allErrorsMu.Unlock()
				}
			case "column":
				if err := db.Where("COLUMN_NAME = ?", q).Find(&query_cos).Error; err == nil {
					finalResultsMu.Lock()
					for _, co := range query_cos {
						finalResults = append(finalResults, Result{
							DBInstance: cfg.Name,
							Schemas:    co.TABLE_SCHEMA,
							TableName:  co.TABLE_NAME,
							Columns:    co.COLUMN_NAME,
						})
					}
					finalResultsMu.Unlock()
				} else {
					allErrorsMu.Lock()
					allErrors = append(allErrors, fmt.Sprintf("%s: %v", cfg.Name, err.Error()))
					allErrorsMu.Unlock()
				}
			}
		}(config)
	}
	wg.Wait()

	return finalResults, allErrors
}

// func queryDatabaseColumn(q string) ([]Result, []string) {

// 	var finalResults []Result
// 	var allErrors []string

// 	for _, config := range dbConfigs {
// 		var tbs []Columns
// 		db, err := gorm.Open(mysql.Open(config.DSN), &gorm.Config{})

// 		if err != nil {
// 			return finalResults, allErrors
// 		}
// 		defer func() {
// 			sqlDB, err := db.DB()
// 			if err == nil {
// 				sqlDB.Close()
// 			}
// 		}()

// 		for _, tb := range tbs {
// 			finalResults = append(finalResults, Result{
// 				DBInstance: config.Name,
// 				Schemas:    tb.TABLE_SCHEMA,
// 				TableName:  tb.TABLE_NAME,
// 				Columns:    tb.COLUMN_NAME,
// 			})
// 		}
// 	}

// 	return finalResults, allErrors
// }

func SearchTableHandler(c *gin.Context) {

	q := c.Param("table")
	finalResults, allErrors := queryDatabase(q, "table")

	c.JSON(http.StatusOK, gin.H{
		"error":   strings.Join(allErrors, "<br>"),
		"results": finalResults,
	})
}

func SearchSchemaHandler(c *gin.Context) {

	q := c.Param("schema")
	finalResults, allErrors := queryDatabase(q, "schema")

	c.JSON(http.StatusOK, gin.H{
		"error":   strings.Join(allErrors, "<br>"),
		"results": finalResults,
	})
}

func SearchColumnHandler(c *gin.Context) {
	q := c.Param("column")
	finalResults, allErrors := queryDatabase(q, "column")

	c.JSON(http.StatusOK, gin.H{
		"error":   strings.Join(allErrors, "<br>"),
		"results": finalResults,
	})

}

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"usage": gin.H{
			"mysql": gin.H{
				"db":      "/search/schema/[库名]",
				"table":   "/search/table/[表名]",
				"column":  "/search/column/[字段名]",
				"example": "http://localhost:8080/search/mysql/userdb",
			},
		},
	})
}
