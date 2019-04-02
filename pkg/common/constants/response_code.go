package constants

type ResponseCode uint

const (
	ResponseCodeSuccess ResponseCode = 0

	ResponseCodeNotEnoughProduct ResponseCode = 1000
	ResponseCodeInvalidArgs      ResponseCode = 1001
	ResponseCodeNotFound         ResponseCode = 1002
	ResponseCodeUnknown         ResponseCode = 1003
)
