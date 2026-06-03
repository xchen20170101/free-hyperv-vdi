package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gin-vue/global"
	"gin-vue/modles/models"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/StackExchange/wmi"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gorm.io/gorm"
)

type GpuInfo struct {
	Prefix   string
	Hardware string
	Position string
	Identify string
	Desc     string
}

type PowerShell struct {
	powerShell string
}

func New() *PowerShell {
	ps, _ := exec.LookPath("powershell.exe")
	return &PowerShell{
		powerShell: ps,
	}
}

func (p *PowerShell) Execute(args ...string) (stdOut string, stdErr string, err error) {
	args = append([]string{"-NoProfile", "-NonInteractive"}, args...)
	cmd := exec.Command(p.powerShell, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	stdOut, stdErr = stdout.String(), stderr.String()
	return
}

func OpenVM(deviceName string, srcStatus string) string {
	posh := New()
	command := fmt.Sprintf("Start-VM -Name %s", deviceName)
	stdout, stderr, err := posh.Execute(command)
	if err != nil {
		global.Logger.Printf("open vm stderr:%+v, err:%+v\n", stderr, err)
		return srcStatus
	}
	global.Logger.Printf("open vm output:%+v\n", stdout)
	return "running"
}

func CloseVM(deviceName string, srcStatus string) string {
	posh := New()
	command := fmt.Sprintf("Stop-VM -Name %s -Force", deviceName)
	stdout, stderr, err := posh.Execute(command)
	if err != nil {
		global.Logger.Printf("close vm stderr:%+v, err:%+v\n", stderr, err)
		return srcStatus
	}
	global.Logger.Printf("close vm output:%+v\n", stdout)
	return "Off"

}

func ImportVM(deviceName string, srcVmName string, srcVmPath string, targetVmPath string, targetDiskPath string, vmSwitch string) error {
	posh := New()
	command := fmt.Sprintf("Import-VM -Path '%s' -Copy -GenerateNewId -VhdDestinationPath '%s' -VirtualMachinePath '%s'", srcVmPath, targetVmPath, targetDiskPath)
	global.Logger.Printf("import vm command:%s\n", command)
	_, stderr, err := posh.Execute(command)
	if err != nil {
		global.Logger.Printf("add vm stderr:%+v, err:%+v", stderr, err)
		return err
	}
	command = fmt.Sprintf("Rename-VM %s -NewName %s", srcVmName, deviceName)
	_, stderr, err = posh.Execute(command)
	if err != nil {
		global.Logger.Printf("set vm name stderr:%+v, err:%+v", stderr, err)
		return err
	}

	command = fmt.Sprintf("Connect-VMNetworkAdapter -VMName '%s' -Name '网络适配器' -SwitchName '%s'", deviceName, vmSwitch)
	_, stderr, err = posh.Execute(command)
	if err != nil {
		global.Logger.Printf("set vm network stderr:%+v, err:%+v", stderr, err)
		return err
	}
	return nil
}

func ModifyMemoryAndCPU(deviceName string, memory string, cpu string) error {
	posh := New()
	iMem, _ := strconv.Atoi(memory)
	command := fmt.Sprintf("Set-VM -VMName %s -MemoryStartupBytes %dMB", deviceName, iMem*1024)
	_, stderr, err := posh.Execute(command)
	if err != nil {
		global.Logger.Printf("set vm memory stderr:%+v, err:%+v", stderr, err)
		return err
	}
	iCpu, _ := strconv.Atoi(cpu)
	command = fmt.Sprintf("Set-VMProcessor -VMName %s -Count %d", deviceName, iCpu)
	_, stderr, err = posh.Execute(command)
	if err != nil {
		global.Logger.Printf("set vm cpu stderr:%+v, err:%+v\n", stderr, err)
		return err
	}
	return nil

}

func GetUserInfoByDevice(deviceId string) string {

	if deviceId == "" {
		return ""
	}
	var bind models.Bind
	global.DB.Where("device_id = ?", deviceId).First(&bind)
	if bind.UserId == "" {
		return ""
	}
	var user models.User
	global.DB.Where("id = ?", bind.UserId).First(&user)
	return user.Name

}

func GetTemplateNameById(templateId string) string {
	if templateId == "" {
		return ""
	}
	var template models.Template
	global.DB.Where("id =?", templateId).First(&template)
	if template.Name == "" {
		return ""
	}
	return template.Name
}

func GetTemplateIdByName(templateName string) string {
	if templateName == "" {
		return ""
	}
	var template models.Template
	global.DB.Where("name =?", templateName).First(&template)
	if template.ID == "" {
		return ""
	}
	return template.ID
}

func GetBindDeviceIds() []string {
	var deviceIds []string
	var binds []*models.Bind
	global.DB.Find(&binds)
	for _, bind := range binds {
		deviceIds = append(deviceIds, bind.DeviceId)
	}
	return deviceIds
}

func GetDeviceIdByName(deviceName string) string {
	if deviceName == "" {
		return ""
	}
	var device models.Device
	global.DB.Where("name = ?", deviceName).First(&device)
	return device.ID
}

func GetDeviceByName(deviceName string) *models.Device {
	if deviceName == "" {
		return nil
	}
	var device models.Device
	global.DB.Where("name =?", deviceName).First(&device)
	return &device
}


func TransferTimeStamp(timeStamp int64) string {
	return time.Unix(timeStamp/1000, 0).Format("2006-01-02 15:04:05")
}


func IsInArray(target string, srcArray []string) bool {
	for _, value := range srcArray {
		if target == value {
			return true
		}
	}
	return false
}


func GetCurrentPath() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}


