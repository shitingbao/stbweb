package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	//这里应用了bson里面的M类型，但是D类型的使用会出错，原因待定
	//所以所有的条件以及结果集反馈都是用了M类型，其实就是map
	//注意多条件查询的对象这么写
	//eg bson.M{"age":"18","number": bson.M{"$in": []string{"1", "2"}}}
	//代表 age=18 and number in (1,2)
	//同理，不同的条件，只要把第二个条件字符串（$in）变成其他的即可，比如 bson.M{"number": bson.M{"$gt": 1}}。就是number大于1的
	//or多条件麻烦点 bson.M{"$or": []bson.M{bson.M{"age": "bb"}, bson.M{"number": 1}}}
)

var defaultTimeout = 10 * time.Second

//Mongodb 一个mongo连接对象
//保存数据库名称（Database），集合对象（CollectionDB用于操作），Client基本连接
type Mongodb struct {
	Database     string
	Client       *mongo.Client
	CollectionDB *mongo.Database
	Ctx          context.Context
}

//OpenMongoDb 连接一个mongodb
//输入连接字符串，待连接数据库，集合名称，反馈该集合对象和error
func OpenMongoDb(driver, database string) (*Mongodb, error) {
	ctx := context.TODO()
	// defer canf()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(driver))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	clo := client.Database(database)
	return &Mongodb{
		Client:       client,
		Database:     database,
		CollectionDB: clo,
		Ctx:          ctx,
	}, nil
}

//InsertOne 插入一条文档，传入一个content,一个集合对象，和待插入的bson的对象数据，反馈该数据的id和error
func (m *Mongodb) InsertOne(collect string, data bson.M) (interface{}, error) {
	res, err := m.CollectionDB.Collection(collect).InsertOne(m.Ctx, data)
	if err != nil {
		return nil, err
	}
	id := res.InsertedID
	return id, nil
}

//InsertMany 插入多条数据
func (m *Mongodb) InsertMany(collect string, data bson.M) (interface{}, error) {
	res, err := m.CollectionDB.Collection(collect).InsertOne(m.Ctx, data)
	if err != nil {
		return nil, err
	}
	id := res.InsertedID
	return id, nil
}

//Selectone 查询一条语句，输入集合，条件（bson.D中包含条件数组）,以及接受的结构所以result得是一个struct的指针，反馈查询出来数据
//opt是限制函数
func (m *Mongodb) Selectone(collect string, where bson.M, opts ...*options.FindOneOptions) (bson.M, error) {
	var result bson.M
	if err := m.CollectionDB.Collection(collect).FindOne(m.Ctx, where, opts...).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

//SelectAll 查询所有，输入集合和条件（bson.D中包含条件数组），反馈所有数据和error
//opt是限制函数
func (m *Mongodb) SelectAll(collect string, where bson.M, opts ...*options.FindOptions) ([]bson.M, error) {
	var result []bson.M
	cur, err := m.CollectionDB.Collection(collect).Find(m.Ctx, where, opts...)
	if err != nil {
		return nil, err
	}
	defer cur.Close(m.Ctx)

	for cur.Next(m.Ctx) {
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

//UpdateOne 更新一条数据
func (m *Mongodb) UpdateOne(collect string, where, update bson.M, opts ...*options.UpdateOptions) (interface{}, error) {
	res, err := m.CollectionDB.Collection(collect).UpdateOne(m.Ctx, where, update, opts...)
	if err != nil {
		return nil, err
	}
	return res.UpsertedID, nil
}

//UpdateMany 更新一条数据
func (m *Mongodb) UpdateMany(collect string, update interface{}, where bson.M, opts ...*options.UpdateOptions) (interface{}, error) {
	res, err := m.CollectionDB.Collection(collect).UpdateMany(m.Ctx, where, update, opts...)
	if err != nil {
		return nil, err
	}
	return res.UpsertedID, nil
}

//DeleteCollecting 删除集合
func (m *Mongodb) DeleteCollecting() error {
	if err := m.CollectionDB.Drop(m.Ctx); err != nil {
		return err
	}
	return nil
}

//DeleteDocument 删除文档
func (m *Mongodb) DeleteDocument(collect string, where bson.M, opts ...*options.FindOneAndDeleteOptions) error {
	if res := m.CollectionDB.Collection(collect).FindOneAndDelete(m.Ctx, where, opts...); res.Err() != nil {
		return res.Err()
	}
	return nil
}

//CloseCtx monggodb释放ctx资源
func (m *Mongodb) CloseCtx() {
	m.Client.Disconnect(m.Ctx)
}
