package storable

/*
 * Parses data in perl Storable format and returns primitives and Go objects
 */

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

type StorableReader struct {
	storable *bytes.Reader   // data buffer
	path     []reflect.Value // result
	aseen    []*interface{}  // values recognized earlier
	aclass   []*interface{}  // classes recognized earlier
	//bless   		 // classes for recognition
	//iconv = iconv  // string converter
	codefunc byte
	err      error // last operation error
}

func read_magic(stream *StorableReader) {
	// Reads the magic number

	//err := binary.Read(buf, binary.BigEndian, &myInt)
	magic_bytes := read(stream, 4)
	if stream.err != nil {
		return
	}
	if !bytes.Equal(magic_bytes, MAGICSTR_BYTES) {
		stream.storable.Seek(0, io.SeekStart)
	}

	use_network_order := readUInt8(stream)
	if stream.err != nil {
		return
	}
	version_major := use_network_order >> 1
	version_minor := byte(0)
	if version_major > 1 {
		version_minor = readUInt8(stream)
		if stream.err != nil {
			return
		}
	}

	if version_major > STORABLE_BIN_MAJOR ||
		version_major == STORABLE_BIN_MAJOR &&
			version_minor > STORABLE_BIN_MINOR {
		stream.err = &StorableError{fmt.Sprintf(
			"Storable version mismatch: v%d.%d required, but data is v%d.%d",
			STORABLE_BIN_MAJOR, STORABLE_BIN_MINOR, version_major, version_minor), stream.path}
		return
	}

	if use_network_order&0x1 == 1 {
		return // /* OK */
	}

	length_magic := readUInt8(stream)
	if stream.err != nil {
		return
	}

	use_NV_size := version_major >= 2 || version_minor >= 2

	buf := read(stream, uint32(length_magic))
	if stream.err != nil {
		return
	}

	if !bytes.Equal(buf, BYTEORDER_BYTES) {
		stream.err = &StorableError{fmt.Sprintf("Magic number is not compatible: %s <> %s", buf, BYTEORDER_BYTES), stream.path}
		return
	}

	if readInt8(stream) != SIZE_OF_INT {
		stream.err = &StorableError{"Integer size is not compatible", stream.path}
		return
	}
	if readInt8(stream) != SIZE_OF_LONG {
		stream.err = &StorableError{"Long integer size is not compatible", stream.path}
		return
	}
	if readInt8(stream) != SIZE_OF_CHAR_PTR {
		stream.err = &StorableError{"Pointer size is not compatible", stream.path}
		return
	}
	if use_NV_size {
		if readInt8(stream) != SIZE_OF_NV {
			stream.err = &StorableError{"Double size is not compatible", stream.path}
			return
		}
	}
}

// // Сохраняет в aseen извлечённое из буфера значение и возвращает его.
// func seen(stream *StorableReader, sv *interface{}) {
// stream.aseen := append(stream.aseen, sv)
// }

func retrieve(stream *StorableReader) {
	// Считывает структуру рекурсивно.
	index := readUInt8(stream)
	stream.codefunc = index
	if int(index) >= len(RETRIVE_METHOD) {
		retrieve_other(stream)
		return
	}
	RETRIVE_METHOD[index](stream)
}

func retrieve_other(stream *StorableReader) {
	stream.err = &StorableError{fmt.Sprintf("Storable structure is corrupted. Handler N°%d", stream.codefunc), stream.path}
}

func retrieve_byte(stream *StorableReader) {
	value := int8(readUInt8(stream) - 128)
	//seen(stream, &value)
	pointer := stream.path[len(stream.path)-1]
	elem := pointer.Elem()
	switch {
	case elem.Kind() == reflect.Interface && elem.NumMethod() == 0:
		elem.Set(reflect.ValueOf(value))
	case elem.CanInt():
		elem.SetInt(int64(value))
	case elem.CanUint():
		elem.SetUint(uint64(value))
	case elem.CanFloat():
		elem.SetFloat(float64(value))
	default:
		stream.err = &StorableError{fmt.Sprintf("The byte corresponds to %v", elem.Type()), stream.path}
	}

}

