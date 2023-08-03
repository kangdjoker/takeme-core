package database

import (
	"os"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/basic"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CommitWithRetry(sctx mongo.SessionContext) error {
	for {
		err := sctx.CommitTransaction(sctx)
		switch e := err.(type) {
		case nil:
			return nil
		case mongo.CommandError:
			if e.HasErrorLabel("UnknownTransactionCommitResult") {
				continue
			}
			return e
		default:
			return e
		}
	}
}

func RunTransactionWithRetry(sessionCtx mongo.SessionContext, function func(mongo.SessionContext) error) error {
	for {
		err := function(sessionCtx)
		if err == nil {
			return nil
		}

		if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.HasErrorLabel("TransientTransactionError") {
			continue
		}

		return err
	}
}

func SessionFindOneByID(colName string, ID string, session mongo.SessionContext) *mongo.SingleResult {
	objectID, _ := primitive.ObjectIDFromHex(ID)

	query := bson.M{"_id": objectID}
	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(colName)
	result := collection.FindOne(session, query)

	return result
}

func SessionFindOne(colName string, query bson.M, session mongo.SessionContext) *mongo.SingleResult {
	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(colName)
	result := collection.FindOne(session, query)

	return result
}

func SessionUpdateOne(paramLog *basic.ParamLog, domain domain.BaseModel, session mongo.SessionContext) error {
	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(domain.CollectionName())
	document, err := toDoc(domain)
	if err != nil {
		return utils.ErrorInternalServer(paramLog, utils.UpdateFailed, err.Error())
	}

	filter := bson.M{"_id": bson.M{"$eq": domain.GetDocumentID()}}
	update := bson.M{"$set": document}

	_, err = collection.UpdateOne(
		session,
		filter,
		update,
	)
	if err != nil {
		return err
	}

	return nil
}

func SessionSaveOne(domain domain.BaseModel, session mongo.SessionContext) error {

	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(domain.CollectionName())
	result, err := collection.InsertOne(session, domain)
	if err != nil {
		return err
	}

	ID := result.InsertedID.(primitive.ObjectID)
	domain.SetDocumentID(ID)

	return nil
}

func SessionDeleteInactive(domain domain.BaseModel, session mongo.SessionContext) error {
	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(domain.CollectionName())

	_, err := collection.DeleteOne(session, bson.M{"_id": domain.GetDocumentID(), "active": false, "pending": false})
	if err != nil {
		return err
	}

	return nil
}
