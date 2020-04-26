package stboutserver

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"stbweb/lib/external_service/stbserver"
	"sync"
)

const (
	//Port 服务端口
	Port = ":5000"
)

//StbServe 外部调用结构体
type StbServe struct{}

//GetSummonerInfo 信息获取
func (s *StbServe) GetSummonerInfo(ctx context.Context, iden *stbserver.Identity) (*stbserver.Character, error) {
	var skillLists []*stbserver.Skill
	var summonerLists []*stbserver.Summoner
	skill := &stbserver.Skill{
		Ordinary: 2.25,
		Qkill:    "cutsteel",
		Wkill:    "windwall",
		Ekill:    "run",
		Rkill:    "yaton",
	}

	summoner := &stbserver.Summoner{
		Dkill: "shan",
		Fkill: "fire",
	}
	return &stbserver.Character{
		Xaxis:    1,
		Yaxis:    1,
		Zaxis:    1,
		Area:     iden.Idcard,
		Name:     iden.Name,
		Skill:    append(skillLists, skill),
		Summoner: append(summonerLists, summoner),
	}, nil
}

//PutSummonerInfo 实时信息发送
func (s *StbServe) PutSummonerInfo(cli stbserver.StbServer_PutSummonerInfoServer) error {
	for {
		// if i > 3 {
		// 	return nil
		// }
		da, err := cli.Recv()
		if err != nil {
			log.Println("err:", err)
			return err
		}
		log.Println("da:", da)
	}
}

//GetAllSummonerInfo 实时信息反馈
func (s *StbServe) GetAllSummonerInfo(iden *stbserver.Identity, req stbserver.StbServer_GetAllSummonerInfoServer) error {
	var skillLists []*stbserver.Skill
	var summonerLists []*stbserver.Summoner
	skill := &stbserver.Skill{
		Ordinary: 2.25,
		Qkill:    "cutsteel",
		Wkill:    "windwall",
		Ekill:    "run",
		Rkill:    "yaton",
	}

	summoner := &stbserver.Summoner{
		Dkill: "shan",
		Fkill: "fire",
	}
	i := 0
	for {
		if i > 3 {
			return nil
		}
		if err := req.Send(&stbserver.Character{
			Xaxis:    1,
			Yaxis:    1,
			Zaxis:    1,
			Area:     "22.5",
			Name:     "yasuo",
			Skill:    append(skillLists, skill),
			Summoner: append(summonerLists, summoner),
		}); err != nil {
			return err
		}
		i++
	}
}

//ShareSummonerInfo 信息共享
func (s *StbServe) ShareSummonerInfo(cli stbserver.StbServer_ShareSummonerInfoServer) error {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for {
			da, err := cli.Recv()
			if err != nil {
				log.Println("get mes err:", err)
				break
			}
			log.Println("da:", da)
		}
		log.Println("接收完成")
		wg.Done()
	}()

	var skillLists []*stbserver.Skill
	var summonerLists []*stbserver.Summoner
	skill := &stbserver.Skill{
		Ordinary: 2.25,
		Qkill:    "cutsteel",
		Wkill:    "windwall",
		Ekill:    "run",
		Rkill:    "yaton",
	}

	summoner := &stbserver.Summoner{
		Dkill: "shan",
		Fkill: "fire",
	}
	go func() {
		i := 0
		for {
			if i > 3 {
				break
			}
			log.Println("发送", i)
			if err := cli.Send(&stbserver.Character{
				Xaxis:    int64(i), //当这里是0的时候，接收方是没有该属性数据的，并不是为0值，而是直接忽略了该属性
				Yaxis:    int64(i),
				Zaxis:    int64(i),
				Area:     "", //同理当这里是“”空字符串的时候，接收方是没有该属性数据的，并不是为空字符串值，而是直接忽略了该属性
				Name:     "yasuo",
				Skill:    append(skillLists, skill),
				Summoner: append(summonerLists, summoner),
			}); err != nil {
				log.Println("err:", err)
				break
			}
			i++
		}
		log.Println("发送完成")
		wg.Done()
	}()
	wg.Wait()
	return nil
}

//SendFile 文件传输
func (s *StbServe) SendFile(cli stbserver.StbServer_SendFileServer) error {

	fDir, err := os.Executable()
	if err != nil {
		panic(err)
	}

	fURL := filepath.Join(filepath.Dir(fDir), "assets")
	mkdir(fURL)
	f, err := os.Create(filepath.Join(fURL, "test.json"))
	if err != nil {
		return err
	}
	defer f.Close()
	for {
		da, err := cli.Recv()
		if err != nil {
			log.Println("err:", err)
			break
		}
		log.Println("name:", da.Filename)
		f.Write(da.Filedata)
	}
	return nil
}

func mkdir(url string) {
	_, err := os.Stat(url)
	if err == nil {
		return
	}
	if os.IsNotExist(err) {
		log.Println("创建目录")
		os.MkdirAll(url, os.ModePerm)
	}
}
