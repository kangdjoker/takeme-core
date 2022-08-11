package transfer_bank

import (
	"context"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/takeme-id/core/domain"
	"github.com/takeme-id/core/service"
	"github.com/takeme-id/core/usecase"
	"github.com/takeme-id/core/utils"
	"github.com/takeme-id/core/utils/database"
	"github.com/takeme-id/core/utils/gateway"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type TransferBank struct {
}

func (self TransferBank) SetupGateway(transaction *domain.Transaction) {

	if transaction.To.AccountNumber == "0000000000" {
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
	} else if transaction.To.InstitutionCode == utils.OCBC || transaction.To.InstitutionCode == utils.DKI || transaction.To.InstitutionCode == utils.JAWA_BARAT {

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

func (self TransferBank) CreateTransferGateway(transaction domain.Transaction) {
	oy := gateway.OYGateway{}
	mmbc := gateway.MMBCGateway{}
	xendit := gateway.XenditGateway{}

	reference := ""
	var err error
	gatewayCode := changeGatewayStrategy(&transaction)

	switch gatewayCode {
	case gateway.OY:
		reference, err = oy.CreateTransfer(transaction)
	case gateway.MMBC:
		reference, err = mmbc.CreateTransfer(transaction)
	case gateway.Xendit:
		reference, err = xendit.CreateTransfer(transaction)
	}

	if gatewayCode == "" {
		transaction.Status = domain.FAILED_STATUS

		rollbackUsecase := RollbackTransferBank{}
		rollbackUsecase.Initialize(transaction)
		rollbackUsecase.ExecuteRollback()
	}

	commitTransactionGateway(transaction.ID.Hex(), transaction.Status, gatewayCode, reference, transaction.GatewayStrategies)

	if err != nil {
		self.CreateTransferGateway(transaction)
		return
	}

	return
}

func (self TransferBank) ProcessCallbackGatewayTransfer(gatewayCode string, transactionCode string, reference string,
	status string) (domain.Transaction, error) {

	var corporate domain.Corporate
	var transaction domain.Transaction
	var nextGateway string
	var err error

	if status == domain.REFUND_STATUS {
		transaction, err = service.TransactionByGatewayReferenceNoSession(reference)
		corporate, err = service.CorporateByIDNoSession(transaction.CorporateID.Hex())
		nextGateway = checkUnexecutedGateway(transaction)
	} else {
		transaction, err = service.TransactionPendingByReferenceNoSession(reference)
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
		go self.CreateTransferGateway(transaction)
		return domain.Transaction{}, nil
	}

	if nextGateway == "" && (status == domain.FAILED_STATUS || status == domain.REFUND_STATUS) {
		transaction.Status = domain.FAILED_STATUS
		commitTransactionGateway(transaction.ID.Hex(), transaction.Status, gatewayCode, reference, transaction.GatewayStrategies)

		rollbackUsecase := RollbackTransferBank{}
		rollbackUsecase.Initialize(transaction)
		rollbackUsecase.ExecuteRollback()

		go usecase.PublishTransferCallback(corporate, transaction)

		return domain.Transaction{}, nil
	}

	transaction.Status = domain.COMPLETED_STATUS
	commitTransactionGateway(transaction.ID.Hex(), transaction.Status, gatewayCode, reference, transaction.GatewayStrategies)

	go usecase.PublishTransferCallback(corporate, transaction)

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

func commitTransactionGateway(transactionID string, status string, gatewayCode string, reference string, gatewayStrategy []domain.GatewayStrategy) {
	function := func(session mongo.SessionContext) error {
		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			session.AbortTransaction(session)
			return utils.ErrorInternalServer(utils.DBStartTransactionFailed, "Initialize commit gateway start transaction failed")
		}

		transaction, err := service.TransactionByID(transactionID, session)
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

		err = service.TransactionUpdateOne(&transaction, session)
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
		log.Error(fmt.Sprintf("Failed commit gateway transaction with id  %v because %v ", transactionID, err.Error()))
	}
}
