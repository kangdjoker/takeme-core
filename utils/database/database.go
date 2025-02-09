package database

import (
	"context"
	"os"
	"strconv"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/basic"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DBClient *mongo.Client

func FindCount(paramLog *basic.ParamLog, colName string, query bson.M) (int64, error) {

	opts := options.CountOptions{}

	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(colName)
	total, err := collection.CountDocuments(
		context.TODO(),
		query,
		&opts,
	)

	if err != nil {
		return 0, utils.ErrorInternalServer(paramLog, utils.QueryFailed, err.Error())
	}

	return total, nil
}

func FindWithJoin(paramLog *basic.ParamLog, colName string, query []bson.M) (*mongo.Cursor, error) {
	a := true
	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(colName)
	cursor, err := collection.Aggregate(
		context.TODO(),
		query,
		&options.AggregateOptions{
			AllowDiskUse: &a,
		},
	)

	if err != nil {
		errorMessage := err.Error()
		return nil, utils.ErrorInternalServer(paramLog, utils.QueryFailed, errorMessage)
	}

	return cursor, nil
}

func FindOneByID(colName string, ID string) *mongo.SingleResult {
	objectID, _ := primitive.ObjectIDFromHex(ID)

	query := bson.M{"_id": objectID}
	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(colName)
	result := collection.FindOne(context.TODO(), query)

	return result
}

func FindOne(colName string, query bson.M) *mongo.SingleResult {
	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(colName)
	result := collection.FindOne(context.TODO(), query)

	return result
}

func Find(paramLog *basic.ParamLog, colName string, query bson.M, page string, limit string) (*mongo.Cursor, error) {

	opts := options.Find()
	opts.SetSort(bson.D{{"time", -1}})

	if page != "" && limit != "" {
		p, _ := strconv.Atoi(page)
		l, _ := strconv.Atoi(limit)
		opts.SetSkip(int64((p - 1) * l))
		opts.SetLimit(int64(l))
	}

	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(colName)
	cursor, err := collection.Find(
		context.TODO(),
		query,
		opts,
	)
	if err != nil {
		return nil, utils.ErrorInternalServer(paramLog, utils.QueryFailed, err.Error())
	}

	return cursor, nil
}

func FindOrderByID(paramLog *basic.ParamLog, colName string, query bson.M, page string, limit string) (*mongo.Cursor, error) {

	opts := options.Find()
	opts.SetSort(bson.D{{"_id", -1}})

	if page != "" && limit != "" {
		p, _ := strconv.Atoi(page)
		l, _ := strconv.Atoi(limit)
		opts.SetSkip(int64((p - 1) * l))
		opts.SetLimit(int64(l))
	}

	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(colName)
	cursor, err := collection.Find(
		context.TODO(),
		query,
		opts,
	)
	if err != nil {
		return nil, utils.ErrorInternalServer(paramLog, utils.QueryFailed, err.Error())
	}

	return cursor, nil
}

func FindAllOrderByID(paramLog *basic.ParamLog, colName string, query bson.M, page string, limit string) (*mongo.Cursor, error) {

	opts := options.Find()
	opts.SetSort(bson.D{{"_id", -1}})

	if page != "" && limit != "" {
		p, _ := strconv.Atoi(page)
		l, _ := strconv.Atoi(limit)
		opts.SetSkip(int64((p - 1) * l))
		opts.SetLimit(int64(l))
	}

	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(colName)
	cursor, err := collection.Find(
		context.TODO(),
		query,
		opts,
	)
	if err != nil {
		return nil, utils.ErrorInternalServer(paramLog, utils.QueryFailed, err.Error())
	}

	return cursor, nil
}

func Aggregate(paramLog *basic.ParamLog, colName string, query []bson.M) (*mongo.Cursor, error) {
	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(colName)
	cursor, err := collection.Aggregate(
		context.TODO(),
		query,
	)

	if err != nil {
		return nil, utils.ErrorInternalServer(paramLog, utils.QueryFailed, err.Error())
	}

	return cursor, nil
}

func IsExist(paramLog *basic.ParamLog, colName string, query bson.M) (*mongo.Cursor, error) {
	opts := options.Find()
	opts.SetLimit(1)

	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(colName)
	cursor, err := collection.Find(
		context.TODO(),
		query,
		opts,
	)
	if err != nil {
		return nil, utils.ErrorInternalServer(paramLog, utils.QueryFailed, err.Error())
	}

	return cursor, nil
}

func SaveOne(paramLog *basic.ParamLog, colName string, domain domain.BaseModel) error {

	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(colName)
	result, err := collection.InsertOne(context.TODO(), domain)
	if err != nil {
		return utils.ErrorInternalServer(paramLog, utils.InsertFailed, err.Error())
	}

	ID := result.InsertedID.(primitive.ObjectID)
	domain.SetDocumentID(ID)

	return nil
}

func UpdateOne(paramLog *basic.ParamLog, colName string, domain domain.BaseModel) error {

	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(colName)
	document, err := toDoc(domain)
	if err != nil {
		return utils.ErrorInternalServer(paramLog, utils.UpdateFailed, err.Error())
	}

	filter := bson.M{"_id": bson.M{"$eq": domain.GetDocumentID()}}
	update := bson.M{"$set": document}

	_, err = collection.UpdateOne(
		context.TODO(),
		filter,
		update,
	)
	if err != nil {
		return utils.ErrorInternalServer(paramLog, utils.UpdateFailed, err.Error())
	}

	return nil
}

func UpdateQuery(paramLog *basic.ParamLog, colName string, id primitive.ObjectID, update bson.M) error {

	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(colName)

	filter := bson.M{"_id": bson.M{"$eq": id}}

	_, err := collection.UpdateOne(
		context.TODO(),
		filter,
		update,
	)
	if err != nil {
		return utils.ErrorInternalServer(paramLog, utils.UpdateFailed, err.Error())
	}

	return nil
}

func UpdateMany(paramLog *basic.ParamLog, domains []domain.BaseModel) error {

	ctx := context.TODO()
	session, err := DBClient.StartSession()
	if err != nil {
		return utils.ErrorInternalServer(paramLog, utils.UpdateFailed, err.Error())
	}

	err = session.StartTransaction()
	if err != nil {
		return utils.ErrorInternalServer(paramLog, utils.UpdateFailed, err.Error())
	}

	var docol []DocAndCollection
	for _, element := range domains {
		document, _ := toDoc(element)
		docol = append(docol, DocAndCollection{
			ID:         element.GetDocumentID(),
			Document:   document,
			Collection: DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(element.CollectionName()),
		})

	}

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		for _, a := range docol {
			filter := bson.M{"_id": bson.M{"$eq": a.ID}}
			update := bson.M{"$set": a.Document}
			_, err = a.Collection.UpdateOne(
				ctx,
				filter,
				update,
			)

			if err != nil {
				return utils.ErrorInternalServer(paramLog, utils.UpdateFailed, err.Error())
			}
		}

		session.CommitTransaction(sc)
		return nil
	}); err != nil {
		session.EndSession(ctx)
		return utils.ErrorInternalServer(paramLog, utils.UpdateFailed, err.Error())
	}
	session.EndSession(ctx)
	return nil
}

