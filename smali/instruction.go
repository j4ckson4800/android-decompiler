package smali

type Opcode byte
type Type int8
type OperandType int8

// ref: https://github.com/google/smali
// https://source.android.com/docs/core/runtime/dalvik-bytecode
const (
	OpNop Opcode = iota
	OpMove
	OpMoveFrom16
	OpMove16
	OpMoveWide
	OpMoveWideFrom16
	OpMoveWide16
	OpMoveObject
	OpMoveObjectFrom16
	OpMoveObject16
	OpMoveResult
	OpMoveResultWide
	OpMoveResultObject
	OpMoveException

	OpReturnVoid
	OpReturnRegular
	OpReturnWide
	OpReturnObject

	OpConst4
	OpConst16
	OpConstRegular
	OpConstHigh16
	OpConstWide16
	OpConstWide32
	OpConstWide
	OpConstWideHigh16
	OpConstString
	OpConstStringJumbo
	OpConstClass

	OpMonitorEnter
	OpMonitorExit
	OpCheckCast
	OpInstanceOf
	OpArrayLength
	OpNewInstance

	OpNewArray
	OpFilledNewArray
	OpFilledNewArrayRange
	OpFilledArrayData

	OpThrowOp

	OpGotoOp
	OpGoto16
	OpGoto32

	OpPackedSwitch
	OpSparseSwitch

	OpCmplFloat
	OpCmpgFloat
	OpCmplDouble
	OpCmpgDouble
	OpCmpLong

	OpIfEq
	OpIfNe
	OpIfLt
	OpIfGe
	OpIfGt
	OpIfLe

	OpIfEqz
	OpIfNez
	OpIfLtz
	OpIfGez
	OpIfGtz
	OpIfLez

	OpAget = iota + 6
	OpAgetWide
	OpAgetObject
	OpAgetBoolean
	OpAgetByte
	OpAgetChar
	OpAgetShort

	OpAput
	OpAputWide
	OpAputObject
	OpAputBoolean
	OpAputByte
	OpAputChar
	OpAputShort

	OpIget
	OpIgetWide
	OpIgetObject
	OpIgetBoolean
	OpIgetByte
	OpIgetChar
	OpIgetShort
	OpIput
	OpIputWide
	OpIputObject
	OpIputBoolean
	OpIputByte
	OpIputChar
	OpIputShort

	OpSget
	OpSgetWide
	OpSgetObject
	OpSgetBoolean
	OpSgetByte
	OpSgetChar
	OpSgetShort
	OpSput
	OpSputWide
	OpSputObject
	OpSputBoolean
	OpSputByte
	OpSputChar
	OpSputShort

	OpInvokeVirtual
	OpInvokeSuper
	OpInvokeDirect
	OpInvokeStatic
	OpInvokeInterface

	_

	OpInvokeVirtualRange
	OpInvokeSuperRange
	OpInvokeDirectRange
	OpInvokeStaticRange
	OpInvokeInterfaceRange

	OpNegInt = iota + 8
	OpNotInt
	OpNegLong
	OpNotLong
	OpNegFloat
	OpNegDouble
	OpIntToLong
	OpIntToFloat
	OpIntToDouble
	OpLongToInt
	OpLongToFloat
	OpLongToDouble
	OpFloatToInt
	OpFloatToLong
	OpFloatToDouble
	OpDoubleToInt
	OpDoubleToLong
	OpDoubleToFloat
	OpIntToByte
	OpIntToChar
	OpIntToShort

	OpAddInt
	OpSubInt
	OpMulInt
	OpDivInt
	OpRemInt
	OpAndInt
	OpOrInt
	OpXorInt
	OpShlInt
	OpShrInt
	OpUshrInt
	OpAddLong
	OpSubLong
	OpMulLong
	OpDivLong
	OpRemLong
	OpAndLong
	OpOrLong
	OpXorLong
	OpShlLong
	OpShrLong
	OpUshrLong
	OpAddFloat
	OpSubFloat
	OpMulFloat
	OpDivFloat
	OpRemFloat
	OpAddDouble
	OpSubDouble
	OpMulDouble
	OpDivDouble
	OpRemDouble

	OpAddInt2addr
	OpSubInt2addr
	OpMulInt2addr
	OpDivInt2addr
	OpRemInt2addr
	OpAndInt2addr
	OpOrInt2addr
	OpXorInt2addr
	OpShlInt2addr
	OpShrInt2addr
	OpUshrInt2addr
	OpAddLong2addr
	OpSubLong2addr
	OpMulLong2addr
	OpDivLong2addr
	OpRemLong2addr
	OpAndLong2addr
	OpOrLong2addr
	OpXorLong2addr
	OpShlLong2addr
	OpShrLong2addr
	OpUshrLong2addr
	OpAddFloat2addr
	OpSubFloat2addr
	OpMulFloat2addr
	OpDivFloat2addr
	OpRemFloat2addr
	OpAddDouble2addr
	OpSubDouble2addr
	OpMulDouble2addr
	OpDivDouble2addr
	OpRemDouble2addr

	OpAddIntLit16
	OpRsubIntLit16
	OpMulIntLit16
	OpDivIntLit16
	OpRemIntLit16
	OpAndIntLit16
	OpOrIntLit16
	OpXorIntLit16

	OpAddIntLit8
	OpRsubIntLit8
	OpMulIntLit8
	OpDivIntLit8
	OpRemIntLit8
	OpAndIntLit8
	OpOrIntLit8
	OpXorIntLit8
	OpShlIntLit8
	OpShrIntLit8
	OpUshrIntLit8

	OpInvokePolymorphic = iota + 31
	OpInvokePolymorphicRange
	OpInvokeCustom
	OpInvokeCustomRange
	OpConstMethodHandle
	OpConstMethodType
)