func CreateVMLocalUser(userName string, pwd string, vmName string, adminUser string, adminPwd string) error {
	return nil
}


func DeleteVMLocalUser(userName string, vmName string, adminUser string, adminPwd string) error {
	return nil
}


func ResetVMUserPwd(userName string, newPassword string, vmName string, adminUser string, adminPwd string) error {
	psFile := fmt.Sprintf("%s\\script\\set_vm_user_pwd.ps1", GetCurrentPath())
	cmd := exec.Command("powershell", "-File", psFile, "-TargetUser", userName, "-NewPassword", newPassword, "-VMName", vmName, "-AdminUser", adminUser, "-AdminPassword", adminPwd)
	out, err := cmd.Output()
	if err != nil {
		global.Logger.Printf("reset user password failed:%s\n", err)
		return err
	}
	global.Logger.Printf("reset user password out:%+v\n", out)
	return nil
}


func ValidVMUserPwd(vmName string, username string, userpwd string) (bool, error) {
	psFile := fmt.Sprintf("%s\\script\\valid_vm_user_pwd.ps1", GetCurrentPath())
	cmd := exec.Command("powershell", "-File", psFile, "-VMName", vmName, "-TargetUser", username, "-TargetUserPassword", userpwd)

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			switch exitError.ExitCode() {
			case 0:

				return true, nil
			case 1:

				return false, nil
			case 2:

				global.Logger.Printf("valid vm user password system error: %s\n", err)
				return false, fmt.Errorf("系统错误：WinRM服务未启动或其他系统错误")
			case 4:

				global.Logger.Printf("valid vm user password vm error: %s\n", err)
				return false, fmt.Errorf("虚拟机不存在或未运行")
			default:
				global.Logger.Printf("valid vm user password unknown error: %s\n", err)
				return false, fmt.Errorf("未知错误")
			}
		}
		global.Logger.Printf("valid vm user password failed: %s\n", err)
		return false, err
	}


	return true, nil
}




























func AddGpuToVm(vmName string, instancePath string) error {
	return nil
}


func IsWindows11() bool {
	cmd := exec.Command("powershell", "Get-CimInstance Win32_OperatingSystem | Select-Object -ExpandProperty Caption")
	stdout, err := cmd.Output()
	if err != nil {
		global.Logger.Printf("is win11 result:%+v\n", err)
		return false
	}
	if strings.Contains(string(stdout), "Windows 11") {
		global.Logger.Println("the system is win11")
		return true
	}
	return false
}


