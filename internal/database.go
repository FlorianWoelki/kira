package internal

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoURI = "mongodb://database:27017"

type Database struct {
	ctx    context.Context
	client *mongo.Client
	db     *mongo.Database
}

func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	d.ctx = ctx
	d.client = client
	return nil
}

func (d *Database) InitDatabase() error {
	db := d.client.Database("kira")
	err := db.CreateCollection(d.ctx, "logs")
	if err != nil {
		return err
	}

	d.db = db
	return nil
}

func (d *Database) Insert(log primitive.E) (interface{}, error) {
	logs := d.db.Collection("logs")
	result, err := logs.InsertOne(d.ctx, bson.D{log})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *Database) Disconnect() {
	d.client.Disconnect(d.ctx)
}
