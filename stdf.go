package stdf

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"time"
)

type C1 byte

// Fixed length character string:
// If a fixed length character string does not fill the entire field, it
// must be left-justified and padded with spaces.
type C12 [12]byte

// Variable length character string:
// first byte = unsigned count of bytes to follow (maximum of 255 bytes)
type CN []byte

// Variable length character string:
// string length is stored in another field
type CF []byte

// One byte unsigned integer
type U1 uint8

// Two byte unsigned integer
type U2 uint16

// Four byte unsigned integer
type U4 uint32

// One byte signed integer
type I1 int8

// Two byte signed integer
type I2 int16

// Four byte signed integer
type I4 int32

// Four byte floating point number
type R4 float32

// Eight byte floating point number
type R8 float64

// Fixed length bit-encoded data
type B6 [6]byte

type KXU1 []U1

type StdfRecordType interface {
	// 将对象转换为字节切片输出
	ToByte() ([]byte, error)

	ToString() string
}

// REC_TYP Code 	Meaning and STDF REC_SUB Codes
// 0 				Information about the STDF file
// 						10 File Attributes Record (FAR)
// 						20 Audit Trail Record (ATR)
// 1 				Data collected on a per lot basis
// 						10 Master Information Record (MIR)
// 						20 Master Results Record (MRR)
// 						30 Part Count Record (PCR)
// 						40 Hardware Bin Record (HBR)
// 						50 Software Bin Record (SBR)
// 						60 Pin Map Record (PMR)
// 						62 Pin Group Record (PGR)
// 						63 Pin List Record (PLR)
// 						70 Retest Data Record (RDR)
// 						80 Site Description Record (SDR)
// 2 				Data collected per wafer
// 						10 Wafer Information Record (WIR)
// 						20 Wafer Results Record (WRR)
// 						30 Wafer Configuration Record (WCR)
// 5 				Data collected on a per part basis
// 						10 Part Information Record (PIR)
// 						20 Part Results Record (PRR)
// 10 				Data collected per test in the test program
// 						30 Test Synopsis Record (TSR)
// 15 				Data collected per test execution
// 						10 Parametric Test Record (PTR)
// 						15 Multiple-Result Parametric Record (MPR)
// 						20 Functional Test Record (FTR)
// 20 				Data collected per program segment
// 						10 Begin Program Section Record (BPS)
// 						20 End Program Section Record (EPS)
// 50 				Generic Data
// 						10 Generic Data Record (GDR)
// 						30 Datalog Text Record (DTR)
// 180 				Reserved for use by Image software
// 181 				Reserved for use by IG900 software
func NewStdfRecord(a []byte) StdfRecordType {
	var t BasicRecordType
	t.Rec_Len = U2(binary.LittleEndian.Uint16(a))
	t.Rec_Type = U1(a[2])
	t.Rec_Sub = U1(a[3])
	switch t.Rec_Type {
	case 0:
		switch t.Rec_Sub {
		case 10:
			var far FAR
			far.BasicRecordType = t
			return &far
		}
	case 1:
		switch t.Rec_Sub {
		case 10:
			var mir MIR
			mir.BasicRecordType = t
			return &mir
		case 80:
			var sdr SDR
			sdr.BasicRecordType = t
			return &sdr
		}
	}
	return nil
}

func TransB2S(s []byte, o1 interface{}) error {
	t := reflect.TypeOf(o1)
	v := reflect.ValueOf(o1)
	m := 0
	for i := 0; i < t.Elem().NumField(); i++ {
		if m >= len(s) {
			break
		}
		field := v.Elem().Field(i)
		fT := fmt.Sprintf("%v", field.Type())
		switch fT {
		case "stdf.U1":
			ii1 := U1(s[m])
			v.Elem().Field(i).Set(reflect.ValueOf(ii1))
			m++
		case "stdf.U2":
			v.Elem().Field(i).Set(reflect.ValueOf(U2(binary.LittleEndian.Uint16(s[m : m+2]))))
			m += 2
		case "stdf.U4":
			v.Elem().Field(i).Set(reflect.ValueOf(U4(binary.LittleEndian.Uint32(s[m : m+4]))))
			m += 4
		case "stdf.C1":
			ii1 := C1(s[m])
			v.Elem().Field(i).Set(reflect.ValueOf(ii1))
			m++
		case "stdf.CN":
			i1 := int(s[m])
			if i1+m > len(s) {
				v.Elem().Field(i).Set(reflect.ValueOf(CN(s[m:])))
			} else {
				v.Elem().Field(i).Set(reflect.ValueOf(CN(s[m : i1+m])))
			}
			m = m + 1 + i1
		case "stdf.KXU1":
			i1 := v.Elem().Field(i - 1).Interface().(int)
			var t2 []U1
			for j := 0; j < int(i1); j++ {
				t2 = append(t2, U1(s[m+j]))
			}
			v.Elem().Field(i).Set(reflect.ValueOf(t2))
			m = m + i1
		}
	}
	return nil
	// var params []reflect.Value
	// params = append(params, reflect.ValueOf(s))

	// v1.MethodByName("ToRecord").Call(params)
	// t.Log("\n 对象对应的字符串内容:", v1.MethodByName("ToString").Call(nil))

	// t.Log("Number of fields", v1.Elem().NumField())
	// for i := 0; i < v1.Elem().NumField(); i++ {

	// 	fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
	// 	t.Logf("Field:%d \t type:%T \t value:%v\n",
	// 		i, v1.Elem().Field(i), v1.Elem().Field(i))
	// }
}

