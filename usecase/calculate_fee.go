package usecase

import (
	"os"
	"strconv"
	"time"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/service"
	"github.com/kangdjoker/takeme-core/utils"
)

// source = user
// corporate_level = corporate
// to bank
// ==============================================
// source transactionStatement - & feeStatement -
// corporate feeStatement + & feeStatement -
// principal feeStatement +

// source = user
// corporate_level = principal
// to bank
// ==============================================
// source transactionalStatement - & feeStatement -
// corporate feeStatement +

// ============================================================== //

// source = corporate
// corporate_level = corporate
// to bank
// ==============================================
// source transactionStatement - & feeStatement -
// principal feeStatement +

// source = corporate
// corporate_level = principal
// to bank
// ==============================================
// corporate transactionalStatement -

type CalculateFee struct {
	result          []domain.Statement
	corporate       domain.Corporate
	balance         domain.Balance
	balanceType     string
	transaction     domain.Transaction
	transactionType string
}

func (self *CalculateFee) Initialize(corporate domain.Corporate, balance domain.Balance, transaction domain.Transaction) {
	self.corporate = corporate
	self.balance = balance
	self.balanceType = balance.Owner.Type
	self.transaction = transaction
	self.transactionType = transaction.Type
}

func (self *CalculateFee) CalculateByOwnerAndTransaction() ([]domain.Statement, error) {
	var err error
	if self.balanceType == domain.ACTOR_TYPE_USER {
		err = self.balanceUser(self.corporate, self.balance, self.transaction)
	} else {
		err = self.balanceCorporate(self.corporate, self.balance, self.transaction)
	}

	return self.result, err
}

func (self *CalculateFee) RollbackFeeStatement(statements []domain.Statement) []domain.Statement {
	var result []domain.Statement
	for _, element := range statements {

		if element.Withdraw != 0 {
			s := service.DepositFeeStatement(
				element.BalanceID,
				time.Now().Format(os.Getenv("TIME_FORMAT")),
				element.Reference,
				element.Withdraw,
			)
			result = append(result, s)
		} else {
			s := service.WithdrawFeeStatement(
				element.BalanceID,
				time.Now().Format(os.Getenv("TIME_FORMAT")),
				element.Reference,
				element.Deposit,
			)
			result = append(result, s)
		}
	}

	return result
}

func (self *CalculateFee) balanceUser(corporate domain.Corporate, userBalance domain.Balance, transaction domain.Transaction) error {
	if self.transactionType == domain.TRANSFER_BANK {
		a, err := balanceUserTransferBank(corporate, userBalance, transaction)
		if err != nil {
			return err
		}

		self.result = a
	} else if self.transactionType == domain.TOPUP {
		a, err := balanceUserTopupBank(corporate, userBalance, transaction)
		if err != nil {
			return err
		}

		self.result = a
	} else if self.transactionType == domain.TRANSFER_WALLET {
		a, err := balanceUserTransferBalance(corporate, userBalance, transaction)
		if err != nil {
			return err
		}

		self.result = a
	} else if self.transactionType == domain.ACCEPT_PAYMENT_CARD {
		a, err := balanceUserAcceptPaymentCard(corporate, userBalance, transaction)
		if err != nil {
			return err
		}

		self.result = a
	}

	return nil
}

func (self *CalculateFee) balanceCorporate(corporate domain.Corporate, corporateBalance domain.Balance, transaction domain.Transaction) error {
	if self.transactionType == domain.TRANSFER_BANK {
		a, err := balanceCorporateTransferBank(corporate, corporateBalance, transaction)
		if err != nil {
			return err
		}

		self.result = a
	} else if self.transactionType == domain.TOPUP {
		a, err := balanceCorporateTopupBank(corporate, corporateBalance, transaction)
		if err != nil {
			return err
		}

		self.result = a
	} else if self.transactionType == domain.TRANSFER_WALLET {
		a, err := balanceCorporateTransferBalance(corporate, corporateBalance, transaction)
		if err != nil {
			return err
		}

		self.result = a
	} else if self.transactionType == domain.DEDUCT {
		a, err := balanceCorporateDeductBalance(corporate, corporateBalance, transaction)
		if err != nil {
			return err
		}

		self.result = a
	} else if self.transactionType == domain.ACCEPT_PAYMENT_CARD {
		a, err := balanceCorporateAcceptPaymentCard(corporate, corporateBalance, transaction)
		if err != nil {
			return err
		}

		self.result = a
	}

	return nil
}

