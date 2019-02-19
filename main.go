package main

import (
	"github.com/abihf/delta"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/guregu/dynamo"
)

func main() {
	r := gin.Default()
	r.Use(errorHandler)
	r.GET("/tasks", listTasks)
	r.POST("/task", addTask)
	r.GET("/task/:taskid", getTask)
	r.PATCH("/task/:taskid", editTask)
	r.DELETE("/task/:taskid", deleteTask)

	// start lambda handling if it runs on lambda
	// otherwise start http server on port 3000
	delta.ServeOrStartLambda(":3000", r)
}

func errorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
			c.JSON(-1, c.Errors) // -1 == not override the current error code
	}
}

var table dynamo.Table

func init() {
	db := dynamo.New(session.New())
	table = db.Table("LambdaGoDemo")
}

type task struct {
	ID   string    `json:"id"`
	Time string `json:"time"`

	Content string    `json:"content" dynamo:"Content"`
	Done    string `json:"done" dynamo:"Done"`
}

func listTasks(c *gin.Context) {
	var tasks []task
	err := table.Scan().AllWithContext(c, &tasks)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
		c.JSON(200, tasks)
}

func getTask(c *gin.Context) {
	id := c.Param("taskid")
	var t task
	err := table.Get("ID", id).OneWithContext(c, &t)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, t)
}

func addTask(c *gin.Context) {
	var t task
	err := c.BindJSON(&t)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	err = table.Put(t).RunWithContext(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, t)
}

func editTask(c *gin.Context) {
	var t task
	err := c.BindJSON(&t)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	id := c.Param("taskid")
	t.ID = id
	err = table.Update("ID", &t).RunWithContext(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, t)
}

func deleteTask(c *gin.Context) {
	id := c.Param("taskid")
	err := table.Delete("ID", id).RunWithContext(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
}