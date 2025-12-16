package behavior

import (
	"errors"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

type UserLoginOptions struct {
	Url      string
	Province string
	City     string
	DomainId string
	Username string
	Password string
}

type UserBehavior struct {
	Options UserLoginOptions

	Home *HomeClient

	log *log.Logger
	jar *cookiejar.Jar
	ua  string
}

type HomeClient struct {
	c       *colly.Collector
	options UserLoginOptions
	log     *log.Logger
}

// 初始化基本信息
func (u *UserBehavior) init() error {
	u.log = log.New(os.Stdout, "ub_", log.LstdFlags)
	u.ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36"

	// 公共 cookie jar
	jar, _ := cookiejar.New(nil)
	u.jar = jar

	homeClient := &HomeClient{
		c:       u.initCollector(),
		options: u.Options,
		log:     u.log,
	}
	u.Home = homeClient

	return nil
}

// 初始化 collector 的公共方法
// 通用的打印
func (u *UserBehavior) initCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{}),
		colly.AllowURLRevisit(),
		colly.MaxDepth(1),
	)

	// 设置 CookieJar
	c.SetCookieJar(u.jar)

	c.WithTransport(&http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // 建立 TCP 连接超时
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   120 * time.Second,
		ResponseHeaderTimeout: 120 * time.Second, // ← 等待响应头的超时（解决你遇到的错误）
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
	}) // 保证 cookie jar 生效

	c.SetRequestTimeout(130 * time.Second)

	c.OnRequest(func(r *colly.Request) {
		u.log.Println(">>> Visiting:", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		u.log.Println("<<< Status:", r.StatusCode)
		u.log.Println("<<< URL:", r.Request.URL.String())
		u.log.Println("<<< Headers:", r.Headers)
		u.log.Println("<<< Body:", string(r.Body))
	})

	c.OnError(func(r *colly.Response, err error) {
		u.log.Println("Error:", err, r.Request.URL)
	})

	return c
}

// 启动行为
func NewUserBehavior(opt UserLoginOptions) *UserBehavior {
	ub := &UserBehavior{
		Options: opt,
	}
	err := ub.init()
	if err != nil {
		panic("ub init fail")
	}

	return ub
}

// 启动行为
func (u *UserBehavior) Start() error {
	if u.Home == nil || u.Home.c == nil {
		return errors.New("home client not initialized")
	}

	return u.Mock()
}

// 示例 Mock 方法
func (u *UserBehavior) Mock() error {
	var err error
	err = u.Home.Login()
	if err != nil {
		return errors.New("home login fail")
	}

	time.Sleep(time.Duration(randomInt(10, 15)) * time.Second)
	err = u.Home.CheckLogin()
	if err != nil {
		return errors.New("home CheckLogin fail")
	}

	time.Sleep(time.Duration(randomInt(6, 8)) * time.Second)
	err = u.Home.GoBackHome()
	if err != nil {
		return errors.New("home GoBackHome fail")
	}

	time.Sleep(time.Duration(randomInt(20, 30)) * time.Second)
	//已完成全部登录前置流程操作
	mockUser(u)

	return nil
}

func mockUser(u *UserBehavior) {
	var err error
	queryStartTime := "2025-01-01"
	pageNo := 5
	pageSize := 999
	//birthDate := 1

	for {
		//去首页尝试保活
		err = u.Home.BackIndex()
		if err != nil {
			u.log.Println("back Index fail")
			break
		}
		time.Sleep(time.Duration(randomInt(8, 16)) * time.Second)

		err = u.Home.GoMemberModule()
		if err != nil {
			u.log.Println("back VisitMemberInfoPage fail")
			break
		}
		time.Sleep(time.Duration(randomInt(20, 30)) * time.Second)
		err = u.Home.DoMemberRequest(queryStartTime, "", pageNo, pageSize)
		if err != nil {
			u.log.Println("back DoMemberRequest fail")
			break
		}

		//if birthDate <= 31 {
		//	time.Sleep(time.Duration(randomInt(3, 4)) * time.Minute)
		//} else {
		//}
		time.Sleep(time.Duration(randomInt(5, 6)) * time.Minute)
		pageNo = pageNo + 1

		//if birthDate <= 31 {
		//	err = u.Home.GoBirthPage()
		//	if err != nil {
		//		u.log.Println("back DoBirthRequest fail")
		//	}
		//
		//	time.Sleep(time.Duration(randomInt(10, 16)) * time.Second)
		//
		//	queryBirthDay := fmt.Sprintf("%d", birthDate)
		//	err = u.Home.DoBirthRequest(queryBirthDay, 1, 600)
		//	if err != nil {
		//		u.log.Println("back DoBirthRequest fail")
		//		break
		//	}
		//	time.Sleep(time.Duration(randomInt(2, 3)) * time.Minute)
		//}

		//birthDate = birthDate + 1

		if pageNo >= 100 {
			break
		}
	}

}

func randomInt(min, max int) int {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min+1) + min
}

//type Score struct {
//	SpaceId int    `json:"space_id"`
//	GoodsId uint   `json:"goods_id"`
//	Value   string `json:"value"`
//}
//
//func (c *HomeClient) CarpetPost() {
//	goods := []uint{575, 576, 577, 578, 579, 580, 581}
//	maxSpaceId := 457
//
//	url1 := "https://phenixai.cn/upload/space-goods/add"
//	headers := map[string][]string{
//		"Content-Type": {"application/json"},
//	}
//
//	for _, goodId := range goods {
//		for spaceId := 1; spaceId <= maxSpaceId; spaceId++ {
//			score := Score{
//				SpaceId: spaceId,
//				GoodsId: goodId,
//				Value:   "0",
//			}
//			body, _ := json.Marshal(score)
//			reader := bytes.NewReader(body)
//			err := c.c.Request("POST", url1, reader, nil, headers)
//			if err != nil {
//				c.log.Println("Failed to send POST request:", err)
//			}
//			// 等待回调执行完毕
//			c.c.Wait()
//			time.Sleep(100 * time.Millisecond)
//
//		}
//	}
//
//}
