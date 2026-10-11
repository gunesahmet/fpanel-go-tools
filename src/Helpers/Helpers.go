package Helpers

import (
	"Logger"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	tld "github.com/jpillora/go-tld"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var Logging Logger.Interface

func Trim(str string) string {
	s1 := strings.TrimSpace(str)
	s2 := strings.Trim(s1, "\t")
	s3 := strings.Trim(s2, "\n")
	return s3
}

func Timer(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func RandFloats(min float64, max float64, n int) []float64 {
	rand.Seed(time.Now().UnixNano())
	res := make([]float64, n)
	var floatn = float64(n)
	max = max / floatn
	for i := range res {
		res[i] = min + rand.Float64()*max
		//time.Sleep(time.Duration(100*time.Millisecond))
	}
	return res
}

func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func RandGuid() string {
	return uuid.New().String()
}

func RemoveSpaces(str string) string {
	return strings.Join(strings.Fields(str), " ")
}
func PrettyPrint(data []byte) {
	var res interface{}

	err := json.Unmarshal(data, &res)

	if err != nil {
		panic(err)
	}
	result, _ := json.MarshalIndent(res, "", "    ")
	fmt.Println(string(result))
}

func CurrentTimestamp() int64 {
	msec := time.Now().UnixNano() / 1000
	return msec
}

func ConvertTimeToTimestamp(date string) int64 {
	//Check the documentation on Go for the const variables!
	//They need to be exactly as they are shown in the documentation to be read correctly!
	format := "2006-01-02 15:04:05"

	t, err := time.Parse(format, date)
	if err != nil {
		HandleError(err, 1000, "Format error")
	} else {
		return t.Unix()
	}
	return 0
}

func GetDateTimeFromString(date string) (bool, time.Time) {
	//date := "2023-10-14 12:01:01"
	format := "2006-01-02 15:04:05"
	t, err := time.Parse(format, date)
	if err == nil {
		return true, t
	}

	format = "2006-01-02.15:04:05"
	t, err = time.Parse(format, date)
	if err == nil {
		return true, t
	}

	format = "2006-01-02T15:04:05Z"
	t, err = time.Parse(format, date)
	if err == nil {
		return true, t
	}

	return false, time.Time{}
}

func GetDateTimeWithFormat(date time.Time, format string) string {
	return date.Format(format)
}

func IsIpAddr(str string) bool {
	addr := net.ParseIP(str)
	if addr != nil {
		return true
	}
	return false
}

func GetCurrentDateTime() time.Time {
	return time.Now()
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	return false
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func DirIsEmpty(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer func() {
		_ = f.Close()
	}()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true
	}
	return false
}

func GetFileList(source string) []string {
	fileList := make([]string, 0)
	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	return fileList
}

func CreateDirectory(path string, recursive bool) bool {
	if recursive {
		err := os.MkdirAll(path, 755)
		if err != nil {
			return false
		}
	} else {
		err := os.Mkdir(path, 755)
		if err != nil {
			return false
		}
	}
	return true
}

func CreateSymLink(source string, target string) {
	_, _ = Exec("/bin/bash", "-c", fmt.Sprintf("ln -s %s %s", source, target))
}

func Compare(first, second string) bool {
	return strings.Compare(first, second) == 0
}

func ReadFile(path string) []string {
	readFile, err := os.Open(path)

	if err != nil {
		Logging.Error(9999, "failed to open file: %s", err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileTextLines []string

	for fileScanner.Scan() {
		if len(fileScanner.Text()) > 0 {
			fileTextLines = append(fileTextLines, fileScanner.Text())
		}
	}

	_ = readFile.Close()

	return fileTextLines
}

func RemoveLastChar(str, suffix string) string {
	if strings.HasSuffix(str, suffix) {
		str = str[:len(str)-len(suffix)]
	}
	return str
}

func TrimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

func WriteFile(content string, filePath string) error {
	if FileExists(filePath) {
		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	}
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	} else {
		_, err := f.WriteString(content)
		if err != nil {
			_ = f.Close()
			return err
		}
		err = f.Close()
		if err != nil {
			//fmt.Println(err)
			return err
		}
	}
	return nil
}

func RemoveFile(filePath string) error {
	err := os.RemoveAll(filePath)
	if err != nil {
		return err
	}
	return nil
}

func DeleteFolder(path string) error {
	dr, err := os.Open(path)  //open the path using os package function
	if err != nil {
		return err
	}
	defer dr.Close()
	names, err := dr.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(path + "/" + name)  //remove the path
		if err != nil {
			return err
		}
	}
	return os.Remove(path)
}

func Exec(command string, args ...string) (error, string) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	_ = cmd.Start()
	err := cmd.Wait()

	if err == nil {
		return nil, Trim(out.String())
	}
	return err, ""
}

func interfaceToStringInterface(in interface{}) (bool, map[string]interface{}) {
	m, ok := in.(map[string]interface{})
	if !ok {
		return false, nil
	}
	result := make(map[string]interface{})
	for k, v := range m {
		result[k] = v
	}
	return true, result
}

func interfaceToIntInterface(in interface{}) (bool, map[int]interface{}) {
	m, ok := in.(map[int]interface{})
	if !ok {
		return false, nil
	}
	result := make(map[int]interface{})
	for k, v := range m {
		result[k] = v
	}
	return true, result
}

func InterfaceToMap(in interface{}) (bool, map[string]interface{}) {
	ok, sm := interfaceToStringInterface(in)
	if !ok {
		ok, im := interfaceToIntInterface(in)
		if !ok {
			return false, nil
		}
		result := make(map[string]interface{})
		for k, v := range im {
			result[strconv.Itoa(k)] = v
		}
		return true, result
	} else {
		result := make(map[string]interface{})
		for k, v := range sm {
			result[k] = v
		}
		return true, result
	}
}

func InterfaceToArrayMap(slice interface{}) (bool, []map[string]interface{}) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return false, nil
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return false, nil
	}

	ret := make([]map[string]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		_, ret[i] = InterfaceToMap(s.Index(i).Interface())
		//ret[i] = s.Index(i).Interface()
	}
	return true, ret
}

