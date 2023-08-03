package service

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/domain/dto"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/basic"
	"github.com/kangdjoker/takeme-core/utils/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserCreate(corporate domain.Corporate, email string, phoneNumber string, fullName string,
	session mongo.SessionContext) (domain.User, error) {

	activationCode := utils.GenerateShortCode()
	verificationCode := utils.GenerateUUID()
	attempt, _ := strconv.Atoi(os.Getenv("SECURITY_ATTEMPT"))

	model := domain.User{
		CorporateID:      corporate.GetDocumentID(),
		Email:            email,
		PhoneNumber:      phoneNumber,
		FullName:         fullName,
		PIN:              "",
		LoginCode:        "-",
		ActivationCode:   activationCode,
		VerificationCode: verificationCode,
		Active:           false,
		Verified:         false,
		AccessAttempt:    int8(attempt), // Default value
		LoginAttempt:     int8(attempt),
		Audit: domain.Audit{
			CreatedTime: time.Now().Format(os.Getenv("TIME_FORMAT")),
			UpdatedTime: time.Now().Format(os.Getenv("TIME_FORMAT")),
		},
		Avatar:    "",
		Pending:   false,
		FaceAsPIN: false,
	}

	err := UserSaveOne(&model, session)
	if err != nil {
		return model, err
	}

	return model, nil
}

func UserCreateUnpending(paramLog *basic.ParamLog, corporate domain.Corporate, userPending domain.User, email string, phoneNumber string, fullName string,
	session mongo.SessionContext) (domain.User, error) {

	activationCode := utils.GenerateShortCode()
	verificationCode := utils.GenerateUUID()
	attempt, _ := strconv.Atoi(os.Getenv("SECURITY_ATTEMPT"))

	userPending.Email = email
	userPending.PhoneNumber = phoneNumber
	userPending.FullName = fullName
	userPending.ActivationCode = activationCode
	userPending.VerificationCode = verificationCode
	userPending.PIN = ""
	userPending.LoginCode = "-"
	userPending.Active = false
	userPending.Verified = false
	userPending.LoginAttempt = int8(attempt)
	userPending.AccessAttempt = int8(attempt)
	userPending.Audit = domain.Audit{
		CreatedTime: time.Now().Format(os.Getenv("TIME_FORMAT")),
		UpdatedTime: time.Now().Format(os.Getenv("TIME_FORMAT")),
	}
	userPending.Avatar = ""
	userPending.FaceAsPIN = false

	err := UserUpdateOne(paramLog, &userPending, session)
	if err != nil {
		return domain.User{}, err
	}

	return userPending, nil
}