func balanceUserTransferBank(corporate domain.Corporate, userBalance domain.Balance, transaction domain.Transaction) ([]domain.Statement, error) {
	userFee := corporate.FeeUser.TransferBank
	userBalanceID := userBalance.ID
	corporateBalanceID := corporate.MainBalance

	var result []domain.Statement

	withdrawUser := service.WithdrawFeeStatement(userBalanceID, transaction.Time, transaction.TransactionCode, userFee)
	depositCorporate := service.DepositFeeStatement(corporateBalanceID, transaction.Time, transaction.TransactionCode, userFee)

	result = append(result, withdrawUser)
	result = append(result, depositCorporate)

	if IsNotPrincipal(corporate) {
		principal, err := service.CorporateByIDNoSession(corporate.Parent.Hex())
		if err != nil {
			return result, err
		}

		corporateFee := corporate.FeeCorporate.TransferBank
		principalBalanceID := principal.MainBalance

		withdrawCorporate := service.WithdrawFeeStatement(corporateBalanceID, transaction.Time, transaction.TransactionCode, corporateFee)
		depositPrincipal := service.DepositFeeStatement(principalBalanceID, transaction.Time, transaction.TransactionCode, corporateFee)

		result = append(result, withdrawCorporate)
		result = append(result, depositPrincipal)
	}

	return result, nil
}

func balanceCorporateTransferBank(corporate domain.Corporate, corporateBalance domain.Balance, transaction domain.Transaction) ([]domain.Statement, error) {
	var result []domain.Statement
	if IsNotPrincipal(corporate) {
		principal, err := service.CorporateByIDNoSession(corporate.Parent.Hex())
		if err != nil {
			return result, err
		}

		corporateFee := corporate.FeeCorporate.TransferBank
		corporateBalanceID := corporate.MainBalance
		principalBalanceID := principal.MainBalance

		withdrawCorporate := service.WithdrawFeeStatement(corporateBalanceID, transaction.Time, transaction.TransactionCode, corporateFee)
		depositPrincipal := service.DepositFeeStatement(principalBalanceID, transaction.Time, transaction.TransactionCode, corporateFee)
		result = append(result, withdrawCorporate)
		result = append(result, depositPrincipal)
	}

	return result, nil
}

func balanceUserTopupBank(corporate domain.Corporate, userBalance domain.Balance, transaction domain.Transaction) ([]domain.Statement, error) {
	userFee := corporate.FeeUser.Topup
	userBalanceID := userBalance.ID
	corporateBalanceID := corporate.MainBalance

	var result []domain.Statement

	withdrawUser := service.WithdrawFeeStatement(userBalanceID, transaction.Time, transaction.TransactionCode, userFee)
	depositCorporate := service.DepositFeeStatement(corporateBalanceID, transaction.Time, transaction.TransactionCode, userFee)

	result = append(result, withdrawUser)
	result = append(result, depositCorporate)

	if IsNotPrincipal(corporate) {
		principal, err := service.CorporateByIDNoSession(corporate.Parent.Hex())
		if err != nil {
			return result, err
		}

		corporateFee := corporate.FeeCorporate.Topup
		principalBalanceID := principal.MainBalance

		withdrawCorporate := service.WithdrawFeeStatement(corporateBalanceID, transaction.Time, transaction.TransactionCode, corporateFee)
		depositPrincipal := service.DepositFeeStatement(principalBalanceID, transaction.Time, transaction.TransactionCode, corporateFee)

		result = append(result, withdrawCorporate)
		result = append(result, depositPrincipal)
	}

	return result, nil
}

