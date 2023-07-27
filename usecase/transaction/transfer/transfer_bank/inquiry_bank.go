package transfer_bank

import (
	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils/gateway"
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
