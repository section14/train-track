package extract

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Element struct {
	NewHTML   []byte
	Js        []byte
	JsRemoved bool
}

type Handler struct {
	AttrName         string
	DataName         []string
	FileName         string
	RelativeFileName string
	FuncName         string
}

type extractState struct {
	extracting     bool
	foundMatch     bool
	currentHandler Handler
	emptyHandler   Handler
}

type MainJs struct {
	Handlers []Handler
}

func MakeJsPath(path string, wd string) (string, string) {
	split := strings.Split(path, fmt.Sprintf("%s/%s", wd, "templates"))
	trimmed := strings.TrimSuffix(split[1], ".html")
	jsPath := fmt.Sprintf("%s%s%s%s", wd, "/static/js/extracted", trimmed, ".js")
	relativeJsPath := fmt.Sprintf("%s%s%s", "/static/js/extracted", trimmed, ".js")

	return jsPath, relativeJsPath
}

func BuildJsFile(currentDir string, data map[string]string) error {
	for jsFileName, jsFileData := range data {

		//create new JS file (plus optional directory) and write to it
		err := os.MkdirAll(filepath.Dir(jsFileName), 0770)
		if err != nil {
			return err
		}

		file, err := os.Create(jsFileName)
		if err != nil {
			return errors.New(fmt.Sprintf("couldn't open extracted.js file %s", err))
		}
		defer file.Close()

		_, err = file.WriteString(jsFileData)
		if err != nil {
			return errors.New(fmt.Sprintf("couldn't write to extracted.js %s", err))
		}
	}

	return nil
}

func dashToCamel(s string) string {
	var builder strings.Builder
	words := strings.Split(s, "-")

	for i, word := range words {
		if i == 0 {
			continue
		} else if i == 1 {
			builder.WriteString(word)
		} else {
			caser := cases.Title(language.AmericanEnglish)
			builder.WriteString(caser.String(word))
		}
	}

	return builder.String()
}

func extractHanlder(path string, line []byte, extractor *extractState) []byte {
	matcher := "//@handle:"

	if !extractor.extracting && strings.Contains(string(line), matcher) {
		//get data-* attribute names
		trimmed := strings.Trim(string(line), " ")
		attr := strings.TrimPrefix(trimmed, matcher)

		rawDataNames := strings.Split(attr, ":")
		dataNames := make([]string, 0)

		for _, n := range rawDataNames {
			dataNames = append(dataNames, dashToCamel(n))
		}

		//get current working directory
		wd, _ := os.Getwd()

		//add everything but the function name
		jsFileName, jsRelativeFileName := MakeJsPath(path, wd)
		extractor.currentHandler.AttrName = attr
		extractor.currentHandler.DataName = dataNames
		extractor.currentHandler.FileName = jsFileName
		extractor.currentHandler.RelativeFileName = jsRelativeFileName
		extractor.extracting = true

		return nil
	}

	if extractor.extracting {
		//regex the func name
		trimmed := strings.Trim(string(line), " ")

		pattern := regexp.MustCompile(`^export\sconst\s(\w+)?[\s=]?`)
		match := pattern.FindStringSubmatch(trimmed)

		//add func name to use in handlers
		extractor.currentHandler.FuncName = match[1]
		extractor.foundMatch = true
	}

	return line
}

func ExtractJs(path string, handlers []Handler) (Element, []Handler, error) {
	elem := Element{NewHTML: nil, Js: nil, JsRemoved: false}

	file, err := os.Open(path)
	if err != nil {
		return elem, handlers, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var sbHTML strings.Builder
	var sbJs strings.Builder

	startExtracting := false

	extractor := extractState{
		extracting:     false,
		foundMatch:     false,
		currentHandler: Handler{},
		emptyHandler:   Handler{},
	}

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Trim(line, " ") == "<script>" {
			startExtracting = true
			elem.JsRemoved = true
			continue
		}

		if strings.Trim(line, " ") == "</script>" {
			startExtracting = false
			continue
		}

		if startExtracting {
			//add handler extraction here
			newLine := extractHanlder(path, scanner.Bytes(), &extractor)

			if !extractor.extracting {
				sbJs.Write(newLine)
				sbJs.Write([]byte("\n"))
			}

			if extractor.foundMatch {
				//add the func name
				handlers = append(handlers, extractor.currentHandler)

				//write func definition line
				sbJs.Write(newLine)
				sbJs.Write([]byte("\n"))

				//reset the func extractor
				extractor.currentHandler = extractor.emptyHandler
				extractor.extracting = false
				extractor.foundMatch = false
			}
		} else {
			sbHTML.Write(scanner.Bytes())
			sbHTML.Write([]byte("\n"))
		}

	}

	elem.NewHTML = []byte(sbHTML.String())
	elem.Js = []byte(sbJs.String())

	return elem, handlers, nil
}

