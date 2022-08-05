package service

import (
	"os"
	"time"

	"github.com/takeme-id/core/domain"
	"github.com/takeme-id/core/utils"
	"github.com/takeme-id/core/utils/database"
)

func CreateBulkInquiry(corporate domain.Corporate, totalBulk int, reference string, banks []domain.Bank,
	actor domain.ActorObject) domain.BulkInquiry {

	bulk := domain.BulkInquiry{
		CorporateID: corporate.ID,
		Reference:   reference,
		Owner:       actor,
		Time:        time.Now().Format(os.Getenv("TIME_FORMAT")),
		Status:      domain.BULK_PROGRESS_STATUS,
		TotalList:   totalBulk,
	}

	number := 1
	for _, a := range banks {
		bulk.List = append(bulk.List, domain.Inquiry{
			Number:        number,
			AccountName:   "",
			AccountNumber: a.AccountNumber,
			BankName:      a.BankCode,
			Valid:         false,
		})
		number++
	}

	return bulk
}

func CreateBulkTransfer(corporate domain.Corporate, totalBulk int, reference string,
	transfers []domain.Transfer, actor domain.ActorObject, balance domain.Balance) (domain.BulkTransfer, error) {

	bulk := domain.BulkTransfer{
		CorporateID: corporate.ID,
		BalanceID:   balance.ID,
		Reference:   reference,
		Owner:       actor,
		Time:        time.Now().Format(os.Getenv("TIME_FORMAT")),
		Status:      domain.BULK_UNEXECUTED_STATUS,
		TotalList:   totalBulk,
	}

	number := 1
	subAmount := 0
	for _, transfer := range transfers {
		if transfer.ExternalID == "" || transfer.ExternalID == " " {
			return domain.BulkTransfer{}, utils.ErrorBadRequest(utils.ExternalIDNotFound, "External ID Not Found")
		}

		transfer.Number = number
		bulk.List = append(bulk.List, transfer)

		number++
		subAmount += transfer.Amount
	}

	bulk.SubAmount = subAmount

	if actor.Type == domain.ACTOR_TYPE_USER {
		bulk.Amount = subAmount + (corporate.FeeUser.TransferBank * totalBulk)
	} else {
		bulk.Amount = subAmount + (corporate.FeeCorporate.TransferBank * totalBulk)
	}

	return bulk, nil
}

func BulkInquiryByID(ID string) (domain.BulkInquiry, error) {
	model := domain.BulkInquiry{}
	cursor := database.FindOneByID(domain.BULK_INQUIRY_COLLECTION, ID)
	err := cursor.Decode(&model)
	if err != nil {
		return domain.BulkInquiry{}, err
	}

	return model, nil
}

func BulkTransferByID(ID string) (domain.BulkTransfer, error) {
	model := domain.BulkTransfer{}
	cursor := database.FindOneByID(domain.BULK_TRANSFER_COLLECTION, ID)
	err := cursor.Decode(&model)
	if err != nil {
		return domain.BulkTransfer{}, err
	}

	return model, nil
}

func SaveBulkInquiry(bulk *domain.BulkInquiry) error {
	err := database.SaveOne(domain.BULK_INQUIRY_COLLECTION, bulk)
	if err != nil {
		return err
	}

	return nil
}

func SaveBulkTransfer(bulk *domain.BulkTransfer) error {
	err := database.SaveOne(domain.BULK_TRANSFER_COLLECTION, bulk)
	if err != nil {
		return err
	}

	return nil
}

func BulkInquiryUpdateOne(model *domain.BulkInquiry) error {
	err := database.UpdateOne(domain.BULK_INQUIRY_COLLECTION, model)
	if err != nil {
		return err
	}

	return nil
}

func BulkTransferUpdateOne(model *domain.BulkTransfer) error {
	err := database.UpdateOne(domain.BULK_TRANSFER_COLLECTION, model)
	if err != nil {
		return err
	}

	return nil
}
