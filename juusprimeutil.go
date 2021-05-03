package juusprime

//juusprime is an engine to find prime Sextuplets and/or Quintuplets and/or Quadruplets
//GNU GPL3
//author: Julius Schön / R. Spicer, 2021

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func init() {
	Basis29Path = os.Getenv("HOME")
	DataPath = Basis29Path
	ReadConfig()
	initBypass = false //I prefer this method instead of a parameter
	isInteger = regexp.MustCompile(`^[-+]?\d+$`)
}

var (
	emptyString = ""
	//Basis29Path : string to hold configuration path to 29basis files
	Basis29Path string
	//DataPath : string to hold configuration path to write tuplet files
	DataPath   string
	initBypass = true
	isInteger  *regexp.Regexp
)

const (
	trimString = " \t\r\n"
	//Basis29PathStr const defining a string for Basis29Path
	Basis29PathStr = "Basis29Path"
	//DataPathStr const defining a string for DataPath
	DataPathStr = "DataPath"
)

const (
	//do not change or re-use the #s for any cfg* consts
	cfgBasis29Path = "0"
	cfgDataPath    = "1"
)

//IsConfigured : check to see if config/paths have been setup,
//if false is returned then recommended to run Configure()
func IsConfigured() bool {
	return FileExists(ConfigFilename())
}

//ConfigFilename : return the expected path of the config
//file, executable name + ".config", dev environment writes jup.config to HOME
func ConfigFilename() string {
	cfgName := os.Args[0]
	if strings.HasPrefix(cfgName, "/tmp/") {
		//adjust for devel. environment and "go run"
		return filepath.Join(os.Getenv("HOME"), "jup.config")
	}
	return cfgName + ".jup.config"
}

