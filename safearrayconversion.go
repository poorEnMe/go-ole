// Helper for converting SafeArray to array of objects.

package ole

import (
	"unsafe"
)

type SafeArrayConversion struct {
	Array *SafeArray
}

func (sac *SafeArrayConversion) ToStringArray() (strings []string) {
	totalElements, _ := sac.TotalElements(0)
	strings = make([]string, totalElements)

	for i := int32(0); i < totalElements; i++ {
		strings[int32(i)], _ = safeArrayGetElementString(sac.Array, i)
	}

	return
}

func (sac *SafeArrayConversion) ToByteArray() (bytes []byte) {
	totalElements, _ := sac.TotalElements(0)
	bytes = make([]byte, totalElements)

	for i := int32(0); i < totalElements; i++ {
		safeArrayGetElement(sac.Array, i, unsafe.Pointer(&bytes[int32(i)]))
	}

	return
}

func (sac *SafeArrayConversion) ToValueArray() (values []interface{}) {
	totalElements, _ := sac.TotalElements(0)
	values = make([]interface{}, totalElements)
	vt, _ := safeArrayGetVartype(sac.Array)

	for i := int32(0); i < totalElements; i++ {
		switch VT(vt) {
		case VT_BOOL:
			var v bool
			safeArrayGetElement(sac.Array, i, unsafe.Pointer(&v))
			values[i] = v
		case VT_I1:
			var v int8
			safeArrayGetElement(sac.Array, i, unsafe.Pointer(&v))
			values[i] = v
		case VT_I2:
			var v int16
			safeArrayGetElement(sac.Array, i, unsafe.Pointer(&v))
			values[i] = v
		case VT_I4:
			var v int32
			safeArrayGetElement(sac.Array, i, unsafe.Pointer(&v))
			values[i] = v
		case VT_I8:
			var v int64
			safeArrayGetElement(sac.Array, i, unsafe.Pointer(&v))
			values[i] = v
		case VT_UI1:
			var v uint8
			safeArrayGetElement(sac.Array, i, unsafe.Pointer(&v))
			values[i] = v
		case VT_UI2:
			var v uint16
			safeArrayGetElement(sac.Array, i, unsafe.Pointer(&v))
			values[i] = v
		case VT_UI4:
			var v uint32
			safeArrayGetElement(sac.Array, i, unsafe.Pointer(&v))
			values[i] = v
		case VT_UI8:
			var v uint64
			safeArrayGetElement(sac.Array, i, unsafe.Pointer(&v))
			values[i] = v
		case VT_R4:
			var v float32
			safeArrayGetElement(sac.Array, i, unsafe.Pointer(&v))
			values[i] = v
		case VT_R8:
			var v float64
			safeArrayGetElement(sac.Array, i, unsafe.Pointer(&v))
			values[i] = v
		case VT_BSTR:
			v, _ := safeArrayGetElementString(sac.Array, i)
			values[i] = v
		case VT_VARIANT:
			var v VARIANT
			safeArrayGetElement(sac.Array, i, unsafe.Pointer(&v))
			values[i] = v.Value()
			v.Clear()
		default:
			// TODO
		}
	}

	return
}

func (sac *SafeArrayConversion) GetType() (varType uint16, err error) {
	return safeArrayGetVartype(sac.Array)
}

func (sac *SafeArrayConversion) GetDimensions() (dimensions *uint32, err error) {
	return safeArrayGetDim(sac.Array)
}

func (sac *SafeArrayConversion) GetSize() (length *uint32, err error) {
	return safeArrayGetElementSize(sac.Array)
}

func (sac *SafeArrayConversion) TotalElements(index uint32) (totalElements int32, err error) {
	if index < 1 {
		index = 1
	}

	// Get array bounds
	var LowerBounds int32
	var UpperBounds int32

	LowerBounds, err = safeArrayGetLBound(sac.Array, index)
	if err != nil {
		return
	}

	UpperBounds, err = safeArrayGetUBound(sac.Array, index)
	if err != nil {
		return
	}

	totalElements = UpperBounds - LowerBounds + 1
	return
}

// Release Safe Array memory
func (sac *SafeArrayConversion) Release() {
	safeArrayDestroy(sac.Array)
}

func (sac *SafeArrayConversion) ToVariantArray() (values []*VARIANT) {
	vt, _ := safeArrayGetVartype(sac.Array)
	if VT(vt) != VT_VARIANT {
		return nil
	}

	totalElements, _ := sac.TotalElements(0)
	values = make([]*VARIANT, totalElements)

	for i := int32(0); i < totalElements; i++ {
		var v VARIANT
		safeArrayGetElement(sac.Array, i, unsafe.Pointer(&v))
		values[i] = &v
	}
	return values
}

// To2DVariantArray invaild，尚未实现
func (sac *SafeArrayConversion) To2DVariantArray() (values [][]*VARIANT) {
	vt, _ := safeArrayGetVartype(sac.Array)
	if VT(vt) != VT_VARIANT {
		return nil
	}

	//dis, _ := sac.GetDimensions()

	firstElements, _ := sac.TotalElements(0)
	secondElements, _ := sac.TotalElements(2)
	values = make([][]*VARIANT, firstElements)

	for i := int32(0); i < firstElements; i++ {
		values[i] = make([]*VARIANT, secondElements)
		for j := int32(0); j < secondElements; j++ {
			var v VARIANT
			safeArray2DGetElement(sac.Array, safeArrayFromInt32Slice([]int32{i, j}), unsafe.Pointer(&v))
			values[i][j] = &v
		}
	}
	return values
}
