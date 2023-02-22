package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(Ctx context.Context, arg TransferTxParams) (TransferResult, error)
}

// Store provides all functions to execute db queries and transactions
type SQLStore struct {
	// to support transaction query -> composition
	// Golang에서는 상속보다 기능적으로 확장(extend)가 선호되는 방법
	// All individual query functions provided by Queries will be available to Store
	*Queries
	db *sql.DB
}

// Create new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// Executes a function within a database transaction
func (Store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := Store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// input the parameters of the transfer transition
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// second bracket means that we're creating a new empty object of that type
var txKey = struct{}{}

// transfer record, add account entries, update accounts' balance
func (Store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferResult, error) {
	var result TransferResult

	err := Store.execTx(ctx, func(q *Queries) error {
		var err error

		// transaction add
		txName := ctx.Value(txKey)

		fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			// update accoounts' balance
			// fmt.Println(txName, "get account 1")
			// account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
			// if err != nil {
			// 	return err
			// }
			// fmt.Println(txName, "update account 1")
			// result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID:     arg.FromAccountID,
			// 	Amount: -arg.Amount,
			// })
			// if err != nil {
			// 	return err
			// }
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			// GetAccount -> GetAccountForUpdate
			// fmt.Println(txName, "get account 2")
			// account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
			// if err != nil {
			// 	return err
			// }
			// fmt.Println(txName, "update account 2")
			// result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID:     arg.ToAccountID,
			// 	Amount: arg.Amount,
			// })
			// if err != nil {
			// 	return err
			// }
		} else {
			// fmt.Println(txName, "update account 2")
			// result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID:     arg.ToAccountID,
			// 	Amount: arg.Amount,
			// })
			// if err != nil {
			// 	return err
			// }
			// fmt.Println(txName, "update account 1")
			// result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID:     arg.FromAccountID,
			// 	Amount: -arg.Amount,
			// })
			// if err != nil {
			// 	return err
			// }
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return // == return account1, account2, error
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})

	return
}
