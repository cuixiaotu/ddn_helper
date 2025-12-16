package ui

import (
	"log"
	"strings"
	"tb_helper/behavior"
	"tb_helper/constant"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type HomePage struct {
	log   *log.Logger
	colly *log.Logger
}

func NewHomePage(logger *log.Logger) *HomePage {
	return &HomePage{
		log: logger,
	}
}

func (h *HomePage) Box() {
	myApp := app.New()
	fyne.SetCurrentApp(myApp) // 确保 CurrentApp 可用（通常自动设置）
	myWindow := myApp.NewWindow("ddn-helper")

	//省
	provinceNames := make([]string, 0, len(constant.ProvinceAndCityArr))
	for _, p := range constant.ProvinceAndCityArr {
		provinceNames = append(provinceNames, p.Name)
	}
	//市
	cityOptions := []string{}
	if len(constant.ProvinceAndCityArr) > 0 {
		cityOptions = append(cityOptions, constant.ProvinceAndCityArr[0].Cities...)
	}

	citySelect := widget.NewSelect(cityOptions, func(selected string) {
		// 可按需处理城市选择事件
		//h.log.Println("city selected:", selected)
	})
	// 创建 Select 控件
	provinceSelect := widget.NewSelect(provinceNames, func(selected string) {
		// 当省份改变时，找到对应城市并更新 citySelect 的 Options
		newCities := []string{}
		for _, p := range constant.ProvinceAndCityArr {
			if p.Name == selected {
				// 复制 slice
				newCities = append(newCities, p.Cities...)
				break
			}
		}
		// 更新 citySelect 的 options & 清空已选中项
		citySelect.Options = newCities
		citySelect.PlaceHolder = "请选择"
		citySelect.Refresh()

		// 自动选中第一个城市，避免显示空白
		if len(newCities) > 0 {
			citySelect.SetSelected(newCities[0])
		} else {
			citySelect.SetSelected("") // 没有城市时清空
		}
	})
	// 默认选择第一个省与其第一个城市（可选）
	if len(provinceNames) > 0 {
		provinceSelect.SetSelected(provinceNames[0])
	}

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("请输入账号")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("请输入密码")

	// 布局
	form := container.NewVBox(
		widget.NewForm(
			&widget.FormItem{Text: "省份：", Widget: provinceSelect},
			&widget.FormItem{Text: "城市：", Widget: citySelect},
		),
	)

	loginButton := widget.NewButton("确定登录", func() {
		h.log.Printf("provinceSelect %s, citySelect %s, ", provinceSelect.Selected, citySelect.Selected)
		domainId := GetDomainId(provinceSelect.Selected, citySelect.Selected)
		opts := behavior.UserLoginOptions{
			Url:      "http://ip:port",
			Province: provinceSelect.Selected,
			City:     citySelect.Selected,
			DomainId: domainId,
			Username: usernameEntry.Text,
			Password: passwordEntry.Text,
		}

		// 显示 loading
		loadingDialog := showLoading(myWindow)
		loadingDialog.dialog.Show()

		//ub := behavior.NewUserBehavior(opts)
		//ub.Start()
		//ub.Home.CarpetPost()

		// 异步执行 Start，避免阻塞 UI
		go func() {
			// 更新进度的回调函数
			//updateProgress := func(val float64) {
			//	fyne.CurrentApp().Driver().DoFromGoroutine(func() {
			//		loadingDialog.progressBar.SetValue(val)
			//	}, false)
			//}

			//// 模拟分步操作
			//time.Sleep(500 * time.Millisecond)
			//updateProgress(30)
			//time.Sleep(800 * time.Millisecond)
			//updateProgress(70)
			//time.Sleep(600 * time.Millisecond)
			//updateProgress(99)

			// 执行耗时操作
			//process, err :=
			ub := behavior.NewUserBehavior(opts)
			ub.Start()

			// 完成后关闭 loading
			fyne.CurrentApp().Driver().DoFromGoroutine(func() {
				loadingDialog.dialog.Hide()
				loadingDialog.dialog.Dismiss()
				dialog.ShowInformation("成功", "导出完成！", myWindow)
			}, false)
			//myWindow.Close()
		}()
	})

	content := container.NewVBox(
		form,
		widget.NewForm(
			&widget.FormItem{Text: "账号：", Widget: usernameEntry},
			&widget.FormItem{Text: "密码：", Widget: passwordEntry},
		),
		loginButton,
	)

	myWindow.Resize(fyne.NewSize(400, 250))
	myWindow.CenterOnScreen()
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

// 根据省份名 城市名获取domainId
func GetDomainId(provinceName, cityName string) string {
	var domainId string
	for _, province := range constant.ProvinceAndCityArr {
		if provinceName == province.Name {
			domainId = province.Code
			break
		}
	}
	if strings.Contains(domainId, "广东") {
		domainId = "4"
		if strings.Contains(cityName, "深圳") && strings.Contains(cityName, "深圳网络出租屋") {
			domainId = "1"
		}
		if strings.Contains(cityName, "东莞") || strings.Contains(cityName, "惠州") {
			domainId = "3"
		}
	}
	return domainId
}

type LoadingDialog struct {
	dialog *dialog.CustomDialog
	//progressBar *widget.ProgressBar
}

// 显示进度条
func showLoading(parent fyne.Window) *LoadingDialog {
	loading := widget.NewLabel("正在导出，请稍候...")
	icon := widget.NewIcon(theme.ViewRefreshIcon())
	//progressBar := widget.NewProgressBar() // 创建可设置值的 ProgressBar
	//progressBar.Min = 0
	//progressBar.Max = 100
	//progressBar.Value = 0

	//content := container.NewVBox(icon, loading, progressBar)
	content := container.NewVBox(icon, loading)

	d := dialog.NewCustomWithoutButtons("请稍等", content, parent)
	//d.SetDismissText("") // 禁止关闭

	return &LoadingDialog{
		dialog: d,
		//progressBar: progressBar,
	}
}