// func retrieve_object(stream *StorableReader) {
// var tag := stream.readInt32BE()
// if tag < 0 || tag >= len(stream.aseen) {
// stream.err = fmt.Errorf('Object //%s out of range 0..%d', tag, len(stream.aseen)-1))
// }
// return stream.aseen[tag]
// }

func retrieve_integer(stream *StorableReader) {
	var value int64
	err := binary.Read(stream.storable, binary.LittleEndian, &value)
	if err != nil {
		stream.err = &StorableError{fmt.Sprintf("read int64: %v", err), stream.path}
		return
	}
	//seen(stream, &value)
	pointer := stream.path[len(stream.path)-1]
	elem := pointer.Elem()
	switch {
	case elem.Kind() == reflect.Interface && elem.NumMethod() == 0:
		elem.Set(reflect.ValueOf(value))
	case elem.CanInt():
		elem.SetInt(int64(value))
	case elem.CanUint():
		elem.SetUint(uint64(value))
	case elem.CanFloat():
		elem.SetFloat(float64(value))
	default:
		stream.err = &StorableError{fmt.Sprintf("Corresponds to the whole %v", elem.Type()), stream.path}
	}
}

func retrieve_double(stream *StorableReader) {
	var value float64
	err := binary.Read(stream.storable, binary.LittleEndian, &value)
	if err != nil {
		stream.err = &StorableError{fmt.Sprintf("read float64: %v", err), stream.path}
		return
	}
	//seen(stream, &value)
	pointer := stream.path[len(stream.path)-1]
	elem := pointer.Elem()
	switch {
	case elem.Kind() == reflect.Interface && elem.NumMethod() == 0:
		elem.Set(reflect.ValueOf(value))
	case elem.CanFloat():
		elem.SetFloat(value)
	default:
		stream.err = &StorableError{fmt.Sprintf("Corresponds to the material %v", elem.Type()), stream.path}
	}
}

func get_string(stream *StorableReader, size int32, is_utf8 bool) {
	if stream.err != nil {
		return
	}
	value := get_lstring(stream, size)
	if stream.err != nil {
		return
	}
	//seen(stream, &value)
	pointer := stream.path[len(stream.path)-1]
	elem := pointer.Elem()
	switch {
	case elem.Kind() == reflect.Interface && elem.NumMethod() == 0:
		if is_utf8 {
			elem.Set(reflect.ValueOf(string(value)))
		} else {
			elem.Set(reflect.ValueOf(value))
		}
	case elem.Kind() == reflect.Slice && elem.Type().Elem().Kind() == reflect.Uint8:
		elem.SetBytes(value)
	case elem.Kind() == reflect.String:
		elem.SetString(string(value))
	default:
		stream.err = &StorableError{fmt.Sprintf("Corresponds to a scalar %v", elem.Type()), stream.path}
	}
}

func retrieve_scalar(stream *StorableReader) {
	get_string(stream, int32(readUInt8(stream)), false)
}

func retrieve_lscalar(stream *StorableReader) {
	get_string(stream, readInt32LE(stream), false)
}

func retrieve_utf8str(stream *StorableReader) {
	get_string(stream, int32(readUInt8(stream)), true)
}

func retrieve_lutf8str(stream *StorableReader) {
	get_string(stream, readInt32LE(stream), true)
}

// func retrieve_array(stream *StorableReader) {
// size = readInt32LE(stream)()
// array = stream.seen([])
// for i in range(0, size):
// array.append(stream.retrieve())

// return array
// }

// func retrieve_ref(stream *StorableReader) {
// """ There is no analogue of a Perl link in Python, in Python everything is links, so we return the value like this """
// stream.seen(nil)
// sv = stream.retrieve()
// stream.aseen[-1] = sv
// return sv
// }

// func retrieve_hash(stream *StorableReader) {
// length = readInt32LE(stream)()
// hash = stream.seen({})
// for i in range(0, length):
// value = stream.retrieve()
// size = readInt32LE(stream)()
// key = get_lstring(stream, size)
// hash[key] = value

// return hash
// }

// func retrieve_flag_hash(stream *StorableReader) {
// hash_flags = readUInt8(stream)
// length = readInt32LE(stream)()
// hash = stream.seen({})

// for i in range(0, length):
// value = stream.retrieve()
// flags = readUInt8(stream)
// key = 0

