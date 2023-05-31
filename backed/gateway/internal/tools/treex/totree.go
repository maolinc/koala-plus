package treex

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// 可能的错误常量
const (
	inSliceErrMustValidSlice       = "参数一(inSlice) 必须是一个不为空值的结构体切片"
	destSlicePtrErrMustPtr         = "参数二(destSlicePtr) 必须是一个指针"
	destSlicePtrErrMustSlice       = "参数二(destSlicePtr) 必须是一个结构体切片的指针"
	structErrMustPrimaryKey        = "每级结构体必须设置 primaryKey 标签，指定每一层结构体的主键"
	structErrMustFid               = "子结构体必须设置 fid 标签，指定的父级键名也必须存在,错误的键名："
	structPrimaryKeyMustUpper      = "结构体主键字段（primaryKey 标签所在键）必须以大写字母开头（扫描原始数据时本程序需要修改主键字段值的权限）"
	destStructFieldNotExists       = "参数二(destSlicePtr) 结构体定义的字段（包括类型）必须在 inSlice 参数传递的结构体中存在，否则无法获取值,请检查字段名称书写是否有误，错误字段名："
	destStructFidFieldNotExists    = "子结构体设置的 fid 标签对应的字段不存在于 inSlice 参数传递的结构体中存在,fid设置的错误字段名："
	allowMaxRows                   = 100000
	overAllowMaxRows               = "程序遍历次数已经超过了 100000 次,可能已经选入了死循环,请检查传入的数据是否符合要求,引起该错误的原因 primaryKey 标签字段和 fid 标签字段值不要出现互相嵌套，此外保证 primaryKey 标签键值唯一"
	primaryKeyDataTypeError        = "主键数据类型必须是 int系列(int、int16、int32、int64) 和 string 类型，其他类型不支持,错误发生的主键为： "
	subKeyDataTypeIsNotIntError    = "原始切片中, 主键、子健数据类型不一致，请确保子键数据类型为 int 系列, 相关子键(外键)："
	subKeyDataTypeIsNotStringError = "原始切片中, 主键、子健数据类型不一致，请确保子键数据类型为 string 系列, 相关子键(外键)："
)

func CreateTreeFactory() *resFormatTree {
	pk := make(map[string]string)
	return &resFormatTree{primaryKey: pk, counts: 0}
}

type resFormatTree struct {
	primaryKey     map[string]string
	counts         int // 统计程序秭归计算的次数
	inSliceValueOf reflect.Value
	inSliceLen     int
}

