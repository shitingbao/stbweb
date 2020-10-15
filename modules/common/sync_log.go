package common

import (
	"errors"
	"io/ioutil"
	"path"
	"stbweb/core"
	"stbweb/lib/excel"
	base "stbweb/lib/file_base"
	"stbweb/lib/task"

	"github.com/sirupsen/logrus"
)

type syncLog struct{}

const lockNx = "log_lock"

func init() {
	core.RegisterFun("synclog", new(syncLog), true)
}

func (*syncLog) Get(p *core.ElementHandleArgs) {
	if p.APIInterceptionGet("sg", nil, logUpdata) {
		return
	}
}

func logUpdata(param interface{}, p *core.ElementHandleArgs) error {
	m := p.Req.URL.Query().Get("starttime")
	if m == "" {
		m = "0 0 1 * * ?"
	}
	ts := task.NewTask(p.Usr, "log_into", m, operaFunc) //默认晚上1点执行
	ts.Run(core.Ddb, core.Rds)
	return nil
}
func operaFunc() error {
	//带守护进程的分布式锁
	if core.DistributeLock("shitingbao", opera) {
		return nil
	}
	return errors.New("执行失败")
}
func opera() {
	fileList, err := ioutil.ReadDir(core.DefaultFilePath)
	if err != nil {
		logrus.WithFields(logrus.Fields{"ReadDir": err.Error()}).Error("opera")
		return
	}
	for _, v := range fileList {
		if v.IsDir() {
			continue
		}
		//解析出数据包中的文件，依次载入
		switch path.Ext(path.Base(v.Name())) {
		case "zip":
			paseZip(v.Name())
		}
	}
}

func paseZip(name string) {
	fList, err := base.ZipParse(path.Join(core.DefaultFilePath, name), "./")
	if err != nil {
		logrus.WithFields(logrus.Fields{"ZipParse": err.Error()}).Error("zip")
		return
	}
	for _, v := range fList {
		loadData(v)
	}
}

func loadData(path string) {
	ddb, err := core.Ddb.Begin()
	if err != nil {
		ddb.Rollback()
		return
	}
	da := excel.PaseCscOrTxt(path, "", false)
	for _, val := range da {
		stmt, err := ddb.Prepare("insert into order_good(order_name,order_info,announcer,receiver,release_time,receive_time,amount,completion) values(?,?,?,?,?,?,?,?,?)")
		if err != nil {
			ddb.Rollback()
			return
		}
		sv := []interface{}{}
		for _, vl := range val {
			sv = append(sv, vl)
		}
		if _, err := stmt.Exec(sv...); err != nil {
			ddb.Rollback()
			return
		}
	}
	ddb.Commit()
}
