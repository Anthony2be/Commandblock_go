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

var function_lock = 0
var datapack_functions []string

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
		//retry 
		e := os.RemoveAll(path + "/.temp")
		if e != nil {
			color.Red("Could not remove the temp folder")
			atexit.Exit(-1)
		}
		err2 := os.Mkdir(path+"/.temp", 0755)
		if err2 != nil {
			color.Red(fmt.Sprintf("ERROR : Oh-oh! An Exception occured while generating a temporary folder at dir \"%s\": %d", path, err))
			atexit.Exit(-1)
		}
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

func (d datapack_struct) Abort (rem_datapack bool) {
	path, _ := os.Getwd()
	color.Red("\n\n##############################\nABORTING THE DATAPACK GENERATION\n##############################\n\n\n")
	_, err := os.Stat(fmt.Sprintf(path+"/.temp/%s", d.datapack_name))
	if rem_datapack && !os.IsNotExist(err) {
		os.RemoveAll(fmt.Sprintf(path+"/.temp/%s", d.datapack_name))
		color.Red("Removed the datapack folder")
	}

	_, err2 := os.Stat(path+ "/.temp")
	if rem_datapack && !os.IsNotExist(err2) {
		os.RemoveAll(path+"/.temp")
		color.Red("Removed temp folder")
	}

	color.Red("Finished removing datapack related files. Shutting down system")
}

func (d datapack_struct) RegisterFunction (name string, content string){
	if (name == d.load_json || name == d.tick_json){
		function_lock ++
	}

	
}
