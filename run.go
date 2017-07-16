package main

import(
	"io/ioutil"
	"os"
	path "path/filepath"
	"runtime"
	"strings"
)

const(
	RunCommandUsage = `Use {{printf "org help %s" .Name | bold}} for more information. {{endline}}`
	RunCommandShort = ``
	RunCommandLong = ``
)

var (
	mainFiles ListOpts
	excludedPaths StrFlags
	buildTags string
	currpath string
	appname string
	exit chan bool
	watchVendor bool
	currentGOPATH string
	env string
)
func init(){
	var runCmd = newRunCommand()
	runCmd.Flag.Var(&mainFiles, "main", "Specify main go files.")
	runCmd.Flag.Var(&excludedPaths, "e", "List of paths to exclude.")
	runCmd.Flag.BoolVar(&watchVendor, "vendor", false, "Enable watch vendor folder.")
	runCmd.Flag.StringVar(&buildTags, "tags", "", "Set the build tags. See: https://golang.org/pkg/go/build/")
	runCmd.Flag.StringVar(&env, "env", "", "Set the Beego run mode.")
	exit = make(chan bool)
	commands[runCmd.name] = runCmd
}

func newRunCommand() *Command{
	return &Command{
		UsageLine: RunCommandUsage,
		Short: RunCommandShort,
		Long: RunCommandLong,
		Run: callRun,
	}
}

// callRun
func callRun(cmd *Command, args []string) int {
	var gopath string
	gopath = os.Getenv("GOPATH") 
	if gopath == ""{
		colorLog("[ERRO] GOPATH is not found")
		os.Exit(1)
	}

	colorLog("[INFO] Running '%s' as 'appliation name.' \n", appname)
	colorLog("[INFO ] Current directory: %s \n", currpath)

	var paths []string
	readAppDirectories(currpath, &paths)

	files := []string{}
	for _, arg := range mainFiles {
		if len(arg) > 0 {
			files = append(files, arg)
		}
	}

	// Start the Reload server (if enabled)
	// if config.Conf.EnableReload {
	// 	startReloadServer()
	// }
	NewWatcher(paths, files, false)
	AutoBuild(files, false)


	for {
		<-exit
		runtime.Goexit()
	}
}

// readAppDirectories
func readAppDirectories(directory string, paths *[]string) {
	fileInfos, err := ioutil.ReadDir(directory)
	if err != nil {
		return
	}

	useDirectory := false
	for _, fileInfo := range fileInfos {

		if !watchVendor && strings.HasSuffix(fileInfo.Name(), "vendor") {
			continue
		}

		if isExcluded(path.Join(directory, fileInfo.Name())) {
			continue
		}

		if fileInfo.IsDir() && fileInfo.Name()[0] != '.' {
			readAppDirectories(directory+"/"+fileInfo.Name(), paths)
			continue
		}

		if useDirectory {
			continue
		}
	}
}


// If a file is excluded
func isExcluded(filePath string) bool {
	for _, p := range excludedPaths {
		absP, err := path.Abs(p)
		if err != nil {
			colorLog("Cannot get absolute path of '%s'", p)
			continue
		}
		absFilePath, err := path.Abs(filePath)
		if err != nil {
			colorLog("Cannot get absolute path of '%s'", filePath)
			break
		}
		if strings.HasPrefix(absFilePath, absP) {
			colorLog("'%s' is not being watched", filePath)
			return true
		}
	}
	return false
}