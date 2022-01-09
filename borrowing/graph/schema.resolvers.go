package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	dbmong "yale/borrowing/db"
	"yale/borrowing/graph/generated"
	"yale/borrowing/graph/model"
	l "yale/borrowing/logger"
	"yale/borrowing/router"
)

func (r *mutationResolver) BorrowCreate(ctx context.Context, data model.BorrowedCreate) (*model.Borrowed, error) {
	l.LogGraph("Mutation Ask: \033[32m[BorrowCreate]\033[0m: with inputs:\nidbook: " + *data.IDBook + "\nidCustomer: " + *data.IDCustomer)
	if router.CheckBook(*data.IDBook) && router.CheckCustomer(*data.IDCustomer) {
		return db.NewBorrow(&data), nil
	} else {
		var r string
		if router.CheckBook(*data.IDBook) {
			r = r + "There is no book with this id:\n" + *data.IDBook + "\n"
		}
		if router.CheckCustomer(*data.IDCustomer) {
			r = r + "There is no book with this id:\n" + *data.IDCustomer + "\n"
		}
		return nil, errors.New(r)
	}
}

func (r *mutationResolver) Returnedbook(ctx context.Context, id *string) (*model.Borrowed, error) {
	l.LogGraph("Mutation Ask: \033[32m[Returnedbook]\033[0m with inputs:\n" + *id)
	return db.Returnedbook(id)
}

func (r *mutationResolver) PostBook(ctx context.Context, title *string, authors *string) (*model.Book, error) {
	//return router.PostBook(*title,*authors)
	payload := map[string]string{"title": *title, "authors": *authors}
	l.LogGraph("Mutation Ask: \033[32m[PostBook]\033[0m with inputs:\ntitle: " + *title + "\nauthors: " + *authors)
	return router.BookRequestHandler(payload, "POST", "")
}

func (r *mutationResolver) PutBook(ctx context.Context, id *string, title *string, authors *string) (*model.Book, error) {
	l.LogGraph("Mutation Ask: \033[32m[PutBook]\033[0m with inputs:\ntitle: " + *title + "\nauthors: " + *authors)
	payload := map[string]string{"title": *title, "authors": *authors}
	return router.BookRequestHandler(payload, "PUT", *id)
}

func (r *mutationResolver) DeleteBook(ctx context.Context, id *string) (*string, error) {
	l.LogGraph("Mutation Ask: \033[32m[DeleteBook]\033[0m with inputs:\n" + *id)
	return router.Eraser(router.BOOKAPI, *id)
}

func (r *mutationResolver) PostCustomer(ctx context.Context, name *string, surname *string, nin *string) (*model.Customer, error) {
	l.LogGraph("Mutation Ask: \033[32m[PostCustomer]\033[0m with inputs:\nname: " + *name + "\nsurname: " + *surname + "\nnin: " + *nin)
	payload := map[string]string{"name": *name, "surname": *surname, "nin": *nin}
	return router.CustomerRequestHandler(payload, "POST", "")
}

func (r *mutationResolver) PutCustomer(ctx context.Context, id *string, name *string, surname *string, nin *string) (*model.Customer, error) {
	l.LogGraph("Mutation Ask: \033[32m[PutCustomer]\033[0m with inputs:\nname: " + *name + "\nsurname: " + *surname + "\nnin: " + *nin)
	payload := map[string]string{"name": *name, "surname": *surname, "nin": *nin}
	return router.CustomerRequestHandler(payload, "PUT", *id)
}

func (r *mutationResolver) DeleteCustomer(ctx context.Context, id *string) (*string, error) {
	l.LogGraph("Mutation Ask: \033[32m[DeleteCustomer]\033[0m with inputs:\n" + *id)
	return router.Eraser(router.CUSTOMERAPI, *id)
}

func (r *queryResolver) Borrows(ctx context.Context) ([]*model.Borrowed, error) {
	return db.Borrows()
}

func (r *queryResolver) Borrow(ctx context.Context, id *string) (*model.Borrowed, error) {
	return db.Borrow(id)
}

func (r *queryResolver) Books(ctx context.Context) ([]*model.Book, error) {
	l.LogGraph("Query Ask: \033[35m[All books]\033[0m")
	return router.GetBooks()
}

func (r *queryResolver) Book(ctx context.Context, id string) (*model.Book, error) {
	l.LogGraph("Query Ask: \033[35m[A book]\033[0m with input:" + id)
	return router.GetBook(id)
}

func (r *queryResolver) Customers(ctx context.Context) ([]*model.Customer, error) {
	l.LogGraph("Query Ask: \033[35m[All Customer]\033[0m")
	return router.GetCustomers()
}

func (r *queryResolver) Customer(ctx context.Context, id string) (*model.Customer, error) {
	l.LogGraph("Query Ask: \033[35m[A Customer]\033[0m with input:" + id)
	return router.GetCustomer(id)
}

func (r *queryResolver) Borrowsnotreturned(ctx context.Context) ([]*model.Borrowed, error) {
	l.LogGraph("Query Ask: \033[35m[Borrowsnotreturned]\033[0m")
	return db.Borrowsnotreturned(), nil
}

func (r *queryResolver) Borrowsforcustomer(ctx context.Context, id string) ([]*model.Borrowed, error) {
	l.LogGraph("Query Ask: \033[35m[Borrowsforcustomer]\033[0m with input:" + id)
	return db.Borrowsforcustomer(id), nil
}

func (r *queryResolver) Borrowsforbook(ctx context.Context, id string) ([]*model.Borrowed, error) {
	l.LogGraph("Query Ask: \033[35m[Borrowsforbook]\033[0m with input:" + id)
	return db.Borrowsforbook(id), nil
}

func (r *queryResolver) GetIDBook(ctx context.Context, title string, authors string) (*model.Book, error) {
	l.LogGraph("Query Ask: \033[35m[GetIDBook]\033[0m with input:\ntitle: " + title + " authors: " + authors)
	return router.GetIDBook(title, authors)
}

func (r *queryResolver) GetIDCustomer(ctx context.Context, name string, surname string, nin string) (*model.Customer, error) {
	l.LogGraph("Query Ask: \033[35m[GetIDBook]\033[0m with input:\nname: " + name + " surname: " + surname + " nin: " + nin)
	return router.GetIDCustomer(name, surname, nin)
}

func (r *queryResolver) GetIDBorrowed(ctx context.Context, idb string, idc string) (*model.Borrowed, error) {
	return db.Get_ID_Borrowed(idb, idc)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
var db = dbmong.Connect()
