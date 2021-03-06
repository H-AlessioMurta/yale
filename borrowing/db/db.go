/*
*This package will define, connect, exec mutations and querry for our borrowed model.
*All borrows will stored in the same collection of mongo
*/


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
	"strconv"
)

const (
	//Different version of trying the connection to our Database on docker's build or k8's pod
	//DB_URI = "mongodb://localhost:27017"
	DB_URI = "mongodb://borrows-mongodb:27017"
	mongoDB = "mongoDB"
	collection =  "Borrows"
)

type DB struct {
	client *mongo.Client
}

func Connect() *DB {	
	credential := options.Credential{
		AuthSource: "admin",
		Username: "root",
		Password: "root",
		//Databases: "Borrows",
	}
	
	opt := options.Client().ApplyURI(DB_URI).SetAuth(credential)
	client, err := mongo.NewClient(opt)
	l.CheckErr(err)
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ctx := context.Background()
	err = client.Connect(ctx)
	l.CheckErr(err)
	err = client.Ping(ctx, nil)
	l.CheckErr(err)
	l.LogInfo("Connected to Mongo db")
	return &DB{
		client: client,
	}
}

// A custom function to connect not to a local/clustered mongo but to atlas mongodb
func ConnectToAtlas() *DB {	

	opt := options.Client().ApplyURI("mongodb+srv://cluster0.gamsi.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	client, err := mongo.NewClient(opt)
	l.CheckErr(err)
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ctx := context.Background()
	err = client.Connect(ctx)
	l.CheckErr(err)
	err = client.Ping(ctx, nil)
	l.CheckErr(err)
	l.LogInfo("Connected to Mongo Atlas")
	return &DB{
		client: client,
	}
}


//Add a new borrow
func (db *DB)NewBorrow(input *model.BorrowedCreate) *model.Borrowed {
	t := time.Now()
	e := t.AddDate(0,0,15)// this line will add 15 days to our time.Time date
	u := uuid.New()// a new random uuid
	newBorrow := model.Borrowed{
		IDBorrowing: u.String(),
		IDCustomer: *input.IDCustomer,
		IDBook: *input.IDBook,
		Starting: t,
		Expiring: e,
		Returned:false,
	}
	collection := db.client.Database(mongoDB).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)//trying to connect to mongo
	defer cancel()
	_, err := collection.InsertOne(ctx, newBorrow) //exec the insert
	l.CheckErr(err)
	l.LogResponseBorrowing(&newBorrow)
	return &newBorrow
}

//Set true returned a specific idborrowing's borrow 
func (db *DB)Returnedbook(id *string) (*model.Borrowed, error) {
	newID:=*id
	filter := bson.M{"idborrowing":newID}// Using mongodriver for creating bson on golang
    update := bson.D{primitive.E{
		Key: "$set", Value: bson.D{ primitive.E{Key: "returned", Value: true},// condition of setting
    }}}
	collection := db.client.Database(mongoDB).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var aborrow model.Borrowed
	_, err := collection.UpdateOne(ctx,filter,update)//exec on collection 
	l.CheckErr(err)
	l.LogResponseBorrowing(&aborrow)
	return &aborrow,err
}
//fetching a borrow field knowing it's idborrowing
func (db *DB)Borrow(id *string) (*model.Borrowed, error) {
	newID:=*id
	filter := bson.M{"idborrowing":newID}
	collection := db.client.Database(mongoDB).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var aborrow model.Borrowed
	err := collection.FindOne(ctx,filter).Decode(&aborrow)
	l.CheckErr(err)
	l.LogResponseBorrowing(&aborrow)
	return &aborrow,err
}


//fetching all borrows
func (db *DB)Borrows() ([]*model.Borrowed,error) {
	collection := db.client.Database(mongoDB).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	b, err := collection.Find(ctx,bson.M{})
	/* Cursor provides a stream of documents through which you can iterate and decode one at a time.
	 Once a Cursor has been exhausted, you should close the Cursor*/
	l.CheckErr(err)	
	var borrows []*model.Borrowed
	err = b.All(ctx,&borrows)
	l.CheckErr(err)
	return borrows,err
}

//fetching all non returned borrows
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

//fetching all borrows for the same customer
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

func (db *DB)RemoveBorrow(id *string) (*string, error) {
	
	newID:=*id
	filter := bson.M{"idborrowing":newID}
	collection := db.client.Database(mongoDB).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.DeleteOne(ctx,filter)
	l.CheckErr(err)
	response:= strconv.Itoa(int(res.DeletedCount))
	l.LogResponse("Deleted "+response+" element/s")
	return &response,err
}

// fetching all borrows with same book
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

//returning a specifi borrowing for the same idbook, idcustomer, this is not intended for be a feauter of yale's users 
func (db *DB)Get_ID_Borrowed(idb string, idc string)(*model.Borrowed,error){
	filter := bson.D{
        primitive.E{Key: "IDBook", Value: idb},
		primitive.E{Key: "IDCustomer", Value: idc},
    }
	collection := db.client.Database(mongoDB).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var aborrow model.Borrowed
	err := collection.FindOne(ctx,filter).Decode(&aborrow)
	l.CheckErr(err)
	l.LogResponseBorrowing(&aborrow)
	return &aborrow,err
}



