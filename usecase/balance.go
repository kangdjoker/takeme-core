package usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/service"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/basic"
	"github.com/kangdjoker/takeme-core/utils/database"
	"github.com/kangdjoker/takeme-core/utils/gateway"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func CreateBalanceUser(paramLog *basic.ParamLog, user domain.ActorAble, corporate domain.Corporate, balanceName string) (domain.Balance, error) {
	var balance domain.Balance

	if utils.IsContainSpecialCharacter(balanceName) {
		return domain.Balance{}, utils.ErrorBadRequest(paramLog, utils.InvalidNameFormat, "Name balance error")
	}

	createBalanceForUser := func(session mongo.SessionContext) error {

		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(paramLog, utils.DBStartTransactionFailed, "Initialize balance start transaction failed")
		}

		balance, err = service.BalanceCreate(corporate.ID,
			user.ToActorObject(),
			balanceName,
			corporate.Currency,
			session,
		)

		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		user, err := service.UserByID(paramLog, user.GetActorID().Hex(), session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		user.ListBalance = append(user.ListBalance, domain.AccessBalance{
			BalanceID: balance.ID,
			Access:    domain.ACCESS_BALANCE_OWNER,
		})

		err = service.UserUpdateOne(paramLog, &user, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		createVABalance(paramLog, &balance, user.FullName)
		err = service.BalanceUpdate(paramLog, balance, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		return database.CommitWithRetry(session)
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, createBalanceForUser)
		},
	)

	if err != nil {
		return balance, utils.ErrorInternalServer(paramLog, utils.DBStartTransactionFailed,
			fmt.Sprintf("Failed initialize balance for user (%v)", user.GetActorID()))
	}

	return balance, nil
}

func CreateBalanceCorporate(paramLog *basic.ParamLog, corp domain.ActorAble, corporate domain.Corporate, balanceName string) (domain.Balance, error) {
	var balance domain.Balance

	if utils.IsContainSpecialCharacter(balanceName) {
		return domain.Balance{}, utils.ErrorBadRequest(paramLog, utils.InvalidNameFormat, "Name balance error")
	}

	createBalanceForCorporate := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(paramLog, utils.DBStartTransactionFailed, "Initialize balance start transaction failed")
		}

		balance, err = service.BalanceCreate(
			corporate.ID,
			corp.ToActorObject(),
			balanceName,
			corporate.Currency,
			session,
		)

		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		corp, err := service.CorporateByID(corporate.ID.Hex(), session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		corp.ListBalance = append(corp.ListBalance, domain.AccessBalance{
			BalanceID: balance.ID,
			Access:    domain.ACCESS_BALANCE_OWNER,
		})

		err = service.CorporateUpdateOne(paramLog, &corp, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		createVABalance(paramLog, &balance, corporate.Name)
		err = service.BalanceUpdate(paramLog, balance, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		return database.CommitWithRetry(session)
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, createBalanceForCorporate)
		},
	)

	if err != nil {
		return balance, utils.ErrorInternalServer(paramLog, utils.DBStartTransactionFailed,
			fmt.Sprintf("Failed initialize balance for corporate (%v)", corporate.GetActorID()))
	}

	return balance, nil
}

func InitializeBalanceUser(paramLog *basic.ParamLog, user domain.ActorAble, corporate domain.Corporate, balanceName string) (domain.Balance, error) {
	var balance domain.Balance
	createBalanceForUser := func(session mongo.SessionContext) error {

		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(paramLog, utils.DBStartTransactionFailed, "Initialize balance start transaction failed")
		}

		balance, err = service.BalanceInitialization(user.GetActorID(),
			corporate.ID,
			user.ToActorObject(),
			balanceName,
			corporate.Currency,
			session,
		)

		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		user, err := service.UserByID(paramLog, user.GetActorID().Hex(), session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		user.MainBalance = balance.ID

		user.ListBalance = append(user.ListBalance, domain.AccessBalance{
			BalanceID: balance.ID,
			Access:    domain.ACCESS_BALANCE_OWNER,
		})

		err = service.UserUpdateOne(paramLog, &user, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		createVABalance(paramLog, &balance, user.FullName)
		err = service.BalanceUpdate(paramLog, balance, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		return database.CommitWithRetry(session)
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, createBalanceForUser)
		},
	)

	if err != nil {
		return balance, utils.ErrorInternalServer(paramLog, utils.DBStartTransactionFailed,
			fmt.Sprintf("Failed initialize balance for user (%v)", user.GetActorID()))
	}

	return balance, nil
}

