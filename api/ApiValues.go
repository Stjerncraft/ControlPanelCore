package api

import "strconv"

type ApiValue interface {
	GetType() ApiType
	GetValueAsString() string
	SetValueFromString(str string) error
}

//--
type ApiValueByte struct {
	Value int8
}

var byteType = ApiTypeByte{}

func (t *ApiValueByte) GetType() ApiType {
	return &byteType
}
func (t *ApiValueByte) GetValueAsString() string {
	return strconv.Itoa(int(t.Value))
}
func (t *ApiValueByte) SetValueFromString(str string) error {
	val, err := strconv.ParseInt(str, 10, 8)
	t.Value = int8(val)
	return err
}

//--
type ApiValueShort struct {
	Value int16
}

var shortType = ApiTypeShort{}

func (t *ApiValueShort) GetType() ApiType {
	return &shortType
}
func (t *ApiValueShort) GetValueAsString() string {
	return strconv.Itoa(int(t.Value))
}
func (t *ApiValueShort) SetValueFromString(str string) error {
	val, err := strconv.ParseInt(str, 10, 16)
	t.Value = int16(val)
	return err
}

//--
type ApiValueInteger struct {
	Value int32
}

var integerType = ApiTypeInteger{}

func (t *ApiValueInteger) GetType() ApiType {
	return &integerType
}
func (t *ApiValueInteger) GetValueAsString() string {
	return strconv.Itoa(int(t.Value))
}
func (t *ApiValueInteger) SetValueFromString(str string) error {
	val, err := strconv.ParseInt(str, 10, 32)
	t.Value = int32(val)
	return err
}

//--
type ApiValueLong struct {
	Value int64
}

var longType = ApiTypeLong{}

func (t *ApiValueLong) GetType() ApiType {
	return &longType
}
func (t *ApiValueLong) GetValueAsString() string {
	return strconv.FormatInt(t.Value, 10)
}
func (t *ApiValueLong) SetValueFromString(str string) error {
	var err error
	t.Value, err = strconv.ParseInt(str, 10, 64)
	return err
}

//--
type ApiValueFloat struct {
	Value float32
}

var floatType = ApiTypeFloat{}

func (t *ApiValueFloat) GetType() ApiType {
	return &floatType
}
func (t *ApiValueFloat) GetValueAsString() string {
	return strconv.FormatFloat(float64(t.Value), 'f', -1, 32)
}
func (t *ApiValueFloat) SetValueFromString(str string) error {
	val, err := strconv.ParseFloat(str, 32)
	t.Value = float32(val)
	return err
}

//--
type ApiValueDouble struct {
	Value float64
}

var doubleType = ApiTypeDouble{}

func (t *ApiValueDouble) GetType() ApiType {
	return &doubleType
}
func (t *ApiValueDouble) GetValueAsString() string {
	return strconv.FormatFloat(t.Value, 'f', -1, 64)
}
func (t *ApiValueDouble) SetValueFromString(str string) error {
	var err error
	t.Value, err = strconv.ParseFloat(str, 64)
	return err
}

//--
type ApiValueString struct {
	Value string
}

var stringType = ApiTypeString{}

func (t *ApiValueString) GetType() ApiType {
	return &stringType
}
func (t *ApiValueString) GetValueAsString() string {
	return t.Value
}
func (t *ApiValueString) SetValueFromString(str string) error {
	t.Value = str
	return nil
}

//--
type ApiValueObject struct {
	Type   *ApiTypeObject
	Fields map[string]ApiValue
}

func (t *ApiValueObject) GetType() ApiType {
	return t.Type
}
func (t *ApiValueObject) GetValueAsString() string {
	//TODO: Turn into JSON object
	//TODO: Only send values, not type names, as they are already known.
	//return t.Value
	return ""
}
func (t *ApiValueObject) SetValueFromString(str string) error {
	//TODO: Parse JSON object
	//t.Value = str

	return nil
}

//--
type ApiValueList struct {
	Type   *ApiTypeList
	Values []ApiValue
}

func (t *ApiValueList) GetType() ApiType {
	return t.Type
}
func (t *ApiValueList) GetValueAsString() string {
	//TODO
	return ""
}
func (t *ApiValueList) SetValueFromString(str string) error {
	//TODO:
	return nil
}