func (s *resFormatTree) ScanToTreeData(inSlice interface{}, destSlicePtr interface{}) (err error) {

	inTypeOf := reflect.TypeOf(inSlice)
	if inTypeOf.Kind() != reflect.Slice {
		return errors.New(inSliceErrMustValidSlice)
	}

	inValueOf := reflect.ValueOf(inSlice)

	s.inSliceValueOf = inValueOf // 原始值的 valueOf 存储起来
	s.inSliceLen = inValueOf.Len()

	inLen := inValueOf.Len()
	if inLen == 0 {
		return
	}

	destTypeOf := reflect.TypeOf(destSlicePtr)
	if destTypeOf.Kind() != reflect.Ptr {
		return errors.New(destSlicePtrErrMustPtr)
	}

	// 根据 dest 参数的指针获取对应的切片元素
	destValueOf := reflect.ValueOf(destSlicePtr).Elem()

	if destValueOf.Type().Kind() != reflect.Slice {
		return errors.New(destSlicePtrErrMustSlice)
	}

	// 创建一个与 dest 切片相同结构的切片，存储中间过程数据，最后载赋值给 dest 切片即可
	destTmpSlice := reflect.MakeSlice(destValueOf.Type(), 0, 0)

	// type 获取的 slice，继续获取内部的结构体元素
	destStructElem := destValueOf.Type().Elem()

	primaryKeyName := s.getCurStructPrimaryKeyName(destStructElem)
	if primaryKeyName == "" {
		return errors.New(structErrMustPrimaryKey)
	} else {
		s.storePrimaryKey(primaryKeyName)
	}
	// 返回结构体指针
	tmpDestStructElem := reflect.New(destStructElem)

	structElemTypeOf := tmpDestStructElem.Elem().Type() //  相当于  typeOf
	structElemValueOf := tmpDestStructElem.Elem()       //  相当于  valueOf
	//分析第一层结构体字段
	fieldNum := structElemTypeOf.NumField()

	// 主键的数据类型
	var primaryKeyDataType int
	var primaryKeyIdInt int64
	var primaryKeyIdStr string
	var primaryKeyIdInterf interface{}

	//遍历结果集的每一行数据
	for rowIndex := 0; rowIndex < inLen; rowIndex++ {
		s.counts++
		// 获取正在遍历的当前行数据
		row := inValueOf.Index(rowIndex)
		if !s.destStructFieldIsExists(row.Type(), primaryKeyName) {
			return errors.New(destStructFieldNotExists + primaryKeyName)
		}
		// 根据dest切片中的元素(结构体),所定义的主键查询 row 中对应的字段数据类型，必须满足 数字系列、字符串系列，否则返回错误
		if primaryKeyDataType, err = s.curPrimaryKeyDataType(row, primaryKeyName); err != nil {
			return err
		}

		mainKeyField := row.FieldByName(primaryKeyName)

		switch primaryKeyDataType {
		case 1:
			primaryKeyIdInt = mainKeyField.Int()
			if primaryKeyIdInt > 0 {
				primaryKeyIdInterf = primaryKeyIdInt
			} else {
				continue
			}
		case 2:
			primaryKeyIdStr = mainKeyField.String()
			if strings.ReplaceAll(primaryKeyIdStr, " ", "") != "" {
				primaryKeyIdInterf = primaryKeyIdStr
			} else {
				continue
			}
		}
		if primaryKeyIdInterf != nil {
			for i := 0; i < fieldNum; i++ {
				if destStructElem.Field(i).Name == "Children" && destStructElem.Field(i).Type.Kind() == reflect.Slice {
					if val, err := s.analysisChildren(int64(rowIndex), row, destStructElem.Field(i).Type); err == nil {
						structElemValueOf.Field(i).Set(val)
					} else {
						return err
					}
				} else {
					// dest 接受字段名称和类型与 切片结果遍历中的某一条，
					// 必须是字段名和数据类型相同，则可以赋值
					if s.destStructFieldIsSame(row.Type(), structElemTypeOf.Field(i)) {
						structElemValueOf.Field(i).Set(row.FieldByName(destStructElem.Field(i).Name))
					} else {
						//如果目的结构体的字段不存在于原始数据结构体中，那么寻找 default 标签对应的默认值进行赋值； 否则跳过
						if val, ok := s.setFieldDefaultValue(structElemTypeOf, structElemTypeOf.Field(i).Name); ok {
							structElemValueOf.Field(i).Set(val)
						}
					}
				}
			}
			destTmpSlice = reflect.Append(destTmpSlice, structElemValueOf)
		}
	}
	destValueOf.Set(destTmpSlice)
	return nil
}

func (s *resFormatTree) storePrimaryKey(keyName string) {
	if s.primaryKey[keyName] != keyName {
		s.primaryKey[keyName] = keyName
	}
}

// 获取正在分析(处理)的结构体主键键名（PrimaryKey）
func (s *resFormatTree) getCurStructPrimaryKeyName(rTypeOf reflect.Type) string {
	numField := rTypeOf.NumField()
	for i := 0; i < numField; i++ {
		if rTypeOf.Field(i).Tag.Get("primaryKey") == "yes" {
			return rTypeOf.Field(i).Name
		}
	}
	return ""
}

// 获取当前结构体对应的父级键名（父子关系键中的父级键,也就是 fid 标签的值）
func (s *resFormatTree) getCurStructParentKeyName(value reflect.Type) string {
	numField := value.NumField()
	for i := 0; i < numField; i++ {
		parentKey := value.Field(i).Tag.Get("fid")
		if parentKey != "" {
			return parentKey
		}
	}
	return ""
}

// 获取当前结构体的 子外键名 （父子关系键中的子外键键,也就是 fid 标签所在的字段名称）
func (s *resFormatTree) getCurStructSubFKeyName(value reflect.Type) string {
	numField := value.NumField()
	for i := 0; i < numField; i++ {
		if _, ok := value.Field(i).Tag.Lookup("fid"); ok {
			return value.Field(i).Name
		}
	}
	return ""
}

// 判断 dest 结构体中的字段是否在 inSlice 参数中的结构体中存在
func (s *resFormatTree) destStructFieldIsExists(inSliceStruct reflect.Type, destFieldStructName string) bool {
	num := inSliceStruct.NumField()
	for i := 0; i < num; i++ {
		if inSliceStruct.Field(i).Name == destFieldStructName {
			return true
		}
	}
	return false
}