func BuildMainJs(path string, handlers []Handler) error {
	wd, _ := os.Getwd()

	//open and create files
	mainFile, err := os.Open(fmt.Sprintf("%s%s", wd, path))
	if err != nil {
		return err
	}
	defer mainFile.Close()

	mainFileName := fmt.Sprintf("%s/static/js/extracted/main.js", wd)
	file, err := os.Create(mainFileName)
	if err != nil {
		return errors.New(fmt.Sprintf("couldn't create extracted main.js file %s", err))
	}
	defer file.Close()

	//write content
	var builder strings.Builder
	var existingImports strings.Builder
	var existingCode strings.Builder

	//add existing js file data
	scanner := bufio.NewScanner(mainFile)

	//add existing code
	for scanner.Scan() {
		line := scanner.Bytes()

		if strings.Contains(string(line), "import") {
			existingImports.Write(line)
			existingImports.WriteString("\n")
		} else {
			existingCode.Write(line)
			existingCode.WriteString("\n")
		}

	}

	//add existing imports
	builder.WriteString(existingImports.String())

	//add generated content
	builder.WriteString("//--------------------------------------------/\n")
	builder.WriteString("//------ Generated Content: DON'T TOUCH ------/\n")
	builder.WriteString("//--------------------------------------------/\n")
	builder.WriteString("\n")

	//build a map with filenames to prevent duplicate import statements
	fileNameMap := make(map[string][]string)
	dataToFuncMap := make(map[string][]string)

	for _, handler := range handlers {
		_, ok := fileNameMap[handler.RelativeFileName]
		if !ok {
			fileNameMap[handler.RelativeFileName] = make([]string, 0)
		}

		dataToFuncMap[handler.FuncName] = handler.DataName
	}

	//add functions
	for _, handler := range handlers {
		arr := fileNameMap[handler.RelativeFileName]
		arr = append(arr, handler.FuncName)
		fileNameMap[handler.RelativeFileName] = arr
	}

	buildImports(&builder, fileNameMap)
	switchStmt := buildGlobalSwitch(dataToFuncMap)
	buildGlobalClickHandler(&builder, switchStmt)

	builder.WriteString("//--------------------------------------------/\n")
	builder.WriteString("//---------- End Generated Content: ----------/\n")
	builder.WriteString("//--------------------------------------------/\n")
	builder.WriteString("\n\n")

	//add existing code
	builder.WriteString(existingCode.String())

	_, err = file.WriteString(builder.String())
	if err != nil {
		return errors.New(fmt.Sprintf("couldn't write to extracted main.js %s", err))
	}

	return nil
}

func buildImports(builder *strings.Builder, fileNameMap map[string][]string) {
	var nameBuilder strings.Builder

	for fileName, funcArr := range fileNameMap {
		nameBuilder.WriteString("import { ")
		importNames := strings.Join(funcArr, ", ")
		nameBuilder.WriteString(importNames)
		nameBuilder.WriteString(" } from \"")
		nameBuilder.WriteString(fileName)
		nameBuilder.WriteString("\"")
		nameBuilder.WriteString("\n")
	}

	nameBuilder.WriteString("\n")
	builder.WriteString(nameBuilder.String())
}

func buildGlobalClickHandler(builder *strings.Builder, switchStmt string) {
	//beginning of function
	builder.WriteString("const globalClickHandler = (e) => {\n")
	builder.WriteString("    const dataSet = e.target.dataset\n")
	builder.WriteString("    const key = Object.keys(dataSet)[0]\n\n")
	builder.WriteString("    const button = e.target.closest(\"button\")\n\n")
	builder.WriteString("    if (button) {\n")
	builder.WriteString("        switch(key) {\n")

	builder.WriteString(switchStmt)

	//end of function
	builder.WriteString("        }\n")
	builder.WriteString("    }\n")
	builder.WriteString("}\n\n")
	builder.WriteString("// global click event listener\n")
	builder.WriteString("document.addEventListener(\"click\", globalClickHandler)\n\n")
}

//build switch case in main.js click handler
func buildGlobalSwitch(dataToFuncMap map[string][]string) string {
	var switchBuilder strings.Builder

	for funcName, cases := range dataToFuncMap {
		for _, caseName := range cases {
			switchBuilder.WriteString(fmt.Sprintf("        case \"%s\": \n", caseName))
		}
		if len(cases) > 1 {
			switchBuilder.WriteString(fmt.Sprintf("            %s(dataSet)\n", funcName))
		} else {
			switchBuilder.WriteString(fmt.Sprintf("            %s(dataSet[key])\n", funcName))
		}
		switchBuilder.WriteString("            break\n")
	}
	switchBuilder.WriteString("         default:\n")
	switchBuilder.WriteString("            console.error(\"unknown dataset item: \", key)\n")

	return switchBuilder.String()
}
