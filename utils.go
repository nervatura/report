package report

import (
	"encoding/json"
	"errors"
	"image/color"
	"strconv"
	"strings"
	"time"
)

// ToString - safe string conversion
func ToString(value interface{}, defValue string) string {
	if stringValue, valid := value.(string); valid {
		if stringValue == "" {
			return defValue
		}
		return stringValue
	}
	if boolValue, valid := value.(bool); valid {
		return strconv.FormatBool(boolValue)
	}
	if intValue, valid := value.(int); valid {
		return strconv.Itoa(intValue)
	}
	if intValue, valid := value.(int32); valid {
		return strconv.Itoa(int(intValue))
	}
	if intValue, valid := value.(int64); valid {
		return strconv.FormatInt(intValue, 10)
	}
	if floatValue, valid := value.(float32); valid {
		return strconv.FormatFloat(float64(floatValue), 'f', -1, 64)
	}
	if floatValue, valid := value.(float64); valid {
		return strconv.FormatFloat(floatValue, 'f', -1, 64)
	}
	if timeValue, valid := value.(time.Time); valid {
		return timeValue.Format("2006-01-02T15:04:05-07:00")
	}
	return defValue
}

// ToFloat - safe float64 conversion
func ToFloat(value interface{}, defValue float64) float64 {
	if floatValue, valid := value.(float64); valid {
		if floatValue == 0 {
			return defValue
		}
		return floatValue
	}
	if boolValue, valid := value.(bool); valid {
		if boolValue {
			return 1
		}
	}
	if intValue, valid := value.(int); valid {
		return float64(intValue)
	}
	if intValue, valid := value.(int32); valid {
		return float64(intValue)
	}
	if intValue, valid := value.(int64); valid {
		return float64(intValue)
	}
	if floatValue, valid := value.(float32); valid {
		return float64(floatValue)
	}
	if stringValue, valid := value.(string); valid {
		floatValue, err := strconv.ParseFloat(stringValue, 64)
		if err == nil {
			return float64(floatValue)
		}
	}
	return defValue
}

// ToRGBA - safe RGBA conversion
func ToRGBA(value interface{}, defValue color.RGBA) color.RGBA {
	if rgbaValue, valid := value.(color.RGBA); valid {
		return rgbaValue
	}
	if c := toRGBAFromString(value); c != nil {
		return *c
	}
	if c := toRGBAFromNumeric(value); c != nil {
		return *c
	}
	return defValue
}

func parseHexColor(v string) (out color.RGBA, err error) {
	if len(v) != 7 {
		return out, errors.New("hex color must be 7 characters")
	}
	red, redError := strconv.ParseUint(v[1:3], 16, 8)
	if redError != nil {
		return out, errors.New("red component invalid")
	}
	out.R = uint8(red)
	green, greenError := strconv.ParseUint(v[3:5], 16, 8)
	if greenError != nil {
		return out, errors.New("green component invalid")
	}
	out.G = uint8(green)
	blue, blueError := strconv.ParseUint(v[5:7], 16, 8)
	if blueError != nil {
		return out, errors.New("blue component invalid")
	}
	out.B = uint8(blue)
	return
}

func toRGBAFromString(value interface{}) *color.RGBA {
	stringValue, valid := value.(string)
	if !valid {
		return nil
	}
	if strings.HasPrefix(stringValue, "#") {
		pvalue, err := parseHexColor(stringValue)
		if err == nil {
			return &pvalue
		}
		return nil
	}
	ivalue := ToInteger(value, -1)
	if ivalue >= 0 && ivalue < 255 {
		c := color.RGBA{uint8(ivalue), uint8(ivalue), uint8(ivalue), 0}
		return &c
	}
	return nil
}

func toRGBAFromNumeric(value interface{}) *color.RGBA {
	v, ok := toUint8FromNumeric(value)
	if !ok {
		return nil
	}
	c := color.RGBA{v, v, v, 0}
	return &c
}

func toUint8FromNumeric(value interface{}) (uint8, bool) {
	var n float64
	switch v := value.(type) {
	case int:
		n = float64(v)
	case int32:
		n = float64(v)
	case int64:
		n = float64(v)
	case float32:
		n = float64(v)
	case float64:
		n = v
	default:
		return 0, false
	}
	if n >= 0 && n < 255 {
		return uint8(n), true
	}
	return 0, false
}

// ToInteger - safe int64 conversion
func ToInteger(value interface{}, defValue int64) int64 {
	if intValue, valid := value.(int64); valid {
		if intValue == 0 {
			return defValue
		}
		return intValue
	}
	if boolValue, valid := value.(bool); valid {
		if boolValue {
			return 1
		}
	}
	if intValue, valid := value.(int); valid {
		return int64(intValue)
	}
	if intValue, valid := value.(int32); valid {
		return int64(intValue)
	}
	if floatValue, valid := value.(float32); valid {
		return int64(floatValue)
	}
	if floatValue, valid := value.(float64); valid {
		return int64(floatValue)
	}
	if stringValue, valid := value.(string); valid {
		intValue, err := strconv.ParseInt(stringValue, 10, 64)
		if err == nil {
			return int64(intValue)
		}
	}
	return defValue
}

// ToBoolean - safe bool conversion
func ToBoolean(value interface{}, defValue bool) bool {
	if boolValue, valid := value.(bool); valid {
		return boolValue
	}
	if intValue, valid := value.(int); valid {
		if intValue == 1 {
			return true
		}
	}
	if intValue, valid := value.(int32); valid {
		if intValue == 1 {
			return true
		}
	}
	if intValue, valid := value.(int64); valid {
		if intValue == 1 {
			return true
		}
	}
	if floatValue, valid := value.(float32); valid {
		if floatValue == 1 {
			return true
		}
	}
	if floatValue, valid := value.(float64); valid {
		if floatValue == 1 {
			return true
		}
	}
	if stringValue, valid := value.(string); valid {
		boolValue, err := strconv.ParseBool(stringValue)
		if err == nil {
			return boolValue
		}
	}
	return defValue
}

// Contains tells whether a contains x.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func ConvertToByte(data interface{}) ([]byte, error) {
	//var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Marshal(data)
}

func ConvertFromByte(data []byte, result interface{}) error {
	//var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal(data, result)
}