func balanceCorporateTopupBank(corporate domain.Corporate, corporateBalance domain.Balance, transaction domain.Transaction) ([]domain.Statement, error) {
	var result []domain.Statement
	if IsNotPrincipal(corporate) {
		principal, err := service.CorporateByIDNoSession(corporate.Parent.Hex())
		if err != nil {
			return result, err
		}

		corporateFee := corporate.FeeCorporate.Topup
		corporateBalanceID := corporate.MainBalance
		principalBalanceID := principal.MainBalance

		withdrawCorporate := service.WithdrawFeeStatement(corporateBalanceID, transaction.Time, transaction.TransactionCode, corporateFee)
		depositPrincipal := service.DepositFeeStatement(principalBalanceID, transaction.Time, transaction.TransactionCode, corporateFee)
		result = append(result, withdrawCorporate)
		result = append(result, depositPrincipal)
	}

	return result, nil
}

func balanceUserTransferBalance(corporate domain.Corporate, userBalance domain.Balance, transaction domain.Transaction) ([]domain.Statement, error) {
	userFee := corporate.FeeUser.TransferBalance
	userBalanceID := userBalance.ID
	corporateBalanceID := corporate.MainBalance

	var result []domain.Statement

	withdrawUser := service.WithdrawFeeStatement(userBalanceID, transaction.Time, transaction.TransactionCode, userFee)
	depositCorporate := service.DepositFeeStatement(corporateBalanceID, transaction.Time, transaction.TransactionCode, userFee)

	result = append(result, withdrawUser)
	result = append(result, depositCorporate)

	if IsNotPrincipal(corporate) {
		principal, err := service.CorporateByIDNoSession(corporate.Parent.Hex())
		if err != nil {
			return result, err
		}

		corporateFee := corporate.FeeCorporate.TransferBalance
		principalBalanceID := principal.MainBalance

		withdrawCorporate := service.WithdrawFeeStatement(corporateBalanceID, transaction.Time, transaction.TransactionCode, corporateFee)
		depositPrincipal := service.DepositFeeStatement(principalBalanceID, transaction.Time, transaction.TransactionCode, corporateFee)

		result = append(result, withdrawCorporate)
		result = append(result, depositPrincipal)
	}

	return result, nil
}

func balanceCorporateTransferBalance(corporate domain.Corporate, corporateBalance domain.Balance, transaction domain.Transaction) ([]domain.Statement, error) {
	var result []domain.Statement
	if IsNotPrincipal(corporate) {
		principal, err := service.CorporateByIDNoSession(corporate.Parent.Hex())
		if err != nil {
			return result, err
		}

		corporateFee := corporate.FeeCorporate.TransferBalance
		corporateBalanceID := corporate.MainBalance
		principalBalanceID := principal.MainBalance

		withdrawCorporate := service.WithdrawFeeStatement(corporateBalanceID, transaction.Time, transaction.TransactionCode, corporateFee)
		depositPrincipal := service.DepositFeeStatement(principalBalanceID, transaction.Time, transaction.TransactionCode, corporateFee)
		result = append(result, withdrawCorporate)
		result = append(result, depositPrincipal)
	}

	return result, nil
}

func balanceCorporateDeductBalance(corporate domain.Corporate, corporateBalance domain.Balance, transaction domain.Transaction) ([]domain.Statement, error) {
	var result []domain.Statement
	if IsNotPrincipal(corporate) {
		principal, err := service.CorporateByIDNoSession(corporate.Parent.Hex())
		if err != nil {
			return result, err
		}

		corporateFee := corporate.FeeCorporate.Deduct
		corporateBalanceID := corporate.MainBalance
		principalBalanceID := principal.MainBalance

		withdrawCorporate := service.WithdrawFeeStatement(corporateBalanceID, transaction.Time, transaction.TransactionCode, corporateFee)
		depositPrincipal := service.DepositFeeStatement(principalBalanceID, transaction.Time, transaction.TransactionCode, corporateFee)
		result = append(result, withdrawCorporate)
		result = append(result, depositPrincipal)
	}

	return result, nil
}

