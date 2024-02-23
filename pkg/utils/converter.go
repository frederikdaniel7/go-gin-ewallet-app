package utils

import (
	"strings"

	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/dto"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/entity"
)

func ConvertUserDetailtoJson(user entity.UserDetail) dto.UserDetail {
	var walletNumSB strings.Builder
	walletNumSB.WriteString("420")
	walletNumSB.WriteString(user.Wallet.WalletNumber)
	converted := dto.UserDetail{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Wallet: &dto.WalletPreview{
			WalletNumber: walletNumSB.String(),
			Balance:      user.Wallet.Balance,
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.DeletedAt.Valid {
		converted.DeletedAt = &user.DeletedAt.Time
	}
	return converted
}

func ConvertTokentoJson(token entity.PasswordToken) dto.PasswordToken {
	converted := dto.PasswordToken{
		Token:     token.Token,
		ExpiredAt: token.ExpiredAt,
	}
	return converted
}
