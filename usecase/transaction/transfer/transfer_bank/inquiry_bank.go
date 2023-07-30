package transfer_bank

import (
	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils/gateway"
)

func InquiryBankAccount(accountNumber string, bankCode string, requestId string) (domain.Bank, error) {
	// gateway := gateway.OYGateway{}
	gateway := gateway.PermataGateway{}

	accountName, err := gateway.Inquiry(bankCode, accountNumber, requestId)
	if err != nil {
		return domain.Bank{}, err
	}

	bank := domain.CreateBank(bankCode, accountName, accountNumber)

	return bank, nil
}