func Update(paramLog *basic.ParamLog, colName string, filter bson.M, changes bson.D) (*mongo.UpdateResult, error) {
	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(colName)
	result, err := collection.UpdateMany(
		context.TODO(),
		filter,
		changes,
	)
	if err != nil {
		return nil, utils.ErrorInternalServer(paramLog, utils.UpdateFailed, err.Error())
	}

	return result, nil
}

func DeleteOne(paramLog *basic.ParamLog, colName string, domain domain.BaseModel) error {
	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(colName)
	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": domain.GetDocumentID()})
	if err != nil {
		return utils.ErrorInternalServer(paramLog, utils.DeleteFailed, err.Error())
	}

	return nil
}

func DeleteInActive(paramLog *basic.ParamLog, colName string, domain domain.BaseModel) error {
	collection := DBClient.Database(os.Getenv("MONGO_DB_NAME")).Collection(colName)

	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": domain.GetDocumentID(), "active": false, "pending": false})
	if err != nil {
		return utils.ErrorInternalServer(paramLog, utils.DeleteFailed, err.Error())
	}

	return nil
}

func SetupDB() error {

	// Set client options
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_CLUSTER_URL"))

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	// Return error if there problem
	if err != nil {
		return err
	}

	DBClient = client

	return nil
}

func CloseDB() error {
	err := DBClient.Disconnect(context.TODO())
	// Return error if there problem
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

type DocAndCollection struct {
	Collection *mongo.Collection
	Document   *bson.D
	ID         primitive.ObjectID
}