func InterfaceToArray(slice interface{}) (bool, []interface{}) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return false, nil
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return false, nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}
	return true, ret
}

func GetParamWithNeedle(param string, defaultVal string, needle string, needleVal string) string {
	if len(param) == 0 {
		param = defaultVal
	} else {
		if param == needle {
			param = needleVal
		} else {
			param = defaultVal
		}
	}
	return param
}

func GetParam(param string, defaultVal string) string {
	if len(param) == 0 {
		param = defaultVal
	} else {
	}
	return param
}

func ExecWithoutOutPut(command string, args ...string) {
	err := exec.Command(command, args...).Run()
	if err != nil {
		_, fn, line, _ := runtime.Caller(1)
		fmt.Println(fmt.Sprintf("Fn: %s Line: %d", fn, line))
		fmt.Println(command)
		fmt.Println(err)
	}
}

func ExecWithoutWait(command string, args ...string) {
	err := exec.Command(command, args...).Start()
	if err != nil {
		_, fn, line, _ := runtime.Caller(1)
		fmt.Println(fmt.Sprintf("Fn: %s Line: %d", fn, line))
		fmt.Println(command)
		fmt.Println(err)
	}
}

func RemoveAll(str string, rem string) string {
	str = strings.ReplaceAll(str, rem, "")
	if strings.Contains(str, rem) {
		return RemoveAll(str, rem)
	}
	return str
}

func ReplaceAllOld(str string, old string, new string) string {
	str = strings.ReplaceAll(str, old, new)
	if strings.Contains(str, old) {
		return ReplaceAllOld(str, old, new)
	}
	return str
}

func ReplaceAll(str string, old string, new string) string {
	re := regexp.MustCompile(old)
	return re.ReplaceAllString(str, new)
	//str = strings.ReplaceAll(str, old, new)
	//if strings.Contains(str, old) {
	//	return ReplaceAll(str, old, new)
	//}
	//return str
}

func ReplaceAllTest(subj *string, old string, new string) string {
	re := regexp.MustCompile(old)
	return re.ReplaceAllString(*subj, new)
	//str = strings.ReplaceAll(str, old, new)
	//if strings.Contains(str, old) {
	//	return ReplaceAll(str, old, new)
	//}
	//return str
}

func IsFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return false
	case mode.IsRegular():
		return true
	}
	return false
}

func Sprintf(format string, params map[string]string) string {
	for key, val := range params {
		format = strings.Replace(format, "%{"+key+"}s", fmt.Sprintf("%s", val), -1)
	}
	return format
}

func IsDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return true
	case mode.IsRegular():
		return false
	}
	return false
}

func ParseVersion(s string, width int) int64 {
	strList := strings.Split(s, ".")
	format := fmt.Sprintf("%%s%%0%ds", width)
	v := ""
	for _, value := range strList {
		v = fmt.Sprintf(format, v, value)
	}
	var result int64
	var err error
	if result, err = strconv.ParseInt(v, 10, 64); err != nil {
		fmt.Printf("parseVersion(%s): error=%s", s, err)
		return 0
	}
	//fmt.Printf("parseVersion: [%s] => [%s] => [%d]\n", s, v, result);
	return result
}

