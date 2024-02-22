package handler

import (
	"net/http"

	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/dto"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/entity"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/usecase"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/constant"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var body dto.CreateUser
	if err := ctx.ShouldBindJSON(&body); err != nil {
		errType := utils.CheckError(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			dto.Response{
				Msg:  errType,
				Data: nil,
			})
		return
	}
	user, err := h.userUseCase.RegisterUser(ctx, &entity.User{
		Email:    body.Email,
		Password: body.Password,
		Name:     body.Name,
	})
	if err != nil {
		ctx.Error(err)
		return
	}
	userJson := utils.ConvertUserDetailtoJson(*user)
	ctx.JSON(http.StatusCreated, dto.Response{
		Msg: constant.ResponseMsgCreated,
		Data: dto.UserObj{
			User: userJson,
		},
	})
}