func GetVMIp(vmName string) string {
	posh := New()
	command := fmt.Sprintf("Get-VM -Name %s | Get-VMNetworkAdapter | %% {Write-Host $_.ipaddresses[0]}", vmName)
	stdout, stderr, err := posh.Execute(command)
	if err != nil {
		global.Logger.Printf("get vm ip stderr:%+v, err:%+v\n", stderr, err)
		return ""
	}
	return strings.ReplaceAll(stdout, "\n", "")

}


func CreateDisk(name string, diskPath string, capacity int) error {
	posh := New()
	path := fmt.Sprintf("%s\\%s.vhdx", diskPath, name)
	capacityBytes := capacity * 1024 * 1024 * 1024
	command := fmt.Sprintf("New-VHD -Path %s -SizeBytes %d -Fixed", path, capacityBytes)
	global.Logger.Printf("create disk command:%+v\n", command)
	_, stderr, err := posh.Execute(command)
	if err != nil {
		global.Logger.Printf("create disk stderr:%+v, err:%+v\n", stderr, err)
		return err
	}
	return nil
}


func GetDeviceNameByDiskId(diskId string) string {
	var diskBind models.DiskBind
	global.DB.Where("disk_id = ?", diskId).First(&diskBind)
	if diskBind.DeviceId == "" {
		return ""
	}
	var device models.Device
	global.DB.Where("id = ?", diskBind.DeviceId).First(&device)
	return device.Name
}


func BindDiskToVM(vmName string, diskName string, diskPath string) error {
	posh := New()
	path := fmt.Sprintf("%s\\%s.vhdx", diskPath, diskName)
	command := fmt.Sprintf("Add-VMHardDiskDrive -VMName %s -Path %s", vmName, path)
	stdout, stderr, err := posh.Execute(command)
	if err != nil {
		global.Logger.Printf("bind disk to vm stderr:%+v, err:%+v\n", stderr, err)
		return err
	}
	global.Logger.Printf("bind disk stdout:%+v\n", stdout)
	return nil
}


func UnBindDiskFromVM(vmName string, diskName string, diskPath string) error {
	posh := New()
	path := fmt.Sprintf("%s\\%s.vhdx", diskPath, diskName)
	command := fmt.Sprintf("Get-VM -Name %s | Get-VMHardDiskDrive", vmName)
	stdout, stderr, err := posh.Execute(command)
	if err != nil {
		global.Logger.Printf("get disk bind relation stderr:%+v, err:%+v\n", stderr, err)
		return err
	}
	temps := strings.Split(stdout, "\n")
	for _, value := range temps[3:] {
		if strings.HasPrefix(value, vmName) {
			subTemps := strings.Fields(value)
			if subTemps[len(subTemps)-1] == path {

				command = fmt.Sprintf("Remove-VMHardDiskDrive -VMName %s -ControllerType %s -ControllerNumber %s -ControllerLocation %s", vmName, subTemps[1], subTemps[2], subTemps[3])
				_, stderr, err = posh.Execute(command)
				if err != nil {
					global.Logger.Printf("unbind disk from vm stderr:%+v, err:%+v\n", stderr, err)
					return err
				}
			}
		}
	}
	return nil
}


func IsVMCreated(vmName string) bool {
	posh := New()
	command := fmt.Sprintf("Get-VM -Name %s", vmName)
	stdout, _, _ := posh.Execute(command)
	return stdout != ""
}


func GetMemoryInfo(vmName string) string {
	posh := New()
	command := fmt.Sprintf("Get-VMMemory -VMName '%s' | Select-Object -ExpandProperty Startup", vmName)
	stdout, _, _ := posh.Execute(command)
	re := regexp.MustCompile("[0-9]+")
	numbers := re.FindAllString(stdout, -1)
	iMem, _ := strconv.ParseInt(numbers[0], 10, 64)
	result := int(iMem / (1024 * 1024 * 1024))
	return fmt.Sprintf("%dGB", result)
}


