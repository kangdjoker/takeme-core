package usecase

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/takeme-id/core/domain"
	"github.com/takeme-id/core/service"
	"github.com/takeme-id/core/usecase/security"
	"github.com/takeme-id/core/utils"
	"github.com/takeme-id/core/utils/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

const (
	SMS_CHANNEL = "sms"
	WA_CHANNEL  = "wa"
)

func UserSignup(fullName string, email string, phoneNumber string, corporate domain.Corporate, OTPChannel string) error {

	userSignup := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(utils.DBStartTransactionFailed, "User prelogin start transaction failed")
		}

		// validate user fullname
		err = service.ValidateUserFullname(fullName)
		if err != nil {
			return err
		}

		// validate is user already exist
		isPending, user, err := service.ValidateUserNotRegisteredYet(corporate, phoneNumber, email, session)
		if err != nil {
			return err
		}

		if isPending == false {
			user, err = service.UserCreate(corporate, email, phoneNumber, fullName, session)
			if err != nil {
				return err
			}
		} else {
			user, err = service.UserCreateUnpending(corporate, user, email, phoneNumber, fullName, session)
			if err != nil {
				return err
			}
		}

		err = database.CommitWithRetry(session)
		if err != nil {
			return err
		}

		// sending SMS verification
		message := fmt.Sprintf(
			"Your signup number  %v",
			user.ActivationCode,
		)

		if OTPChannel == SMS_CHANNEL {
			go utils.SendSMS(phoneNumber, message)
		} else {
			go utils.SendWAHubungi(phoneNumber, user.ActivationCode)
		}

		go deleteInactiveUser(user.ID.Hex())

		return nil
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, userSignup)
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func UserActivation(phoneNumber string, corporate domain.Corporate, code string) (string, error) {
	token := ""

	userActivation := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(utils.DBStartTransactionFailed, "User login start transaction failed")
		}

		user, err := service.UserByPhoneNumber(corporate.ID, phoneNumber, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = service.ValidateIsUserAlreadyActive(user)
		if err != nil {
			return err
		}

		err = service.ValidateUserActivationCode(user, code)
		if err != nil {
			return err
		}

		err = service.UserActivate(&user, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		// Generate JWT
		tokenString, err := utils.JWTEncode(user, corporate)
		if err != nil {
			return err
		}

		err = database.CommitWithRetry(session)
		if err != nil {
			return err
		}

		go InitializeBalanceUser(user, corporate, "Main")

		token = tokenString

		return nil
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, userActivation)
		},
	)

	if err != nil {
		return "", err
	}

	return token, nil
}

func UserPrelogin(phoneNumber string, corporate domain.Corporate, OTPChannel string) error {

	userPrelogin := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(utils.DBStartTransactionFailed, "User prelogin start transaction failed")
		}

		user, err := service.UserByPhoneNumber(corporate.ID, phoneNumber, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = service.ValidateUserLocked(user)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = service.ValidateUserLoginAttempt(user)
		if err != nil {
			go security.LockUser(user.ID.Hex())

			session.AbortTransaction(session)
			return err
		}

		err = service.UserGenerateLoginCode(&user, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = database.CommitWithRetry(session)
		if err != nil {
			return err
		}

		// sending SMS login code
		message := fmt.Sprintf(
			"Your login number %v",
			user.LoginCode,
		)

		if OTPChannel == SMS_CHANNEL {
			go utils.SendSMS(phoneNumber, message)
		} else {
			go utils.SendWAHubungi(phoneNumber, user.LoginCode)
		}

		go userRemoveLoginCode(user.ID.Hex(), user.LoginCode)

		return nil
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, userPrelogin)
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func UserLogin(phoneNumber string, corporate domain.Corporate, code string) (string, error) {

	token := ""

	userLogin := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(utils.DBStartTransactionFailed, "User login start transaction failed")
		}

		user, err := service.UserByPhoneNumber(corporate.ID, phoneNumber, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = service.ValidateUserLocked(user)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = service.ValidateUserLoginCode(user, code)
		if err != nil {
			// Reduce user access attempt
			go security.InvalidUserAuth(user)
			session.AbortTransaction(session)
			return err
		}

		// Generate JWT
		tokenString, err := utils.JWTEncode(user, corporate)
		if err != nil {
			return err
		}

		err = service.UserRefreshAttempt(&user, session)
		if err != nil {
			return err
		}

		err = database.CommitWithRetry(session)
		if err != nil {
			return err
		}

		token = tokenString

		return nil
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, userLogin)
		},
	)

	if err != nil {
		return "", err
	}

	return token, nil
}

func UserFaceLogin(phoneNumber string, corporate domain.Corporate, faceImage string) (string, error) {

	token := ""

	userLogin := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(utils.DBStartTransactionFailed, "User login start transaction failed")
		}

		user, err := service.UserByPhoneNumber(corporate.ID, phoneNumber, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		err = service.ValidateUserLocked(user)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		_, err = utils.EKYCVerifyUser(user.NIK, faceImage, user.DigitalID)
		if err != nil {
			go security.InvalidUserAuth(user)
			session.AbortTransaction(session)
			return err
		}

		// Generate JWT
		tokenString, err := utils.JWTEncode(user, corporate)
		if err != nil {
			return err
		}

		err = service.UserRefreshAttempt(&user, session)
		if err != nil {
			return err
		}

		err = database.CommitWithRetry(session)
		if err != nil {
			return err
		}

		token = tokenString

		return nil
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, userLogin)
		},
	)

	if err != nil {
		return "", err
	}

	return token, nil
}

func userRemoveLoginCode(userID string, loginCode string) {
	time.Sleep(120 * time.Second)
	userRemoveLogin := func(session mongo.SessionContext) error {

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

		if user.LoginCode == loginCode {
			log.Info(fmt.Sprintf("Remove login code for user (%v)", user.PhoneNumber))
			user.LoginCode = "-"
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
			return database.RunTransactionWithRetry(sctx, userRemoveLogin)
		},
	)

	if err != nil {
		log.Error(fmt.Sprintf("Failed remove login code for userID (%v)", userID))
	}
}

func deleteInactiveUser(userID string) {
	time.Sleep(120 * time.Second)
	userRemoveLogin := func(session mongo.SessionContext) error {

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

		err = service.UserDeleteInactive(&user, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		if err != nil {
			log.Info("Delete unactive user failed because user already active")
		} else {
			log.Info("Delete unactive user success")
		}

		return database.CommitWithRetry(session)
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, userRemoveLogin)
		},
	)

	if err != nil {
		log.Error(fmt.Sprintf("Failed remove login code for userID (%v)", userID))
	}
}
