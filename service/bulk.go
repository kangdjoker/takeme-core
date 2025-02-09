package service

import (
	"os"
	"time"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/basic"
	"github.com/kangdjoker/takeme-core/utils/database"
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

func CreateBulkTransfer(paramLog *basic.ParamLog, corporate domain.Corporate, totalBulk int, reference string,
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
			return domain.BulkTransfer{}, utils.ErrorBadRequest(paramLog, utils.ExternalIDNotFound, "External ID Not Found")
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

func SaveBulkInquiry(paramLog *basic.ParamLog, bulk *domain.BulkInquiry) error {
	err := database.SaveOne(paramLog, domain.BULK_INQUIRY_COLLECTION, bulk)
	if err != nil {
		return err
	}

	return nil
}

func SaveBulkTransfer(paramLog *basic.ParamLog, bulk *domain.BulkTransfer) error {
	err := database.SaveOne(paramLog, domain.BULK_TRANSFER_COLLECTION, bulk)
	if err != nil {
		return err
	}

	return nil
}

func BulkInquiryUpdateOne(paramLog *basic.ParamLog, model *domain.BulkInquiry) error {
	err := database.UpdateOne(paramLog, domain.BULK_INQUIRY_COLLECTION, model)
	if err != nil {
		return err
	}

	return nil
}

func BulkTransferUpdateOne(paramLog *basic.ParamLog, model *domain.BulkTransfer) error {
	err := database.UpdateOne(paramLog, domain.BULK_TRANSFER_COLLECTION, model)
	if err != nil {
		return err
	}

	return nil
}
