package behavior

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// 用户行为策略
var SubUrl = "http://ip:port"

func (c *HomeClient) BackIndex() error {
	indexUrl := fmt.Sprintf("%s%s", SubUrl, "/netbar/index.jsp")
	return c.c.Visit(indexUrl)
}

/** VisitMemberInfoPage
* 访问成员页面首页 /netbar/log.jsp?moduleId=memberinfo&moduleGroupId=1
*  http://ip:port/netbar/memberinfo/main.jsp?moduleName=%D3%C3%BB%A7%D0%C5%CF%A2%B2%E9%D1%AF
*
 */
func (c *HomeClient) GoMemberModule() error {
	memberLogUrl := fmt.Sprintf("%s%s", SubUrl, "/netbar/log.jsp?moduleId=memberinfo&moduleGroupId=1")
	memberInfoUrl := fmt.Sprintf("%s%s", SubUrl, "/netbar/memberinfo/memberinfo.jsp")
	var respBody string
	// 设置 OnResponse 处理跳转
	c.c.OnResponse(func(r *colly.Response) {
		body := string(r.Body)
		// ub 日志
		c.log.Printf("会员模块登录开始")
		respBody = body
		c.log.Printf(respBody)

	})

	err := c.c.Visit(memberLogUrl)
	if err != nil {
		return err
	}

	// 等待所有回调执行完
	c.c.Wait()

	c.log.Println("✅ 后台首页登录完成")
	//这里缺一个
	return c.c.Visit(memberInfoUrl)
}

/** VisitMemberInfoPage
* 访问成员页面首页 /netbar/log.jsp?moduleId=memberinfo&moduleGroupId=1
*  http://ip:port/netbar/memberinfo/main.jsp?moduleName=%D3%C3%BB%A7%D0%C5%CF%A2%B2%E9%D1%AF
*
 */
func (c *HomeClient) VisitMemberInfoPage() error {
	MemberInfoUrl := fmt.Sprintf("%s%s", SubUrl, "/netbar/memberinfo/main.jsp?moduleName=%D3%C3%BB%A7%D0%C5%CF%A2%B2%E9%D1%AF")
	return c.c.Visit(MemberInfoUrl)
}

// DoMemberRequest 使用 colly 发送 POST 请求获取会员数据
func (c *HomeClient) DoMemberRequest(queryStartTime, queryEndTime string, pageNo, pageSize int) error {
	memberPath := "/netbar/memberinfo/memberinfo.jsp"
	baseUrl := SubUrl + memberPath

	// 构建表单数据
	formData := url.Values{}
	formData.Set("disuseId", "")
	formData.Set("classId", "0")
	formData.Set("name", "")
	formData.Set("idNumber", "")
	formData.Set("serialNum", "")
	formData.Set("memberId", "")
	formData.Set("minBalance", "")
	formData.Set("status", "-1")
	formData.Set("queryStartTime", queryStartTime)
	formData.Set("queryEndTime", queryEndTime)
	formData.Set("pageNo", strconv.Itoa(pageNo))
	formData.Set("pageSize", strconv.Itoa(pageSize))

	c.c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-US;q=0.7,zh-TW;q=0.6")
		r.Headers.Set("Accept-Encoding", "gzip, deflate")

		r.Headers.Set("Host", "113.106.48.76:9080")
		r.Headers.Set("Origin", "http://ip:port")
		r.Headers.Set("Content-Type", "application/x-www-form-urlencoded")

		r.Headers.Set("Referer", "http://ip:port/netbar/memberinfo/tree.jsp?moduleName=%D3%C3%BB%A7%D0%C5%CF%A2%B2%E9%D1%AF")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")

		// 你抓包里的 Cookie 必须也带上
		c.log.Println(">>> Sending Request:", r.URL.String())
		c.log.Println(">>>r.Headers", r.Headers)
	})

	// 存储响应 body
	var responseBody string
	var queryUrl string

	// 成功接收到响应时触发
	c.c.OnResponse(func(r *colly.Response) {
		responseBody = string(r.Body)
		queryUrl = r.Request.URL.String()
	})

	// 错误处理
	c.c.OnError(func(r *colly.Response, err error) {
		c.log.Printf("Request failed with status %d: %s\n", r.StatusCode, err)
	})

	err := c.c.Request("POST", baseUrl, strings.NewReader(formData.Encode()), nil, map[string][]string{})
	if err != nil {
		c.log.Println("Failed to send POST request:", err)
		return err
	}
	// 等待回调执行完毕
	c.c.Wait()

	if responseBody == "" {
		c.log.Print("empty response from server")
		return errors.New("empty response from server")
	}

	if c.IsTimeOut(responseBody, queryUrl) {
		return errors.New("time out")
	}

	if err := formMemberHtml(responseBody); err != nil {
		c.log.Println("formMemberHtml error:", err)
		return err
	}

	return nil
}

// DoBirthRequest 使用 colly 发送 POST 请求获取会员数据
func (c *HomeClient) GoBirthPage() error {
	memberPath := "/netbar/memberinfo/memberbirthday.jsp"
	birthPage := SubUrl + memberPath
	return c.c.Visit(birthPage)
}