func SetMemoryInfo(vmName string, memoryInfo string) error {
	re := regexp.MustCompile("[0-9]+")
	numbers := re.FindAllString(memoryInfo, -1)
	iMem, _ := strconv.ParseInt(numbers[0], 10, 64)
	target := fmt.Sprintf("%dMB", iMem*1024)
	posh := New()
	command := fmt.Sprintf("Set-VM -VMName %s -MemoryStartupBytes %s", vmName, target)
	_, stderr, err := posh.Execute(command)
	if err != nil {
		global.Logger.Printf("set memory info stderr:%+v, err:%+v\n", stderr, err)
		return err
	}
	return nil
}


func GetCpuInfo(vmName string) string {
	posh := New()
	command := fmt.Sprintf("Get-VM -VMName '%s' | Select-Object -ExpandProperty ProcessorCount", vmName)
	stdout, _, _ := posh.Execute(command)
	return stdout
}


func SetCpuInfo(vmName string, cpuInfo string) error {
	posh := New()
	command := fmt.Sprintf("Set-VMProcessor -VMName %s -Count %s", vmName, cpuInfo)
	_, stderr, err := posh.Execute(command)
	if err != nil {
		global.Logger.Printf("set cpu info stderr:%+v, err:%+v\n", stderr, err)
		return err
	}
	return nil
}


func DeleteVMInForce(vmName string) error {

	posh := New()
	command := fmt.Sprintf("Remove-VM -Name %s -Force", vmName)
	global.Logger.Printf("DeleteVM command:%+v\n", command)
	_, stderr, err := posh.Execute(command)
	if err != nil {
		global.Logger.Printf("delete vm stderr:%+v, err:%+v\n", stderr, err)
	}
	vmPath := fmt.Sprintf("%s\\%s", global.Config.Vm.Path, vmName)
	err = os.RemoveAll(vmPath)
	if err != nil {
		global.Logger.Printf("delete vm path failed:%+v\n", err)
		return err
	}
	diskPath := fmt.Sprintf("%s\\%s_disk", global.Config.Vm.DiskPath, vmName)
	err = os.RemoveAll(diskPath)
	if err != nil {
		global.Logger.Printf("delete vm disk path failed:%+v\n", err)
		return err
	}
	return nil
}


func GetSwitchs() []string {
	var result []string
	posh := New()
	command := "Get-VMSwitch | Select-Object -Property Name"
	stdout, _, err := posh.Execute(command)
	if err != nil {
		global.Logger.Printf("get switchs failed:%+v\n", err)
		return result
	}
	temps := strings.Split(stdout, "\n")
	for _, value := range temps[3:] {
		strValue := strings.TrimSpace(value)
		if strValue == "Default Switch" || strValue == "" {
			continue
		}
		resTarget, err := DecodeGBK(strValue)
		if err != nil {
			continue
		}
		global.Logger.Printf("GetSwitchs value:%v\n", resTarget)
		result = append(result, resTarget)
	}
	return result
}