const (
	TypeUnknown Type = iota - 1
	TypeMove
	TypeCond
	TypeGoto
	TypeNoop
	TypeCast
	TypeConst
	TypeReturn
	TypeArrayOp
	TypeStaticOp
	TypeSwitchOp
	TypeException
	TypeInvocation
	TypeComparison
	TypeInstanceOp
	TypeMoveResult
	TypeArithmetics
)

const (
	OperandTypeNone OperandType = iota
	OperandTypeReg
	OperandType2reg
	OperandTypeRegShort
	OperandTypeShort
	OperandType2short
	OperandType2regShort
	OperandTypeUint
	OperandTypeRegUint
	OperandTypeRegUlong
	OperandType3reg
	OperandRegisterArray
	OperandRegisterArrayRange
	OperandRegHigh32
	OperandRegHigh64
	OperandRegWide16
	OperandRegWide32
)

func getInstructionOperandsType(opcode Opcode) OperandType {
	switch opcode {
	case OpNop, OpReturnVoid:
		return OperandTypeNone
	case OpAddInt2addr, OpSubInt2addr, OpMulInt2addr, OpDivInt2addr, OpRemInt2addr,
		OpAndInt2addr, OpOrInt2addr, OpXorInt2addr, OpShlInt2addr, OpShrInt2addr, OpUshrInt2addr,
		OpAddLong2addr, OpSubLong2addr, OpMulLong2addr, OpDivLong2addr, OpRemLong2addr,
		OpAndLong2addr, OpOrLong2addr, OpXorLong2addr, OpShlLong2addr, OpShrLong2addr, OpUshrLong2addr,
		OpAddFloat2addr, OpSubFloat2addr, OpMulFloat2addr, OpDivFloat2addr, OpRemFloat2addr,
		OpAddDouble2addr, OpSubDouble2addr, OpMulDouble2addr, OpDivDouble2addr, OpRemDouble2addr,
		OpNegInt, OpNotInt, OpNegLong, OpNotLong, OpNegFloat, OpNegDouble,
		OpIntToLong, OpIntToFloat, OpIntToDouble, OpLongToInt, OpLongToFloat, OpLongToDouble,
		OpFloatToInt, OpFloatToLong, OpFloatToDouble, OpDoubleToInt, OpDoubleToLong, OpDoubleToFloat,
		OpIntToByte, OpIntToChar, OpIntToShort, OpArrayLength, OpConst4,
		OpMove, OpMoveObject, OpMoveWide:
		return OperandType2reg
	case OpConstMethodType, OpConstMethodHandle,
		OpSget, OpSgetWide, OpSgetObject, OpSgetBoolean, OpSgetByte, OpSgetChar, OpSgetShort,
		OpSput, OpSputWide, OpSputObject, OpSputBoolean, OpSputByte, OpSputChar, OpSputShort,
		OpConst16, OpConstString, OpConstClass,
		OpMoveFrom16, OpMoveWideFrom16, OpMoveObjectFrom16,
		OpNewInstance,
		OpCheckCast,
		OpIfEqz, OpIfNez, OpIfLtz, OpIfGez, OpIfGtz, OpIfLez:
		return OperandTypeRegShort
	case OpMove16, OpMoveWide16, OpMoveObject16:
		return OperandType2short
	case OpThrowOp,
		OpMonitorEnter, OpMonitorExit,
		OpReturnWide, OpReturnObject, OpReturnRegular,
		OpMoveResult, OpMoveResultWide, OpMoveResultObject,
		OpGotoOp,
		OpMoveException:
		return OperandTypeReg
	case OpConstRegular, OpConstStringJumbo:
		return OperandTypeRegUint
	case OpAddIntLit16, OpRsubIntLit16, OpMulIntLit16, OpDivIntLit16, OpRemIntLit16,
		OpAndIntLit16, OpOrIntLit16, OpXorIntLit16,
		OpIget, OpIgetWide, OpIgetObject, OpIgetBoolean, OpIgetByte, OpIgetChar, OpIgetShort,
		OpIput, OpIputWide, OpIputObject, OpIputBoolean, OpIputByte, OpIputChar, OpIputShort,
		OpInstanceOf,
		OpNewArray,
		OpIfEq, OpIfNe, OpIfLt, OpIfGe, OpIfGt, OpIfLe:
		return OperandType2regShort
	case OpAddIntLit8, OpRsubIntLit8, OpMulIntLit8, OpDivIntLit8, OpRemIntLit8,
		OpAndIntLit8, OpOrIntLit8, OpXorIntLit8, OpShlIntLit8, OpShrIntLit8, OpUshrIntLit8,
		OpAddInt, OpSubInt, OpMulInt, OpDivInt, OpRemInt,
		OpAndInt, OpOrInt, OpXorInt, OpShlInt, OpShrInt, OpUshrInt,
		OpAddLong, OpSubLong, OpMulLong, OpDivLong, OpRemLong,
		OpAndLong, OpOrLong, OpXorLong, OpShlLong, OpShrLong, OpUshrLong,
		OpAddFloat, OpSubFloat, OpMulFloat, OpDivFloat, OpRemFloat,
		OpAddDouble, OpSubDouble, OpMulDouble, OpDivDouble, OpRemDouble,
		OpAget, OpAgetWide, OpAgetObject, OpAgetBoolean, OpAgetByte, OpAgetChar, OpAgetShort,
		OpAput, OpAputWide, OpAputObject, OpAputBoolean, OpAputByte, OpAputChar, OpAputShort,
		OpCmpLong, OpCmplFloat, OpCmpgFloat, OpCmplDouble, OpCmpgDouble:
		return OperandType3reg
	case OpGoto16:
		return OperandTypeShort
	case OpGoto32:
		return OperandTypeUint
	case OpConstWide:
		return OperandTypeRegUlong
	case OpFilledNewArray,
		OpInvokePolymorphic, OpInvokeCustom, OpInvokeVirtual, OpInvokeSuper, OpInvokeDirect, OpInvokeStatic, OpInvokeInterface:
		return OperandRegisterArray
	case OpFilledNewArrayRange,
		OpInvokePolymorphicRange, OpInvokeCustomRange, OpInvokeVirtualRange, OpInvokeSuperRange, OpInvokeDirectRange, OpInvokeStaticRange, OpInvokeInterfaceRange:
		return OperandRegisterArrayRange
	case OpConstHigh16:
		return OperandRegHigh32
	case OpConstWideHigh16:
		return OperandRegHigh64
	case OpConstWide16:
		return OperandRegWide16
	case OpConstWide32:
		return OperandRegWide32
	case OpFilledArrayData,
		OpPackedSwitch,
		OpSparseSwitch:
		// We don't want to deal with packed payloads at this moment
		// Because I don't want to add packed payloads to the instruction type
		return OperandTypeRegUint
	}

	return OperandTypeNone
}

