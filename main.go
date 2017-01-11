package main

import (
	"fenqiwanh5/models"
	"fenqiwanh5/multi"
	//"fmt"
	_ "wechatmenu/routers"

	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/chanxuehong/wechat/mp/menu"
)

//var redisx *lib.RedisxCache

func main() {

	//initCache()
	orm.Debug, _ = beego.AppConfig.Bool("orm_debug")
	beego.SetLogger("file", `{"filename":"logs/wechatmenu.log"}`)
	beego.SetLogFuncCall(true)

	logLevel, _ := beego.AppConfig.Int("log_level")
	beego.SetLevel(logLevel)

	if logLevel < 7 {
		beego.BeeLogger.DelLogger("console")
	}

	models.Connect()

	SetupWechatMenus()

	beego.Run()
}

//func InitMenu() {

//	var mn menu.Menu
//	mn.Buttons = make([]menu.Button, 2)

//	subButtons := make([]menu.Button, 1)

//	siteUrl := base.SiteUrl

//	beego.Info(siteUrl)

//	subButtons[0].SetAsViewButton("5126充值", siteUrl+"/merchant/entry?mk=aafabf9052bec2464daab9078b5514ce01e61")

//	mn.Buttons[0].SetAsViewButton("查信用", siteUrl+"/merchant/jdb")

//	mn.Buttons[1].SetAsSubMenuButton("信用支付", subButtons)

//	//mn.Buttons[1].SetAsScanCodePushButton("扫订单", "0e772eb5b0")

//	menuClient := (*menu.Client)(base.MpClient)
//	if err := menuClient.CreateMenu(mn); err != nil {
//		beego.Error(err)
//		return
//	}

//	beego.Info("======Menu Created=====")
//}

func SetupWechatMenus() {

	beego.Debug("enter setup menu")

	rootMenus, err := models.GetWechatMenus()

	if err != nil {
		panic(err.Error())
	}

	err = multi.InitMultipleWechat("")

	if err == nil {

		//		beego.Debug(fmt.Sprintf("%+v", rootMenus))
		//				beego.Debug(fmt.Sprintf("%+v", rootMenus[0].SubMenus))
		//				beego.Debug(fmt.Sprintf("%+v", rootMenus[0].SubMenus[0].MenuName))
		//wechat level
		for _, rootMenu := range rootMenus {
			var mn menu.Menu
			mn.Buttons = make([]menu.Button, rootMenu.MenuSize)
			siteUrl := rootMenu.SiteUrl
			beego.Debug(rootMenu.MenuSize)

			if rootMenu.MenuSize > 0 {

				//get wechat client
				instanceServer, ok := multi.InstanceMap[rootMenu.InstanceName]

				if !ok {
					panic("wechat client not fountd for:" + rootMenu.InstanceName)
				}

				//root menu
				for i, menuItem := range rootMenu.SubMenus {
					beego.Debug(menuItem.MenuName)

					if strings.Compare("submenu", menuItem.MenuType) == 0 {
						subButtons := make([]menu.Button, menuItem.MenuSize)

						beego.Debug(menuItem.MenuSize)

						for j, sumMenuItem := range menuItem.SubMenus {

							if strings.Compare("view", sumMenuItem.MenuType) == 0 {

								subButtons[j].SetAsViewButton(sumMenuItem.MenuName, siteUrl+sumMenuItem.MenuUrl)
							}
						}

						mn.Buttons[i].SetAsSubMenuButton(menuItem.MenuName, subButtons)

					} else if strings.Compare("view", menuItem.MenuType) == 0 {
						mn.Buttons[i].SetAsViewButton(menuItem.MenuName, siteUrl+menuItem.MenuUrl)

					}
				}

				beego.Debug(mn)

				menuClient := (*menu.Client)(instanceServer.MpClient)
				if err := menuClient.CreateMenu(mn); err != nil {
					beego.Error(err)
				} else {
					beego.Info("Menu Created for:" + rootMenu.InstanceName)
				}
			}

		}

	} else {
		panic(err.Error())
	}

}