func InitializeBalanceCorporate(paramLog *basic.ParamLog, corp domain.ActorAble, corporate domain.Corporate, balanceName string) (domain.Balance, error) {
	var balance domain.Balance
	createBalanceForCorporate := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return utils.ErrorInternalServer(paramLog, utils.DBStartTransactionFailed, "Initialize balance start transaction failed")
		}

		balance, err = service.BalanceInitialization(corp.GetActorID(),
			corporate.ID,
			corp.ToActorObject(),
			balanceName,
			corporate.Currency,
			session,
		)

		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		corp, err := service.CorporateByID(corporate.ID.Hex(), session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		corp.MainBalance = balance.ID

		corp.ListBalance = append(corp.ListBalance, domain.AccessBalance{
			BalanceID: balance.ID,
			Access:    domain.ACCESS_BALANCE_OWNER,
		})

		err = service.CorporateUpdateOne(paramLog, &corp, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		createVABalance(paramLog, &balance, corporate.Name)
		err = service.BalanceUpdate(paramLog, balance, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		return database.CommitWithRetry(session)
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, createBalanceForCorporate)
		},
	)

	if err != nil {
		return balance, utils.ErrorInternalServer(paramLog, utils.DBStartTransactionFailed,
			fmt.Sprintf("Failed initialize balance for corporate (%v)", corporate.GetActorID()))
	}

	return balance, nil
}

func WithdrawBalance(paramLog *basic.ParamLog, statement domain.Statement, session mongo.SessionContext) error {
	balanceID := statement.BalanceID.Hex()
	amount := statement.Withdraw

	balance, err := service.BalanceByID(balanceID, session)
	if err != nil {
		return err
	}
	basic.LogInformation(paramLog, "balance.Amount: "+strconv.Itoa(balance.Amount)+", "+strconv.Itoa(amount))
	if balance.Amount < amount {
		return utils.ErrorBadRequest(paramLog, utils.InsufficientBalance, "Insufficient balance")
	}

	balance.Amount = balance.Amount - amount

	err = service.BalanceUpdate(paramLog, balance, session)
	if err != nil {
		return err
	}

	statement.Balance = balance.Amount

	err = service.StatementSaveOne(statement, session)
	if err != nil {
		return err
	}

	return nil
}

func DepositBalance(paramLog *basic.ParamLog, statement domain.Statement, session mongo.SessionContext) error {
	balanceID := statement.BalanceID.Hex()
	amount := statement.Deposit

	balance, err := service.BalanceByID(balanceID, session)
	if err != nil {
		return err
	}
	basic.LogInformation2(paramLog, "BalanceByID", balance)

	balance.Amount = balance.Amount + amount

	err = service.BalanceUpdate(paramLog, balance, session)
	if err != nil {
		return err
	}

	basic.LogInformation2(paramLog, "BalanceUpdate", balance)
	statement.Balance = balance.Amount

	err = service.StatementSaveOne(statement, session)
	if err != nil {
		return err
	}

	basic.LogInformation(paramLog, "StatementSaveOne")
	return nil
}

func StatementByBalanceID(paramLog *basic.ParamLog, balanceID string, page string, limit string) ([]domain.Statement, error) {
	balance, err := service.BalanceByIDNoSession(balanceID)
	if err != nil || balance.Owner.Type == "" {
		return []domain.Statement{}, utils.ErrorBadRequest(paramLog, utils.InvalidBalanceID, "Balance not found")

	}

	statements, err := service.StatementsByBalanceID(paramLog, balance.ID, page, limit)
	if err != nil {
		return []domain.Statement{}, err
	}

	return statements, nil
}

