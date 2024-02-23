package utils

import (
	"strings"

	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/constant"
)

func CheckError(errMap string) string {
	if strings.Contains(errMap, constant.ErrEmailField) {
		return constant.ResponseMsgInvalidEmailInput
	} else if strings.Contains(errMap, constant.ErrPasswordField) {
		return constant.ResponseMsgInvalidPasswordInput
	} else if strings.Contains(errMap, constant.ErrAmountField) || strings.Contains(errMap, constant.ErrAmountUpercaseField) {
		return constant.ResponseMsgInvalidAmountInput
	} else if strings.Contains(errMap, constant.ErrRecipientWalletNumberField) {
		return constant.ResponseMsgInvalidWalletNumberInput
	} else if strings.Contains(errMap, constant.ErrDescriptionField) || strings.Contains(errMap, "description") {
		return constant.ResponsemsgInvalidDescriptionField
	}
	return constant.ResponseMsgInvalidBody
}
