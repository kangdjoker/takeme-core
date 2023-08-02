package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/domain/dto"
	"github.com/kangdjoker/takeme-core/service"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/basic"
	"github.com/kangdjoker/takeme-core/utils/database"
	"github.com/kangdjoker/takeme-core/utils/storage"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func UserCheck(user domain.User) (dto.User, error) {

	result, err := service.UserDTOByID(user.ID.Hex())
	if err != nil {
		return dto.User{}, err
	}

	return result, nil
}

func UserUpgrade(user domain.User, nik string, faceImage string, deviceID string) error {
	body, err := utils.EKYCEnrollUser(nik, faceImage)
	if err != nil {
		return err
	}

	userUpgrade := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(utils.DBStartTransactionFailed, "User login start transaction failed")
		}

		user, err = service.UserByID(user.ID.Hex(), session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = service.UserVerify(&user, body.DeviceID, body.NIK, body.DigitalID, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = database.CommitWithRetry(session)
		if err != nil {
			return err
		}

		return nil

	}

	err = database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, userUpgrade)
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func UserSaveBankAccount(user domain.User, name string, bankCode string, accountNumber string) error {
	account := domain.Bank{
		Name:          name,
		BankCode:      bankCode,
		AccountNumber: accountNumber,
	}

	userSaveBank := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(utils.DBStartTransactionFailed, "User login start transaction failed")
		}

		user, err = service.UserByID(user.ID.Hex(), session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = service.UserAddBankAccount(&user, account, session)
		if err != nil {
			return err
		}

		err = database.CommitWithRetry(session)
		if err != nil {
			return err
		}

		return nil
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, userSaveBank)
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func UserDeleteBankAccount(user domain.User, name string, bankCode string, accountNumber string) error {
	account := domain.Bank{
		Name:          name,
		BankCode:      bankCode,
		AccountNumber: accountNumber,
	}

	userDeleteBank := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(utils.DBStartTransactionFailed, "User login start transaction failed")
		}

		user, err = service.UserByID(user.ID.Hex(), session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = service.UserRemoveBankAccount(&user, account, session)
		if err != nil {
			return err
		}

		err = database.CommitWithRetry(session)
		if err != nil {
			return err
		}

		return nil
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, userDeleteBank)
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func CheckUserPhonebook(corporate domain.Corporate, phonebook []domain.Contact) []domain.Contact {
	var members []domain.Contact
	for _, contact := range phonebook {
		isExist, name := isPhoneNumberAlreadyExist(corporate, contact.Number)

		if isExist {
			members = append(members, domain.Contact{
				Number: contact.Number,
				Name:   name,
			})
		}
	}

	return members
}

func isPhoneNumberAlreadyExist(corporate domain.Corporate, phoneNumber string) (bool, string) {

	user, err := service.UserByPhoneNumberWithoutSession(corporate.ID, phoneNumber)
	if err != nil {
		return false, ""
	}

	return true, user.FullName
}

func UserSavePIN(user domain.User, encryptedPIN string) error {
	userSavePIN := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(utils.DBStartTransactionFailed, "User login start transaction failed")
		}

		user, err = service.UserByID(user.ID.Hex(), session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = service.UserSavePIN(&user, encryptedPIN, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = database.CommitWithRetry(session)
		if err != nil {
			return err
		}

		return nil
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, userSavePIN)
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func UserPreForgotPIN(paramLog basic.ParamLog, user domain.User, encryptedPIN string, OTPChannel string) error {

	preForgotPIN := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(utils.DBStartTransactionFailed, "User login start transaction failed")
		}

		user, err = service.UserByID(user.ID.Hex(), session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		changePINCode, err := service.UserGenerateChangePINCode(&user, encryptedPIN, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = database.CommitWithRetry(session)
		if err != nil {
			return err
		}

		// sending SMS code
		phoneNumber := user.PhoneNumber
		userID := user.ID.Hex()
		message := fmt.Sprintf("Your forgot number %v",
			changePINCode)

		if OTPChannel == SMS_CHANNEL {
			go utils.SendSMS(paramLog, phoneNumber, message)
		} else {
			go utils.SendWAHubungi(paramLog, phoneNumber, changePINCode)
		}

		go removeForgotPINCode(userID, changePINCode)

		return nil
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, preForgotPIN)
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func UserForgotPIN(user domain.User, code string) error {

	forgotPIN := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(utils.DBStartTransactionFailed, "User login start transaction failed")
		}

		user, err = service.UserByID(user.ID.Hex(), session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = service.ValidateUserChangePINCode(user, code)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = service.UserChangePIN(&user, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = database.CommitWithRetry(session)
		if err != nil {
			return err
		}

		return nil
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, forgotPIN)
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func UserMainBalanceVA(user domain.User) ([]domain.VirtualAccount, error) {
	var va []domain.VirtualAccount

	userMainBalanceVA := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(utils.DBStartTransactionFailed, "User login start transaction failed")
		}

		balance, err := service.BalanceByID(user.MainBalance.Hex(), session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = database.CommitWithRetry(session)
		if err != nil {
			return err
		}

		va = balance.VA

		return nil
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, userMainBalanceVA)
		},
	)

	if err != nil {
		return []domain.VirtualAccount{}, err
	}

	return va, nil
}