type BasicRecordType struct {
	// Bytes of data following header
	Rec_Len U2
	// Record type (0)
	Rec_Type U1
	// REC_SUB U*1 Record sub-type (10)
	Rec_Sub U1
}

// File Attributes Record (FAR)
// Function: Contains the information necessary to determine how to decode the STDF data contained in the file.
//
// Data Fields:
// Field Name         Data Type 	Field Description  	 				Missing/Invalid Data Flag
// REC_LEN 				U*2 		Bytes of data following header
// REC_TYP 				U*1 		Record type (0)
// REC_SUB				U*1 		Record sub-type (10)
// CPU_TYPE 			U*1 		CPU type that wrote this file
// STDF_VER 			U*1 		STDF version number
//
// Notes on Specific Fields:
// ---------------------------------------------
// CPU_TYPE Indicates which type of CPU wrote this STDF file. This information is useful for
// determining the CPU-dependent data representation of the integer and floating point
// fields in the file’s records. The valid values are:
// 0 = DEC PDP-11 and VAX processors. F and D floating point formats
// will be used. G and H floating point formats will not be used.
// 1 = Sun 1, 2, 3, and 4 computers.
// 2 = Sun 386i computers, and IBM PC, IBM PC-AT, and IBM PC-XT
// computers.
// 3-127 = Reserved for future use by Teradyne.
// 128-255 = Reserved for use by customers.
// Acode defined heremay also be valid for otherCPUtypeswhose data formats are fully
// compatible with that of the type listed here. Before using one of these codes for a CPU
// type not listed here, please check with the Teradyne hotline, which can provide
// additional information on CPU compatibility.
// STDF_VER Identifies the version number of the STDF specification used in generating the data.
// This allows data analysis programs to handle STDF specification enhancements.
// ---------------------------------------------
//
// Location: Required as the first record of the file.
type FAR struct {
	BasicRecordType
	// CPU type that wrote this file
	Cpu_Type U1
	// STDF version number
	Stdf_Ver U1
}

func (f FAR) ToByte() ([]byte, error) {
	// var b []byte
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(f.Rec_Len))
	b = append(b, byte(f.Rec_Type))
	b = append(b, byte(f.Rec_Sub))
	b = append(b, byte(f.Cpu_Type))
	b = append(b, byte(f.Stdf_Ver))
	return b, nil
}

func (f FAR) ToString() string {
	return fmt.Sprintf("Rec Len=%v, Rec Type=%v, Rec sub=%v, Cpu Type=%v, Stdf Ver=%v", f.Rec_Len, f.Rec_Type, f.Rec_Sub, f.Cpu_Type, f.Stdf_Ver)
}

// Audit Trail Record (ATR)
// Function: Used to record any operation that alters the contents of the STDF file. The name of the
// program and all its parameters should be recorded in the ASCII field provided in this
// record. Typically, this record will be used to track filter programs that have been
// applied to the data.
// Data Fields:
// Field Data Field Missing/Invalid
// Name Type Description Data Flag
// REC_LEN U*2 Bytes of data following header
// REC_TYP U*1 Record type (0)
// REC_SUB U*1 Record sub-type (20)
// MOD_TIM U*4 Date and time of STDF file modification
// CMD_LINE C*n Command line of program
// Frequency: Optional. One for each filter or other data transformation program applied to the STDF
// data.
// Location: Between the File Attributes Record (FAR) and the Master Information Record (MIR).
// The filter program that writes the altered STDF file must write its ATR immediately
// after the FAR (and hence before any other ATRs that may be in the file). In this way,
// multiple ATRs will be in reverse chronological order.
// Possible Use: Determining whether a particular filter has been applied to the data.
// type ATR struct {
// 	BasicRecordType
// 	// Date and time of STDF file modification
// 	MOD_TIM U4
// 	// Command line of program
// 	CMD_LINE CN
// }

