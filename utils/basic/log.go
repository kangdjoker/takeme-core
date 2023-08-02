package basic

import (
	"context"
	"errors"
	"io"
	"os"
	"time"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

const LOG_COLLECTION string = "log"

var DBClient *mongo.Client

const (
	ACCESS_LOG_VIEW_ONLY = "View Only"
	ACCESS_LOG_SHARED    = "Shared"
	ACCESS_LOG_OWNER     = "Owner"
)

func init() {
	SetupDB()
}

type Log struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TimeCreate string             `json:"time_create" bson:"time_create,omitempty"`
	Tag        string             `json:"tag" bson:"tag,omitempty"`
	IsError    bool               `json:"is_error" bson:"is_error"`
	Data       interface{}        `json:"data" bson:"data"`
}

func (domain *Log) SetDocumentID(ID primitive.ObjectID) {
	domain.ID = ID
}

func (domain *Log) GetDocumentID() primitive.ObjectID {
	return domain.ID
}

func (domain *Log) CollectionName() string {
	return LOG_COLLECTION
}

func LogCreate(isError bool, paramLog ParamLog, data interface{}, session mongo.SessionContext) (Log, error) {
	now := time.Now().Format("2006-01-02 15:04:05")

	model := Log{
		Data:       data,
		TimeCreate: now,
		IsError:    isError,
	}

	err := LogSaveOne(&model, session)
	if err != nil {
		return model, err
	}

	return model, nil
}

func LogSaveOne(model *Log, session mongo.SessionContext) error {
	err := SessionSaveOne(model, session)
	if err != nil {
		return err
	}

	return nil
}

func LogSaveOneNoSession(model *Log) error {
	err := SaveOne(LOG_COLLECTION, model)
	if err != nil {
		return err
	}

	return nil
}

func LogById(ID string, session mongo.SessionContext) (Log, error) {
	model := Log{}
	cursor := SessionFindOneByID(LOG_COLLECTION, ID, session)
	err := cursor.Decode(&model)
	if err != nil {
		return Log{}, err
	}

	return model, nil
}

func LogByIDNoSession(ID string) (Log, error) {
	model := Log{}
	cursor := FindOneByID(LOG_COLLECTION, ID)
	err := cursor.Decode(&model)
	if err != nil {
		return Log{}, err
	}

	return model, nil
}

func LogUpdate(model Log, session mongo.SessionContext) error {
	err := SessionUpdateOne(&model, session)
	if err != nil {
		return err
	}

	return nil
}

func logInformation(isError bool, paramLog ParamLog, data interface{}) (Log, error) {
	var log Log

	createLog := func(session mongo.SessionContext) error {

		err := session.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)

		if err != nil {
			return errors.New("log db error")
		}
		log, err = LogCreate(isError, paramLog, data, session)

		if err != nil {
			session.AbortTransaction(session)
			return err
		}
		return CommitWithRetry(session)
	}

	err := DBClient.UseSessionWithOptions(
		context.TODO(), options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return RunTransactionWithRetry(sctx, createLog)
		},
	)

	if err != nil {
		return log, errors.New("fail to initialize log")
	}

	return log, nil
}
func LogError(paramLog ParamLog, data interface{}) (Log, error) {
	if DBClient == nil {
		return Log{}, errors.New("No DB Client")
	}
	return logInformation(true, paramLog, data)
}
func LogInformation(paramLog ParamLog, data interface{}) (Log, error) {
	if DBClient == nil {
		return Log{}, errors.New("No DB Client")
	}
	return logInformation(false, paramLog, data)
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
func SaveOne(colName string, domain domain.BaseModel) error {

	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(colName)
	result, err := collection.InsertOne(context.TODO(), domain)
	if err != nil {
		return err
	}

	ID := result.InsertedID.(primitive.ObjectID)
	domain.SetDocumentID(ID)

	return nil
}
func SessionFindOneByID(colName string, ID string, session mongo.SessionContext) *mongo.SingleResult {
	objectID, _ := primitive.ObjectIDFromHex(ID)

	query := bson.M{"_id": objectID}
	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(colName)
	result := collection.FindOne(session, query)

	return result
}
func FindOneByID(colName string, ID string) *mongo.SingleResult {
	objectID, _ := primitive.ObjectIDFromHex(ID)

	query := bson.M{"_id": objectID}
	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(colName)
	result := collection.FindOne(context.TODO(), query)

	return result
}
func SessionUpdateOne(domain domain.BaseModel, session mongo.SessionContext) error {
	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(domain.CollectionName())
	document, err := toDoc(domain)
	if err != nil {
		return err
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
func toDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}
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
func SetupDB() error {

	// Set client options
	clusterUrl := os.Getenv("MONGO_CLUSTER_URL")
	logrus.Info("clusterUrl", clusterUrl)
	clientOptions := options.Client().ApplyURI(clusterUrl)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	// Return error if there problem
	if err != nil {
		logrus.Info("Unable to connect MONGO ", err.Error(), clusterUrl)
		return err
	}

	DBClient = client

	return nil
}

type ParamLog struct {
	TrCloser *io.Closer
	Span     *opentracing.Span
	Tag      string
}