//Configure : should be called at least once or routines will have to
//default to HOME folder for writing/reading data, the answers will be written to
//config file, for now simple setup instead of full blown ini file
func Configure() {
	fmt.Println("\njuusprime Configuration:")
	fmt.Println("")
	cfgName := ConfigFilename()

	setPath("29Basis", &Basis29Path)
	setPath("Output data", &DataPath)
	cfg, err := FileOpen(cfgName, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer FileClose(cfg)
	fmt.Fprintln(cfg, fmt.Sprintf("%s^%s^%s", cfgBasis29Path, Basis29PathStr, Basis29Path))
	fmt.Fprintln(cfg, fmt.Sprintf("%s^%s^%s", cfgDataPath, DataPathStr, DataPath))
	fmt.Println("Settings written to config file: ", cfgName)

}

//ReadConfig : open the configuration file and retrieve settings
func ReadConfig() {

	cfgName := ConfigFilename()
	if !FileExists(cfgName) {
		fmt.Println("ReadConfig(): config file not found.")
		if !initBypass {
			fmt.Println("Opening Configure() routine:")
			Configure()
		}
		return
	}

	failure := false
	//fmt.Println("Reading from config file: ", cfgName)

	setVar := func(idx, path string) {
		fi, _ := os.Stat(path)
		switch idx {
		case cfgBasis29Path:
			if !FileExists(path) || !fi.Mode().IsDir() {
				fmt.Println(fmt.Sprintf("'%s' is not a valid path to a Folder", path))
				failure = true
				return
			}
			Basis29Path = path
		case cfgDataPath:
			if !FileExists(path) || !fi.Mode().IsDir() {
				fmt.Println(fmt.Sprintf("'%s' is not a valid path to a Folder", path))
				failure = true
				return
			}
			DataPath = path
		}
	}

	cfg, err := FileOpen(cfgName, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer FileClose(cfg)

	cfgScan := bufio.NewScanner(cfg)

	lineNo := 0
	for cfgScan.Scan() {
		lineNo++
		sl := strings.Split(cfgScan.Text(), "^")
		if len(sl) != 3 {
			fmt.Println(fmt.Sprintf("config file line #%v '%v' is invalid", lineNo, sl))
			continue
		}
		setVar(sl[0], sl[2])
	}
	if failure {
		fmt.Println("There were failures reading config file.")
		if !initBypass {
			fmt.Println("Running Configure() routine")
			Configure()
		}
	}
}

//setPath : hint is a short description of which path is being
//set and currPath the variable current setting, Only accepts folders
func setPath(hint string, currPath *string) {
	var (
		path        string
		wascanceled bool
	)

	for {
		if path, wascanceled = GetUserInput(fmt.Sprintf("Path for %s files:", hint), *currPath, "xxx"); wascanceled {
			fmt.Println(fmt.Sprintf("\ncanceled, %s path remains as '%s'", hint, *currPath))
			fmt.Println("")
			break
		}
		fi, _ := os.Stat(path)
		if !FileExists(path) || !fi.IsDir() {
			fmt.Println(fmt.Sprintf("'%s' is not a valid path to a Folder", path))
			continue
			//return
		}
		*currPath = path
		fmt.Println(fmt.Sprintf("\n%s path set to '%s'", hint, *currPath))
		fmt.Println("")
		break
	}
}

//DisplayProgressBookend : Display the msg, isStarting is true for
//the first call (a header) and false for the end call (a footer), it
//returns a Time, meant to be used to wrap DisplayProgress
func DisplayProgressBookend(msg string, isStarting bool) time.Time {
	var theTime time.Time
	switch isStarting {
	case true:
		theTime = time.Now()
		fmt.Println("start:", theTime)
		fmt.Println(msg)
		fmt.Print("| ")
	default:
		theTime = time.Now()
		fmt.Print("100% | " + msg + " \n")
		fmt.Println("end:", theTime)
		fmt.Println("")
	}
	return theTime
}

//DisplayProgress : from and to are usually ranges of TNumbers, but it
//depends on what you want to measure, part is how fine to divide the
//notifications: for example, 20 will partition it into 5% blocks, 100
//into 1% blocks, this is setup as a closure func, its virtue is being
//able to accurately measure progress even when the index vs TNum do
//not match, ie. "jerky" progress, primarily used when generating "small"
//int/int64 compatible ranges, see also DisplayProgressBig
func DisplayProgress(from, to int64, part int) func() {
	calcRange := to - from + 1
	if part > 100 {
		part = 100
	}
	progressCnt := int(math.Floor(float64(calcRange) / float64(part)))
	calcRangeCnt := 0
	progressCntr := 0
	return func() {
		calcRangeCnt++
		progressCntr++
		if progressCntr > progressCnt {
			progressCntr = 0
			fmt.Printf("%v%% ", math.Floor((float64(calcRangeCnt)/float64(calcRange))*100))
		}
	}
}

//DisplayProgressBig : closure for displaying progress for data where the range
//is not evenly, ie. 1 to 1, probed, for example, a range from 234 to 2345
//which is actually indexed randomly in that range, used when processing GTE 31's
//against the Basis29 file, works for huge numbers
func DisplayProgressBig(from, to *big.Int, part int64) func(*big.Int, *big.Int) {
	printResult := big.NewInt(0)
	calcRange := big.NewInt(0).Sub(to, from)
	calcRange.Add(calcRange, big1)
	if part > 100 {
		part = 100
	}
	prec := uint(to.BitLen()) + 20
	bf100 := big.NewFloat(0).SetPrec(prec).SetInt64(100)
	numerator := big.NewFloat(0).SetPrec(prec)
	denominator := big.NewFloat(0).SetPrec(prec).SetInt(calcRange)
	bigFloat := big.NewFloat(0).SetPrec(prec).Quo(
		big.NewFloat(0).SetPrec(prec).SetInt(calcRange),
		big.NewFloat(0).SetPrec(prec).SetInt64(part))
	progressCnt, _ := bigFloat.Int(nil)
	//progressCnt := int(math.Floor(float64(calcRange) / float64(part)))
	calcRangeCnt := big.NewInt(0)
	progressCntr := big.NewInt(0)

	return func(currPos, lastPos *big.Int) {
		//keep track of the uneven change through the range
		printResult.Sub(currPos, lastPos)
		lastPos.Set(currPos)

		calcRangeCnt.Add(calcRangeCnt, printResult)
		progressCntr.Add(progressCntr, printResult)
		//we've reached a position greater than the precalculated chunk size
		if progressCntr.Cmp(progressCnt) > 0 {
			progressCntr.Set(big0)
			//at this chunk, calculate the percent progress (eg if split into 20 should get around 5%)
			bigFloat.Quo(numerator.SetInt(calcRangeCnt), denominator)
			printResult, _ = bigFloat.Mul(bigFloat, bf100).Int(nil)
			fmt.Printf("%v%% ", printResult)
		}
	}
}

//FileExists : Returns true if file or directory (?) exists
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//ChangeFileExt : change the extension of fName to newExt
func ChangeFileExt(fName, newExt string) string {
	return strings.TrimSuffix(fName, filepath.Ext(fName)) + newExt
}

//FileOpen : wrapper for opening files, if append is true
//the file is opened for reading/writing/appending, if false
//the file is created or overwritten
func FileOpen(p string, append bool) (*os.File, error) {
	var (
		f   *os.File
		err error
	)

	switch append {
	case true:
		f, err = os.Open(p)
	case false:
		f, err = os.Create(p)
	}

	if err != nil {
		return nil, err
		/*switch {
			case errors.Is(err, os.ErrInvalid) : return nil, err
			default :	panic(err)
		}*/
	}
	return f, nil
}

//FileClose : close a file previously opened with FileOpen
//with error checking, usually called with defer
func FileClose(f *os.File) {
	if err := f.Close(); err != nil {
		if errors.Is(err, os.ErrClosed) {
			fmt.Println(fmt.Sprintf("File: %s is already closed..", f.Name()))
			return
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		panic(err)
	}
}

//valueClean : simple helper function to trim and (semi)validate parameters sent to func()s
//afunc will be run if it is a non-empty func().
func valueClean(toClean, theDefault *string) {
	//this is mostly for programmer errors and/or dynamic menu errors.
	//if menus are properly constructed, one will never see this displayed
	//except when vAfterBlock is specified.
	*toClean = strings.Trim(*toClean, trimString)
	if *toClean == "" {
		*toClean = *theDefault
	}
}

//GetUserBoolChoice : wraps GetUserInput to allow various true values
func GetUserBoolChoice(prompt, deflt string, result *bool) (canceled bool) {
	var input string
	*result = false
	if input, canceled = GetUserInput(prompt, deflt, "x"); canceled {
		return canceled
	}
	switch strings.ToUpper(input) {
	case "T", "TRUE", "Y", "YES":
		*result = true
	}
	return canceled
}

//GetUserInput : Good for basic input, returns the string the User enters.
//The calling func must deal with type verification
func GetUserInput(prompt, defaultTo, cancel string) (string, bool) {
	fmt.Println("")
	scanner := bufio.NewScanner(os.Stdin)
	valueClean(&prompt, &emptyString)
	valueClean(&defaultTo, &emptyString)
	valueClean(&cancel, &emptyString)
	if prompt == emptyString {
		prompt = "Enter data: "
	}
	if cancel == emptyString {
		cancel = "q"
	}
	fmt.Println(fmt.Sprintf("%s (<RET> for default '%s')  ['%s' %s]:", prompt, defaultTo, cancel, "cancels"))
	fmt.Print("=> ")
	scanner.Scan()
	input := strings.Trim(scanner.Text(), trimString)
	if input == "" && defaultTo != "" {
		input = defaultTo
		//fmt.Println("<canceled>")
	}
	if input == cancel {
		fmt.Println("<canceled>")
		return "", true
	}
	return input, false
}

//GetUserInputInteger : Get input and validate that answer is a Number. Loops
//forever unless canceled.
func GetUserInputInteger(prompt, defaultTo, cancel string) (string, bool) {
	for {
		input, wasCanceled := GetUserInput(prompt, defaultTo, cancel)
		if isInteger.MatchString(input) || wasCanceled {
			return input, wasCanceled
		}
		fmt.Println(fmt.Sprintf("'%s' is not an integer, it has spaces or characters other than digits", input))
	}
}

//GetUserConfirmation : basic confirmation "dialog"
func GetUserConfirmation(prompt, answer string) bool {
	fmt.Println("")
	scanner := bufio.NewScanner(os.Stdin)

	const trueAnswer = "yes"
	aPrompt := strings.Trim(prompt, trimString)
	anAnswer := strings.Trim(answer, trimString)
	if anAnswer == "" {
		anAnswer = trueAnswer
		fmt.Println("GetUserInput: No answer sent in, User doesn't know what to Type!!")
		fmt.Printf("setting the 'right' answer to '%s'\n", trueAnswer)
		//they sent in a prompt, it may have the wrong word to type
		//add the newly programatically created "answer"
		if aPrompt != "" {
			aPrompt = aPrompt + fmt.Sprintf(" [ignore prompt, answer set to '%s']", anAnswer)
		}
	}
	if aPrompt == "" {
		aPrompt = fmt.Sprintf("Type '%s' for returning true", anAnswer)
	}
	fmt.Println(fmt.Sprintf("%s ['%s' confirms, anything else cancels]:", aPrompt, anAnswer))
	fmt.Print(">>: ")
	scanner.Scan()
	input := strings.Trim(scanner.Text(), " \t\r")
	if input == anAnswer {
		return true
	}
	fmt.Println("<canceled>")
	return false
}

func waitForInput() {
	fmt.Println("Press <RET> to continue...")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	scanner.Text()
}

//GetFileInfosFromFolder : a func that gets a list of all files
//in given path
func GetFileInfosFromFolder(path string) ([]os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fileInfos, err := f.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return fileInfos, nil
}

//Choose29BasisFile : Routine that prepares a list of 29basis files
//from the Basis29Path and presents them to the user.
func Choose29BasisFile() (string, bool) {

	fis, err := GetFileInfosFromFolder(Basis29Path)
	if err != nil {
		fmt.Println(err)
		return "", false
	}

	var (
		input       string
		wasCanceled bool
		result      string = ""
	)

	sort.Slice(fis, func(i, j int) bool { return fis[i].Name() < fis[j].Name() })

	//prepare and show choices
	choice := -1
	listMap := make(map[int]int)

	fmt.Println(fmt.Sprintf("Searching path '%s' for files of form: %s---%s", Basis29Path, filePrefix29, fileExtRaw29))
	fmt.Println("")
	for idx, fi := range fis {
		if filepath.Ext(fi.Name()) == fileExtRaw29 && strings.HasPrefix(fi.Name(), filePrefix29) {
			choice++
			listMap[choice] = idx
			fmt.Println(fmt.Sprintf("%v : %s", choice, fi.Name()))
		}
	}

	if choice < 0 {
		fmt.Println("<no standard files found>")
	}

	//add a manual choice just in case
	optManual := strconv.Itoa(choice + 1)
	fmt.Println("")

	//get users returned choice
	if choice >= 0 {
		//other files were found, otherwise don't make them have to chose manaul option
		fmt.Println(fmt.Sprintf("%s : Enter path manually (advanced users, recommended is to use default basis file)", optManual))
		if input, wasCanceled = GetUserInputInteger("Enter choice for 29Basis rawdata file:", "0", "x"); wasCanceled {
			return result, false
		}
	} else {
		input = optManual
		fmt.Println("\n(if you are confused by this, or don't know where it could be, then just do/re-do your 29 basis file)")
	}

	switch input {
	case optManual:
		if input, wasCanceled = GetUserInput("Manual entry: enter full path to 29Basis raw data file:", "<none>", "xxx"); wasCanceled {
			fmt.Println("\ncanceled")
			fmt.Println("")
			return result, false
		}
		if !FileExists(input) {
			fmt.Println(fmt.Sprintf("\nfile %s does not exist, quitting.", input))
			fmt.Println("")
			return result, false
		}
		result = input
	default:
		selected, _ := strconv.Atoi(input)

		found := false
		for key := range listMap {
			//fmt.Println(key)
			if selected == key {
				found = true
				break
			}
		}

		if !found {
			fmt.Println(fmt.Sprintf("\n%v is an invalid choice in the list...", selected))
			fmt.Println("")
			return result, false
		}
		result = filepath.Join(Basis29Path, fis[listMap[selected]].Name())
	}

	return result, true
}