// 判断 dest 结构体中的字段是否在 inSlice 参数中的结构体中存在(字段名称+数据类型相同)
func (s *resFormatTree) destStructFieldIsSame(inSliceStruct reflect.Type, destFieldStruct reflect.StructField) bool {
	num := inSliceStruct.NumField()
	for i := 0; i < num; i++ {
		if inSliceStruct.Field(i).Name == destFieldStruct.Name && inSliceStruct.Field(i).Type == destFieldStruct.Type {
			return true
		}
	}
	return false
}

// 继续分析 children 结构体
// 参数解释
// 1.parentRowIndex ： 正在遍历的结果集的当前行号
// 2.parentField ：正在遍历的结果集的当前行结构体（valueOf）
// 3.childrenType : dest结构体中的 Children 字段类型(typeOf)，本质上就是一个切片类型

func (s *resFormatTree) analysisChildren(parentRowIndex int64, parentField reflect.Value, childrenType reflect.Type) (reflect.Value, error) {
	s.counts++
	resChildren := reflect.MakeSlice(childrenType, 0, 0)

	if s.counts > allowMaxRows {
		return resChildren, errors.New(overAllowMaxRows)
	}

	vType := childrenType.Elem()
	newStruct := reflect.New(vType)

	// 分析 Children 切片中的元素(结构体)
	newTypeOf := newStruct.Elem().Type()
	newValueOf := newStruct.Elem()
	fieldNum := newTypeOf.NumField()

	// 获取当前结构体中，某个字段定义的fid对应的父键名称
	parentKeyName := s.getCurStructParentKeyName(newTypeOf)
	if parentKeyName == "" {
		return reflect.Value{}, errors.New(structErrMustFid)
	}

	// 获取当前结构体的主键（primaryKey 标签）所在的字段名称
	curStructPrimaryKeyName := s.getCurStructPrimaryKeyName(newTypeOf)
	if curStructPrimaryKeyName == "" {
		return reflect.Value{}, errors.New(structErrMustPrimaryKey)
	}

	// 子级结构体中定义的外键（fid标签设置的父键），必须在父级字段中存在
	// 这样才能形成  父---子 数据关联关系
	if !s.destStructFieldIsExists(parentField.Type(), parentKeyName) {
		return reflect.Value{}, errors.New(destStructFidFieldNotExists + parentKeyName)
	}

	s.storePrimaryKey(curStructPrimaryKeyName)

	var ParentDataType int
	var ParentIdInt int64
	var ParentIdStr string
	var err error
	// 判断主键的数据类型
	if ParentDataType, err = s.curPrimaryKeyDataType(parentField, parentKeyName); err == nil {
		switch ParentDataType {
		// 主键为 int 系列
		case 1:
			ParentIdInt = parentField.FieldByName(parentKeyName).Int()
			if ParentIdInt > 0 {
				for subRowIndex := int(parentRowIndex); subRowIndex < s.inSliceLen; subRowIndex++ {
					subRow := s.inSliceValueOf.Index(subRowIndex)

					// 获取children切片中的结构体元素 fid 所在标签的字段名
					// 对于上层结构体来说，就是外键字段名
					subFKeyName := s.getCurStructSubFKeyName(newTypeOf)
					if subFKeyName == "" {
						return reflect.Value{}, errors.New(structErrMustFid + subFKeyName)
					}

					subFKeyField := subRow.FieldByName(subFKeyName)
					// 获取children切片中的结构体元素中， primaryKey 所在的标签对应的字段名，即 主键字段名
					subPrimaryKeyName := s.getCurStructParentKeyName(newTypeOf)
					if subPrimaryKeyName == "" {
						return reflect.Value{}, errors.New(structErrMustPrimaryKey)
					}

					// 子集 数据中的主键
					subPrimaryKeyField := subRow.FieldByName(subPrimaryKeyName)

					//判断子健的数据类型必须为 int 系列
					if val, _ := s.curPrimaryKeyDataType(subRow, subFKeyName); val != 1 {
						return reflect.Value{}, errors.New(subKeyDataTypeIsNotIntError + subFKeyName)
					}
					//相对父级行来说，就是子外键的值
					subFKeyId := subFKeyField.Int()

					s.storePrimaryKey(subPrimaryKeyName)

					if subFKeyId > 0 && subFKeyId == ParentIdInt && !subPrimaryKeyField.IsZero() {
						if tmpChildren, err := s.getLevelGe2Children(fieldNum, resChildren, newTypeOf, parentRowIndex, subRowIndex, ParentIdInt, subFKeyName, subPrimaryKeyName, subRow, newValueOf); err != nil {
							return reflect.Value{}, err
						} else {
							resChildren = tmpChildren
						}
					}
				}
			}
		//  字符串系列
		case 2:
			ParentIdStr = parentField.FieldByName(parentKeyName).String()
			if ParentIdStr != "" {
				for subRowIndex := int(parentRowIndex); subRowIndex < s.inSliceLen; subRowIndex++ {
					subRow := s.inSliceValueOf.Index(subRowIndex)

					// 获取children切片中的结构体元素 fid 所在标签的字段名
					// 对于上层结构体来说，就是外键字段名
					subFKeyName := s.getCurStructSubFKeyName(newTypeOf)
					if subFKeyName == "" {
						return reflect.Value{}, errors.New(structErrMustFid + subFKeyName)
					}

					subFKeyField := subRow.FieldByName(subFKeyName)
					// 获取children切片中的结构体元素中， primaryKey 所在的标签对应的字段名，即 主键字段名
					subPrimaryKeyName := s.getCurStructPrimaryKeyName(newTypeOf)
					if subPrimaryKeyName == "" {
						return reflect.Value{}, errors.New(structErrMustPrimaryKey)
					}

					subPrimaryKeyField := subRow.FieldByName(subPrimaryKeyName)

					//判断子健的数据类型必须为 string 系列
					if val, _ := s.curPrimaryKeyDataType(subRow, subFKeyName); val != 2 {
						return reflect.Value{}, errors.New(subKeyDataTypeIsNotStringError + subFKeyName)
					}

					//相对父级行来说，就是子外键的值
					subFKeyId := subFKeyField.String()

					s.storePrimaryKey(subPrimaryKeyName)

					if subFKeyId != "" && subFKeyId == ParentIdStr && !subPrimaryKeyField.IsZero() {
						if tmpChildren, err := s.getLevelGe2Children(fieldNum, resChildren, newTypeOf, parentRowIndex, subRowIndex, ParentIdStr, subFKeyName, subPrimaryKeyName, subRow, newValueOf); err != nil {
							return reflect.Value{}, err
						} else {
							resChildren = tmpChildren
						}
					}
				}
			}
		}
	}

	return resChildren, nil
}

