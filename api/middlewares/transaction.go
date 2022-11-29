package middlewares

import (
	"fmt"
	"github/abinav-07/mcq-test/constants"
	"github/abinav-07/mcq-test/infrastructure"
	"github/abinav-07/mcq-test/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DBTransactionMW struct {
	db infrastructure.Database
}

// Constructor
func NewDBTransactionMiddleware(
	db infrastructure.Database,
) DBTransactionMW {
	return DBTransactionMW{
		db: db,
	}
}

// Handle Transaction middleware
func (m DBTransactionMW) HandleDBTransaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		trxHandle := m.db.DB.Begin()
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Trx Commit  error", trxHandle.Error)
				trxHandle.Rollback()
			}
		}()

		ctx.Set(constants.DBTransaction, trxHandle)
		ctx.Next()

		//After Request
		//Checks if Response Status is Ok
		if utils.StatusInList(ctx.Writer.Status(), []int{http.StatusOK, http.StatusCreated}) {

			if err := trxHandle.Commit().Error; err != nil {
				fmt.Println("trx commit error: ", err)
			}
		} else {
			fmt.Println("rolling back transaction due to status code: ", ctx.Writer.Status())
			trxHandle.Rollback()
		}
	}
}
