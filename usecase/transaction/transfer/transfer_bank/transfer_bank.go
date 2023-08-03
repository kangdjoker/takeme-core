package transfer_bank

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/service"
	"github.com/kangdjoker/takeme-core/usecase"
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

type TransferBank struct {
}

func (self TransferBank) SetupGateway(transaction *domain.Transaction) {
	if true {
		transaction.GatewayStrategies = []domain.GatewayStrategy{
			{
				Code:       gateway.Permata,
				IsExecuted: false,
			},
		}
	} else if transaction.To.AccountNumber == "8691577392" && transaction.To.InstitutionCode == utils.BCA {
		transaction.GatewayStrategies = []domain.GatewayStrategy{
			{
				Code:       gateway.Xendit,
				IsExecuted: false,
			},
			{
				Code:       gateway.MMBC,
				IsExecuted: false,
			},
			{
				Code:       gateway.OY,
				IsExecuted: false,
			},
		}
	} else if transaction.To.InstitutionCode == utils.ALADIN {
		transaction.GatewayStrategies = []domain.GatewayStrategy{
			{
				Code:       gateway.Xendit,
				IsExecuted: false,
			},
		}
	} else if transaction.To.InstitutionCode == utils.OCBC || transaction.To.InstitutionCode == utils.DKI ||
		transaction.To.InstitutionCode == utils.JAWA_BARAT || transaction.To.InstitutionCode == utils.BTN {

		transaction.GatewayStrategies = []domain.GatewayStrategy{
			{
				Code:       gateway.MMBC,
				IsExecuted: false,
			},
			{
				Code:       gateway.OY,
				IsExecuted: false,
			},
			{
				Code:       gateway.Xendit,
				IsExecuted: false,
			},
		}
	} else {

		transaction.GatewayStrategies = []domain.GatewayStrategy{
			{
				Code:       gateway.OY,
				IsExecuted: false,
			},
			{
				Code:       gateway.MMBC,
				IsExecuted: false,
			},
			{
				Code:       gateway.Xendit,
				IsExecuted: false,
			},
		}
	}
}

func (self TransferBank) CreateTransferGateway(paramLog *basic.ParamLog, transaction domain.Transaction, requestID string) {
	oy := gateway.OYGateway{}
	mmbc := gateway.MMBCGateway{}
	xendit := gateway.XenditGateway{}
	permata := gateway.PermataGateway{}

	reference := ""
	var err error
	gatewayCode := changeGatewayStrategy(&transaction)
	basic.LogInformation(paramLog, "gatewayCode:"+gatewayCode)

	switch gatewayCode {
	case gateway.OY:
		reference, err = oy.CreateTransfer(paramLog, transaction)
	case gateway.Permata:
		reference, err = permata.CreateTransfer(paramLog, transaction, requestID)
	case gateway.MMBC:
		reference, err = mmbc.CreateTransfer(paramLog, transaction)
	case gateway.Xendit:
		reference, err = xendit.CreateTransfer(paramLog, transaction)
	}

	if gatewayCode == "" {
		transaction.Status = domain.FAILED_STATUS

		rollbackUsecase := RollbackTransferBank{}
		rollbackUsecase.Initialize(transaction)
		rollbackUsecase.ExecuteRollback(paramLog)
	}

	commitTransactionGateway(paramLog, transaction.ID.Hex(), transaction.Status, gatewayCode, reference, transaction.GatewayStrategies)

	if err != nil {
		self.CreateTransferGateway(paramLog, transaction, requestID)
		return
	}

	return
}

func (self TransferBank) ProcessCallbackGatewayTransfer(paramLog *basic.ParamLog, gatewayCode string, transactionCode string, reference string,
	status string, requestId string) (domain.Transaction, error) {

	var corporate domain.Corporate
	var transaction domain.Transaction
	var nextGateway string
	var err error

	if status == domain.REFUND_STATUS {
		transaction, err = service.TransactionByGatewayReferenceNoSession(paramLog, reference)
		corporate, err = service.CorporateByIDNoSession(transaction.CorporateID.Hex())
		nextGateway = checkUnexecutedGateway(transaction)
	} else {
		transaction, err = service.TransactionPendingByReferenceNoSession(paramLog, reference)
		corporate, err = service.CorporateByIDNoSession(transaction.CorporateID.Hex())
		nextGateway = checkUnexecutedGateway(transaction)
	}

	if err != nil {
		return domain.Transaction{}, err
	}

	if status == domain.PENDING_STATUS {
		return domain.Transaction{}, nil
	}

	if nextGateway != "" && (status == domain.FAILED_STATUS || status == domain.REFUND_STATUS) {
		go self.CreateTransferGateway(paramLog, transaction, requestId)
		return domain.Transaction{}, nil
	}

	if nextGateway == "" && (status == domain.FAILED_STATUS || status == domain.REFUND_STATUS) {
		transaction.Status = domain.FAILED_STATUS
		commitTransactionGateway(paramLog, transaction.ID.Hex(), transaction.Status, gatewayCode, reference, transaction.GatewayStrategies)

		rollbackUsecase := RollbackTransferBank{}
		rollbackUsecase.Initialize(transaction)
		rollbackUsecase.ExecuteRollback(paramLog)

		go usecase.PublishTransferCallback(paramLog, corporate, transaction)

		return domain.Transaction{}, nil
	}

	transaction.Status = domain.COMPLETED_STATUS
	commitTransactionGateway(paramLog, transaction.ID.Hex(), transaction.Status, gatewayCode, reference, transaction.GatewayStrategies)

	go usecase.PublishTransferCallback(paramLog, corporate, transaction)

	return transaction, nil
}

func changeGatewayStrategy(transaction *domain.Transaction) string {
	gatewayCode := ""

	for index, element := range transaction.GatewayStrategies {

		if element.IsExecuted == false {
			gatewayCode = element.Code
			transaction.GatewayStrategies[index].IsExecuted = true
			break
		}

	}

	return gatewayCode
}

func checkUnexecutedGateway(transaction domain.Transaction) string {
	gatewayCode := ""

	for _, element := range transaction.GatewayStrategies {

		if element.IsExecuted == false {
			gatewayCode = element.Code
			break
		}

	}

	return gatewayCode
}

func commitTransactionGateway(paramLog *basic.ParamLog, transactionID string, status string, gatewayCode string, reference string, gatewayStrategy []domain.GatewayStrategy) {
	function := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			session.AbortTransaction(session)
			return utils.ErrorInternalServer(paramLog, utils.DBStartTransactionFailed, "Initialize commit gateway start transaction failed")
		}

		transaction, err := service.TransactionByID(paramLog, transactionID, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		if transaction.GatewayReference != reference && transaction.Gateway != gatewayCode {
			history := domain.GatewayHistory{
				Code:      gatewayCode,
				Reference: reference,
				Time:      time.Now().Format(os.Getenv("TIME_FORMAT")),
			}
			transaction.GatewayHistories = append(transaction.GatewayHistories, history)
		}

		transaction.Status = status
		transaction.Gateway = gatewayCode
		transaction.GatewayReference = reference
		transaction.GatewayStrategies = gatewayStrategy

		err = service.TransactionUpdateOne(paramLog, &transaction, session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}

		return database.CommitWithRetry(session)
	}

	err := database.DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return database.RunTransactionWithRetry(sctx, function)
		},
	)

	if err != nil {
		basic.LogError(paramLog, fmt.Sprintf("Failed commit gateway transaction with id  %v because %v ", transactionID, err.Error()))
	}
}
