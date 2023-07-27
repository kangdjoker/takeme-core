package service

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/domain/dto"
	"github.com/kangdjoker/takeme-core/utils"
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

func UserCreateUnpending(corporate domain.Corporate, userPending domain.User, email string, phoneNumber string, fullName string,
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

	err := UserUpdateOne(&userPending, session)
	if err != nil {
		return domain.User{}, err
	}

	return userPending, nil
}

func UserActivate(user *domain.User, session mongo.SessionContext) error {

	user.Active = true
	user.ActivationCode = ""
	user.Pending = false

	err := UserUpdateOne(user, session)
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

func UserUpdateOne(model *domain.User, session mongo.SessionContext) error {
	err := database.SessionUpdateOne(model, session)
	if err != nil {
		return err
	}

	return nil
}

func UserUpdateOneNoSession(model *domain.User) error {
	err := database.UpdateOne(domain.USER_COLLECTION, model)
	if err != nil {
		return err
	}

	return nil
}

func UserByID(ID string, session mongo.SessionContext) (domain.User, error) {
	model := domain.User{}
	cursor := database.SessionFindOneByID(domain.USER_COLLECTION, ID, session)
	err := cursor.Decode(&model)
	if err != nil {
		return domain.User{}, utils.ErrorInternalServer(utils.QueryFailed, "Query failed or cannot decode")
	}

	return model, nil
}

func UserByIDNoSession(ID string) (domain.User, error) {
	model := domain.User{}
	cursor := database.FindOneByID(domain.USER_COLLECTION, ID)
	err := cursor.Decode(&model)
	if err != nil {
		return domain.User{}, utils.ErrorInternalServer(utils.QueryFailed, "Query failed or cannot decode")
	}

	return model, nil
}

func UserByPhoneNumber(corporateID primitive.ObjectID, phoneNumber string, session mongo.SessionContext) (domain.User, error) {
	model := domain.User{}
	query := bson.M{"phone_number": phoneNumber, "corporate_id": corporateID, "pending": false}
	cursor := database.SessionFindOne(domain.USER_COLLECTION, query, session)
	cursor.Decode(&model)

	if model.PhoneNumber == "" {
		return domain.User{}, utils.ErrorBadRequest(utils.UserNotFound, "User not found")
	}

	return model, nil
}

func UserByPhoneNumberWithoutSession(corporateID primitive.ObjectID, phoneNumber string) (domain.User, error) {
	model := domain.User{}
	query := bson.M{"phone_number": phoneNumber, "corporate_id": corporateID, "pending": false}
	cursor := database.FindOne(domain.USER_COLLECTION, query)
	cursor.Decode(&model)

	if model.PhoneNumber == "" {
		return domain.User{}, utils.ErrorBadRequest(utils.UserNotFound, "User not found")
	}

	return model, nil
}

func ValidateUserNotRegisteredYet(corporate domain.Corporate, phoneNumber string, email string,
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
		return isPending, model, utils.ErrorBadRequest(utils.UserAlreadyExist, "User already exist")
	}

	return isPending, model, nil
}

func UserByIDWithValidation(ID string, validations []func(user domain.User) error) (domain.User, error) {
	model := domain.User{}
	cursor := database.FindOneByID(domain.USER_COLLECTION, ID)
	err := cursor.Decode(&model)
	if err != nil {
		return domain.User{}, utils.ErrorInternalServer(utils.QueryFailed, "Query failed or cannot decode")
	}

	for _, element := range validations {
		err := element(model)
		if err != nil {
			return domain.User{}, err
		}
	}

	return model, nil
}

