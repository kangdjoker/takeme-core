package transfer_bank

import (
	"github.com/takeme-id/core/domain"
	"github.com/takeme-id/core/utils/gateway"
)

func InquiryBankAccount(accountNumber string, bankCode string) (domain.Bank, error) {
	gateway := gateway.OYGateway{}

	accountName, err := gateway.Inquiry(bankCode, accountNumber)
	if err != nil {
		return domain.Bank{}, err
	}

	bank := domain.CreateBank(bankCode, accountName, accountNumber)

	return bank, nil
}