func ShareBalance(paramLog *basic.ParamLog, corporate domain.Corporate, balanceID string, access string, actorID string, pin string) error {

	err := ValidateActorPIN(paramLog, corporate, pin)
	if err != nil {
		return err
	}

	balance, err := service.BalanceByIDNoSession(balanceID)
	if err != nil || balance.Owner.Type == "" {
		return utils.ErrorBadRequest(paramLog, utils.InvalidBalanceID, "Balance not found")
	}

	actor, err := ActorByID(paramLog, actorID)
	if err != nil {
		return err
	}

	if access != domain.ACCESS_BALANCE_SHARED && access != domain.ACCESS_BALANCE_VIEW_ONLY {
		return utils.ErrorBadRequest(paramLog, utils.InvalidAccessType, "Invalid request access type")
	}

	if balance.CorporateID != corporate.ID {
		return utils.ErrorBadRequest(paramLog, utils.InvalidAccessType, "Invalid balance scope")
	}

	if IsAccessBalanceAlreadyHave(actor, balanceID) {
		return utils.ErrorBadRequest(paramLog, utils.AccessBalanceAlreadyHave, "Access balance already have")
	}

	err = ActorAddBalance(paramLog, actor, domain.AccessBalance{
		BalanceID: balance.ID,
		Access:    access,
	})

	if err != nil {
		return err
	}

	return nil
}

func RevokeBalance(paramLog *basic.ParamLog, corporate domain.Corporate, balanceID string, revokeFrom string, pin string) error {

	err := ValidateActorPIN(paramLog, corporate, pin)
	if err != nil {
		return err
	}

	balance, err := service.BalanceByIDNoSession(balanceID)
	if err != nil || balance.Owner.Type == "" {
		return utils.ErrorBadRequest(paramLog, utils.InvalidBalanceID, "Balance not found")
	}

	actor, err := ActorByID(paramLog, revokeFrom)
	if err != nil {
		return err
	}

	if balance.CorporateID != corporate.ID {
		return utils.ErrorBadRequest(paramLog, utils.InvalidAccessType, "Invalid balance scope")
	}

	err = ActorRemoveBalance(paramLog, actor, balanceID)
	if err != nil {
		return err
	}

	return nil
}

func CreateRequestAccesssBalance(paramLog *basic.ParamLog, corporate domain.Corporate, requester domain.ActorAble,
	balanceID string, access string) (domain.RequestAccessBalance, error) {

	balance, err := service.BalanceByIDNoSession(balanceID)
	if err != nil || balance.Owner.Type == "" {
		return domain.RequestAccessBalance{}, utils.ErrorBadRequest(paramLog, utils.InvalidBalanceID, "Balance not found")
	}

	if access != domain.ACCESS_BALANCE_SHARED && access != domain.ACCESS_BALANCE_VIEW_ONLY {
		return domain.RequestAccessBalance{}, utils.ErrorBadRequest(paramLog, utils.InvalidAccessType, "Invalid request access type")
	}

	if IsAccessBalanceAlreadyHave(requester, balanceID) {
		return domain.RequestAccessBalance{}, utils.ErrorBadRequest(paramLog, utils.AccessBalanceAlreadyHave, "Access balance already have")
	}

	request, err := service.CreateRAB(paramLog, corporate, balance, requester.ToActorObject(), balance.Owner, access)
	if err != nil {
		return domain.RequestAccessBalance{}, err
	}

	return request, nil
}

func ListRequesterAccesssBalance(paramLog *basic.ParamLog, actor domain.ActorAble, status string) ([]domain.RequestAccessBalance, error) {
	result, err := service.RABByRequsterID(paramLog, actor.GetActorID().Hex(), status)
	if err != nil {
		return []domain.RequestAccessBalance{}, err
	}

	return result, nil
}

func ListOwnerAccesssBalance(paramLog *basic.ParamLog, actor domain.ActorAble, status string) ([]domain.RequestAccessBalance, error) {
	result, err := service.RABByOwnerID(paramLog, actor.GetActorID().Hex(), status)
	if err != nil {
		return []domain.RequestAccessBalance{}, err
	}

	return result, nil
}