func RemoveAfter(s string, d string) string {
	parts := strings.Split(s, d)
	return parts[0]
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func CopyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer func() {
		_ = sourcefile.Close()
	}()
	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer func() {
		_ = destfile.Close()
	}()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}
	}

	return
}

func CopyDir(source string, dest string) (err error) {

	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)
	objects, err := directory.Readdir(-1)

	for _, obj := range objects {
		sourcefilepointer := source + "/" + obj.Name()
		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			err = CopyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// perform copy
			err = CopyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return
}

func FindFile(directory string, fileName string, path *string) {
	directory = ReplaceAll(directory, "//", "/")
	files, err := ioutil.ReadDir(directory)

	if err == nil {
		for _, file := range files {
			if file.Mode().IsDir() {
				FindFile(directory+file.Name()+"/", fileName, path)
			} else {
				if strings.Compare(file.Name(), fileName) == 0 {
					*path = directory
				}
			}
		}
	}
}

func IsNumeric(val any) bool {
	typeOfVar := reflect.TypeOf(val)
	if strings.Compare(typeOfVar.Name(), "string") == 0 {
		return regexp.MustCompile(`\d`).MatchString(val.(string))
	} else if strings.Compare(typeOfVar.Name(), "int") == 0 {
		return true
	}
	return false
}

func IsEmpty(str string) bool {
	if len(str) == 0 {
		return true
	}
	return false
}

func HandleError(err error, code int, msg string) {
	_, fn, line, _ := runtime.Caller(1)
	fn = ReplaceAll(fn, ".go", "")
	fn = ReplaceAll(fn, "/usr/local/src/sources/fpanel/tools/src/", "")
	fn = ReplaceAll(fn, "H:/Linux/fpanel-tools-linux/src/", "")
	if err == nil {
		err = errors.New("nullptr")
	}
	Logging.Error(code, "[%s] %v:%d %v", msg, fn, line, err)
}

func AppendFile(str string, filePath string, perm os.FileMode) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	if _, err := f.WriteString(str + "\n"); err != nil {
		return err
	}
	return nil
}

func SearchStringInFile(needle string, filePath string) bool {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.Index(line, needle) > -1 {
			return true
		}
	}
	return true
}

func SearchStringInFileRegexp(rx string, filePath string) bool {
	status := false
	readFile, err := os.Open(filePath)

	if err != nil {
		Logging.Error(9999, "failed to open file: %s", err)
		return false
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		if len(fileScanner.Text()) > 0 {
			result, _ := regexp.MatchString(rx, fileScanner.Text())
			if result {
				status = true
				break
			}
		}
	}
	_ = readFile.Close()
	return status
}

func ExtractString(str string, rgx string) []string {
	re := regexp.MustCompile(rgx)
	return re.FindStringSubmatch(str)
}

func Equals(str string, compare string, ignoreCase bool) bool {
	if ignoreCase {
		str = strings.ToLower(str)
		compare = strings.ToLower(compare)
	}
	if strings.Compare(str, compare) == 0 {
		return true
	}
	return false
}

func ClearEmptyLinesInFile(filePath string) {
	ExecWithoutOutPut("/bin/sh", "-c", "sed -i '/^$/d' "+filePath)
}

func Split(s string, d string) []string {
	a := strings.Split(s, d)
	return a
}

func ExtractDomainName(domain string) string {
	arr := Split(domain, ".")
	if len(arr) > 2 {
		u, _ := tld.Parse(domain)
		return fmt.Sprintf("%s.%s", u.Domain, u.TLD)
	}
	return domain
}

func InSlice(val interface{}, array interface{}) bool {
	exists := false

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				exists = true
				break
			}
		}
	}
	return exists
}

func TrimLeftChars(s string, n int) string {
	m := 0
	for i := range s {
		if m >= n {
			return s[i:]
		}
		m++
	}
	return s[:0]
}

func TrimRight(s string, n int) string {
	m := 0
	for i := range s {
		if m >= n {
			return s[:i]
		}
		m++
	}
	return s[:0]
}

func Sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func Find(slice []string, key int) bool {
	for i, _ := range slice {
		if i == key {
			return true
		}
	}
	return false
}

func CheckIPV4(ip string) bool {
	if net.ParseIP(ip) == nil {
		return false
	} else {
		return true
	}
}

func StringIsExists(rx string, str string) bool {
	status := false
	result, _ := regexp.MatchString(rx, str)
	if result {
		status = true
	}
	return status
}