// 获取 ≥ 2层级的 children 切片的内容
// 参数解释
// fieldNum 正在等待填充的 chilren 切片元素中的结构体字段数
// resChildren 正在等待填充的 chilren 切片初始化后的变量
// newTypeOf 正在等待填充的 chilren 切片元素中的结构体的 typeof
// parentRowIndex 正在提取 chilren中数据时的当前行号
// subRowIndex 正在遍历的子级数据行号
// ParentId 正在提取 chilren中数据时的主键id(interface类型，主要是 int、string类型)
// subFKeyName 正在遍历的子级数据外键名
// subPrimaryKeyName 正在遍历的子级数据主键名
// subRow 正在遍历的子级数据行
// newValueOf 正在等待填充的 chilren 切片元素中的结构体的 valueof
func (s *resFormatTree) getLevelGe2Children(fieldNum int, resChildren reflect.Value, newTypeOf reflect.Type, parentRowIndex int64, subRowIndex int, ParentId interface{}, subFKeyName, subPrimaryKeyName string, subRow, newValueOf reflect.Value) (reflect.Value, error) {
	for j := 0; j < fieldNum; j++ {
		if newTypeOf.Field(j).Type.Kind() == reflect.Slice && newTypeOf.Field(j).Name == "Children" {
			if s.curItemHasSubLists(parentRowIndex, ParentId, subFKeyName) {
				if dataType, err := s.curPrimaryKeyDataType(subRow, subPrimaryKeyName); err == nil {
					switch dataType {
					case 1:
						if val, err := s.analysisChildren(int64(subRowIndex), subRow, newTypeOf.Field(j).Type); err == nil {
							newValueOf.Field(j).Set(val)
						} else {
							return reflect.Value{}, err
						}
					case 2:
						if val, err := s.analysisChildren(int64(subRowIndex), subRow, newTypeOf.Field(j).Type); err == nil {
							newValueOf.Field(j).Set(val)
						} else {
							return reflect.Value{}, err
						}
					}
				}
			} else {
				return resChildren, nil
			}
		} else {
			if s.destStructFieldIsExists(subRow.Type(), newTypeOf.Field(j).Name) {
				newValueOf.Field(j).Set(subRow.FieldByName(newTypeOf.Field(j).Name))
			} else {
				//如果目的结构体的字段不存在于原始数据结构体中，那么寻找 default 标签对应的默认值进行赋值； 否则跳过
				if val, ok := s.setFieldDefaultValue(newTypeOf, newTypeOf.Field(j).Name); ok {
					newValueOf.Field(j).Set(val)
				}
			}
		}
	}
	if err := s.setUsedKeyInvalid(subRowIndex); err != nil {
		return reflect.Value{}, err
	}
	resChildren = reflect.Append(resChildren, newValueOf)
	return resChildren, nil
}

