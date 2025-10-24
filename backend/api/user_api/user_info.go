package user_api

import (
	"gin-vue/api/utils"
	"gin-vue/global"
	"gin-vue/modles/form"
	"gin-vue/modles/models"
	"gin-vue/modles/res"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (UsersApi) UserCount(c *gin.Context) {
	var user []models.User
	rows := global.DB.Find(&user).RowsAffected
	res.OkWithData(rows, c)
}


func (UsersApi) PasswordModify(c *gin.Context) {
	var myUser models.User

	result := global.DB.First(&myUser, "name", "admin")
	if result.Error != nil {
		res.FailWithMsg("failed to find user", c)
		return
	}

	userName := "admin"
	password := c.PostForm("password")
	userRole := "1"
	myNewUser := models.User{
		Name:     userName,
		Password: password,
		Role:     userRole,
	}

	global.DB.Model(myUser).Updates(myNewUser)


	res.OkWithData(myNewUser, c)
}


func (UsersApi) AndroidResetPassword(c *gin.Context) {
	var user models.User
	userName := c.PostForm("username")
	global.DB.Where("name =?", userName).First(&user)
	if user.Name == "" {
		res.FailWithMsg("User.NotExist", c)
		return
	}
	if user.Status != "启用" {
		res.FailWithMsg("User.Disable", c)
		return
	}
	password := c.PostForm("oldpassword")
	if user.Password != password {
		res.FailWithMsg("User.PasswordIsWrong", c)
		return
	}
	newPassword := c.PostForm("newpassword")
	myNewUser := models.User{
		Password: newPassword,
	}
	global.DB.Model(user).Updates(myNewUser)
	res.OkWithData(userName, c)
}

func (UsersApi) UserAdd(c *gin.Context) {
	var myUser models.User
	var form form.AddUserForm
	if err := c.ShouldBind(&form); err != nil {
		global.Logger.Printf("add user failed:%v\n", err.Error())
		res.FailWithMsg("Common.InvalidParam", c)
		return
	}
	userName := form.Name
	password := form.Password
	userRole := form.Role
	status := form.Status

	global.DB.Where("name = ?", userName).First(&myUser)

	if myUser.Name != "" {
		res.FailWithMsg("User.Exists", c)
		return
	}
	myUser.ID = models.NewUUID()
	myUser.Name = userName
	myUser.Password = password
	myUser.Role = userRole
	myUser.Status = status
	global.DB.Create(&myUser)
	global.Logger.Println("user:", myUser)

	res.OkWithData(myUser, c)
}

func (UsersApi) UserUpdate(c *gin.Context) {
	var myUser models.User
	id := c.Param("id")
	result := global.DB.First(&myUser, "id", id)
	if result.Error != nil {
		res.FailWithMsg("User.NotExist", c)
		return
	}
	var form form.UpdateUserForm

	if err := c.ShouldBind(&form); err != nil {
		global.Logger.Printf("update user failed:%v\n", err.Error())
		res.FailWithMsg("Common.InvalidParam", c)
		return
	}

	password := form.Password
	userRole := form.Role
	status := form.Status


	if myUser.Name == "vmadmin" && userRole == "云桌面用户" {
		res.FailWithMsg("Common.InvalidParam", c)
		return
	}


	if myUser.Name == "vmadmin" && status == "禁用" {
		res.FailWithMsg("Common.InvalidParam", c)
		return
	}



	myNewUser := models.User{
		Name:     myUser.Name,
		Password: password,
		Role:     userRole,
		Status:   status,
	}

	global.DB.Model(myUser).Updates(myNewUser)
	global.Logger.Println("user:", myNewUser)
	res.OkWithData(myNewUser, c)
}

func (UsersApi) UserUpdateSelfPassword(c *gin.Context) {
	var myUser models.User
	// 从上下文中获取userId（由中间件设置）
	userId, exists := c.Get("userId")
	if !exists {
		res.FailWithMsg("Common.InvalidParam", c)
		return
	}
	
	result := global.DB.First(&myUser, "id", userId)
	if result.Error != nil {
		res.FailWithMsg("User.NotExist", c)
		return
	}

	password := c.PostForm("password")

	myNewUser := models.User{
		Name:     myUser.Name,
		Password: password,
	}

	global.DB.Model(myUser).Updates(myNewUser)
	res.OkWithData(myNewUser, c)
}

func (UsersApi) UserGet(c *gin.Context) {
	var users []models.User
	var totalNum int64
	count, _ := strconv.Atoi(c.Query("count"))
	pageNum, _ := strconv.Atoi(c.Query("index"))
	name := c.Query("name")
	keyword := "%" + name + "%"
	offset := (pageNum - 1) * count
	if name == "" {
		global.DB.Model(&models.User{}).Count(&totalNum)
		global.DB.Limit(count).Offset(offset).Find(&users)
	} else {
		global.DB.Model(&models.User{}).Where("name LIKE ?", keyword).Count(&totalNum)
		global.DB.Limit(count).Offset(offset).Where("name LIKE ?", keyword).Find(&users)
	}
	data := make(map[string]interface{})
	data["users"] = users
	data["totalNum"] = totalNum
	res.OkWithData(data, c)
}


func (UsersApi) UserAllCountGet(c *gin.Context) {
	var users []models.User
	global.DB.Find(&users)
	data := make(map[string]interface{})
	data["num"] = len(users)
	res.OkWithData(data, c)
}


func (UsersApi) UserProfileGet(c *gin.Context) {
	// 从上下文中获取userId（由中间件设置）
	userId, exists := c.Get("userId")
	if !exists {
		res.FailWithMsg("Common.InvalidParam", c)
		return
	}
	
	var user models.User
	result := global.DB.Where("id = ?", userId).First(&user)
	if result.Error != nil {
		res.FailWithMsg("User.NotExist", c)
		return
	}
	
	data := &res.GetUserProfileReponse{
		UserName: user.Name,
		Password: user.Password,
	}
	res.OkWithData(data, c)
}


func (UsersApi) UserLicenseGet(c *gin.Context) {
	var license models.License
	var isChecked bool
	global.DB.First(&license)
	if license.LicenseCode == "" {
		isChecked = false
	}
	machineCode := utils.GetMachineCode()
	isChecked, _ = utils.CheckLicense(machineCode, license.LicenseCode)

	now := time.Now().UnixNano() / int64(time.Millisecond)


	if license.ExpiredTime != -1 && license.ExpiredTime != 0 && license.ExpireFlag != 2 {
		timeInterval := now - license.ExpiredTime

		if timeInterval > 0 {

			myLicense := models.License{
				ExpireFlag: 2,
			}
			global.DB.Model(license).Updates(myLicense)
			res.FailWithMsg("User.LicenseExpired", c)
			return
		}
	}

	if license.ExpireFlag == 2 {
		isChecked = false
	}

	data := &res.GetUserLicenseResponse{
		MachineCode: machineCode,
		IsChecked:   isChecked,
	}
	res.OkWithData(data, c)
}


func (UsersApi) LicenseActive(c *gin.Context) {
	var license models.License
	var form form.ActiveLicenseForm
	if err := c.ShouldBind(&form); err != nil {
		global.Logger.Printf("license active failed:%v\n", err.Error())
		res.FailWithMsg("Common.InvalidParam", c)
		return
	}
	if form.MachineCode != utils.GetMachineCode() {
		res.FailWithMsg("Common.InvalidParam", c)
		return
	}
	global.DB.Where("license_code = ?", form.LicenseCode).First(&license)
	if license.LicenseCode != "" {
		res.FailWithMsg("User.LicenseExists", c)
		return
	}
	isChecked, licenseType := utils.CheckLicense(utils.GetMachineCode(), form.LicenseCode)
	if !isChecked {
		res.FailWithMsg("User.LicenseActiveFailed", c)
		return
	}
	now := time.Now().UnixNano() / int64(time.Millisecond)
	license.ID = models.NewUUID()
	license.LicenseCode = form.LicenseCode
	license.LicenseType = licenseType
	license.ActiveTime = now
	var expireFlag int64
	var expiredTime int64
	if licenseType == "1" {
		expiredTime = now + 30*24*3600*1000
		expireFlag = 1
	} else {
		expiredTime = -1
		expireFlag = 0
	}
	license.ExpiredTime = expiredTime
	license.ExpireFlag = expireFlag
	global.DB.Create(&license)
	res.OkWithData(license, c)
}


func (UsersApi) GetLicenses(c *gin.Context) {
	var licenses []models.License
	var myLicenses []*res.GetLicensesResponse
	global.DB.Where("expire_flag < ?", "2").Find(&licenses)
	data := make(map[string]interface{})
	for _, value := range licenses {
		var expireTime string
		if value.ExpiredTime == -1 {
			expireTime = "-"
		} else {
			expireTime = utils.TransferTimeStamp(value.ExpiredTime)
		}

		var licenseTypeTarget string
		if value.LicenseType == "1" {
			licenseTypeTarget = "测试授权"
		} else {
			licenseTypeTarget = "正式授权"
		}
		temp := &res.GetLicensesResponse{
			MachineCode: utils.GetMachineCode(),
			LicenseCode: value.LicenseCode,
			LicenseType: licenseTypeTarget,
			ExpireTime:  expireTime,
		}
		myLicenses = append(myLicenses, temp)
	}
	data["licenses"] = myLicenses
	res.OkWithData(data, c)
}


func (UsersApi) UserBindDevices(c *gin.Context) {


	var user models.User
	var form form.BindUserForm
	if err := c.ShouldBind(&form); err != nil {
		global.Logger.Printf("user bind device failed:%v\n", err.Error())
		res.FailWithMsg("Common.InvalidParam", c)
		return
	}
	userId := form.UserId
	deviceName := form.StrDeviceName
	deviceNames := strings.Split(deviceName, ",")
	global.DB.Where("id = ?", userId).First(&user)
	if user.Name == "" {
		res.FailWithMsg("User.NotExist", c)
		return
	}

	for _, value := range deviceNames {
		device := utils.GetDeviceByName(value)
		if device == nil {
			continue
		}
		var template models.Template
		global.DB.Where("id =?", device.TemplateId).First(&template)
		err := utils.CreateVMLocalUser(user.Name, user.Password, value, template.UserName, template.UserPwd)
		if err != nil {
			continue
		}
		myBind := &models.Bind{
			ID:       models.NewUUID(),
			DeviceId: device.ID,
			UserId:   userId,
		}
		global.DB.Create(&myBind)
	}


	res.OkWithData(nil, c)
}

func (UsersApi) UserDel(c *gin.Context) {
	var myUser models.User
	var myBinds []*models.Bind
	id := c.Param("id")
	result := global.DB.First(&myUser, "id", id)
	if result.Error != nil {
		res.FailWithMsg("User.NotExist", c)
		return
	}

	global.DB.Find(&myBinds, "user_id", id)
	if len(myBinds) > 0 {
		res.FailWithMsg("User.HasBind", c)
		return
	}

	result = global.DB.Delete(&myUser)
	global.Logger.Println("userDel error:", result.Error)
	if result.Error != nil {
		res.FailWithMsg("User.NotExist", c)
		return
	}

	global.DB.Where("user_id = ?", myUser.ID).Delete(&models.Bind{})

	res.OkWithData(myUser, c)
}


func (UsersApi) UserLogin(c *gin.Context) {



	global.Logger.Println("UserAgent:", c.Request.UserAgent())

	var user models.User
	userName := c.PostForm("username")
	pwd := c.PostForm("password")
	global.DB.Where("name = ?", userName).First(&user)


	requestFrom := c.Request.UserAgent()
	if strings.Contains(requestFrom, "Windows") && user.Role == "云桌面用户" {
		global.Logger.Println("User.NoPermission")
		res.FailWithMsg("User.NoPermission", c)
		return
	}

	if user.Password != pwd {
		global.Logger.Println("User.PasswordIsWrong")
		res.FailWithMsg("User.PasswordIsWrong", c)
		return
	}

	if user.Status != "启用" {
		global.Logger.Println("User.Disable")
		res.FailWithMsg("User.Disable", c)
		return
	}

	// 生成新的token（30分钟有效期）
	tokenValue := models.NewUUID()
	now := time.Now().Unix()
	cryptoperiod := int64(1800) // 30分钟 = 1800秒

	// 查询该用户是否已有token
	var existingToken models.Token
	result := global.DB.Where("user_id = ?", user.ID).First(&existingToken)

	var responseToken models.Token
	
	if result.Error != nil {
		// 不存在，创建新token
		newToken := models.Token{
			ID:           models.NewUUID(),
			Value:        tokenValue,
			Cryptoperiod: cryptoperiod,
			CreatedTime:  now,
			UserId:       user.ID,
		}
		global.DB.Create(&newToken)
		responseToken = newToken
	} else {
		// 已存在，更新token（这会使之前的token失效，实现单点登录）
		existingToken.Value = tokenValue
		existingToken.CreatedTime = now
		existingToken.Cryptoperiod = cryptoperiod
		global.DB.Save(&existingToken)
		responseToken = existingToken
	}
	
	// 设置Cookie（关键：让浏览器自动携带token）
	// 参数：name, value, maxAge(秒), path, domain, secure, httpOnly
	c.SetCookie("accessToken", tokenValue, int(cryptoperiod), "/", "", false, true)
	c.SetCookie("userId", user.ID, int(cryptoperiod), "/", "", false, true)
	
	global.Logger.Printf("用户 %s 登录成功，Token: %s", user.Name, tokenValue)
	
	res.Ok(responseToken, "登录成功", c)
}

func (UsersApi) UserLogout(c *gin.Context) {
	// 从上下文中获取userId（由中间件设置）
	userId, exists := c.Get("userId")
	if !exists {
		// 如果中间件没有设置，尝试从Cookie获取
		userIdStr, err := c.Cookie("userId")
		if err != nil {
			res.FailWithMsg("User.LogoutFailed", c)
			return
		}
		userId = userIdStr
	}
	
	// 删除数据库中的token
	global.DB.Where("user_id = ?", userId).Delete(&models.Token{})
	
	// 清除Cookie
	c.SetCookie("accessToken", "", -1, "/", "", false, true)
	c.SetCookie("userId", "", -1, "/", "", false, true)
	
	global.Logger.Printf("用户 %s 退出成功", userId)
	
	res.Ok(userId, "退出成功", c)
}