func UserChangePIN(user domain.User, encryptedOldPIN string, encryptedNewPIN string) error {

	userChangePIN := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(utils.DBStartTransactionFailed, "User login start transaction failed")
		}

		newPIN, err := utils.RSADecrypt(encryptedNewPIN)
		if err != nil {
			return err
		}

		err = ValidateFormatPIN(newPIN)
		if err != nil {
			return err
		}

		user, err = service.UserByID(user.ID.Hex(), session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = ValidateActorPIN(user, encryptedOldPIN)
		if err != nil {
			return err
		}

		err = service.UserChangeNewPIN(&user, newPIN, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = database.CommitWithRetry(session)
		if err != nil {
			return err
		}

		return nil
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, userChangePIN)
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func UserChangeFaceAsPIN(user domain.User, isFaceAsPIN bool) error {
	userChangeFaceAsPIN := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(utils.DBStartTransactionFailed, "User login start transaction failed")
		}

		user, err = service.UserByID(user.ID.Hex(), session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = service.UserChangeFaceAsPIN(&user, isFaceAsPIN, session)
		if err != nil {
			return err
		}

		err = database.CommitWithRetry(session)
		if err != nil {
			return err
		}

		return nil
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, userChangeFaceAsPIN)
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func UserTransactions(user domain.User, page string, limit string) ([]domain.Transaction, error) {
	var ownBalances []primitive.ObjectID

	for _, index := range user.ListBalance {
		if index.Access == domain.ACCESS_BALANCE_OWNER {
			ownBalances = append(ownBalances, index.BalanceID)
		}
	}

	transactions, err := service.TransactionsByActorNoSession(user.PhoneNumber, ownBalances, page, limit)
	if err != nil {
		return []domain.Transaction{}, err
	}

	return transactions, nil
}

func UserTemporaryPIN(faceImage string, user domain.User) (string, error) {
	_, err := utils.EKYCVerifyUser(user.NIK, faceImage, user.DigitalID)
	if err != nil {
		return "", err
	}

	temporaryPIN := generateTemporaryPIN(user)
	go removeTemporaryPIN(user)

	return temporaryPIN, nil
}

func UserVerify(aktaImage multipart.File, aktaHeader *multipart.FileHeader,
	npwpImage multipart.File, npwpHeader *multipart.FileHeader, nibImage multipart.File,
	nibHeader *multipart.FileHeader, identityImage multipart.File, identityHeader *multipart.FileHeader,
	nik string, legalName string, legalAddress string, userID string, verifyType string) error {

	user, err := service.UserByIDNoSession(userID)
	if err != nil {
		return err
	}

	err, akta := storage.SaveFile(aktaImage, *aktaHeader)
	if err != nil {
		return err
	}

	err, npwp := storage.SaveFile(npwpImage, *npwpHeader)
	if err != nil {
		return err
	}

	err, nib := storage.SaveFile(nibImage, *nibHeader)
	if err != nil {
		return err
	}

	err, identity := storage.SaveFile(identityImage, *identityHeader)
	if err != nil {
		return err
	}

	user.VerifyData.AktaImage = akta
	user.VerifyData.NPWPImage = npwp
	user.VerifyData.NIBImage = nib
	user.VerifyData.IdentityImage = identity
	user.VerifyData.NIK = nik
	user.VerifyData.LegalName = legalName
	user.VerifyData.LegalAddress = legalAddress
	user.VerifyData.Type = verifyType
	user.Verified = true

	err = service.UserUpdateOneNoSession(&user)
	if err != nil {
		return err
	}

	return nil
}

func removeForgotPINCode(userID string, code string) {
	time.Sleep(120 * time.Second)
	userRemoveForgot := func(session mongo.SessionContext) error {

		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return err
		}

		user, err := service.UserByID(userID, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		if user.ChangePINCode == code {
			log.Info(fmt.Sprintf("Remove login code for user (%v)", user.PhoneNumber))
			user.ChangePINCode = "-"
			err := service.UserUpdateOne(&user, session)
			if err != nil {
				session.AbortTransaction(session)
				return err
			}
		}

		return database.CommitWithRetry(session)
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, userRemoveForgot)
		},
	)

	if err != nil {
		log.Error(fmt.Sprintf("Failed remove preforgot code for userID (%v)", userID))
	}
}

func generateTemporaryPIN(user domain.User) string {
	temporaryPIN := utils.GenerateShortCode()
	user.TemporaryPIN = temporaryPIN
	database.UpdateOne(domain.USER_COLLECTION, &user)

	return temporaryPIN
}

func removeTemporaryPIN(user domain.User) {
	time.Sleep(600 * time.Second)
	user.TemporaryPIN = " "
	database.UpdateOne(domain.USER_COLLECTION, &user)
}