func ProccedRAB(paramLog *basic.ParamLog, requestID string, status string, owner domain.ActorAble, pin string) (domain.RequestAccessBalance, error) {

	request, err := service.RABByID(paramLog, requestID)
	if err != nil {
		return domain.RequestAccessBalance{}, utils.ErrorBadRequest(paramLog, utils.RequestAccessBalanceNotFound, "Balance not found")
	}

	requester, err := ActorObjectToActor(paramLog, request.BalanceRequester)
	if err != nil {
		return domain.RequestAccessBalance{}, err
	}

	balanceID := request.BalanceID.Hex()

	if request.Status != domain.REQUEST_ACCESS_BALANCE_STATUS_PENDING {
		return domain.RequestAccessBalance{}, utils.ErrorBadRequest(paramLog, utils.RequestAccessBalanceAlreadyProcced, "Request already procced")
	}

	if IsAccessBalanceAlreadyHave(requester, balanceID) {
		return domain.RequestAccessBalance{}, utils.ErrorBadRequest(paramLog, utils.AccessBalanceAlreadyHave, "Access balance already have")
	}

	if !IsBalanceOwner(owner, balanceID) {
		return domain.RequestAccessBalance{}, utils.ErrorBadRequest(paramLog, utils.InvalidBalanceOwner, "Invalid balance access")
	}

	err = ValidateActorPIN(paramLog, owner, pin)
	if err != nil {
		return domain.RequestAccessBalance{}, err
	}

	request.Status = status

	if status == domain.REQUEST_ACCESS_BALANCE_STATUS_APPROVE {

		err = ActorAddBalance(paramLog, requester, domain.AccessBalance{
			BalanceID: request.BalanceID,
			Access:    request.Access,
		})

		if err != nil {
			return domain.RequestAccessBalance{}, err
		}

	}

	err = service.RABUpdateOne(paramLog, &request)
	if err != nil {
		return domain.RequestAccessBalance{}, err
	}

	return request, nil
}

func createVABalance(paramLog *basic.ParamLog, balance *domain.Balance, ownerName string) {

	balanceName := ownerName + " " + balance.Name
	gatewayXendit := gateway.XenditGateway{}

	mandiriAccountNumber, err := gatewayXendit.CreateVA(paramLog, balance.ID.Hex(), balanceName, "MANDIRI")
	if err != nil {
		balance.VA = append(balance.VA, domain.VirtualAccount{
			BankCode:      "MANDIRI",
			AccountNumber: "Call administrator for fix this",
		})
	} else {
		balance.VA = append(balance.VA, domain.VirtualAccount{
			BankCode:      "MANDIRI",
			AccountNumber: mandiriAccountNumber,
		})
	}

	bniAccountNumber, err := gatewayXendit.CreateVA(paramLog, balance.ID.Hex(), balanceName, "BNI")
	if err != nil {
		balance.VA = append(balance.VA, domain.VirtualAccount{
			BankCode:      "BNI",
			AccountNumber: "Call administrator for fix this",
		})
	} else {
		balance.VA = append(balance.VA, domain.VirtualAccount{
			BankCode:      "BNI",
			AccountNumber: bniAccountNumber,
		})
	}

	briAccountNumber, err := gatewayXendit.CreateVA(paramLog, balance.ID.Hex(), balanceName, "BRI")
	if err != nil {
		balance.VA = append(balance.VA, domain.VirtualAccount{
			BankCode:      "BRI",
			AccountNumber: "Call administrator for fix this",
		})
	} else {
		balance.VA = append(balance.VA, domain.VirtualAccount{
			BankCode:      "BRI",
			AccountNumber: briAccountNumber,
		})
	}

	permataAccountNumber, err := gatewayXendit.CreateVA(paramLog, balance.ID.Hex(), balanceName, "PERMATA")
	if err != nil {
		balance.VA = append(balance.VA, domain.VirtualAccount{
			BankCode:      "PERMATA",
			AccountNumber: "Call administrator for fix this",
		})
	} else {
		balance.VA = append(balance.VA, domain.VirtualAccount{
			BankCode:      "PERMATA",
			AccountNumber: permataAccountNumber,
		})
	}
}
