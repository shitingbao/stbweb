package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	//这里应用了bson里面的M类型，但是D类型的使用会出错，原因待定
	//所以所有的条件以及结果集反馈都是用了M类型，其实就是map
)

//Mongodb 一个mongo连接对象
//保存数据库名称（Database），集合名称（Collect），集合对象（CollectionDB用于操作），过期时间（OutTime，默认10S,可以重新指定），Cancelf取消函数（释放资源函数）
type Mongodb struct {
	Database     string
	Collect      string
	CollectionDB *mongo.Collection
	OutTime      time.Duration
	Cancelf      context.CancelFunc
	Ctx          context.Context
}

//OpenMongoDb 连接一个mongodb
//输入连接字符串，待连接数据库，集合名称，反馈该集合对象和error
func OpenMongoDb(driver, database, collect string) (*Mongodb, error) {
	ctx, canf := context.WithTimeout(context.Background(), 10*time.Second)
	// defer canf()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(driver))
	if err != nil {
		return nil, err
	}
	clo := client.Database(database).Collection(collect)
	return &Mongodb{
		Database:     database,
		Collect:      collect,
		CollectionDB: clo,
		OutTime:      10 * time.Second,
		Cancelf:      canf,
		Ctx:          ctx,
	}, nil
}

//InsertOne 插入一条文档，传入一个content,一个集合对象，和待插入的bson的对象数据，反馈该数据的id和error
func (m *Mongodb) InsertOne(data bson.M) (interface{}, error) {
	res, err := m.CollectionDB.InsertOne(m.Ctx, data)
	if err != nil {
		return nil, err
	}
	id := res.InsertedID
	return id, nil
}

//InsertMany 插入多条数据
func (m *Mongodb) InsertMany(data bson.M) (interface{}, error) {
	res, err := m.CollectionDB.InsertOne(m.Ctx, data)
	if err != nil {
		return nil, err
	}
	id := res.InsertedID
	return id, nil
}

//Selectone 查询一条语句，输入集合，条件（bson.D中包含条件数组）,以及接受的结构所以result得是一个struct的指针，反馈查询出来数据
//opt是限制函数
func (m *Mongodb) Selectone(where bson.M, opts ...*options.FindOneOptions) (bson.M, error) {
	var result bson.M
	if err := m.CollectionDB.FindOne(m.Ctx, where, opts...).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

//SelectAll 查询所有，输入集合和条件（bson.D中包含条件数组），反馈所有数据和error
//opt是限制函数
func (m *Mongodb) SelectAll(where bson.M, opts ...*options.FindOptions) ([]bson.M, error) {
	var result []bson.M
	cur, err := m.CollectionDB.Find(m.Ctx, where, opts...)
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
func (m *Mongodb) UpdateOne(update interface{}, where bson.M, opts ...*options.UpdateOptions) (interface{}, error) {
	res, err := m.CollectionDB.UpdateOne(m.Ctx, where, update, opts...)
	if err != nil {
		return nil, err
	}
	return res.UpsertedID, nil
}

//UpdateMany 更新一条数据
func (m *Mongodb) UpdateMany(update interface{}, where bson.M, opts ...*options.UpdateOptions) (interface{}, error) {
	res, err := m.CollectionDB.UpdateMany(m.Ctx, where, update, opts...)
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
func (m *Mongodb) DeleteDocument(where bson.M, opts ...*options.FindOneAndDeleteOptions) error {
	if res := m.CollectionDB.FindOneAndDelete(m.Ctx, where, opts...); res.Err() != nil {
		return res.Err()
	}
	return nil
}

//UpOutTime 更新过期时间
func (m *Mongodb) UpOutTime(outime time.Duration) {
	m.Cancelf()
	ctx, canf := context.WithTimeout(context.Background(), outime)
	m.Ctx = ctx
	m.OutTime = outime
	m.Cancelf = canf
}

//CloseCtx monggodb释放ctx资源
func (m *Mongodb) CloseCtx() {
	m.Cancelf()
}
