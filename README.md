<p align="center"><img src="https://github.com/mongodb/mongo-go-driver/raw/master/etc/assets/mongo-gopher.png" width="250"></p>

<p align="center">
 <a href="https://goreportcard.com/report/go.mongodb.org/mongo-driver"><img src="https://goreportcard.com/badge/github.com/podanypepa/connpool"></a>
</p>

# connpool

Golang package for creating a connection pool to MongoDB.

## installation

```bash
go get github.com/podanypepa/connpool/mongodb
```

## usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/podanypepa/connpool/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {

	// create mongo connection pool with 20 independent connections
	cp, err := mongodb.Create("mongodb://localhost:27017", 20)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("connections pool length:", cp.Length())

	// get single connection from connection pool
	c, i := cp.GetRandom()
	fmt.Println("connection id:", i)

	// get data through geted connection
	getData(c, "test", "users")

	// close all created connections in pool
	cp.Destroy()
}

// getData is helper function for getting examples data from db
func getData(c *mongo.Client, db, coll string) {
	cur, err := c.
		Database(db).
		Collection(coll).
		Find(context.Background(), bson.D{})

	if err != nil {
		fmt.Println("ERR")
		fmt.Println(err)
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var result interface{}
		err := cur.Decode(&result)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result)
	}

}

```
