package handler

import (
	"net/http"

	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/dto"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/entity"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/usecase"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/constant"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/utils"
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
		Msg:  constant.ResponseMsgCreated,
		Data: transactionJson,
	})
}
