package generators

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	dTypes "codeGen.com/app/dataTypes"
)

func CreateInterface(entity dTypes.Entity) {
	importsInterface := "import { Document } from 'mongoose';\n\n\n"
	var interfaceElements string
	for i, s := range entity.TypeDetails {
		if !s.Array {
			if i != len(entity.TypeDetails)-1 {
				
				interfaceElements += fmt.Sprintf("\t%s: %s,\n", s.Name, s.DataType)
			} else {
				
				interfaceElements += fmt.Sprintf("\t%s: %s\n", s.Name, s.DataType)
			}

		} else {
			if i != len(entity.TypeDetails)-1 {
				
				interfaceElements += fmt.Sprintf("\t%s: %s[],\n", s.Name, s.DataType)
			} else {
				
				if s.Mandatory {
					interfaceElements += fmt.Sprintf("\t%s: %s[]\n", s.Name, s.DataType)
				} else {
					interfaceElements += fmt.Sprintf("\t%s?: %s[]\n", s.Name, s.DataType)
				}

			}

		}

	}
	modelInterface := fmt.Sprintf("export interface I%s extends Document {\n%s}\n\n", strings.Title(entity.ModelName), interfaceElements)
	writeInterface := fmt.Sprintf("%s%s", importsInterface, modelInterface)
	interfaceFileName := fmt.Sprintf("%s.interface.ts", entity.ModelName)
	ioutil.WriteFile(interfaceFileName, []byte(writeInterface), 0755)
}

func UpdateInterface (entity dTypes.Entity) {
	interfaceFile:= fmt.Sprintf("%s.interface.ts", entity.ModelName)
	interfaceDat, err := ioutil.ReadFile(interfaceFile)
	if err != nil {
		fmt.Println(err)
	}
	interfaceRe, _ := regexp.Compile(`Document (?s)(\{.*\})`)
	interfaceString := string(interfaceDat)
	interfaceMatch := interfaceRe.FindStringSubmatch(interfaceString)[1]
	wordReSpecial := regexp.MustCompile(`(\w+\[\]|\w+)`)
	interfaceJson := wordReSpecial.ReplaceAllString(interfaceMatch, `"$1"`)
	var interfaceSchema map[string]interface{}
	err = json.Unmarshal([]byte(interfaceJson), &interfaceSchema)
	for _, s := range entity.TypeDetails {
		if !s.Array {
			interfaceSchema[s.Name] = s.DataType
		} else {
			interfaceSchema[s.Name] = fmt.Sprintf(`%s[]`, s.DataType)
		}
	}
	b2, _ := json.Marshal(interfaceSchema)
	reg := regexp.MustCompile(`"([^"]*)"`)
	unquotedInterface := reg.ReplaceAllString(string(b2), `${1}`)
	reg = regexp.MustCompile(`,`)
	formattedInterface := reg.ReplaceAllString(unquotedInterface, `,\n`)
	reg = regexp.MustCompile(`{`)
	formattedInterface = reg.ReplaceAllString(formattedInterface, `{\n`)
	reg = regexp.MustCompile(`}`)
	formattedInterface = reg.ReplaceAllString(formattedInterface, `\n}`)
	replacementInterface := fmt.Sprintf("Document %s", formattedInterface)
	newDataInterface := interfaceRe.ReplaceAllString(interfaceString, replacementInterface)
	formattedOutInterface := strings.Replace(newDataInterface, `\n`, "\n\t", -1)
	formattedOutInterface = strings.Replace(formattedOutInterface, `,`, ", ", -1)
	formattedOutInterface = strings.Replace(formattedOutInterface, `:`, ": ", -1)
	ioutil.WriteFile(interfaceFile, []byte(formattedOutInterface), 0755)
}