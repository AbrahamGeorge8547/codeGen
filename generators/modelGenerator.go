package generators

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	dTypes "codeGen.com/app/dataTypes"
)

func CreateModel(entity dTypes.Entity) {
	// const fileData =
	imports := "import * as mongoose from 'mongoose';\n\n\n"
	var elements string
	var interfaceElements string
	for i, s := range entity.TypeDetails {
		if !s.Array {
			if i != len(entity.TypeDetails)-1 {
				elements += fmt.Sprintf("\t%s: {type: %s, required: %t},\n", s.Name, strings.Title(s.DataType), s.Mandatory)
				interfaceElements += fmt.Sprintf("\t%s: %s,\n", s.Name, s.DataType)
			} else {
				elements += fmt.Sprintf("\t%s: {type: %s, required: %t}\n", s.Name, strings.Title(s.DataType), s.Mandatory)
				interfaceElements += fmt.Sprintf("\t%s: %s\n", s.Name, s.DataType)
			}

		} else {
			if i != len(entity.TypeDetails)-1 {
				elements += fmt.Sprintf("\t%s: {type: Array, required: %t},\n", s.Name, s.Mandatory)
				interfaceElements += fmt.Sprintf("\t%s: %s[],\n", s.Name, s.DataType)
			} else {
				elements += fmt.Sprintf("\t%s: {type: Array, required: %t}\n", s.Name, s.Mandatory)
				if s.Mandatory {
					interfaceElements += fmt.Sprintf("\t%s: %s[]\n", s.Name, s.DataType)
				} else {
					interfaceElements += fmt.Sprintf("\t%s?: %s[]\n", s.Name, s.DataType)
				}

			}

		}

	}
	schema := fmt.Sprintf("const %sSchema = new mongoose.Schema({\n%s})\n\n", strings.Title(entity.ModelName), elements)
	modelInterface := fmt.Sprintf("export interface I%s extends mongoose.Document {\n%s}\n\n", strings.Title(entity.ModelName), interfaceElements)
	exports := fmt.Sprintf("export default mongoose.model<I%s>(\"%s\", %sSchema)", strings.Title(entity.ModelName), strings.Title(entity.ModelName), strings.Title(entity.ModelName))
	writeFile := fmt.Sprintf("%s%s%s%s", imports, schema, modelInterface, exports)
	fileName := fmt.Sprintf("%s.model.ts", entity.ModelName)
	ioutil.WriteFile(fileName, []byte(writeFile), 0755)
}

func UpdateModel(entity dTypes.Entity) {
	fileName := fmt.Sprintf("%s.model.ts", entity.ModelName)
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
	}
	modelRe, _ := regexp.Compile(`new mongoose.Schema\((\{[^]]+})`)
	interfaceRe, _ := regexp.Compile(`mongoose.Document (?s)(\{.*\})`)
	dataString := string(dat)
	modelMatch := modelRe.FindStringSubmatch(dataString)[1]
	interfaceMatch := interfaceRe.FindStringSubmatch(dataString)[1]
	var wordRe = regexp.MustCompile(`(\w+)`)
	modelJson := wordRe.ReplaceAllString(modelMatch, `"$1"`)
	wordReSpecial := regexp.MustCompile(`(\w+\[\]|\w+)`)
	interfaceJson := wordReSpecial.ReplaceAllString(interfaceMatch, `"$1"`)
	if err != nil {
		fmt.Println(err)
	}
	var model map[string]map[string]interface{}
	var interfaceSchema map[string]interface{}
	err = json.Unmarshal([]byte(modelJson), &model)
	err = json.Unmarshal([]byte(interfaceJson), &interfaceSchema)
	if err != nil {
		fmt.Println(err)
	}

	for _, s := range entity.TypeDetails {
		if !s.Array {
			if val, ok := model[s.Name]; ok {
				interfaceSchema[s.Name] = s.DataType
				val["type"] = strings.Title(s.DataType)
				val["required"] = s.Mandatory
			} else {
				model[s.Name] = map[string]interface{}{
					"required": s.Mandatory,
					"type":     strings.Title(s.DataType),
				}
			}
		} else {
			interfaceSchema[s.Name] = fmt.Sprintf(`%s[]`, s.DataType)
			if val, ok := model[s.Name]; ok {
				val["type"] = "Array"
				val["required"] = s.Mandatory
			} else {
				model[s.Name] = map[string]interface{}{
					"required": s.Mandatory,
					"type":     strings.Title(s.DataType),
				}
			}
		}
	}
	fmt.Println(interfaceSchema)
	b, _ := json.Marshal(model)
	b2, _ := json.Marshal(interfaceSchema)
	reg := regexp.MustCompile(`"([^"]*)"`)
	fmt.Println(string(b))
	unquotedModel := reg.ReplaceAllString(string(b), `${1}`)
	unquotedInterface := reg.ReplaceAllString(string(b2), `${1}`)
	reg = regexp.MustCompile(`},`)
	formattedModel := reg.ReplaceAllString(unquotedModel, `}, \n`)
	reg = regexp.MustCompile(`,`)
	formattedInterface := reg.ReplaceAllString(unquotedInterface, `,\n`)
	reg = regexp.MustCompile(`{`)
	formattedInterface = reg.ReplaceAllString(formattedInterface, `{\n`)
	reg = regexp.MustCompile(`}`)
	formattedInterface = reg.ReplaceAllString(formattedInterface, `\n}`)
	reg = regexp.MustCompile("^(.*?){")
	formattedModel = reg.ReplaceAllString(formattedModel, `{ \n`)
	reg = regexp.MustCompile(`}}`)
	formattedModel = reg.ReplaceAllString(formattedModel, `} \n }`)
	replacementModel := fmt.Sprintf("new mongoose.Schema(%s", formattedModel)
	replacementInterface := fmt.Sprintf("mongoose.Document %s", formattedInterface)
	newData := modelRe.ReplaceAllString(dataString, replacementModel)
	newData = interfaceRe.ReplaceAllString(newData, replacementInterface)
	formattedOut := strings.Replace(newData, `\n`, "\n", -1)
	ioutil.WriteFile(fileName, []byte(formattedOut), 0755)

	// jsonParsed, err := gabs.ParseJSON([]byte(s))
	// fmt.Println(jsonParsed)
}
