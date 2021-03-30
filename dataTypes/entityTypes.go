package dTypes

type VarType string

const (
	ObjectId VarType = "objectId"
	Number   VarType = "number"
	DateTime VarType = "dateTime"
	String   VarType = "string"
	Array    VarType = "array"
)

type TypeDetails struct {
	DataType  string `json:"dType"`
	Mandatory bool   `json:"mandatory"`
	Name      string `json:"name"`
	Array     bool   `json:"array"`
}
type Entity struct {
	ModelName   string        `json:"modelName"`
	TypeDetails []TypeDetails `json:"types"`
	Desc        string        `json:"desc"`
}