// func (f *ATR) ToRecord(b []byte) error {
// 	f.MOD_TIM = U1(b[4])
// 	f.CMD_LINE = U1(b[5])
// 	return nil
// }

// func (f ATR) ToByte() ([]byte, error) {
// 	// var b []byte
// 	b := make([]byte, 2)
// 	binary.LittleEndian.PutUint16(b, uint16(f.Rec_Len))
// 	b = append(b, byte(f.Rec_Type))
// 	b = append(b, byte(f.Rec_Sub))
// 	b = append(b, byte(f.Cpu_Type))
// 	b = append(b, byte(f.Stdf_Ver))
// 	return b, nil
// }

// func (f ATR) ToString() string {
// 	return fmt.Sprintf("Rec Len=%v, Rec Type=%v, Rec Type=%v, Cpu Type=%v, Stdf Ver=%v", f.Rec_Len, f.Rec_Type, f.Rec_Type, f.Cpu_Type, f.Stdf_Ver)
// }

// Master Information Record (MIR)
// Function: The MIR and the MRR (Master Results Record) contain all the global information that
// is to be stored for a tested lot of parts. Each data stream must have exactly one MIR,
// immediately after the FAR (and the ATRs, if they are used). This will allow any data
// reporting or analysis programs access to this information in the shortest possible
// amount of time.
// Data Fields:
// Field Data Field Missing/Invalid
// Name Type Description Data Flag
// REC_LEN U*2 Bytes of data following header
// REC_TYP U*1 Record type (1)
// REC_SUB U*1 Record sub-type (10)
// SETUP_T U*4 Date and time of job setup
// START_T U*4 Date and time first part tested
// STAT_NUM U*1 Tester station number
// MODE_COD C*1 Test mode code (e.g. prod, dev) space
// RTST_COD C*1 Lot retest code space
// PROT_COD C*1 Data protection code space
// BURN_TIM U*2 Burn-in time (in minutes) 65,535
// CMOD_COD C*1 Command mode code space
// LOT_ID C*n Lot ID (customer specified)
// PART_TYP C*n Part Type (or product ID)
// NODE_NAM C*n Name of node that generated data
// TSTR_TYP C*n Tester type
// JOB_NAM C*n Job name (test program name)
// JOB_REV C*n Job (test program) revision number length byte = 0
// SBLOT_ID C*n Sublot ID length byte = 0
// OPER_NAM C*n Operator name or ID (at setup time) length byte = 0
// EXEC_TYP C*n Tester executive software type length byte = 0
// EXEC_VER C*n Tester exec software version number length byte = 0
// TEST_COD C*n Test phase or step code length byte = 0
// TST_TEMP C*n Test temperature length byte = 0
// USER_TXT C*n Generic user text length byte = 0
// AUX_FILE C*n Name of auxiliary data file length byte = 0
// PKG_TYP C*n Package type length byte = 0
// FAMLY_ID C*n Product family ID length byte = 0
// DATE_COD C*n Date code length byte = 0
// FACIL_ID C*n Test facility ID length byte = 0
// FLOOR_ID C*n Test floor ID length byte = 0
// PROC_ID C*n Fabrication process ID length byte = 0
// OPER_FRQ C*n Operation frequency or step length byte = 0
// SPEC_NAM C*n Test specification name length byte = 0
// SPEC_VER C*n Test specification version number length byte = 0
// FLOW_ID C*n Test flow ID length byte = 0
// SETUP_ID C*n Test setup ID length byte = 0
// DSGN_REV C*n Device design revision length byte = 0
// ENG_ID C*n Engineering lot ID length byte = 0
// ROM_COD C*n ROM code ID length byte = 0
// SERL_NUM C*n Tester serial number length byte = 0
// SUPR_NAM C*n Supervisor name or ID length byte = 0

