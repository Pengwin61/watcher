package core

import (
	"bufio"
	"log"
	"os"
	"strings"
	"watcher/internal/auth"
	"watcher/internal/utils"

	"github.com/spf13/viper"
)

type TempUser struct {
	Username    string
	ProjectName string
	ProjectID   int
}

const (
	PATH_ROOT_QUOTA string = "/"

	PATH_PROJECT    string = "./test_env/etc/projects"
	PATH_PROJECT_ID string = "./test_env/etc/projid"

	TMPL_PROJECT    string = "{{.ProjectName}}:/export/{{.ProjectName}}/{{.Username}}\n"
	TMPL_PROJECT_ID string = "{{.ProjectName}}_{{.Username}}:{{.ProjectName}}\n"
)

var data = make([]TempUser, 0)
var count = 0

var fsType = ""
var fileExists = false
var fileIsEmpty = true

func InitQuota(hard, soft string) error {
	var err error
	var list []string

	if fsType == "" && !fileExists {
		fsType = utils.CheckFS(PATH_ROOT_QUOTA)
		log.Printf("current filesystem type [%s] mount path: %s", fsType, PATH_ROOT_QUOTA)
		fsType = "xfs"

		list, fileExists, err = checkProjectFile()
		if err != nil {
			return err
		}
	}

	if fileExists && fsType == "xfs" {
		preparingQuotaFile()
		defaultQuotaXFS(list)
		for _, projname := range data {
			err := setQuotaXFS(hard, projname.ProjectName)
			if err != nil {
				return err
			}
		}
	}

	if fsType == "ext4" {
		log.Printf("current fsType: %s", fsType)
		setQuotaEXT4()
	}

	data = nil
	count = 0
	return nil
}

func preparingQuotaFile() {
	var projects = map[string]string{
		PATH_PROJECT:    TMPL_PROJECT,
		PATH_PROJECT_ID: TMPL_PROJECT_ID,
	}

	if !fileIsEmpty {
		log.Printf("files is created %s, %s\n", PATH_PROJECT, PATH_PROJECT_ID)
	} else {
		log.Printf("files is edit %s, %s\n", PATH_PROJECT, PATH_PROJECT_ID)
	}

	for path, tmpl := range projects {
		createQuotaFile(path, tmpl)
	}
}

func GenerationListQuota(userListID map[string]auth.Employee, group string) {

	for _, v := range userListID {

		tmp := &TempUser{
			Username:    v.Username,
			ProjectID:   count,
			ProjectName: group,
		}
		data = append(data, *tmp)
		count++
	}
}

func createQuotaFile(filePath string, tmp string) {

	projert, err := utils.CreateFile(filePath)
	if err != nil {
		log.Println("Ошибка при открытии файла:", err)
	}

	templa, err := utils.CreateTemlpate(tmp)
	if err != nil {
		log.Println("Ошибка при создании шаблона:", err)
	}

	for _, v := range data {
		// Записываем данные в файл с использованием шаблона
		err = templa.Execute(projert, v)
		if err != nil {
			log.Println("Ошибка при записи данных:", err)
			return
		}
	}

	defer projert.Close()

}

func defaultQuotaXFS(list []string) (ok bool, err error) {
	defQuota := "0G"

	for _, projname := range list {
		err = setQuotaXFS(defQuota, projname)
		if err != nil {
			return false, err
		}
	}
	// log.Fatalf("default_xfs unknown file system type [%s], or invalid mount path: %s", fsType, PATH_ROOT_QUOTA)
	return false, err

}

func setQuotaXFS(bhard, project string) error {
	// command := "xfs_quota"
	// arg0 := "-x"
	// arg1 := "-c"

	// cmd := exec.Command(command, arg0, arg1, fmt.Sprintf("limit -p bhard=%s %s", bhard, project), PATH_ROOT_QUOTA)
	// log.Printf("User %s has a quota applied %s \n", project, bhard)
	// err := cmd.Run()
	// if err != nil {
	// 	return err
	// }
	log.Printf("quota %s for project %s is set %s", fsType, project, bhard)
	return nil

}

func setQuotaEXT4() error {

	for _, users := range data {

		// _, err := exec.Command("setquota", "-u", users.Username, viper.GetString("userQuota.soft"), viper.GetString("userQuota.hard"), "0", "0", "/").Output()
		// if err != nil {
		// 	return err
		// }
		log.Printf("user quota ext4 for %s is set hard:%s, soft:%s", users.Username, viper.GetString("userQuota.hard"), viper.GetString("userQuota.soft"))
	}
	return nil
}

func checkProjectFile() (list []string, ok bool, err error) {

	_, err = os.Stat(PATH_PROJECT_ID)
	if err != nil {
		log.Printf("err: %s", err)
		return nil, false, nil
	}

	f, err := utils.ReadFile(PATH_PROJECT_ID)
	if err != nil {
		log.Printf("err: %s", err)
		return nil, false, err
	}

	defer f.Close()

	list = parseFile(f)
	if len(list) == 0 {
		fileIsEmpty = true
		return nil, false, err
	}
	return list, true, err
}

func parseFile(file *os.File) (list []string) {

	scanner := bufio.NewScanner(file)
	// Читаем файл построчно
	for scanner.Scan() {
		line := scanner.Text()
		arr := strings.Split(line, ":")

		list = append(list, arr[0])
	}

	defer file.Close()

	return list
}
