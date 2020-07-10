package opcua

import (
	"time"

	ua "github.com/awcullen/opcua"
	uuid "github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// toUaVariant coerces the local Variant to ua.Variant
func toUaVariant(value *Variant) *ua.Variant {
	switch x := value.GetValue().(type) {
	case *Variant_Null:
		return &ua.NilVariant
	case *Variant_Boolean:
		return ua.NewVariantBoolean(x.Boolean)
	case *Variant_SByte:
		return ua.NewVariantSByte(int8(x.SByte))
	case *Variant_Byte:
		return ua.NewVariantByte(uint8(x.Byte))
	case *Variant_Int16:
		return ua.NewVariantInt16(int16(x.Int16))
	case *Variant_UInt16:
		return ua.NewVariantUInt16(uint16(x.UInt16))
	case *Variant_Int32:
		return ua.NewVariantInt32(x.Int32)
	case *Variant_UInt32:
		return ua.NewVariantUInt32(x.UInt32)
	case *Variant_Int64:
		return ua.NewVariantInt64(x.Int64)
	case *Variant_UInt64:
		return ua.NewVariantUInt64(x.UInt64)
	case *Variant_Float:
		return ua.NewVariantFloat(x.Float)
	case *Variant_Double:
		return ua.NewVariantDouble(x.Double)
	case *Variant_String_:
		return ua.NewVariantString(x.String_)
	case *Variant_DateTime:
		return ua.NewVariantDateTime(time.Unix(x.DateTime.Seconds, int64(x.DateTime.Nanos)).UTC())
	case *Variant_Guid:
		y, _ := uuid.FromBytes(x.Guid)
		return ua.NewVariantGUID(y)
	case *Variant_ByteString:
		return ua.NewVariantByteString(ua.ByteString(x.ByteString))
	case *Variant_BooleanArray:
		return ua.NewVariantBooleanArray(x.BooleanArray.Value)
	case *Variant_SByteArray:
		a := make([]int8, len(x.SByteArray.Value))
		for i, v := range x.SByteArray.Value {
			a[i] = int8(v)
		}
		return ua.NewVariantSByteArray(a)
	case *Variant_ByteArray:
		a := make([]uint8, len(x.ByteArray.Value))
		for i, v := range x.ByteArray.Value {
			a[i] = uint8(v)
		}
		return ua.NewVariantByteArray(a)
	case *Variant_Int16Array:
		a := make([]int16, len(x.Int16Array.Value))
		for i, v := range x.Int16Array.Value {
			a[i] = int16(v)
		}
		return ua.NewVariantInt16Array(a)
	case *Variant_UInt16Array:
		a := make([]uint16, len(x.UInt16Array.Value))
		for i, v := range x.UInt16Array.Value {
			a[i] = uint16(v)
		}
		return ua.NewVariantUInt16Array(a)
	case *Variant_Int32Array:
		return ua.NewVariantInt32Array(x.Int32Array.Value)
	case *Variant_UInt32Array:
		return ua.NewVariantUInt32Array(x.UInt32Array.Value)
	case *Variant_Int64Array:
		return ua.NewVariantInt64Array(x.Int64Array.Value)
	case *Variant_UInt64Array:
		return ua.NewVariantUInt64Array(x.UInt64Array.Value)
	case *Variant_FloatArray:
		return ua.NewVariantFloatArray(x.FloatArray.Value)
	case *Variant_DoubleArray:
		return ua.NewVariantDoubleArray(x.DoubleArray.Value)
	case *Variant_StringArray:
		return ua.NewVariantStringArray(x.StringArray.Value)
	case *Variant_DateTimeArray:
		a := make([]time.Time, len(x.DateTimeArray.Value))
		for i, v := range x.DateTimeArray.Value {
			a[i] = time.Unix(v.Seconds, int64(v.Nanos)).UTC()
		}
		return ua.NewVariantDateTimeArray(a)
	case *Variant_GuidArray:
		a := make([]uuid.UUID, len(x.GuidArray.Value))
		for i, v := range x.GuidArray.Value {
			a[i], _ = uuid.FromBytes(v)
		}
		return ua.NewVariantGUIDArray(a)
	case *Variant_ByteStringArray:
		a := make([]ua.ByteString, len(x.ByteStringArray.Value))
		for i, v := range x.ByteStringArray.Value {
			a[i] = ua.ByteString(v)
		}
		return ua.NewVariantByteStringArray(a)
	default:
		return &ua.NilVariant
	}
}

// toVariant coerces the ua.Variant to Variant
func toVariant(value *ua.Variant) *Variant {
	if len(value.ArrayDimensions()) == 0 {
		switch value.Type() {
		case ua.VariantTypeNull:
			return &Variant{Value: &Variant_Null{}}
		case ua.VariantTypeBoolean:
			return &Variant{Value: &Variant_Boolean{Boolean: value.Value().(bool)}}
		case ua.VariantTypeSByte:
			return &Variant{Value: &Variant_SByte{SByte: int32(value.Value().(int8))}}
		case ua.VariantTypeByte:
			return &Variant{Value: &Variant_Byte{Byte: uint32(value.Value().(uint8))}}
		case ua.VariantTypeInt16:
			return &Variant{Value: &Variant_Int16{Int16: int32(value.Value().(int16))}}
		case ua.VariantTypeUInt16:
			return &Variant{Value: &Variant_UInt16{UInt16: uint32(value.Value().(uint16))}}
		case ua.VariantTypeInt32:
			return &Variant{Value: &Variant_Int32{Int32: value.Value().(int32)}}
		case ua.VariantTypeUInt32:
			return &Variant{Value: &Variant_UInt32{UInt32: value.Value().(uint32)}}
		case ua.VariantTypeInt64:
			return &Variant{Value: &Variant_Int64{Int64: value.Value().(int64)}}
		case ua.VariantTypeUInt64:
			return &Variant{Value: &Variant_UInt64{UInt64: value.Value().(uint64)}}
		case ua.VariantTypeFloat:
			return &Variant{Value: &Variant_Float{Float: value.Value().(float32)}}
		case ua.VariantTypeDouble:
			return &Variant{Value: &Variant_Double{Double: value.Value().(float64)}}
		case ua.VariantTypeString:
			return &Variant{Value: &Variant_String_{String_: value.Value().(string)}}
		case ua.VariantTypeDateTime:
			x := value.Value().(time.Time)
			return &Variant{Value: &Variant_DateTime{DateTime: &timestamppb.Timestamp{Seconds: int64(x.Second()), Nanos: int32(x.Nanosecond())}}}
		case ua.VariantTypeGUID:
			x := value.Value().(uuid.UUID)
			return &Variant{Value: &Variant_Guid{Guid: x[:]}}
		case ua.VariantTypeByteString:
			return &Variant{Value: &Variant_ByteString{ByteString: value.Value().([]byte)}}
		default:
			return &Variant{Value: &Variant_Null{}}
		}
	}
	/*
		case ua.VariantTypeBooleanArray:
			return ua.NewVariantWithType(x.BooleanArray.Value, ua.VariantTypeBoolean, []int32{int32(len(x.BooleanArray.Value))})
		case ua.VariantTypeSByteArray:
			a := make([]int8, len(x.SByteArray.Value))
			for i, v := range x.SByteArray.Value {
				a[i] = int8(v)
			}
			return ua.NewVariantWithType(a, ua.VariantTypeSByte, []int32{int32(len(x.SByteArray.Value))})
		case *Variant_ByteArray:
			a := make([]uint8, len(x.ByteArray.Value))
			for i, v := range x.ByteArray.Value {
				a[i] = uint8(v)
			}
			return ua.NewVariantWithType(a, ua.VariantTypeByte, []int32{int32(len(x.ByteArray.Value))})
		case *Variant_Int16Array:
			a := make([]int16, len(x.Int16Array.Value))
			for i, v := range x.Int16Array.Value {
				a[i] = int16(v)
			}
			return ua.NewVariantWithType(a, ua.VariantTypeInt16, []int32{int32(len(x.Int16Array.Value))})
		case *Variant_UInt16Array:
			a := make([]uint16, len(x.UInt16Array.Value))
			for i, v := range x.UInt16Array.Value {
				a[i] = uint16(v)
			}
			return ua.NewVariantWithType(a, ua.VariantTypeUInt16, []int32{int32(len(x.UInt16Array.Value))})
		case *Variant_Int32Array:
			return ua.NewVariantWithType(x.Int32Array.Value, ua.VariantTypeInt32, []int32{int32(len(x.Int32Array.Value))})
		case *Variant_UInt32Array:
			return ua.NewVariantWithType(x.UInt32Array.Value, ua.VariantTypeUInt32, []int32{int32(len(x.UInt32Array.Value))})
		case *Variant_Int64Array:
			return ua.NewVariantWithType(x.Int64Array.Value, ua.VariantTypeInt64, []int32{int32(len(x.Int64Array.Value))})
		case *Variant_UInt64Array:
			return ua.NewVariantWithType(x.UInt64Array.Value, ua.VariantTypeUInt64, []int32{int32(len(x.UInt64Array.Value))})
		case *Variant_FloatArray:
			return ua.NewVariantWithType(x.FloatArray.Value, ua.VariantTypeFloat, []int32{int32(len(x.FloatArray.Value))})
		case *Variant_DoubleArray:
			return ua.NewVariantWithType(x.DoubleArray.Value, ua.VariantTypeDouble, []int32{int32(len(x.DoubleArray.Value))})
		case *Variant_StringArray:
			return ua.NewVariantWithType(x.StringArray.Value, ua.VariantTypeString, []int32{int32(len(x.StringArray.Value))})
		case *Variant_DateTimeArray:
			a := make([]time.Time, len(x.DateTimeArray.Value))
			for i, v := range x.DateTimeArray.Value {
				a[i] = time.Unix(v.Seconds, int64(v.Nanos)).UTC()
			}
			return ua.NewVariantWithType(a, ua.VariantTypeDateTime, []int32{int32(len(x.DateTimeArray.Value))})
		case *Variant_GuidArray:
			a := make([]uuid.UUID, len(x.GuidArray.Value))
			for i, v := range x.GuidArray.Value {
				a[i], _ = uuid.FromBytes(v)
			}
			return ua.NewVariantWithType(a, ua.VariantTypeGuid, []int32{int32(len(x.GuidArray.Value))})
		case *Variant_ByteStringArray:
			a := make([]ua.ByteString, len(x.ByteStringArray.Value))
			for i, v := range x.ByteStringArray.Value {
				a[i] = ua.ByteString(v)
			}
			return ua.NewVariantWithType(a, ua.VariantTypeByteString, []int32{int32(len(x.ByteStringArray.Value))})
		default:
			return &ua.NilVariant
		}
	*/
	return &Variant{Value: &Variant_Null{}}

}
