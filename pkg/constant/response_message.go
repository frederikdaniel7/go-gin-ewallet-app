package constant

const (
	ResponseMsgOK            = "OK"
	ResponseMsgErrorInternal = "Internal Server Error"
	ResponseMsgBadRequest    = "Bad Request"

	ResponseMsgInvalidBody              = "Invalid Body"
	ResponseMsgInvalidAmountInput       = "Invalid Amount Input"
	ResponseMsgInvalidWalletNumberInput = "Invalid Wallet Number Input"
	ResponsemsgInvalidDescriptionField  = "Invalid Description Field"
	ResponseMsgInvalidToken             = "Invalid Token"
	ResponseMsgInvalidPasswordInput     = "Password minimum length must be at least 8 characters"
	ResponseMsgInvalidEmailInput        = "Put the correct email format"

	ResponseMsgInsufficientFunds = "Insufficient Balance for Transfer"
	ResponseMsgUnauthorized      = "Unauthorized"
	ResponseMsgCreated           = "Created"

	ResponseMsgUserAlreadyExists = "user already exists"
	ResponseMsgUserDoesNotExist  = "user does not exist"

	ResponseMsgRecordDoesNotExist = "record does not exist"
	ResponseMsgErrorCredentials   = "email or password incorrect"
	ResponseMsgSQLError           = "bad request database error"
)