func getInstructionType(opcode Opcode) Type {
	switch opcode {
	case OpNop:
		return TypeNoop
	case OpReturnVoid,
		OpReturnWide, OpReturnObject, OpReturnRegular:
		return TypeReturn
	case OpAddInt2addr, OpSubInt2addr, OpMulInt2addr, OpDivInt2addr, OpRemInt2addr,
		OpAndInt2addr, OpOrInt2addr, OpXorInt2addr, OpShlInt2addr, OpShrInt2addr, OpUshrInt2addr,
		OpAddLong2addr, OpSubLong2addr, OpMulLong2addr, OpDivLong2addr, OpRemLong2addr,
		OpAndLong2addr, OpOrLong2addr, OpXorLong2addr, OpShlLong2addr, OpShrLong2addr, OpUshrLong2addr,
		OpAddFloat2addr, OpSubFloat2addr, OpMulFloat2addr, OpDivFloat2addr, OpRemFloat2addr,
		OpAddDouble2addr, OpSubDouble2addr, OpMulDouble2addr, OpDivDouble2addr, OpRemDouble2addr,
		OpNegInt, OpNotInt, OpNegLong, OpNotLong, OpNegFloat, OpNegDouble,
		OpAddIntLit16, OpRsubIntLit16, OpMulIntLit16, OpDivIntLit16, OpRemIntLit16,
		OpAndIntLit16, OpOrIntLit16, OpXorIntLit16,
		OpAddIntLit8, OpRsubIntLit8, OpMulIntLit8, OpDivIntLit8, OpRemIntLit8,
		OpAndIntLit8, OpOrIntLit8, OpXorIntLit8, OpShlIntLit8, OpShrIntLit8, OpUshrIntLit8,
		OpAddInt, OpSubInt, OpMulInt, OpDivInt, OpRemInt,
		OpAndInt, OpOrInt, OpXorInt, OpShlInt, OpShrInt, OpUshrInt,
		OpAddLong, OpSubLong, OpMulLong, OpDivLong, OpRemLong,
		OpAndLong, OpOrLong, OpXorLong, OpShlLong, OpShrLong, OpUshrLong,
		OpAddFloat, OpSubFloat, OpMulFloat, OpDivFloat, OpRemFloat,
		OpAddDouble, OpSubDouble, OpMulDouble, OpDivDouble, OpRemDouble:
		return TypeArithmetics
	case OpIntToLong, OpIntToFloat, OpIntToDouble, OpLongToInt, OpLongToFloat, OpLongToDouble,
		OpFloatToInt, OpFloatToLong, OpFloatToDouble, OpDoubleToInt, OpDoubleToLong, OpDoubleToFloat,
		OpIntToByte, OpIntToChar, OpIntToShort:
		return TypeCast
	case OpArrayLength,
		OpNewArray,
		OpAget, OpAgetWide, OpAgetObject, OpAgetBoolean, OpAgetByte, OpAgetChar, OpAgetShort,
		OpAput, OpAputWide, OpAputObject, OpAputBoolean, OpAputByte, OpAputChar, OpAputShort,
		OpFilledArrayData,
		OpFilledNewArray,
		OpFilledNewArrayRange:
		return TypeArrayOp
	case OpConst4,
		OpConst16, OpConstClass, OpConstString,
		OpConstMethodType, OpConstMethodHandle,
		OpConstRegular, OpConstStringJumbo,
		OpConstWide,
		OpConstHigh16,
		OpConstWideHigh16,
		OpConstWide16,
		OpConstWide32:
		return TypeConst
	case OpMove, OpMoveObject, OpMoveWide,
		OpMoveFrom16, OpMoveWideFrom16, OpMoveObjectFrom16,
		OpMove16, OpMoveWide16, OpMoveObject16:
		return TypeMove
	case OpMoveResult, OpMoveResultWide, OpMoveResultObject:
		return TypeMoveResult
	case OpSget, OpSgetWide, OpSgetObject, OpSgetBoolean, OpSgetByte, OpSgetChar, OpSgetShort,
		OpSput, OpSputWide, OpSputObject, OpSputBoolean, OpSputByte, OpSputChar, OpSputShort:
		return TypeStaticOp
	case OpNewInstance,
		OpIget, OpIgetWide, OpIgetObject, OpIgetBoolean, OpIgetByte, OpIgetChar, OpIgetShort,
		OpIput, OpIputWide, OpIputObject, OpIputBoolean, OpIputByte, OpIputChar, OpIputShort,
		OpInstanceOf:
		return TypeInstanceOp
	case OpCheckCast,
		OpCmpLong,
		OpCmplFloat, OpCmpgFloat,
		OpCmplDouble, OpCmpgDouble:
		return TypeComparison
	case OpThrowOp,
		OpMonitorEnter, OpMonitorExit,
		OpMoveException:
		return TypeException
	case OpIfEq, OpIfNe, OpIfLt, OpIfGe, OpIfGt, OpIfLe,
		OpIfEqz, OpIfNez, OpIfLtz, OpIfGez, OpIfGtz, OpIfLez:
		return TypeCond
	case OpGotoOp,
		OpGoto16,
		OpGoto32:
		return TypeGoto
	case OpPackedSwitch,
		OpSparseSwitch:
		return TypeSwitchOp
	case OpInvokePolymorphic, OpInvokeCustom, OpInvokeVirtual, OpInvokeSuper, OpInvokeDirect, OpInvokeStatic, OpInvokeInterface,
		OpInvokePolymorphicRange, OpInvokeCustomRange, OpInvokeVirtualRange, OpInvokeSuperRange, OpInvokeDirectRange, OpInvokeStaticRange, OpInvokeInterfaceRange:
		return TypeInvocation
	}

	return TypeUnknown
}

type Instruction struct {
	Opcode      Opcode      // 0x0
	Type        Type        // 0x1
	OperandType OperandType // 0x2
	Pad         byte        // 0x3
	Operands    []int64     // 0x4
} // size: 0x14
