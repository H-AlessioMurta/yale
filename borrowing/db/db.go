package db

import(
	"context"
	"time"
	"github.com/google/uuid"
	"yale/borrowing/graph/model"
	l "yale/borrowing/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DB_URI = "mongodb://localhost:27017"
	mongoDB = "mongoDB"
	collection =  "Borrows"
)

type DB struct {
	client *mongo.Client
}

func Connect() *DB {	
	client, err := mongo.NewClient(options.Client().ApplyURI(DB_URI))
	l.CheckErr(err)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	err = client.Ping(ctx, nil)
	l.CheckErr(err)
	l.LogInfo("Connected to Mongo db")
	return &DB{
		client: client,
	}
}

func (db *DB)NewBorrow(input *model.BorrowedCreate) *model.Borrowed {
	t := time.Now()
	e := t.AddDate(0,0,15)
	u := uuid.New()
	newBorrow := model.Borrowed{
		IDBorrowing: u.String(),
		IDCustomer: *input.IDCustomer,
		IDBook: *input.IDBook,
		Starting: t,
		Expiring: e,
		Returned:false,
	}
	collection := db.client.Database(mongoDB).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, newBorrow)
	l.CheckErr(err)
	l.LogResponseBorrowing(&newBorrow)
	return &newBorrow
}

func (db *DB)Returnedbook(id *string) (*model.Borrowed, error) {
	newID:=*id
	filter := bson.M{"idborrowing":newID}
    update := bson.D{primitive.E{
		Key: "$set", Value: bson.D{ primitive.E{Key: "returned", Value: true},
    }}}
	collection := db.client.Database(mongoDB).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var aborrow model.Borrowed
	_, err := collection.UpdateOne(ctx,filter,update)
	l.CheckErr(err)
	l.LogResponseBorrowing(&aborrow)
	return &aborrow,err
}

func (db *DB)Borrowsnotreturned() []*model.Borrowed {
	filter := bson.D{
        primitive.E{Key: "returned", Value: false},
    }
	collection := db.client.Database(mongoDB).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	notReturned, err := collection.Find(ctx, filter)
	/* Cursor provides a stream of documents through which you can iterate and decode one at a time.
	 Once a Cursor has been exhausted, you should close the Cursor*/
	l.CheckErr(err)	
	var borrows []*model.Borrowed
	for notReturned.Next(ctx) {
		var aborrow model.Borrowed
		err := notReturned.Decode(&aborrow)
		l.CheckErr(err)
		borrows = append(borrows, &aborrow)
		l.LogResponseBorrowing(&aborrow)
	}
	return borrows
}

func (db *DB) Borrowsexpriring() []*model.Borrowed {
	t:= time.Now()
	filter := bson.D{
        primitive.E{Key: "returned", Value: false},
    }
	collection := db.client.Database(mongoDB).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	notReturned, err := collection.Find(ctx,filter)
	/* Cursor provides a stream of documents through which you can iterate and decode one at a time.
	 Once a Cursor has been exhausted, you should close the Cursor*/
	l.CheckErr(err)	
	var borrows []*model.Borrowed
	for notReturned.Next(ctx) {
		var aborrow model.Borrowed
		err := notReturned.Decode(&aborrow)
		l.CheckErr(err)
		if t.Before(aborrow.Expiring){
			borrows = append(borrows, &aborrow)
			l.LogResponseBorrowing(&aborrow)
		}
	}
	return borrows
}

func (db *DB) Borrowsforcustomer(id string) []*model.Borrowed {
	
	filter := bson.D{
        primitive.E{Key: "IDCustomer", Value: id},
    }
	collection := db.client.Database(mongoDB).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	bxc, err := collection.Find(ctx,filter)
	/* Cursor provides a stream of documents through which you can iterate and decode one at a time.
	 Once a Cursor has been exhausted, you should close the Cursor*/
	l.CheckErr(err)	
	var borrows []*model.Borrowed
	for bxc.Next(ctx) {
		var aborrow model.Borrowed
		err := bxc.Decode(&aborrow)
		l.CheckErr(err)
		borrows = append(borrows, &aborrow)
		l.LogResponseBorrowing(&aborrow)
	}
	return borrows
}

func (db *DB) Borrowsforbook(id string) []*model.Borrowed {

	filter := bson.D{
        primitive.E{Key: "IDBook", Value: id},
    }
	collection := db.client.Database(mongoDB).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	bxb, err := collection.Find(ctx,filter)
	/* Cursor provides a stream of documents through which you can iterate and decode one at a time.
	 Once a Cursor has been exhausted, you should close the Cursor*/
	l.CheckErr(err)	
	var borrows []*model.Borrowed
	for bxb.Next(ctx) {
		var aborrow *model.Borrowed
		err := bxb.Decode(&aborrow)
		l.CheckErr(err)
		l.LogInfo(aborrow.IDBorrowing)
		borrows = append(borrows, aborrow)
		l.LogResponseBorrowing(aborrow)
	}
	return borrows
}



