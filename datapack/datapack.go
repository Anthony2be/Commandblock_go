package datapack

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/tebeka/atexit"
)

type datapack_struct struct {
	datapack_name        string
	namespace_id         string
	pack_version         int
	load_json            string
	tick_json            string
	datapack_description string
}

func New(Datapack_name string, Namespace_id string, Pack_version int, Load_json string, Tick_json string, Datapack_description string) datapack_struct {
	path, err := os.Getwd()
	atexit.Register(stopAbrupt)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(path)

	d := datapack_struct{Datapack_name, Namespace_id, Pack_version, Load_json, Tick_json, Datapack_description}

	if len([]rune(Namespace_id)) > 5 {
		warning := fmt.Sprintf("WARNING : The name space id given (%s) is longer than 5 characters\n          Datapacks are *usually* made with 2-4 namespace characters\n          Nothing will be changed in the datapack", Namespace_id)
		color.Yellow(warning)
	}

	if Pack_version < 4 {
		color.Yellow(fmt.Sprintf("WARNING : The pack version given (\"%d\") is less than 4. This may cause errors in the datapack", Pack_version))
	}

	err = os.Mkdir(path+"/.temp", 0755)
	if err != nil {
		//os.RemoveAll(path + "/.temp")
		//fmt.Println(err)
		color.Red(fmt.Sprintf("ERROR : Oh-oh! An Exception occured while generating a temporary folder at dir \"%s\": %d", path, err))
		atexit.Exit(-1)
	}

	return d
}

func stopAbrupt() {
	path, e := os.Getwd()
	if e != nil {
		fmt.Println(e)
	}

	folderInfo, err := os.Stat(path + "/.temp")
	if !os.IsNotExist(err) {
		_ = folderInfo
		color.Red("Detected a temp folder, and erasing it")
		os.RemoveAll(path + "/.temp")
	}
}

//func RegisterFunction()

func Abort(rem_datapack bool){
	color.Red("\n\n##############################\nABORTING THE DATAPACK GENERATION\n##############################\n\n\n")
}