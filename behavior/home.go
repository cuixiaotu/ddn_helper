package behavior

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

var ddnLoginUrl = "https://www.ddn.cn/front/login"
var ddnIndexUrl = "https://www.ddn.cn/front/index"
var ddnBackUrl = "https://www.ddn.cn/back/netBar"

func (c *HomeClient) GoBackHome() error {
	var respBody string
	// 设置 OnResponse 处理跳转
	c.c.OnResponse(func(r *colly.Response) {
		body := string(r.Body)
		// ub 日志
		c.log.Printf("后台登录真正开始")
		respBody = body
	})

	err := c.c.Visit(ddnBackUrl)
	if err != nil {
		return err
	}

	// 等待所有回调执行完
	c.c.Wait()

	// 正则提取 location.href 中的 URL
	re := regexp.MustCompile(`location(?:\.href)?\s*=\s*["']([^"']+)["']`)
	matches := re.FindStringSubmatch(respBody)
	if len(matches) < 2 {
		c.log.Println("❌ 未匹配到跳转后台登录 URL")
		return errors.New("未匹配到跳转后台登录 URL")
	}

	targetURL := matches[1]
	c.log.Println("🔗 检测到跳转 URL:", targetURL)

	// 执行跳转
	if err := c.c.Visit(targetURL); err != nil {
		c.log.Println("❌ 跳转失败:", err)
	} else {
		c.log.Println("➡️ 已跳转到页面")
	}
	// 等待所有回调执行完
	c.c.Wait()

	c.log.Println("✅ 后台首页登录完成")
	return nil
}

// http://ip:port/netbar/common/ddnforward.jsp?netBarId=6681&timestamp=1764156233244&userId=6681&isddnBusiness=9999&mac=8086782b984f0af8ff5c59437e3529ee&targetPage=../index.jsp
func (c *HomeClient) Login() error {
	// 构造 form-urlencoded body
	values := url.Values{}
	values.Set("loginIdType", "userId")
	values.Set("province", c.options.Province)
	values.Set("city", c.options.City)
	values.Set("domainId", c.options.DomainId)
	values.Set("loginType", "userId")
	values.Set("userId", c.options.Username)
	values.Set("password", c.options.Password)
	values.Set("charPwd", c.options.Password)
	values.Set("mianze", "on")
	sendBody := values.Encode()
	c.log.Println("Send Body:", sendBody)
	// 初始化 CookieJar
	jar, _ := cookiejar.New(nil)
	c.c.SetCookieJar(jar)

	// ----------------------------
	// 2. 登录请求所需 Headers（强制补齐所有浏览器行为）
	// ----------------------------
	c.c.OnRequest(func(r *colly.Request) {
		setDefaultHeaders(r)
		c.log.Println(">>> Sending Request:", r.URL.String())
	})

	// 打印响应
	c.c.OnResponse(func(r *colly.Response) {
		saved := c.c.Cookies("https://www.ddn.cn")
		for _, sc := range saved {
			c.log.Printf("Collector cookie: %s=%s\n", sc.Name, sc.Value)
		}

	})

	c.c.OnError(func(r *colly.Response, err error) {
		c.log.Println("Error:", err)
	})

	// 发起 POST 登录请求
	err := c.c.Request("POST", ddnLoginUrl, bytes.NewBufferString(sendBody), nil, http.Header{
		"Content-Type": []string{"application/x-www-form-urlencoded"},
	})
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}

	c.log.Println("✅ 登录完成")
	return nil
}

func (c *HomeClient) CheckLogin() error {
	var resp string

	c.c.OnResponse(func(r *colly.Response) {
		resp = string(r.Body)
		saved := c.c.Cookies("https://www.ddn.cn")
		for _, sc := range saved {
			c.log.Printf("Collector cookie: %s=%s\n", sc.Name, sc.Value)
		}
	})

	c.c.OnRequest(func(r *colly.Request) {
		setDefaultHeaders(r)
	})

	// 登录完成后访问首页保持登录态
	err := c.c.Visit(ddnIndexUrl)
	if err != nil {
		return fmt.Errorf("visit index failed: %w", err)
	}

	// 等待所有回调执行完
	c.c.Wait()

	// 用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp))
	if err != nil {
		return err
	}

	// 找到表格
	loginFlag := doc.Find(".login_t")
	if loginFlag.Length() > 0 {
		return nil
	}

	return errors.New("login failed")
}

func setDefaultHeaders(r *colly.Request) {
	r.Headers.Set("Origin", "https://www.ddn.cn")
	r.Headers.Set("Referer", "https://www.ddn.cn/front/index")
	r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")

	r.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
	r.Headers.Set("Accept-Encoding", "gzip, deflate, br")

	// Fetch Metadata（非常关键）
	r.Headers.Set("Sec-Fetch-Dest", "document")
	r.Headers.Set("Sec-Fetch-Mode", "navigate")
	r.Headers.Set("Sec-Fetch-Site", "same-origin")
	r.Headers.Set("Sec-Fetch-User", "?1")

	// Client-Hints
	r.Headers.Set("Sec-Ch-Ua", `"Chromium";v="142", "Google Chrome";v="142", "Not_A Brand";v="99"`)
	r.Headers.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	r.Headers.Set("Sec-Ch-Ua-Mobile", "?0")

	r.Headers.Set("Upgrade-Insecure-Requests", "1")
}
