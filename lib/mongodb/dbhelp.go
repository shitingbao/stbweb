package mongodb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

//连接一个mongodb，输入连接字符串，待连接数据库，集合名称，反馈该集合对象和error
func openMongoDb(driver, database, collect string) (*mongo.Collection, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(driver))
	if err != nil {
		return nil, err
	}
	clo := client.Database(database).Collection(collect)
	return clo, nil
}

//插入一条文档，传入一个content,一个集合对象，和待插入的bson的对象数据，反馈该数据的id和error
func insertOne(ctx context.Context, clo *mongo.Collection, data bson.M) (interface{}, error) {
	res, err := clo.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	id := res.InsertedID
	return id, nil
}

//插入多条数据
func insertMany(ctx context.Context, clo *mongo.Collection, data bson.M) (interface{}, error) {
	res, err := clo.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	id := res.InsertedID
	return id, nil
}

//查询一条语句，输入集合，条件（bson.D中包含条件数组）,以及接受的结构所以result得是一个struct的指针，反馈查询出来数据
//opt是限制函数
func selectone(ctx context.Context, clo *mongo.Collection, where bson.D, opts ...*options.FindOneOptions) (bson.M, error) {
	var result bson.M
	if err := clo.FindOne(ctx, where, opts...).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

//查询所有，输入集合和条件（bson.D中包含条件数组），反馈所有数据和error
//opt是限制函数
func selectAll(ctx context.Context, clo *mongo.Collection, where bson.D, opts ...*options.FindOptions) ([]bson.M, error) {
	var result []bson.M
	cur, err := clo.Find(ctx, where, opts...)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var res bson.M
		if err := cur.Decode(&res); err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

//更新一条数据
func updateOne(ctx context.Context, clo *mongo.Collection, update interface{}, where bson.D, opts ...*options.UpdateOptions) {
	res, err := clo.UpdateOne(ctx, where, update, opts...)
	if err != nil {
		return
	}
	log.Println(res.UpsertedID)
}

//更新一条数据
func updateMany(ctx context.Context, clo *mongo.Collection, update interface{}, where bson.D, opts ...*options.UpdateOptions) {
	res, err := clo.UpdateMany(ctx, where, update, opts...)
	if err != nil {
		return
	}
	log.Println(res.UpsertedID)
}

//删除集合
func deleteCollecting(ctx context.Context, clo *mongo.Collection) {
	if err := clo.Drop(ctx); err != nil {
		return
	}
}

//删除文档
func deleteDocument(ctx context.Context, clo *mongo.Collection, where bson.D, opts ...*options.FindOneAndDeleteOptions) {
	if res := clo.FindOneAndDelete(ctx, where, opts...); res.Err() != nil {
		return
	}
}
