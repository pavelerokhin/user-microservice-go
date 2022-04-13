package repository

//
//import (
//	"context"
//	"log"
//	"time"
//
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"
//)
//
//func NewMongodbRepo(l *log.Logger) (UserRepository, error) {
//	//l.Println("preparing MongoDB database")
//	//
//	//client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
//	//err = client.Connect(ctx)
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//defer client.Disconnect(ctx)
//	//
//	///*
//	//   List databases
//	//*/
//	//dbs, err := client.ListDatabaseNames(ctx, bson.M{})
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//
//	//l.Println("MongoDB database is ready")
//	//return &repo{DB: client, Logger: l}, nil
//}