func InitGpuInfo(db *gorm.DB) {


	posh := New()
	command := "Get-VMPartitionableGpu | Select-Object -ExpandProperty Name"
	stdout, _, err := posh.Execute(command)
	if err != nil {
		global.Logger.Printf("init gpu info failed:%+v\n", err)
		return
	}
	instancePaths := strings.Split(stdout, "\n")
	for _, value := range instancePaths {
		if value == "" {
			continue
		}
		var gpu models.Gpu
		_ = db.Where("instance_path = ?", value).First(&gpu)
		if gpu.InstancePath != "" {

			continue
		}

		global.Logger.Println("value:", value)
		if !strings.HasPrefix(value, `\`) {
			continue
		}
		temps := strings.Split(value, "#")
		id := models.NewUUID()
		name := fmt.Sprintf("GPU_%s", temps[2])
		desc := ParseInstancePath(value)

		myGpu := &models.Gpu{
			ID:           id,
			Name:         name,
			InstancePath: value,
			Desc:         desc,
		}
		db.Create(&myGpu)
	}
}


func ParseInstancePath(instancePath string) string {
	if instancePath == "" {
		return ""
	}

	temps := strings.Split(instancePath, "#")
	subtemps := strings.Split(temps[3], `\`)
	gpuInfo := &GpuInfo{
		Prefix:   "PCI",
		Hardware: temps[1],
		Position: temps[2],
		Identify: subtemps[0],
		Desc:     subtemps[1],
	}
	jsonData, err := json.Marshal(gpuInfo)
	if err != nil {
		return ""
	}
	return string(jsonData)
}


func GetInstancePathFromDB(desc string) string {
	byteData := []byte(desc)
	var gpuInfo GpuInfo
	err := json.Unmarshal(byteData, &gpuInfo)
	if err != nil {
		return ""
	}




	target := fmt.Sprintf(`\\?\%v#%v#%v#%v\%v`, gpuInfo.Prefix, gpuInfo.Hardware, gpuInfo.Position, gpuInfo.Identify, "GPUPARAV")

	return target
}


func GetAdapterIdByInstancePath(vmName string, instancePath string) string {





















	idList := GetGpuAdpterIdList(vmName)
	for _, value := range idList {

		instPath := GetInstancePathByVMAndId(vmName, value)
		if strings.Contains(instancePath, instPath) {
			return value
		}
	}
	return ""
}


func GetGpuAdpterIdList(vmName string) []string {
	var result []string
	posh := New()
	command := fmt.Sprintf("Get-VMGpuPartitionAdapter -VMName '%s' | Select-Object -ExpandProperty Id", vmName)
	stdout, _, err := posh.Execute(command)
	if err != nil {
		return result
	}
	temps := strings.Split(stdout, "\n")
	return temps
}


func GetInstancePathByVMAndId(vmName string, adapterId string) string {
	posh := New()
	command := fmt.Sprintf("Get-VMGpuPartitionAdapter -VMName '%s' -AdapterId '%v' | Select-Object -ExpandProperty InstancePath", vmName, adapterId)
	stdout, _, err := posh.Execute(command)
	if err != nil {
		return ""
	}
	return stdout
}


func RemoveGpuAdpater(vmName string, instancePath string) error {
	adapterId := GetAdapterIdByInstancePath(vmName, instancePath)
	global.Logger.Println("adapterId:", adapterId)
	posh := New()
	command := fmt.Sprintf("Remove-VMGpuPartitionAdapter -VMName '%v' -AdapterId ", vmName)
	target := command + adapterId
	_, stderr, err := posh.Execute(target)
	if err != nil {
		global.Logger.Printf("RemoveGpuAdpater failed, stderr:%+v, err:%+v\n", stderr, err)
		return err
	}
	return nil
}


func DecodeGBK(input string) (string, error) {
	decoder := simplifiedchinese.GBK.NewDecoder()
	inputReader := strings.NewReader(input)
	transformReader := transform.NewReader(inputReader, decoder)

	out, err := ioutil.ReadAll(transformReader)
	if err != nil {
		return "", err
	}

	return string(out), nil
}


func getCPUSerialNumber() (string, error) {
	cmd := exec.Command("wmic", "cpu", "get", "ProcessorId")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}


	cpuInfo := strings.Split(string(output), "\n")
	if len(cpuInfo) >= 2 {
		return strings.TrimSpace(cpuInfo[1]), nil
	}

	return "", nil
}

type Win32_BaseBoard struct {
	SerialNumber string
}


func getBoardSerialNumber() (string, error) {
	var baseBoard []Win32_BaseBoard
	query := "select SerialNumber from Win32_BaseBoard"
	if err := wmi.Query(query, &baseBoard); err != nil {
		global.Logger.Printf("WMI query failed:%+v\n", err)
		return "", err
	}
	serialNumber := ""
	for _, board := range baseBoard {
		serialNumber += board.SerialNumber
	}
	return serialNumber, nil
}


func GetMachineCode() string {

	serialNum, err := getCPUSerialNumber()
	if err != nil {
		global.Logger.Printf("get cpu serial number failed:%+v\n", err)
		return ""
	}
	boardNum, err := getBoardSerialNumber()
	if err != nil {
		global.Logger.Printf("get board serial number failed:%+v\n", err)
	}
	target := serialNum + boardNum
	md5Value := calculateMD5(target)
	return md5Value
}


func CheckLicense(machineCode string, licenseCode string) (bool, string) {
	if licenseCode == "" {
		return false, ""
	}
	value, err := decrypt(licenseCode)
	if err != nil {
		global.Logger.Printf("decrypt failed:%+v\n", err)
		return false, ""
	}
	temps := strings.Split(value, "&")
	if len(temps) != 2 {
		return false, ""
	}
	if temps[0] != machineCode {
		global.Logger.Println("license check failed")
		return false, ""
	}
	global.Logger.Println("license check pass")
	return true, temps[1]
}

func CheckSystemTimeModify() bool {

	nowDay := time.Now().Format("2006-01-02")
	var licenseDate models.LicenseDate
	global.DB.Where("day_date = ?", nowDay).First(&licenseDate)
	if licenseDate.DayDate != "" {
		global.Logger.Println("system time is modify, license invalid")
		return false
	}
	return true
}


func InsertOrUpdateTemplate(name string, desc string) error {
	var template models.Template
	global.DB.Where("name =?", name).First(&template)
	if template.Name == "" {
		template.ID = models.NewUUID()
		template.Name = name
		template.Desc = desc
		template.UserName = "vmadmin"
		template.UserPwd = "vmadmin"
		global.DB.Create(&template)
	} else {
		template.Desc = desc
		global.DB.Save(&template)
	}
	return nil
}


func DeleteTemplateFile(name string) {
	dir := global.Config.Vm.Template + "\\" + name
	err := os.RemoveAll(dir)
	if err != nil {
		global.Logger.Println("Error removing template:", err)
	}
	global.Config.Vm.DeleteTemplateMap(name)
}


func GetVMState(vmName string) string {
	posh := New()
	command := fmt.Sprintf("Get-VM -Name '%s' | Select-Object -ExpandProperty State", vmName)
	stdout, stderr, err := posh.Execute(command)
	if err != nil {
		global.Logger.Printf("get vm state stderr:%+v, err:%+v\n", stderr, err)
		return "Unknown"
	}
	return strings.TrimSpace(stdout)
}


func GetAllVMStates(vmNames []string) map[string]string {
	if len(vmNames) == 0 {
		return make(map[string]string)
	}
	

	vmNamesStr := "'" + strings.Join(vmNames, "','") + "'"
	posh := New()
	command := fmt.Sprintf("Get-VM -Name %s | Select-Object Name,State | Format-Table -HideTableHeaders", vmNamesStr)
	stdout, stderr, err := posh.Execute(command)
	
	result := make(map[string]string)
	
	if err != nil {
		global.Logger.Printf("get all vm states stderr:%+v, err:%+v\n", stderr, err)

		for _, vmName := range vmNames {
			result[vmName] = GetVMState(vmName)
		}
		return result
	}
	

	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 {
			vmName := fields[0]
			state := fields[1]
			result[vmName] = state
		}
	}
	

	for _, vmName := range vmNames {
		if _, exists := result[vmName]; !exists {
			result[vmName] = "Unknown"
		}
	}
	
	return result
}


func GetAllVMIps(vmNames []string) map[string]string {
	if len(vmNames) == 0 {
		return make(map[string]string)
	}
	
	result := make(map[string]string)
	

	type vmIPResult struct {
		vmName string
		ip     string
	}
	
	resultChan := make(chan vmIPResult, len(vmNames))
	

	for _, vmName := range vmNames {
		go func(name string) {
			ip := GetVMIp(name)
			resultChan <- vmIPResult{vmName: name, ip: ip}
		}(vmName)
	}
	

	for i := 0; i < len(vmNames); i++ {
		res := <-resultChan
		result[res.vmName] = res.ip
	}
	
	return result
}



func ValidateVMName(name string) bool {
	if name == "" {
		return false
	}
	

	if len(name) > 64 {
		return false
	}
	


	matched, _ := regexp.MatchString("^[a-zA-Z0-9][a-zA-Z0-9_]{0,63}$", name)
	return matched
}
