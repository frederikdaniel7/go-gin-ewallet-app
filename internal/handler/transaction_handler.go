package handler

import (
	"net/http"
	"strings"

	"github.com/frederikdaniel7/go-gin-ewallet-app/internal/dto"
	"github.com/frederikdaniel7/go-gin-ewallet-app/internal/entity"
	"github.com/frederikdaniel7/go-gin-ewallet-app/internal/usecase"
	"github.com/frederikdaniel7/go-gin-ewallet-app/pkg/constant"
	"github.com/frederikdaniel7/go-gin-ewallet-app/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type TransactionHandler struct {
	transactionUseCase usecase.TransactionUseCase
}

func NewTransactionHandler(transactionUseCase usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{
		transactionUseCase: transactionUseCase,
	}
}

var mapSourceFunds = map[int]string{
	1: constant.SourceBankTransfer,
	2: constant.SourceCreditCard,
	3: constant.SourceCash,
	4: constant.SourceReward,
	5: constant.SourceWallet,
}

func (h *TransactionHandler) Transfer(ctx *gin.Context) {
	var body dto.Transfer
	userId := ctx.GetFloat64("id")
	if err := ctx.ShouldBindJSON(&body); err != nil {
		errType := utils.CheckError(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			dto.Response{
				Msg:  errType,
				Data: nil,
			})
		return
	}
	amountDecimal := decimal.NewFromFloat(body.Amount)
	transaction, err := h.transactionUseCase.MakeTransfer(ctx, &entity.Transfer{
		RecipientWalletNumber: utils.ConvertWalletNumber(body.RecipientWalletNumber),
		Amount:                amountDecimal,
		Descriptions:          body.Descriptions,
	}, int64(userId))
	if err != nil {
		ctx.Error(err)
		return
	}
	transactionJson := utils.ConvertTransactiontoJson(*transaction)
	ctx.JSON(http.StatusCreated, dto.Response{
		Msg: constant.ResponseMsgCreated,
		Data: dto.TransactionObj{
			Transaction: transactionJson,
		},
	})
}

func (h *TransactionHandler) TopUpBalance(ctx *gin.Context) {
	var body dto.TopUpBody
	userId := ctx.GetFloat64("id")
	if err := ctx.ShouldBindJSON(&body); err != nil {
		errType := utils.CheckError(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			dto.Response{
				Msg:  errType,
				Data: nil,
			})
		return
	}
	var sbDesc strings.Builder
	sbDesc.WriteString(constant.DescTopUp)
	sbDesc.WriteString(mapSourceFunds[body.SourceOfFunds])
	amountDecimal := decimal.NewFromFloat(body.Amount)
	transaction, err := h.transactionUseCase.TopUpBalance(ctx, &entity.Transfer{
		SourceOfFunds: mapSourceFunds[body.SourceOfFunds],
		Amount:        amountDecimal,
		Descriptions:  sbDesc.String(),
	}, int64(userId))
	if err != nil {
		ctx.Error(err)
		return
	}
	transactionJson := utils.ConvertTransactiontoJson(*transaction)
	ctx.JSON(http.StatusCreated, dto.Response{
		Msg: constant.ResponseMsgCreated,
		Data: dto.TransactionObj{
			Transaction: transactionJson},
	})
}

func (h *TransactionHandler) GetTransactions(ctx *gin.Context) {
	var params dto.TransactionFilter
	userId := ctx.GetFloat64("id")
	if err := ctx.ShouldBindQuery(&params); err != nil {
		errType := utils.CheckError(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			dto.Response{
				Msg:  errType,
				Data: nil,
			})
		return
	}
	convertedParams := utils.ConvertQueryJsonToObject(params)
	transactionsPage, err := h.transactionUseCase.GetTransaction(ctx, convertedParams, int64(userId))
	if err != nil {
		ctx.Error(err)
		return
	}

	transactionsJson := utils.ConvertTransactionstoJson(transactionsPage.Transactions)
	ctx.JSON(http.StatusOK, dto.Response{
		Msg: constant.ResponseMsgOK,
		Data: dto.TransactionPage{
			Transactions: dto.Transactions{
				Transactions: transactionsJson,
			},
			ItemCount:   transactionsPage.ItemCount,
			PageCount:   transactionsPage.PageCount,
			CurrentPage: transactionsPage.CurrentPage,
		},
	})

}