// Notes on Specific Fields:
// MODE_COD Indicates the station mode under which the parts were tested. Currently defined
// values for the MODE_COD field are:
// A = AEL(AutomaticEdgeLock)mode
// C = Checkermode
// D = Development / Debug test mode
// E = Engineering mode (same as Development mode)
// M = Maintenancemode
// P = Production test mode
// Q = Quality Control
// All other alphabetic codes are reserved for future use by Teradyne. The characters 0 -
// 9 are available for customer use.
// RTST_COD Indicates whether the lot of parts has been previously tested under the same test
// conditions. Suggested values are:
// Y = Lot was previously tested.
// N = Lot has not been previously tested.
// space = Not known if lot has been previously tested.
// 0 - 9 = Number of times lot has previously been tested.
// PROT_COD User-defined field indicating the protection desired for the test data being stored. Valid
// values are the ASCII characters 0 - 9 and A - Z. A space in this field indicates a missing
// value (default protection).
// CMOD_COD Indicates the command mode of the tester during testing of the parts. The user or the
// tester executive software defines command mode values. Valid values are the ASCII
// characters 0 - 9 and A - Z. A space indicates a missing value.
// STDF Record Types Master Information Record (MIR)
// STDF Specification V4 Page 20
// Main Menu
// Frequency: Always required. One per data stream.
// Location: Immediately after the File Attributes Record (FAR) and the Audit Trail Records (ATR),
// if ATRs are used.
// Possible Use: Header information for all reports
type MIR struct {
	BasicRecordType
	// Date and time of job setup
	SETUP_T U4
	// Date and time first part tested
	START_T U4
	// Tester station number
	STAT_NUM U1
	// Test mode code (e.g. prod, dev) space
	MODE_COD C1
	// Lot retest code space
	RTST_COD C1
	// Data protection code space
	PROT_COD C1
	// Burn-in time (in minutes) 65,535
	BURN_TIM U2
	// Command mode code space
	CMOD_COD C1
	// Lot ID (customer specified)
	LOT_ID CN
	// Part Type (or product ID)
	PART_TYP CN
	// Name of node that generated data
	NODE_NAM CN
	// Tester type
	TSTR_TYP CN
	// Job name (test program name)
	JOB_NAM CN
	// Job (test program) revision number length byte = 0
	JOB_REV CN
	// Sublot ID length byte = 0
	SBLOT_ID CN
	// Operator name or ID (at setup time) length byte = 0
	OPER_NAM CN
	// Tester executive software type length byte = 0
	EXEC_TYP CN
	// Tester exec software version number length byte = 0
	EXEC_VER CN
	// Test phase or step code length byte = 0
	TEST_COD CN
	// Test temperature length byte = 0
	TST_TEMP CN
	// Generic user text length byte = 0
	USER_TXT CN
	// Name of auxiliary data file length byte = 0
	AUX_FILE CN
	// Package type length byte = 0
	PKG_TYP CN
	// Product family ID length byte = 0
	FAMLY_ID CN
	// Date code length byte = 0
	DATE_COD CN
	// Test facility ID length byte = 0
	FACIL_ID CN
	// Test floor ID length byte = 0
	FLOOR_ID CN
	// Fabrication process ID length byte = 0
	PROC_ID CN
	// Operation frequency or step length byte = 0
	OPER_FRQ CN
	// Test specification name length byte = 0
	SPEC_NAM CN
	// Test specification version number length byte = 0
	SPEC_VER CN
	// Test flow ID length byte = 0
	FLOW_ID CN
	// Test setup ID length byte = 0
	SETUP_ID CN
	// Device design revision length byte = 0
	DSGN_REV CN
	// Engineering lot ID length byte = 0
	ENG_ID CN
	// ROM code ID length byte = 0
	ROM_COD CN
	// Tester serial number length byte = 0
	SERL_NUM CN
	// Supervisor name or ID length byte = 0
	SUPR_NAM CN
}

func (f MIR) ToByte() ([]byte, error) {
	// var b []byte
	b := make([]byte, 2)

	// b = append(b, byte(f.Rec_Type))
	// b = append(b, byte(f.Rec_Sub))
	// binary.LittleEndian.PutUint16(b, uint16(f.SETUP_T))
	// b = append(b, byte(f.SETUP_T))
	// b = append(b, byte(f.Stdf_Ver))
	return b, nil
}

func (f MIR) ToString() string {
	// return ""
	return fmt.Sprintf("Rec Len=%v, Rec Type=%v, Rec Sub=%v, SETUP_T=%v, LOT_ID=%v",
		f.Rec_Len, f.Rec_Type, f.Rec_Sub, time.Unix(int64(f.SETUP_T), 0), string(f.LOT_ID))
}

