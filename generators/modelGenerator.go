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
	importsModel := "import {Schema, model} from 'mongoose';\n\n\n"
	var elements string
	for i, s := range entity.TypeDetails {
		if !s.Array {
			if i != len(entity.TypeDetails)-1 {
				elements += fmt.Sprintf("\t%s: {type: %s, required: %t},\n", s.Name, strings.Title(s.DataType), s.Mandatory)	
			} else {
				elements += fmt.Sprintf("\t%s: {type: %s, required: %t}\n", s.Name, strings.Title(s.DataType), s.Mandatory)	
			}
		} else {
			if i != len(entity.TypeDetails)-1 {
				elements += fmt.Sprintf("\t%s: {type: Array, required: %t},\n", s.Name, s.Mandatory)	
			} else {
				elements += fmt.Sprintf("\t%s: {type: Array, required: %t}\n", s.Name, s.Mandatory)
			}
		}
	}
	schema := fmt.Sprintf("const %sSchema = new Schema({\n%s})\n\n", strings.Title(entity.ModelName), elements)
	exports := fmt.Sprintf("export default model<I%s>(\"%s\", %sSchema)", strings.Title(entity.ModelName), strings.Title(entity.ModelName), strings.Title(entity.ModelName))
	writeModel := fmt.Sprintf("%s%s%s", importsModel, schema, exports)
	modelFileName := fmt.Sprintf("%s.model.ts", entity.ModelName)
	ioutil.WriteFile(modelFileName, []byte(writeModel), 0755)
}

func UpdateModel(entity dTypes.Entity) {
	modelFile := fmt.Sprintf("%s.model.ts", entity.ModelName)
	modelDat, err := ioutil.ReadFile(modelFile)
	if err != nil {
		fmt.Println(err)
	}
	modelRe, _ := regexp.Compile(`new Schema\((\{[^]]+})`)
	modelString := string(modelDat)
	modelMatch := modelRe.FindStringSubmatch(modelString)[1]
	var wordRe = regexp.MustCompile(`(\w+)`)
	modelJson := wordRe.ReplaceAllString(modelMatch, `"$1"`)
	if err != nil {
		fmt.Println(err)
	}
	var model map[string]map[string]interface{}
	
	err = json.Unmarshal([]byte(modelJson), &model)
	
	if err != nil {
		fmt.Println(err)
	}

	for _, s := range entity.TypeDetails {
		if !s.Array {
			if val, ok := model[s.Name]; ok {
				
				val["type"] = strings.Title(s.DataType)
				val["required"] = s.Mandatory
			} else {
				model[s.Name] = map[string]interface{}{
					"required": s.Mandatory,
					"type":     strings.Title(s.DataType),
				}
			}
		} else {
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

	b, _ := json.Marshal(model)
	reg := regexp.MustCompile(`"([^"]*)"`)
	unquotedModel := reg.ReplaceAllString(string(b), `${1}`)
	reg = regexp.MustCompile(`},`)
	formattedModel := reg.ReplaceAllString(unquotedModel, `}, \n`)
	reg = regexp.MustCompile("^(.*?){")
	formattedModel = reg.ReplaceAllString(formattedModel, `{ \n`)
	reg = regexp.MustCompile(`}}`)
	formattedModel = reg.ReplaceAllString(formattedModel, `} \n}`)
	replacementModel := fmt.Sprintf("new Schema(%s", formattedModel)
	newDataModel := modelRe.ReplaceAllString(modelString, replacementModel)
	formattedOutModel := strings.Replace(newDataModel, `\n`, "\n\t", -1)
	formattedOutModel = strings.Replace(formattedOutModel, `,`, ", ", -1)
	formattedOutModel = strings.Replace(formattedOutModel, `:`, ": ", -1)
	ioutil.WriteFile(modelFile, []byte(formattedOutModel), 0755)
	// jsonParsed, err := gabs.ParseJSON([]byte(s))
	// fmt.Println(jsonParsed)
}