// TODO PROVIDE LOGIC FOR PRINCIPAL CAN ACCEPT MONEY FROM MULTICURRENCY TRANSACTION
func balanceCorporateAcceptPaymentCard(corporate domain.Corporate, corporateBalance domain.Balance, transaction domain.Transaction) ([]domain.Statement, error) {
	var result []domain.Statement
	if IsNotPrincipal(corporate) && IsNotIDRCurrency(transaction.Currency) {
		principal, err := service.CorporateByIDNoSession(corporate.Parent.Hex())
		if err != nil {
			return result, err
		}

		percentageFee, err := strconv.ParseFloat(corporate.FeeCorporate.AcceptPaymentCard, 64)
		if err != nil {
			return result, utils.ErrorBadRequest(utils.WrongAcceptCardFee, "Cannot convert accept payment card fee")
		}

		corporateFee := int(float64(transaction.SubAmount) * percentageFee)

		corporateBalanceID := corporate.MainBalance
		principalBalanceID := principal.MainBalance

		withdrawCorporate := service.WithdrawFeeStatement(corporateBalanceID, transaction.Time, transaction.TransactionCode, corporateFee)
		depositPrincipal := service.DepositFeeStatement(principalBalanceID, transaction.Time, transaction.TransactionCode, corporateFee)
		result = append(result, withdrawCorporate)
		result = append(result, depositPrincipal)
	}

	return result, nil
}

// TODO PROVIDE LOGIC FOR PRINCIPAL CAN ACCEPT MONEY FROM MULTICURRENCY TRANSACTION
func balanceUserAcceptPaymentCard(corporate domain.Corporate, userBalance domain.Balance, transaction domain.Transaction) ([]domain.Statement, error) {
	var result []domain.Statement
	percentageFee, err := strconv.ParseFloat(corporate.FeeUser.AcceptPaymentCard, 64)
	if err != nil {
		return result, utils.ErrorBadRequest(utils.WrongAcceptCardFee, "Cannot convert accept payment card fee")
	}

	userFee := int(float64(transaction.SubAmount) * percentageFee)

	userBalanceID := userBalance.ID
	corporateBalanceID := corporate.MainBalance

	withdrawUser := service.WithdrawFeeStatement(userBalanceID, transaction.Time, transaction.TransactionCode, userFee)
	depositCorporate := service.DepositFeeStatement(corporateBalanceID, transaction.Time, transaction.TransactionCode, userFee)

	result = append(result, withdrawUser)
	result = append(result, depositCorporate)

	if IsNotPrincipal(corporate) && IsNotIDRCurrency(transaction.Currency) {

		percentageFee, err := strconv.ParseFloat(corporate.FeeCorporate.AcceptPaymentCard, 64)
		if err != nil {
			return result, utils.ErrorBadRequest(utils.WrongAcceptCardFee, "Cannot convert accept payment card fee")
		}

		principal, err := service.CorporateByIDNoSession(corporate.Parent.Hex())
		if err != nil {
			return result, err
		}

		corporateFee := int(float64(transaction.SubAmount) * percentageFee)
		principalBalanceID := principal.MainBalance

		withdrawCorporate := service.WithdrawFeeStatement(corporateBalanceID, transaction.Time, transaction.TransactionCode, corporateFee)
		depositPrincipal := service.DepositFeeStatement(principalBalanceID, transaction.Time, transaction.TransactionCode, corporateFee)

		result = append(result, withdrawCorporate)
		result = append(result, depositPrincipal)
	}

	return result, nil
}

func IsNotPrincipal(corporate domain.Corporate) bool {
	if corporate.Parent.Hex() == "000000000000000000000000" {
		return false
	}

	return true
}

func IsNotIDRCurrency(currency string) bool {
	if currency != "idr" {
		return true
	}

	return false
}
