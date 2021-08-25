package datapack

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/otiai10/copy"
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

//makes a new datapack
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

//Aborts the datapack generation
func (d datapack_struct) Abort(rem_datapack bool) {
	path, _ := os.Getwd()
	color.Red("\n\n##############################\nABORTING THE DATAPACK GENERATION\n##############################\n\n\n")
	_, err := os.Stat(fmt.Sprintf(path+"/.temp/%s", d.datapack_name))
	if rem_datapack && !os.IsNotExist(err) {
		os.RemoveAll(fmt.Sprintf(path+"/.temp/%s", d.datapack_name))
		color.Red("Removed the datapack folder")
	}

	_, err2 := os.Stat(path + "/.temp")
	if rem_datapack && !os.IsNotExist(err2) {
		os.RemoveAll(path + "/.temp")
		color.Red("Removed temp folder")
	}

	color.Red("Finished removing datapack related files. Shutting down system")
}

//Registers a function
func (d datapack_struct) RegisterFunction(name string, content string) {
	path, _ := os.Getwd()
	if name == d.load_json || name == d.tick_json {
		function_lock++
	}

	datapack_functions = append(datapack_functions, name)

	err := os.MkdirAll(fmt.Sprintf("%s/.temp/%s/%s/functions/", path, d.datapack_name, d.namespace_id /*name*/), 0755)
	if err != nil {
		//fmt.Println(err)
		color.Red(fmt.Sprintf("ERROR : %s", err.Error()))
		atexit.Exit(-1)
	}

	f, err := os.Create(fmt.Sprintf("%s/.temp/%s/%s/functions/%s.mcfunction", path, d.datapack_name, d.namespace_id, name))
	if err != nil {
		color.Red(fmt.Sprintf("ERROR : %s", err.Error()))
		atexit.Exit(-1)
	}
	_, err2 := f.WriteString(content)
	if err2 != nil {
		color.Red(fmt.Sprintf("ERROR : %s", err.Error()))
		f.Close()
		atexit.Exit(-1)
	}
	err = f.Close()
	if err != nil {
		color.Red(fmt.Sprintf("ERROR : %s", err.Error()))
		atexit.Exit(-1)
	}
}

func (d datapack_struct) json_value(load bool) string {
	value := "{\"values\":[\"%s:%s\"]}"

	if load {
		json := fmt.Sprintf(value, d.namespace_id, d.load_json)
		return json
	} else {
		json := fmt.Sprintf(value, d.namespace_id, d.tick_json)
		return json
	}
}

//Generates the datapack
func (d datapack_struct) Generate( /*zip bool*/ ) {
	path, _ := os.Getwd()

	color.Cyan("Commandblock_go is still in its Beta stages. If You find any bugs, please send a screenshot of the terminal to \"Terroid#0490\" or \"Anthony2be#1900\" on Discord")
	if function_lock < 2 {
		color.Red(fmt.Sprintf("ERROR : Your Datapack doesn't contain a \"%s\" or a \"%s\" function or both", d.load_json, d.tick_json))
		atexit.Exit(-1)
	}

	err := os.Mkdir(fmt.Sprintf("%s/%s", path, d.datapack_name), 0755)
	if err != nil {
		color.Yellow("WARNING : A folder with this name already exists. This action is going to replace the old datapack folder")

		err := os.RemoveAll(fmt.Sprintf("%s/%s", path, d.datapack_name))
		if err != nil {
			color.Red(fmt.Sprintf("ERROR : Oh-oh! An Exception occured while deleting a previous file during generation : %s", err.Error()))
			d.Abort(true)
		}

		err2 := os.Mkdir(fmt.Sprintf("%s/%s", path, d.datapack_name), 0755)
		if err2 != nil {
			color.Red(fmt.Sprintf("ERROR : Oh-oh! An Exception occured while making the datapack folder: %s", err2.Error()))
			d.Abort(false)
		}
	}

	err2 := os.MkdirAll(fmt.Sprintf("%s/%s/data/minecraft/tags/functions", path, d.datapack_name), 0755)
	if err2 != nil {
		color.Red("ERROR : Oh-oh! An Exception occured while generating the template of the datapack: " + err.Error())
		d.Abort(true)
	}

	err3 := os.MkdirAll(fmt.Sprintf("%s/%s/data/%s/functions", path, d.datapack_name, d.namespace_id), 0755)
	if err3 != nil {
		color.Red("ERROR : Oh-oh! An Exception occured while generating the template of the datapack: " + err.Error())
		d.Abort(true)
	}

	color.Green("Generating function files...")

	//load.json
	f, err := os.Create(fmt.Sprintf("%s/%s/data/minecraft/tags/functions/load.json", path, d.datapack_name))
	if err != nil {
		color.Red(fmt.Sprintf("ERROR : %s", err.Error()))
		atexit.Exit(-1)
	}
	_, err = f.WriteString(d.json_value(true))
	if err != nil {
		color.Red(fmt.Sprintf("ERROR : %s", err.Error()))
		f.Close()
		atexit.Exit(-1)
	}
	err = f.Close()
	if err != nil {
		color.Red(fmt.Sprintf("ERROR : %s", err.Error()))
		atexit.Exit(-1)
	}

	//tick.json
	f, err = os.Create(fmt.Sprintf("%s/%s/data/minecraft/tags/functions/tick.json", path, d.datapack_name))
	if err != nil {
		color.Red(fmt.Sprintf("ERROR : %s", err.Error()))
		atexit.Exit(-1)
	}
	_, err = f.WriteString(d.json_value(false))
	if err != nil {
		color.Red(fmt.Sprintf("ERROR : %s", err.Error()))
		f.Close()
		atexit.Exit(-1)
	}
	err = f.Close()
	if err != nil {
		color.Red(fmt.Sprintf("ERROR : %s", err.Error()))
		atexit.Exit(-1)
	}

	//pack.mcmeta
	f, err = os.Create(fmt.Sprintf("%s/%s/pack.mcmeta", path, d.datapack_name))
	if err != nil {
		color.Red(fmt.Sprintf("ERROR : %s", err.Error()))
		atexit.Exit(-1)
	}
	_, err = f.WriteString(fmt.Sprintf("{\"pack\": {\"pack_format\": %d,\"description\": \"%s\"}}", d.pack_version, d.datapack_description))
	if err != nil {
		color.Red(fmt.Sprintf("ERROR : %s", err.Error()))
		f.Close()
		atexit.Exit(-1)
	}
	err = f.Close()
	if err != nil {
		color.Red(fmt.Sprintf("ERROR : %s", err.Error()))
		atexit.Exit(-1)
	}

	color.Green("Generating MCFUNCTION files")
	e := copy.Copy(fmt.Sprintf("%s/.temp/%s", path, d.datapack_name), fmt.Sprintf("%s/%s/data", path, d.datapack_name))
	if e != nil {
		color.Red("ERROR: Something went wrong while generating function files")
		atexit.Exit(-1)
	}

	color.Green(fmt.Sprintf("SUCCESS : The datapack \"%s\" has been generated at \"%s\"!", d.datapack_name, path))
	err = os.RemoveAll(fmt.Sprintf("%s/.temp/", path))
	if err != nil {
		color.Red("ERROR: Something went wrong while removing temp folder")
	}
}
