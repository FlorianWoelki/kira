package internal

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const mongoURI = "mongodb://database:27017"

type Database struct {
	client         *mongo.Client
	db             *mongo.Database
	collectionName string
}

// NewDatabase creates an empty database struct.
func NewDatabase() *Database {
	return &Database{}
}

// Connect connects the database to `mongodb://database:27017` and sets up the client to
// this database. It also initializes the `client` in the database struct.
func (d *Database) Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	// Check if the client is reachable.
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return err
	}

	d.client = client
	return nil
}

// CreateCollection creates a new collection in the database `kira`. It also sets the `db`
// database struct field. When the creation of the collection returns an error, it will
// return this error in this function.
func (d *Database) CreateCollection(collectionName string) error {
	d.collectionName = collectionName

	db := d.client.Database("kira")
	coll := db.Collection(d.collectionName)
	if coll == nil {
		err := db.CreateCollection(context.Background(), d.collectionName)
		if err != nil {
			return err
		}
	}

	d.db = db
	return nil
}

// Insert inserts a specific log to the database. If something went wrong while inserting
// the entry into the database, it will return the error.
func (d *Database) Insert(log interface{}) (interface{}, error) {
	logs := d.db.Collection(d.collectionName)
	result, err := logs.InsertOne(context.Background(), log)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Disconnect closes the connection to the database.
func (d *Database) Disconnect() {
	d.client.Disconnect(context.Background())
}