// 设置已经遍历过的所有主键失效
func (s *resFormatTree) setUsedKeyInvalid(subRowIndex int) error {
	for primaryKeyName := range s.primaryKey {
		if s.inSliceValueOf.Index(subRowIndex).FieldByName(primaryKeyName).CanSet() {
			if dataType, err := s.curPrimaryKeyDataType(s.inSliceValueOf.Index(subRowIndex), primaryKeyName); err == nil {
				switch dataType {
				case 1:
					s.inSliceValueOf.Index(subRowIndex).FieldByName(primaryKeyName).SetInt(0)
				case 2:
					s.inSliceValueOf.Index(subRowIndex).FieldByName(primaryKeyName).SetString("")
				}
			} else {
				return errors.New(primaryKeyDataTypeError + primaryKeyName)
			}
		} else {
			return errors.New(structPrimaryKeyMustUpper)
		}

	}
	return nil
}

// 针对目的结构体中不存在的字段，根据tag标签设置的 default值进行默认赋值
func (s *resFormatTree) setFieldDefaultValue(fieldType reflect.Type, fieldName string) (reflect.Value, bool) {
	if f, ok := fieldType.FieldByName(fieldName); ok {
		if val, ok2 := f.Tag.Lookup("default"); ok2 {
			switch f.Type.Kind() {
			case reflect.String:
				return reflect.ValueOf(val), true
			case reflect.Float32:
				if tmp, err := strconv.ParseFloat(val, 32); err == nil {
					return reflect.ValueOf(tmp), true
				}
			case reflect.Float64:
				if tmp, err := strconv.ParseFloat(val, 64); err == nil {
					return reflect.ValueOf(tmp), true
				}
			case reflect.Int:
				if tmp, err := strconv.Atoi(val); err == nil {
					return reflect.ValueOf(tmp), true
				}
			case reflect.Int32:
				if tmp, err := strconv.ParseInt(val, 10, 32); err == nil {
					return reflect.ValueOf(tmp), true
				}
			case reflect.Int64:
				if tmp, err := strconv.ParseInt(val, 10, 64); err == nil {
					return reflect.ValueOf(tmp), true
				}
			case reflect.Bool:
				if tmp, err := strconv.ParseBool(val); err == nil {
					return reflect.ValueOf(tmp), true
				}
			default:
				break
			}
		}
	}
	return reflect.Value{}, false
}

// 判断当前行底下是否还有挂接子级数据
// 参数解释
// 1.curIndex ：    结果集循环的当前行号
// 2.curMainId ：   结果集循环的当前结构体的主键ID
// 3.subFKeyName：  结果集循环的当前结构体的主键ID对应的子级外键名称

func (s *resFormatTree) curItemHasSubLists(curIndex int64, curMainId interface{}, subFKeyName string) (res bool) {
	for i := int(curIndex); i <= s.inSliceLen-1; i++ {
		tmpField := s.inSliceValueOf.Index(i)
		if pDataType, err := s.curPrimaryKeyDataType(tmpField, subFKeyName); err == nil {
			switch pDataType {
			case 1:
				subFKeyValue := s.inSliceValueOf.Index(i).FieldByName(subFKeyName).Int()
				if curMainId.(int64) == subFKeyValue {
					return true
				}
			case 2:
				subFKeyValue := s.inSliceValueOf.Index(i).FieldByName(subFKeyName).String()
				if curMainId.(string) == subFKeyValue {
					return true
				}
			}
		}

	}
	return res
}

// 判断当前主键数据类型是否为 数字类型 ( int unit  int16  int32  int64 )
func (s *resFormatTree) curPrimaryKeyDataType(rValue reflect.Value, keyName string) (int, error) {
	switch rValue.FieldByName(keyName).Kind() {
	case reflect.Int64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return 1, nil
	case reflect.String:
		return 2, nil
	default:

	}
	return 0, errors.New(primaryKeyDataTypeError + keyName)
}
