package nixpacks

/*


// This file was autogenerated by some hot garbage in the `uniffi` crate.
// Trust me, you don't want to mess with it!

#include <stdbool.h>
#include <stdint.h>

// The following structs are used to implement the lowest level
// of the FFI, and thus useful to multiple uniffied crates.
// We ensure they are declared exactly once, with a header guard, UNIFFI_SHARED_H.
#ifdef UNIFFI_SHARED_H
	// We also try to prevent mixing versions of shared uniffi header structs.
	// If you add anything to the #else block, you must increment the version suffix in UNIFFI_SHARED_HEADER_V4
	#ifndef UNIFFI_SHARED_HEADER_V4
		#error Combining helper code from multiple versions of uniffi is not supported
	#endif // ndef UNIFFI_SHARED_HEADER_V4
#else
#define UNIFFI_SHARED_H
#define UNIFFI_SHARED_HEADER_V4
// ⚠️ Attention: If you change this #else block (ending in `#endif // def UNIFFI_SHARED_H`) you *must* ⚠️
// ⚠️ increment the version suffix in all instances of UNIFFI_SHARED_HEADER_V4 in this file.           ⚠️

typedef struct RustBuffer {
	int32_t capacity;
	int32_t len;
	uint8_t *data;
} RustBuffer;

typedef int32_t (*ForeignCallback)(uint64_t, int32_t, RustBuffer, RustBuffer *);

typedef struct ForeignBytes {
	int32_t len;
	const uint8_t *data;
} ForeignBytes;

// Error definitions
typedef struct RustCallStatus {
	int8_t code;
	RustBuffer errorBuf;
} RustCallStatus;

// ⚠️ Attention: If you change this #else block (ending in `#endif // def UNIFFI_SHARED_H`) you *must* ⚠️
// ⚠️ increment the version suffix in all instances of UNIFFI_SHARED_HEADER_V4 in this file.           ⚠️
#endif // def UNIFFI_SHARED_H

RustBuffer nixpacks_65e0_get_plan_providers(
	RustBuffer path,
	RustBuffer envs,
	RustCallStatus* out_status
);

RustBuffer ffi_nixpacks_65e0_rustbuffer_alloc(
	int32_t size,
	RustCallStatus* out_status
);

RustBuffer ffi_nixpacks_65e0_rustbuffer_from_bytes(
	ForeignBytes bytes,
	RustCallStatus* out_status
);

void ffi_nixpacks_65e0_rustbuffer_free(
	RustBuffer buf,
	RustCallStatus* out_status
);

RustBuffer ffi_nixpacks_65e0_rustbuffer_reserve(
	RustBuffer buf,
	int32_t additional,
	RustCallStatus* out_status
);


*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"unsafe"
)

type rustBuffer struct {
	capacity int
	length   int
	data     unsafe.Pointer
	self     C.RustBuffer
}

func fromCRustBuffer(crbuf C.RustBuffer) rustBuffer {
	return rustBuffer{
		capacity: int(crbuf.capacity),
		length:   int(crbuf.len),
		data:     unsafe.Pointer(crbuf.data),
		self:     crbuf,
	}
}

// asByteBuffer reads the full rust buffer and then converts read bytes to a new reader which makes
// it quite inefficient
// TODO: Return an implementation which reads only when needed
func (rb rustBuffer) asReader() *bytes.Reader {
	b := C.GoBytes(rb.data, C.int(rb.length))
	return bytes.NewReader(b)
}

func (rb rustBuffer) asCRustBuffer() C.RustBuffer {
	return C.RustBuffer{
		capacity: C.int(rb.capacity),
		len:      C.int(rb.length),
		data:     (*C.uchar)(unsafe.Pointer(rb.data)),
	}
}

func stringToCRustBuffer(str string) C.RustBuffer {
	b := []byte(str)
	cs := C.CString(str)
	return C.RustBuffer{
		capacity: C.int(len(b)),
		len:      C.int(len(b)),
		data:     (*C.uchar)(unsafe.Pointer(cs)),
	}
}

func (rb rustBuffer) free() {
	rustCall(func(status *C.RustCallStatus) bool {
		C.ffi_nixpacks_65e0_rustbuffer_free(rb.self, status)
		return false
	})
}

type bufLifter[GoType any] interface {
	lift(value C.RustBuffer) GoType
}

type bufLowerer[GoType any] interface {
	lower(value GoType) C.RustBuffer
}

type ffiConverter[GoType any, FfiType any] interface {
	lift(value FfiType) GoType
	lower(value GoType) FfiType
}

type bufReader[GoType any] interface {
	read(reader io.Reader) GoType
}

type bufWriter[GoType any] interface {
	write(writer io.Writer, value GoType)
}

type ffiRustBufConverter[GoType any, FfiType any] interface {
	ffiConverter[GoType, FfiType]
	bufReader[GoType]
}

func lowerIntoRustBuffer[GoType any](bufWriter bufWriter[GoType], value GoType) C.RustBuffer {
	// This might be not the most efficient way but it does not require knowing allocation size
	// beforehand
	var buffer bytes.Buffer
	bufWriter.write(&buffer, value)

	bytes, err := io.ReadAll(&buffer)
	if err != nil {
		panic(fmt.Errorf("reading written data: %w", err))
	}

	return stringToCRustBuffer(string(bytes))
}

func liftFromRustBuffer[GoType any](bufReader bufReader[GoType], rbuf rustBuffer) GoType {
	defer rbuf.free()
	reader := rbuf.asReader()
	item := bufReader.read(reader)
	if reader.Len() > 0 {
		// TODO: Remove this
		leftover, _ := io.ReadAll(reader)
		panic(fmt.Errorf("Junk remaining in buffer after lifting: %s", string(leftover)))
	}
	return item
}

func rustCallWithError[U any](converter bufLifter[error], callback func(*C.RustCallStatus) U) (U, error) {
	var status C.RustCallStatus
	returnValue := callback(&status)
	switch status.code {
	case 0:
		return returnValue, nil
	case 1:
		return returnValue, converter.lift(status.errorBuf)
	case 2:
		// when the rust code sees a panic, it tries to construct a rustbuffer
		// with the message.  but if that code panics, then it just sends back
		// an empty buffer.
		if status.errorBuf.len > 0 {
			panic(fmt.Errorf("%s", FfiConverterstringINSTANCE.lift(status.errorBuf)))
		} else {
			panic(fmt.Errorf("Rust panicked while handling Rust panic"))
		}
	default:
		return returnValue, fmt.Errorf("unknown status code: %d", status.code)
	}
}

func rustCall[U any](callback func(*C.RustCallStatus) U) U {
	returnValue, err := rustCallWithError(nil, callback)
	if err != nil {
		panic(err)
	}
	return returnValue
}

func writeInt8(writer io.Writer, value int8) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeUint8(writer io.Writer, value uint8) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeInt16(writer io.Writer, value int16) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeUint16(writer io.Writer, value uint16) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeInt32(writer io.Writer, value int32) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeUint32(writer io.Writer, value uint32) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeInt64(writer io.Writer, value int64) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeUint64(writer io.Writer, value uint64) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeFloat32(writer io.Writer, value float32) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func writeFloat64(writer io.Writer, value float64) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		panic(err)
	}
}

func readInt8(reader io.Reader) int8 {
	var result int8
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readUint8(reader io.Reader) uint8 {
	var result uint8
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readInt16(reader io.Reader) int16 {
	var result int16
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readUint16(reader io.Reader) uint16 {
	var result uint16
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readInt32(reader io.Reader) int32 {
	var result int32
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readUint32(reader io.Reader) uint32 {
	var result uint32
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readInt64(reader io.Reader) int64 {
	var result int64
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readUint64(reader io.Reader) uint64 {
	var result uint64
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readFloat32(reader io.Reader) float32 {
	var result float32
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func readFloat64(reader io.Reader) float64 {
	var result float64
	if err := binary.Read(reader, binary.BigEndian, &result); err != nil {
		panic(err)
	}
	return result
}

func init() {

}

type FfiConverterstring struct{}

var FfiConverterstringINSTANCE = FfiConverterstring{}

func (FfiConverterstring) lift(cRustBuf C.RustBuffer) string {
	reader := fromCRustBuffer(cRustBuf).asReader()
	b, err := io.ReadAll(reader)
	if err != nil {
		panic(fmt.Errorf("reading reader: %w", err))
	}
	return string(b)
}

func (FfiConverterstring) read(reader io.Reader) string {
	length := readInt32(reader)
	buffer := make([]byte, length)
	read_length, err := reader.Read(buffer)
	if err != nil {
		panic(err)
	}
	if read_length != int(length) {
		panic(fmt.Errorf("bad read length when reading string, expected %d, read %d", length, read_length))
	}
	return string(buffer)
}

func (FfiConverterstring) lower(value string) C.RustBuffer {
	return stringToCRustBuffer(value)
}

func (FfiConverterstring) write(writer io.Writer, value string) {
	if len(value) > math.MaxInt32 {
		panic("String is too large to fit into Int32")
	}

	writeInt32(writer, int32(len(value)))
	write_length, err := io.WriteString(writer, value)
	if err != nil {
		panic(err)
	}
	if write_length != len(value) {
		panic(fmt.Errorf("bad write length when writing string, expected %d, written %d", len(value), write_length))
	}
}

type FfiDestroyerstring struct{}

func (FfiDestroyerstring) destroy(_ string) {}

type FfiConverterSequencestring struct{}

var FfiConverterSequencestringINSTANCE = FfiConverterSequencestring{}

func (c FfiConverterSequencestring) lift(cRustBuf C.RustBuffer) []string {
	return liftFromRustBuffer[[]string](c, fromCRustBuffer(cRustBuf))
}

func (c FfiConverterSequencestring) read(reader io.Reader) []string {
	length := readInt32(reader)
	if length == 0 {
		return nil
	}
	result := make([]string, 0, length)
	for i := int32(0); i < length; i++ {
		result = append(result, FfiConverterstringINSTANCE.read(reader))
	}
	return result
}

func (c FfiConverterSequencestring) lower(value []string) C.RustBuffer {
	return lowerIntoRustBuffer[[]string](c, value)
}

func (c FfiConverterSequencestring) write(writer io.Writer, value []string) {
	if len(value) > math.MaxInt32 {
		panic("[]string is too large to fit into Int32")
	}

	writeInt32(writer, int32(len(value)))
	for _, item := range value {
		FfiConverterstringINSTANCE.write(writer, item)
	}
}

type FfiDestroyerSequencestring struct{}

func (FfiDestroyerSequencestring) destroy(sequence []string) {
	for _, value := range sequence {
		FfiDestroyerstring{}.destroy(value)
	}
}

func GetPlanProviders(path string, envs []string) []string {

	return FfiConverterSequencestringINSTANCE.lift(rustCall(func(_uniffiStatus *C.RustCallStatus) C.RustBuffer {
		return C.nixpacks_65e0_get_plan_providers(FfiConverterstringINSTANCE.lower(path), FfiConverterSequencestringINSTANCE.lower(envs), _uniffiStatus)
	}))

}
