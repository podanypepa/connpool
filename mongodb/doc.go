// Package mongodb provide connection pool to MongoDB
//
// Typycal use case:
//
//	mcp := mongodb.Create("mongodb://localhost:27017", 20)
//	conn := mcp.GetRandom()
//	collection := conn.Database("test").Collection("users")
//	mcp.Destroy()
package mongodb