func UserActivate(paramLog *basic.ParamLog, user *domain.User, session mongo.SessionContext) error {

	user.Active = true
	user.ActivationCode = ""
	user.Pending = false

	err := UserUpdateOne(paramLog, user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserSaveOne(model *domain.User, session mongo.SessionContext) error {
	err := database.SessionSaveOne(model, session)
	if err != nil {
		return err
	}

	return nil
}

func UserUpdateOne(paramLog *basic.ParamLog, model *domain.User, session mongo.SessionContext) error {
	err := database.SessionUpdateOne(paramLog, model, session)
	if err != nil {
		return err
	}

	return nil
}

func UserUpdateOneNoSession(paramLog *basic.ParamLog, model *domain.User) error {
	err := database.UpdateOne(paramLog, domain.USER_COLLECTION, model)
	if err != nil {
		return err
	}

	return nil
}

func UserByID(paramLog *basic.ParamLog, ID string, session mongo.SessionContext) (domain.User, error) {
	model := domain.User{}
	cursor := database.SessionFindOneByID(domain.USER_COLLECTION, ID, session)
	err := cursor.Decode(&model)
	if err != nil {
		return domain.User{}, utils.ErrorInternalServer(paramLog, utils.QueryFailed, "Query failed or cannot decode")
	}

	return model, nil
}

func UserByIDNoSession(paramLog *basic.ParamLog, ID string) (domain.User, error) {
	model := domain.User{}
	cursor := database.FindOneByID(domain.USER_COLLECTION, ID)
	err := cursor.Decode(&model)
	if err != nil {
		return domain.User{}, utils.ErrorInternalServer(paramLog, utils.QueryFailed, "Query failed or cannot decode")
	}

	return model, nil
}

func UserByPhoneNumber(paramLog *basic.ParamLog, corporateID primitive.ObjectID, phoneNumber string, session mongo.SessionContext) (domain.User, error) {
	model := domain.User{}
	query := bson.M{"phone_number": phoneNumber, "corporate_id": corporateID, "pending": false}
	cursor := database.SessionFindOne(domain.USER_COLLECTION, query, session)
	cursor.Decode(&model)

	if model.PhoneNumber == "" {
		return domain.User{}, utils.ErrorBadRequest(paramLog, utils.UserNotFound, "User not found")
	}

	return model, nil
}

func UserByPhoneNumberWithoutSession(paramLog *basic.ParamLog, corporateID primitive.ObjectID, phoneNumber string) (domain.User, error) {
	model := domain.User{}
	query := bson.M{"phone_number": phoneNumber, "corporate_id": corporateID, "pending": false}
	cursor := database.FindOne(domain.USER_COLLECTION, query)
	cursor.Decode(&model)

	if model.PhoneNumber == "" {
		return domain.User{}, utils.ErrorBadRequest(paramLog, utils.UserNotFound, "User not found")
	}

	return model, nil
}

func ValidateUserNotRegisteredYet(paramLog *basic.ParamLog, corporate domain.Corporate, phoneNumber string, email string,
	session mongo.SessionContext) (bool, domain.User, error) {

	isPending := false
	var model domain.User

	query := bson.M{
		"corporate_id": corporate.GetDocumentID(),
		"$or":          bson.A{bson.M{"email": email}, bson.M{"phone_number": phoneNumber}},
	}
	cursor := database.SessionFindOne(domain.USER_COLLECTION, query, session)
	cursor.Decode(&model)

	if model.Pending == true {
		isPending = true
		return isPending, model, nil
	}

	if model.PhoneNumber != "" {
		return isPending, model, utils.ErrorBadRequest(paramLog, utils.UserAlreadyExist, "User already exist")
	}

	return isPending, model, nil
}

func UserByIDWithValidation(paramLog *basic.ParamLog, ID string, validations []func(paramLog *basic.ParamLog, user domain.User) error) (domain.User, error) {
	model := domain.User{}
	cursor := database.FindOneByID(domain.USER_COLLECTION, ID)
	err := cursor.Decode(&model)
	if err != nil {
		return domain.User{}, utils.ErrorInternalServer(paramLog, utils.QueryFailed, "Query failed or cannot decode")
	}

	for _, element := range validations {
		err := element(paramLog, model)
		if err != nil {
			return domain.User{}, err
		}
	}

	return model, nil
}

func UserReduceAccessAttempt(paramLog *basic.ParamLog, userID string, session mongo.SessionContext) error {

	user, err := UserByID(paramLog, userID, session)
	if err != nil {
		return err
	}

	err = ValidateUserLocked(paramLog, user)
	if err != nil {
		return err
	}

	remaining := user.AccessAttempt
	if remaining <= 0 {
		user.AccessAttempt = 0
		user.Active = false
	} else {
		user.AccessAttempt -= 1
	}

	err = UserUpdateOne(paramLog, &user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserGenerateLoginCode(paramLog *basic.ParamLog, user *domain.User, session mongo.SessionContext) error {
	user.LoginAttempt -= 1
	user.LoginCode = utils.GenerateShortCode()

	err := UserUpdateOne(paramLog, user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserRefreshAttempt(paramLog *basic.ParamLog, user *domain.User, session mongo.SessionContext) error {
	attempt, _ := strconv.Atoi(os.Getenv("SECURITY_ATTEMPT"))
	user.LoginAttempt = int8(attempt)
	user.AccessAttempt = int8(attempt)
	user.LoginCode = "-"

	err := UserUpdateOne(paramLog, user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserLock(paramLog *basic.ParamLog, user *domain.User, session mongo.SessionContext) error {

	user.Active = false
	user.AccessAttempt = 0
	user.LoginAttempt = 0

	err := UserUpdateOne(paramLog, user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserUnlock(paramLog *basic.ParamLog, user *domain.User, session mongo.SessionContext) error {
	attempt, _ := strconv.Atoi(os.Getenv("SECURITY_ATTEMPT"))
	user.Active = true
	user.AccessAttempt = int8(attempt)
	user.LoginAttempt = int8(attempt)

	err := UserUpdateOne(paramLog, user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserSavePIN(paramLog *basic.ParamLog, user *domain.User, pin string, session mongo.SessionContext) error {

	pin, err := utils.RSADecrypt(pin)
	if err != nil {
		return utils.ErrorInternalServer(paramLog, utils.DecryptError, err.Error())
	}

	user.PIN = pin

	err = UserUpdateOne(paramLog, user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserGenerateChangePINCode(paramLog *basic.ParamLog, user *domain.User, pin string, session mongo.SessionContext) (string, error) {

	pin, err := utils.RSADecrypt(pin)
	if err != nil {
		return "", utils.ErrorInternalServer(paramLog, utils.DecryptError, err.Error())
	}

	user.ChangePIN = pin
	user.ChangePINCode = utils.GenerateShortCode()

	err = UserUpdateOne(paramLog, user, session)
	if err != nil {
		return "", err
	}

	return user.ChangePINCode, nil
}

func UserChangePIN(paramLog *basic.ParamLog, user *domain.User, session mongo.SessionContext) error {
	user.PIN = user.ChangePIN
	user.ChangePIN = " "
	user.ChangePINCode = " "

	err := UserUpdateOne(paramLog, user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserChangeNewPIN(paramLog *basic.ParamLog, user *domain.User, newPIN string, session mongo.SessionContext) error {
	user.PIN = newPIN

	err := UserUpdateOne(paramLog, user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserChangeFaceAsPIN(paramLog *basic.ParamLog, user *domain.User, isFaceAsPIN bool, session mongo.SessionContext) error {
	user.FaceAsPIN = isFaceAsPIN

	err := UserUpdateOne(paramLog, user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserAddBankAccount(paramLog *basic.ParamLog, user *domain.User, bankAccount domain.Bank, session mongo.SessionContext) error {
	existinglistBank := user.SavedBankAccount
	existinglistBank = append(existinglistBank, bankAccount)

	user.SavedBankAccount = existinglistBank

	err := UserUpdateOne(paramLog, user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserRemoveBankAccount(paramLog *basic.ParamLog, user *domain.User, bankAccount domain.Bank, session mongo.SessionContext) error {
	existListBank := user.SavedBankAccount
	var newListBank = existListBank
	for index, bank := range existListBank {
		if bank.Name == bankAccount.Name && bank.BankCode == bankAccount.BankCode &&
			bank.AccountNumber == bankAccount.AccountNumber {
			newListBank = removeBankAccountByIndex(existListBank, index)
		}
	}

	user.SavedBankAccount = newListBank

	err := UserUpdateOne(paramLog, user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserVerify(paramLog *basic.ParamLog, user *domain.User, deviceID string, nik string, digitalID string, session mongo.SessionContext) error {
	user.Verified = true
	user.DeviceID = deviceID
	user.NIK = nik
	// user.ImageUpgrade = payload.Image
	user.DigitalID = digitalID

	err := UserUpdateOne(paramLog, user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserDTOByID(paramLog *basic.ParamLog, userID string) (dto.User, error) {

	objectID, err := primitive.ObjectIDFromHex(userID)

	query := []bson.M{
		{"$match": bson.M{"_id": objectID}},
		{"$unwind": "$list_balance"},
		{
			"$lookup": bson.M{
				"from":         "balance",
				"localField":   "list_balance.balance_id",
				"foreignField": "_id",
				"as":           "list_balance.detail",
			},
		},
		{"$unwind": "$list_balance.detail"},
		{
			"$group": bson.M{
				"_id": "$_id",
				"list_balance": bson.M{
					"$push": "$list_balance",
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "user",
				"localField":   "_id",
				"foreignField": "_id",
				"as":           "user_detail",
			},
		},
		{"$unwind": "$user_detail"},
		{
			"$addFields": bson.M{
				"user_detail.list_balance": "$list_balance",
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$user_detail",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "balance",
				"localField":   "main_balance",
				"foreignField": "_id",
				"as":           "main_balance",
			},
		},
		{"$unwind": "$main_balance"},
	}

	var result []dto.User
	cursor, err := database.Aggregate(paramLog, domain.USER_COLLECTION, query)
	cursor.All(context.TODO(), &result)
	if err != nil {
		return dto.User{}, utils.ErrorInternalServer(paramLog, utils.QueryFailed, "Query failed or cannot decode")
	}

	return result[0], nil
}

func UsersDTOFindByID(paramLog *basic.ParamLog, usersID []string) ([]dto.User, error) {

	var queryUsers []bson.M
	for _, element := range usersID {
		objectID, _ := primitive.ObjectIDFromHex(element)
		queryUsers = append(queryUsers, bson.M{"_id": bson.M{"$eq": objectID}})
	}

	query := []bson.M{
		{
			"$match": bson.M{
				"$or": queryUsers,
			},
		},
		{"$unwind": "$list_balance"},
		{
			"$lookup": bson.M{
				"from":         "balance",
				"localField":   "list_balance.balance_id",
				"foreignField": "_id",
				"as":           "list_balance.detail",
			},
		},
		{"$unwind": "$list_balance.detail"},
		{
			"$group": bson.M{
				"_id": "$_id",
				"list_balance": bson.M{
					"$push": "$list_balance",
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "user",
				"localField":   "_id",
				"foreignField": "_id",
				"as":           "user_detail",
			},
		},
		{"$unwind": "$user_detail"},
		{
			"$addFields": bson.M{
				"user_detail.list_balance": "$list_balance",
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$user_detail",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "balance",
				"localField":   "main_balance",
				"foreignField": "_id",
				"as":           "main_balance",
			},
		},
		{"$unwind": "$main_balance"},
	}

	var result []dto.User
	cursor, err := database.Aggregate(paramLog, domain.USER_COLLECTION, query)
	cursor.All(context.TODO(), &result)
	if err != nil {
		return []dto.User{}, utils.ErrorInternalServer(paramLog, utils.QueryFailed, "Query failed or cannot decode")
	}

	return result, nil
}

func UserDeleteInactive(user *domain.User, session mongo.SessionContext) error {
	err := database.SessionDeleteInactive(user, session)
	if err != nil {
		return err
	}

	return nil
}

func ValidateUserLocked(paramLog *basic.ParamLog, user domain.User) error {
	if user.Active == false {
		return utils.ErrorBadRequest(paramLog, utils.UserLocked, "User Locked")
	}

	return nil
}

func ValidateUserExist(paramLog *basic.ParamLog, user domain.User) error {
	if user.PhoneNumber == "" {
		return utils.ErrorUnauthorized(paramLog)
	}

	return nil
}

func ValidateUserLoginCode(paramLog *basic.ParamLog, user domain.User, loginCode string) error {

	if user.LoginCode != loginCode {
		return utils.ErrorBadRequest(paramLog, utils.InvalidLoginCode, "Invalid login code")
	}

	return nil
}

func ValidateUserFullname(paramLog *basic.ParamLog, fullName string) error {
	if utils.IsContainSpecialCharacter(fullName) {
		return utils.ErrorBadRequest(paramLog, utils.InvalidNameFormat, "Fullname error")
	}

	return nil
}

func ValidateUserActivationCode(paramLog *basic.ParamLog, user domain.User, activationCode string) error {

	if user.ActivationCode != activationCode {
		return utils.ErrorBadRequest(paramLog, utils.InvalidActivationCode, "Invalid activation code")
	}

	return nil
}

func ValidateIsUserAlreadyActive(paramLog *basic.ParamLog, user domain.User) error {
	if user.Active {
		return utils.ErrorBadRequest(paramLog, utils.UserAlreadyActive, "User already active")
	}

	return nil
}

func ValidateUserLoginAttempt(paramLog *basic.ParamLog, user domain.User) error {
	if user.LoginAttempt <= 0 {
		return utils.ErrorBadRequest(paramLog, utils.UserLocked, "User already locked")
	}

	return nil
}

func removeBankAccountByIndex(s []domain.Bank, index int) []domain.Bank {
	result := append(s[:index], s[index+1:]...)
	if len(result) == 0 {
		return []domain.Bank{}
	}

	return result
}

func ValidateUserChangePINCode(paramLog *basic.ParamLog, user domain.User, changePINCode string) error {

	if user.ChangePINCode != changePINCode {
		return utils.ErrorBadRequest(paramLog, utils.InvalidCode, "Invalid code")
	}

	return nil
}

func ValidateUserPIN(paramLog *basic.ParamLog, user domain.User, pin string) error {

	if pin != user.PIN {
		return utils.ErrorBadRequest(paramLog, utils.InvalidPIN, "Invalid Old PIN")
	}

	return nil
}

func UserbyIDNoSession(ID string) (domain.User, error) {
	model := domain.User{}
	cursor := database.FindOneByID(domain.USER_COLLECTION, ID)
	err := cursor.Decode(&model)
	if err != nil {
		return model, err
	}

	return model, nil
}