// if flags & SHV_K_ISSV:
// """ XXX you can't set a placeholder with an SV key.
// Then again, you can't get an SV key.
// Without messing around beyond what the API is supposed to do.
// """
// key = stream.retrieve()
// else:
// size = readInt32LE(stream)()
// key = get_lstring(stream, size, flags & (SHV_K_UTF8 | SHV_K_WASUTF8))

// hash[key] = value

// // if (hash_flags & SHV_RESTRICTED) or (flags & SHV_K_LOCKED):
// //     Object.defineProperty(hash, key, {
// //         value,
// //         writable: false,
// //         configurable: false,
// //         enumerable: true,
// //     })
// // else:
// //     hash[key] = value

// return hash
// }

// func retrieve_weakref(stream *StorableReader) {
// return stream.retrieve_ref()
// }

// func retrieve_undef(stream *StorableReader) {
// return stream.seen(None)
// }

// func retrieve_sv_undef(stream *StorableReader) {
// return stream.seen(None)

// func make_obj(stream, sv, classname):
// classname_python = classname.replace('::', '__')

// // make class F the same as class stream.bless[classname]
// // objects of class F will be "instanceof A"
// a_class = stream.bless[classname] if classname in stream.bless else type(classname_python, (
// type(sv),), {})

// // переписываем свойства
// if isinstance(sv, list):
// obj = a_class.__new__(a_class)
// for v in sv:
// obj.append(v)
// elif isinstance(sv, dict):
// obj = a_class.__new__(a_class)
// for key, val in sv.items():
// setattr(obj, key, val)
// else:
// obj = a_class(sv)
// //setattr(obj, "scalar", sv)

// return obj
// }

// func retrieve_blessed(stream *StorableReader) {
// length = readUInt8(stream)

// if length & 0x80:
// length = readInt32LE(stream)()

// classname = get_lstring(stream, length)
// stream.aclass.append(classname)
// sv = stream.retrieve()
// return stream.make_obj(sv, classname)
// }

// func retrieve_idx_blessed(stream *StorableReader) {
// idx = readUInt8(stream)
// if idx & 0x80:
// idx = this.readInt32LE()
// if idx<0 or idx>=len(stream.aclass):
// raise PerlStorableException("Storable structure is corrupted: broken index in aclass: " + idx)
// classname = this.aclass[idx]
// sv = stream.retrieve()
// return stream.make_obj(sv, classname)
// }

// func retrieve_overloaded(stream *StorableReader) interface{} {
// return stream.retrieve_ref()
// }

func readUInt8(stream *StorableReader) byte {
	var result byte
	err := binary.Read(stream.storable, binary.LittleEndian, &result)
	if err != nil {
		stream.err = &StorableError{fmt.Sprintf("readUInt8: %v", err), stream.path}
		return 0
	}
	return result
}

func readInt32LE(stream *StorableReader) int32 {
	var result int32
	err := binary.Read(stream.storable, binary.LittleEndian, &result)
	if err != nil {
		stream.err = &StorableError{fmt.Sprintf("readInt32LE: %v", err), stream.path}
		return 0
	}
	return result
}

// func readInt32BE(stream *StorableReader) int32 {
// stream.pos += SIZE_OF_INT
// return struct.unpack_from(">i", stream.storable, stream.pos - SIZE_OF_INT)[0]
// }

func readInt8(stream *StorableReader) int8 {
	var result int8
	err := binary.Read(stream.storable, binary.LittleEndian, &result)
	if err != nil {
		stream.err = &StorableError{fmt.Sprintf("readInt8: %v", err), stream.path}
		return 0
	}
	return result
}

func read(stream *StorableReader, length uint32) []byte {
	result := make([]byte, length)
	n, err := stream.storable.Read(result)
	if err != nil {
		stream.err = &StorableError{fmt.Sprintf("read: %v", err), stream.path}
		return []byte{}
	}
	if int64(n) != int64(length) {
		stream.err = &StorableError{fmt.Sprintf("Unexpected end of data: %d bytes read when %d were required.", n, length), stream.path}
		return []byte{}
	}
	return result
}