func UserReduceAccessAttempt(userID string, session mongo.SessionContext) error {

	user, err := UserByID(userID, session)
	if err != nil {
		return err
	}

	err = ValidateUserLocked(user)
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

	err = UserUpdateOne(&user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserGenerateLoginCode(user *domain.User, session mongo.SessionContext) error {
	user.LoginAttempt -= 1
	user.LoginCode = utils.GenerateShortCode()

	err := UserUpdateOne(user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserRefreshAttempt(user *domain.User, session mongo.SessionContext) error {
	attempt, _ := strconv.Atoi(os.Getenv("SECURITY_ATTEMPT"))
	user.LoginAttempt = int8(attempt)
	user.AccessAttempt = int8(attempt)
	user.LoginCode = "-"

	err := UserUpdateOne(user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserLock(user *domain.User, session mongo.SessionContext) error {

	user.Active = false
	user.AccessAttempt = 0
	user.LoginAttempt = 0

	err := UserUpdateOne(user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserUnlock(user *domain.User, session mongo.SessionContext) error {
	attempt, _ := strconv.Atoi(os.Getenv("SECURITY_ATTEMPT"))
	user.Active = true
	user.AccessAttempt = int8(attempt)
	user.LoginAttempt = int8(attempt)

	err := UserUpdateOne(user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserSavePIN(user *domain.User, pin string, session mongo.SessionContext) error {

	pin, err := utils.RSADecrypt(pin)
	if err != nil {
		return utils.ErrorInternalServer(utils.DecryptError, err.Error())
	}

	user.PIN = pin

	err = UserUpdateOne(user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserGenerateChangePINCode(user *domain.User, pin string, session mongo.SessionContext) (string, error) {

	pin, err := utils.RSADecrypt(pin)
	if err != nil {
		return "", utils.ErrorInternalServer(utils.DecryptError, err.Error())
	}

	user.ChangePIN = pin
	user.ChangePINCode = utils.GenerateShortCode()

	err = UserUpdateOne(user, session)
	if err != nil {
		return "", err
	}

	return user.ChangePINCode, nil
}

func UserChangePIN(user *domain.User, session mongo.SessionContext) error {
	user.PIN = user.ChangePIN
	user.ChangePIN = " "
	user.ChangePINCode = " "

	err := UserUpdateOne(user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserChangeNewPIN(user *domain.User, newPIN string, session mongo.SessionContext) error {
	user.PIN = newPIN

	err := UserUpdateOne(user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserChangeFaceAsPIN(user *domain.User, isFaceAsPIN bool, session mongo.SessionContext) error {
	user.FaceAsPIN = isFaceAsPIN

	err := UserUpdateOne(user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserAddBankAccount(user *domain.User, bankAccount domain.Bank, session mongo.SessionContext) error {
	existinglistBank := user.SavedBankAccount
	existinglistBank = append(existinglistBank, bankAccount)

	user.SavedBankAccount = existinglistBank

	err := UserUpdateOne(user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserRemoveBankAccount(user *domain.User, bankAccount domain.Bank, session mongo.SessionContext) error {
	existListBank := user.SavedBankAccount
	var newListBank = existListBank
	for index, bank := range existListBank {
		if bank.Name == bankAccount.Name && bank.BankCode == bankAccount.BankCode &&
			bank.AccountNumber == bankAccount.AccountNumber {
			newListBank = removeBankAccountByIndex(existListBank, index)
		}
	}

	user.SavedBankAccount = newListBank

	err := UserUpdateOne(user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserVerify(user *domain.User, deviceID string, nik string, digitalID string, session mongo.SessionContext) error {
	user.Verified = true
	user.DeviceID = deviceID
	user.NIK = nik
	// user.ImageUpgrade = payload.Image
	user.DigitalID = digitalID

	err := UserUpdateOne(user, session)
	if err != nil {
		return err
	}

	return nil
}

func UserDTOByID(userID string) (dto.User, error) {

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
	cursor, err := database.Aggregate(domain.USER_COLLECTION, query)
	cursor.All(context.TODO(), &result)
	if err != nil {
		return dto.User{}, utils.ErrorInternalServer(utils.QueryFailed, "Query failed or cannot decode")
	}

	return result[0], nil
}

func UsersDTOFindByID(usersID []string) ([]dto.User, error) {

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
	cursor, err := database.Aggregate(domain.USER_COLLECTION, query)
	cursor.All(context.TODO(), &result)
	if err != nil {
		return []dto.User{}, utils.ErrorInternalServer(utils.QueryFailed, "Query failed or cannot decode")
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

func ValidateUserLocked(user domain.User) error {
	if user.Active == false {
		return utils.ErrorBadRequest(utils.UserLocked, "User Locked")
	}

	return nil
}

func ValidateUserExist(user domain.User) error {
	if user.PhoneNumber == "" {
		return utils.ErrorUnauthorized()
	}

	return nil
}

func ValidateUserLoginCode(user domain.User, loginCode string) error {

	if user.LoginCode != loginCode {
		return utils.ErrorBadRequest(utils.InvalidLoginCode, "Invalid login code")
	}

	return nil
}

func ValidateUserFullname(fullName string) error {
	if utils.IsContainSpecialCharacter(fullName) {
		return utils.ErrorBadRequest(utils.InvalidNameFormat, "Fullname error")
	}

	return nil
}

func ValidateUserActivationCode(user domain.User, activationCode string) error {

	if user.ActivationCode != activationCode {
		return utils.ErrorBadRequest(utils.InvalidActivationCode, "Invalid activation code")
	}

	return nil
}

func ValidateIsUserAlreadyActive(user domain.User) error {
	if user.Active {
		return utils.ErrorBadRequest(utils.UserAlreadyActive, "User already active")
	}

	return nil
}

func ValidateUserLoginAttempt(user domain.User) error {
	if user.LoginAttempt <= 0 {
		return utils.ErrorBadRequest(utils.UserLocked, "User already locked")
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

func ValidateUserChangePINCode(user domain.User, changePINCode string) error {

	if user.ChangePINCode != changePINCode {
		return utils.ErrorBadRequest(utils.InvalidCode, "Invalid code")
	}

	return nil
}

func ValidateUserPIN(user domain.User, pin string) error {

	if pin != user.PIN {
		return utils.ErrorBadRequest(utils.InvalidPIN, "Invalid Old PIN")
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