// DoBirthRequest 使用 colly 发送 POST 请求获取会员数据
func (c *HomeClient) DoBirthRequest(birthDay string, pageNo, pageSize int) error {
	memberPath := "/netbar/memberinfo/memberbirthday.jsp"
	baseUrl := SubUrl + memberPath

	// 构建表单数据 todo 补偿其他参数
	formData := url.Values{}
	formData.Set("verb", "calc")
	formData.Set("classId", "0")
	formData.Set("month", "12")
	formData.Set("day", birthDay)
	formData.Set("pageNo", strconv.Itoa(pageNo))
	formData.Set("pageSize", strconv.Itoa(pageSize))
	//verb=calc&classId=0&month=12&day=1&pageSize=25&pageNo=1

	c.c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-US;q=0.7,zh-TW;q=0.6")

		r.Headers.Set("Host", "113.106.48.76:9080")
		r.Headers.Set("Origin", "http://ip:port")
		r.Headers.Set("Content-Type", "application/x-www-form-urlencoded")

		r.Headers.Set("Referer", "http://ip:port/netbar/memberinfo/memberbirthday.jsp")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")

		// 你抓包里的 Cookie 必须也带上
		c.log.Println(">>> Sending Request:", r.URL.String())
	})

	// 存储响应 body
	var responseBody string
	var queryUrl string

	// 成功接收到响应时触发
	c.c.OnResponse(func(r *colly.Response) {
		responseBody = string(r.Body)
		queryUrl = r.Request.URL.String()
	})

	// 错误处理
	c.c.OnError(func(r *colly.Response, err error) {
		c.log.Printf("Request failed with status %d: %s\n", r.StatusCode, err)
	})

	// 手动发起 POST 请求
	err := c.c.Request("POST", baseUrl, strings.NewReader(formData.Encode()), nil, map[string][]string{})
	if err != nil {
		c.log.Println("Failed to send POST request:", err)
		return err
	}

	// 等待回调执行完毕
	c.c.Wait()

	if responseBody == "" {
		c.log.Print("empty response from server")
		return errors.New("empty response from server")
	}

	if c.IsTimeOut(responseBody, queryUrl) {
		return errors.New("time out")
	}

	if err := formBirthHtml(birthDay, responseBody); err != nil {
		c.log.Println("formMemberHtml error:", err)
		return err
	}

	return nil
}

func (c *HomeClient) IsTimeOut(respBody string, url string) bool {
	// 1. URL 直接判断
	if strings.Contains(url, "/timeout.html") {
		return true
	}

	// 2. 页面标题判断
	if strings.Contains(respBody, "<title>会话超时") {
		return true
	}

	// 3. 页面内容判断
	if strings.Contains(respBody, "当前会话超时") ||
		strings.Contains(respBody, "请重新登录") {
		return true
	}
	return false
}

func formMemberHtml(resp string) error {
	// 用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp))
	if err != nil {
		fmt.Println("formHtml : ", err.Error())
		return err
	}

	// 找到表格
	table := doc.Find("#tblData")

	var rows [][]string

	// 空行统计
	emptyCount := 0
	maxEmptyLog := 10
	maxEmptyErr := 100
	isEnd := false

	table.Find("tr").NextAll().Each(func(i int, tr *goquery.Selection) {

		// 跳过合计行
		if strings.Contains(tr.Text(), "合计") {
			return
		}

		var row []string
		allEmpty := true

		tr.Find("td").Each(func(_ int, td *goquery.Selection) {
			text := strings.TrimSpace(td.Text())

			// &nbsp; 识别成空
			if text == "" || text == "&nbsp;" {
				row = append(row, "")
			} else {
				row = append(row, text)
				allEmpty = false
			}
		})

		// 处理空行逻辑
		if allEmpty {
			emptyCount++

			// Log：超过 10 行才打一次日志（避免刷屏）
			if emptyCount == maxEmptyLog {
				log.Printf("⚠️ 检测到超过 %d 行连续空数据", maxEmptyLog)
			}

			return
		}

		// 如果当前行是有效数据，重置计数
		emptyCount = 0

		// 保存有效行
		if len(row) > 0 {
			rows = append(rows, row)
		}
	})

	// 空行超限报错
	if emptyCount > 100 {
		fmt.Printf("❌ 连续空行超过 %d 行，源页面可能异常 或页面完结", maxEmptyErr)
		isEnd = true
	}

	// ===== 写 CSV =====

	out, err := os.OpenFile("output.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	w := csv.NewWriter(out)
	defer w.Flush()

	for _, row := range rows {
		if err := w.Write(row); err != nil {
			return err
		}
	}

	fmt.Println("✅ 已追加写入 output.csv")
	if isEnd {
		return errors.New("已结束")
	}

	return nil
}

func formBirthHtml(birthDay string, resp string) error {
	// 用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp))
	if err != nil {
		fmt.Println("formHtml : ", err.Error())
		return err
	}
	// 找到表格
	table := doc.Find("#tblData")

	// 结果存储
	var rows [][]string

	//// 取表头
	//table.Find("tr").First().Find("td").Each(func(i int, s *goquery.Selection) {
	//	rows = append(rows, []string{}) // 初始化第一行
	//	rows[0] = append(rows[0], strings.TrimSpace(s.Text()))
	//})

	// 取表格内容
	table.Find("tr").NextAll().Each(func(i int, tr *goquery.Selection) {

		// 跳过合计行
		if strings.Contains(tr.Text(), "合计") {
			return
		}

		var row []string
		tr.Find("td").Each(func(_ int, td *goquery.Selection) {
			// 拿到纯文本
			text := strings.TrimSpace(td.Text())
			row = append(row, text)
		})

		if len(row) > 0 {
			rows = append(rows, row)
		}
	})

	// 写 CSV 文件
	csvNanme := fmt.Sprintf("birth_12_%s.csv", birthDay)
	out, err := os.Create(csvNanme)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	w := csv.NewWriter(out)
	defer w.Flush()

	for _, row := range rows {
		if err := w.Write(row); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("✅ 已生成 " + csvNanme)
	return nil
}