func get_lstring(stream *StorableReader, length int32) []byte {
	if length == 0 {
		return []byte("")
	}
	s := read(stream, uint32(length))

	// if stream.iconv:
	// return stream.iconv(s)
	return s
}

func end(stream *StorableReader) {
	buf := make([]byte, 1)
	_, err := stream.storable.Read(buf)
	if err == nil {
		stream.err = &StorableError{"The structure is not completely disassembled", stream.path}
	} else if len(stream.path) == 0 {
		stream.err = &StorableError{"No result", stream.path}
	} else if len(stream.path) > 1 {
		stream.err = &StorableError{"There are several results left on the way", stream.path}
	}
}

var RETRIVE_METHOD = []func(*StorableReader){
	retrieve_other,    // retrieve_object,  // /* SX_OBJECT -- entry unused dynamically */
	retrieve_lscalar,  // /* SX_LSCALAR */
	retrieve_other,    // retrieve_array,  // /* SX_ARRAY */
	retrieve_other,    // retrieve_hash,  // /* SX_HASH */
	retrieve_other,    // retrieve_ref,  // /* SX_REF */
	retrieve_other,    // retrieve_undef,  // /* SX_UNDEF */
	retrieve_integer,  // /* SX_INTEGER */
	retrieve_double,   // /* SX_DOUBLE */
	retrieve_byte,     // /* SX_BYTE */
	retrieve_other,    // retrieve_netint,  // /* SX_NETINT */
	retrieve_scalar,   // /* SX_SCALAR */
	retrieve_other,    // retrieve_tied_array,  // /* SX_TIED_ARRAY */
	retrieve_other,    // retrieve_tied_hash,  // /* SX_TIED_HASH */
	retrieve_other,    // retrieve_tied_scalar,  // /* SX_TIED_SCALAR */
	retrieve_other,    // retrieve_sv_undef,  // /* SX_SV_UNDEF */
	retrieve_other,    // retrieve_sv_yes,  // /* SX_SV_YES */
	retrieve_other,    // retrieve_sv_no,  // /* SX_SV_NO */
	retrieve_other,    // retrieve_blessed,  // /* SX_BLESS */
	retrieve_other,    // retrieve_idx_blessed,  // /* SX_IX_BLESS */
	retrieve_other,    // retrieve_hook,  // /* SX_HOOK */
	retrieve_other,    // retrieve_overloaded,  // /* SX_OVERLOAD */
	retrieve_other,    // retrieve_tied_key,  // /* SX_TIED_KEY */
	retrieve_other,    // retrieve_tied_idx,  // /* SX_TIED_IDX */
	retrieve_utf8str,  // /* SX_UTF8STR  */
	retrieve_lutf8str, // /* SX_LUTF8STR */
	retrieve_other,    // retrieve_flag_hash,  // /* SX_FLAG_HASH */
	retrieve_other,    // retrieve_code,  // /* SX_CODE */
	retrieve_other,    // retrieve_weakref,  // /* SX_WEAKREF */
	retrieve_other,    // retrieve_weakoverloaded,  // /* SX_WEAKOVERLOAD */
	retrieve_other,    // retrieve_vstring,  // /* SX_VSTRING */
	retrieve_other,    // retrieve_lvstring,  // /* SX_LVSTRING */
	retrieve_other,    // retrieve_svundef_elem,  // /* SX_SVUNDEF_ELEM */
	retrieve_other,    // retrieve_regexp,  // /* SX_REGEXP */
	retrieve_other,    // retrieve_lobject,  // /* SX_LOBJECT */
	retrieve_other,    // /* SX_LAST */
}

func Unmarshal(storable []byte, result any) error {
	if result == nil {
		return &StorableError{
			message: "Pointer to structure or interface required, but nil returned",
		}
	}

	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &StorableError{
			fmt.Sprintf("Pointer to structure or interface required, but %v was returned", rv.Type()),
			[]reflect.Value{rv},
		}
	}

	reader := bytes.NewReader(storable)
	stream := &StorableReader{
		storable: reader,
		path:     []reflect.Value{rv},
	}
	read_magic(stream)
	if stream.err != nil {
		return stream.err
	}
	retrieve(stream)
	if stream.err != nil {
		return stream.err
	}
	end(stream)
	return stream.err
}
