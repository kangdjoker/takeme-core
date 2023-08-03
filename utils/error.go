package utils

import (
	"fmt"
	"net/http"
	"runtime"

	log "github.com/sirupsen/logrus"
)

const (
	NotFoundCode     = 701
	ForbiddenCode    = 702
	UnauthorizedCode = 703

	// Bad request error
	InvalidHeader                      = 801
	InvalidRequestPayload              = 802
	UserAlreadyExist                   = 803
	UserNotFound                       = 804
	InvalidActivationCode              = 805
	InvalidSignature                   = 806
	InvalidLoginCode                   = 807
	InsufficientBalance                = 808
	InvalidTransactionType             = 809
	MinimumAmountTransaction           = 810
	RequestAlreadySubmit               = 811
	QRReadError                        = 812
	InvalidCode                        = 813
	InvalidPhoneNUmber                 = 814
	InvalidPIN                         = 815
	UpgradeAccountFirst                = 816
	InquiryAccountHolderNameNotFound   = 817
	FraudTransaction                   = 818
	UserLocked                         = 819
	BillerCodeNotFound                 = 820
	BillerBadRequest                   = 821
	CodeNotFoundInFusindo              = 822
	CorporateLocked                    = 823
	MaximumAmountTransaction           = 824
	TransactionNotFound                = 825
	RequestNotFound                    = 826
	MerchantNotFound                   = 827
	PreforgotAlreadyProcced            = 828
	InvalidCorporateKey                = 829
	BulkNotFound                       = 830
	BulkListEmpty                      = 831
	NationalIdentityRequired           = 832
	InvalidReceiverType                = 833
	ExternalIDNotFound                 = 834
	InvalidLevelAccessRevoke           = 835
	SprintParamError                   = 880
	SprintDeletedVA                    = 881
	InvalidNameFormat                  = 882
	UserAlreadyActive                  = 883
	InvalidBalanceAccess               = 884
	InvalidBalanceID                   = 885
	BankCodeNotFound                   = 886
	InvalidAccessType                  = 887
	AccessBalanceAlreadyHave           = 888
	RequestAccessBalanceNotFound       = 889
	RequestAccessBalanceAlreadyProcced = 890
	InvalidDeductTarget                = 891
	InvalidBalanceOwner                = 892
	InvalidFormatPIN                   = 893
	AccountNotFound                    = 894
	InvalidBalanceScope                = 895
	FaceNotRecognize                   = 896
	CurrencyError                      = 897
	OnlySupportOnIDR                   = 898
	WrongAcceptCardFee                 = 899

	// Internal server
	QueryFailed               = 901
	InsertFailed              = 902
	UpdateFailed              = 903
	DeleteFailed              = 904
	ReadEnvironmentFailed     = 905
	TwilioApiCallFailed       = 906
	SendGridApiCallFailed     = 907
	MongoDBConnectFailed      = 908
	DecodeJSONError           = 909
	BsonUnmarshalFailed       = 910
	DecodeTokenFailed         = 911
	EncodeTokenFailed         = 912
	XenditApiCallFailed       = 913
	TopupXenditCallbackFailed = 914
	CorporateNotFound         = 915
	DecryptError              = 916
	EKYCCallError             = 917
	BiometricFail             = 918
	StorageCloudFail          = 919
	InquiryBankAccountFail    = 920
	FusindoApiCallFailed      = 921
	CallbackError             = 922
	SprintApiCallFailed       = 923
	OYApiCallFailed           = 924
	APIIndonesiaRegionFailed  = 925
	NIKRequired               = 926
	InvalidCashoutCode        = 927
	InvalidFace               = 928
	MandatoryFieldIsEmpty     = 929
	TransactionAlreadyClaim   = 930
	MMBCApiCallFailed         = 931
	MMBCRetryTransctionFailed = 932
	MMBCBankNoutFound         = 933
	QontakAPICallFailed       = 934
	DBStartTransactionFailed  = 935
	StripeAPICallFail         = 936
	SaveFileFailed            = 937
	PermataApiCallFailed      = 938
)

type CustomError struct {
	HttpStatus  int    `json:"-"`
	Code        int    `json:"code"`
	Description string `json:"description"`
	Time        string `json:"time"`
}

func (error CustomError) Error() string {
	return error.Description
}

// This method have dynamic description message according to code and language
func ErrorBadRequest(errorCode int, logMessage string) error {
	_, fn, line, _ := runtime.Caller(1)
	log.Error(fmt.Sprintf("Bad request on %v at line %v (%v)", fn, line, logMessage))

	return CustomError{
		HttpStatus: http.StatusBadRequest,
		Code:       errorCode,
		Time:       TimestampNow(),
	}
}

func ErrorInternalServer(errorCode int, logMessage string) error {
	_, fn, line, _ := runtime.Caller(1)
	log.Error(fmt.Sprintf("Internal server error on %v at line %v (%v)", fn, line, logMessage))

	return CustomError{
		HttpStatus:  http.StatusInternalServerError,
		Code:        errorCode,
		Description: "Internal Server Error",
		Time:        TimestampNow(),
	}
}

func ErrorForbidden() error {
	log.Error("Forbidden operation")

	return CustomError{
		HttpStatus:  http.StatusForbidden,
		Code:        ForbiddenCode,
		Description: "Forbidden operation",
		Time:        TimestampNow(),
	}
}

func ErrorUnauthorized() error {
	log.Error("Unauthorized or invalid token")

	return CustomError{
		HttpStatus:  http.StatusUnauthorized,
		Code:        UnauthorizedCode,
		Description: "Unauthorized or invalid token",
		Time:        TimestampNow(),
	}
}