// Site Description Record (SDR)
// Function: Contains the configuration information for one or more test sites, connected to one test
// head, that compose a site group.
// Data Fields:
// Field Data Field Missing/Invalid
// Name Type Description Data Flag
// REC_LEN U*2 Bytes of data following header
// REC_TYP U*1 Record type (1)
// REC_SUB U*1 Record sub-type (80)
// HEAD_NUM U*1 Test head number
// SITE_GRP U*1 Site group number
// SITE_CNT U*1 Number (k) of test sites in site group
// SITE_NUM kxU*1 Array of test site numbers
// HAND_TYP C*n Handler or prober type length byte = 0
// HAND_ID C*n Handler or prober ID length byte = 0
// CARD_TYP C*n Probe card type length byte = 0
// CARD_ID C*n Probe card ID length byte = 0
// LOAD_TYP C*n Load board type length byte = 0
// LOAD_ID C*n Load board ID length byte = 0
// DIB_TYP C*n DIB board type length byte = 0
// DIB_ID C*n DIB board ID length byte = 0
// CABL_TYP C*n Interface cable type length byte = 0
// CABL_ID C*n Interface cable ID length byte = 0
// CONT_TYP C*n Handler contactor type length byte = 0
// CONT_ID C*n Handler contactor ID length byte = 0
// LASR_TYP C*n Laser type length byte = 0
// LASR_ID C*n Laser ID length byte = 0
// EXTR_TYP C*n Extra equipment type field length byte = 0
// EXTR_ID C*n Extra equipment ID length byte = 0
// Notes on Specific Fields:
// SITE_GRP Specifies a site group number (called a station number on some testers) for the group
// of sites whose configuration is defined by this record. Note that this is different from
// the station number specified in the MIR, which refers to a software station only.
// The value in this field must be unique within the STDF file.
// SITE_CNT,
// SITE_NUM
// SITE_CNT tells how many sites are in the site group that the current SDR configuration
// applies to. SITE_NUM is an array of those site numbers.
// STDF Record Types Site Description Record (SDR)
// STDF Specification V4 Page 34
// Main Menu
// Frequency: One for each site or group of sites that is differently configured.
// Location: Immediately after the MIR and RDR (if an RDR is used).
// Possible Use: Correlation of yield to interface or peripheral equipment

type SDR struct {
	BasicRecordType
	// Test head number
	HEAD_NUM U1
	// Site group number
	SITE_GRP U1
	// Number (k) of test sites in site group
	SITE_CNT U1
	// Array of test site numbers
	SITE_NUM KXU1
	// Handler or prober type length byte = 0
	HAND_TYP CN
	// Handler or prober ID length byte = 0
	HAND_ID CN
	// Probe card type length byte = 0
	CARD_TYP CN
	// Probe card ID length byte = 0
	CARD_ID CN
	// Load board type length byte = 0
	LOAD_TYP CN
	// Load board ID length byte = 0
	LOAD_ID CN
	// DIB board type length byte = 0
	DIB_TYP CN
	// DIB board ID length byte = 0
	DIB_ID CN
	// Interface cable type length byte = 0
	CABL_TYP CN
	// Interface cable ID length byte = 0
	CABL_ID CN
	// Handler contactor type length byte = 0
	CONT_TYP CN
	// Handler contactor ID length byte = 0
	CONT_ID CN
	// Laser type length byte = 0
	LASR_TYP CN
	// Laser ID length byte = 0
	LASR_ID CN
	// Extra equipment type field length byte = 0
	EXTR_TYP CN
	// Extra equipment ID length byte = 0
	EXTR_ID CN
}

func (f SDR) ToByte() ([]byte, error) {
	// var b []byte
	b := make([]byte, 2)

	// b = append(b, byte(f.Rec_Type))
	// b = append(b, byte(f.Rec_Sub))
	// binary.LittleEndian.PutUint16(b, uint16(f.SETUP_T))
	// b = append(b, byte(f.SETUP_T))
	// b = append(b, byte(f.Stdf_Ver))
	return b, nil
}

func (f SDR) ToString() string {
	// return ""
	return fmt.Sprintf("Rec Len=%v, Rec Type=%v, Rec Sub=%v, SETUP_T=%v, LOT_ID=%v",
		f.Rec_Len, f.Rec_Type, f.Rec_Sub, f.HEAD_NUM, f.SITE_GRP)
}
