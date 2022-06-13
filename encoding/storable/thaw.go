package storable

/*
 * Парсит данные в формате perl Storable и возвращает примитивы и объекты Go
 */

import (
	"bytes"
	"encoding/binary"
	"io"
	"fmt"
	"reflect"
)

type StorableReader struct {
	storable *bytes.Reader  // буффер с данными
	path []reflect.Value // результат
	aseen []*interface{}   // значения распознанные раньше
	aclass []*interface{}  // классы распознанные раньше
	//bless   // классы для распознавания
	//iconv = iconv  // конвертер строк
	codefunc byte
	err error		// ошибка последней операции
}

func read_magic(stream *StorableReader) () {
	// Считывает магическое число

	//err := binary.Read(buf, binary.BigEndian, &myInt)
	magic_bytes := read(stream, 4)
	if stream.err != nil {return}
	if !bytes.Equal(magic_bytes, MAGICSTR_BYTES) {
		stream.storable.Seek(0, io.SeekStart)
	}

	use_network_order := readUInt8(stream)
	if stream.err != nil {return}
	version_major := use_network_order >> 1
	version_minor := byte(0)
	if version_major > 1 {
		version_minor = readUInt8(stream)
		if stream.err != nil {return}
	}

	if version_major > STORABLE_BIN_MAJOR || 
		version_major == STORABLE_BIN_MAJOR &&
		version_minor > STORABLE_BIN_MINOR {
		stream.err = &StorableError{fmt.Sprintf(
			"Версия Storable не совпадает: требуется v%d.%d, а данные имеют версию v%d.%d", 
				STORABLE_BIN_MAJOR, STORABLE_BIN_MINOR, version_major, version_minor), stream.path}
		return
	}

	if use_network_order & 0x1 == 1 {
		return // /* OK */
	}

	length_magic := readUInt8(stream)
	if stream.err != nil {return}
	
	use_NV_size := version_major >= 2 || version_minor >= 2

	buf := read(stream, uint32(length_magic))
	if stream.err != nil {return}

	if !bytes.Equal(buf, BYTEORDER_BYTES) {
		stream.err = &StorableError{fmt.Sprintf("Magic number is not compatible: %s <> %s",	buf, BYTEORDER_BYTES), stream.path}
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
// func seen(stream *StorableReader, sv *interface{}) any {
	// append(stream.aseen, sv)
	// return sv
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
	stream.err = &StorableError{fmt.Sprintf("Структура Storable повреждена. Обработчик № %d", stream.codefunc), stream.path}
}

func retrieve_byte(stream *StorableReader) {
	value := int8(readUInt8(stream) - 128)
	pointer := stream.path[len(stream.path)-1]
	elem := pointer.Elem()
	switch {
	case elem.Kind() == reflect.Interface && elem.NumMethod() == 0:
		//append(stream.aseen, sv)
		elem.Set(reflect.ValueOf(value))
	case elem.CanInt():
		elem.SetInt(int64(value))
	case elem.CanUint():
		elem.SetUint(uint64(value))
	case elem.CanFloat():
		elem.SetFloat(float64(value))
		
	default:
		stream.err = &StorableError{fmt.Sprintf("Байту соответствует %v", elem.Type()), stream.path}
	}
}

// func retrieve_object(stream *StorableReader) {
	// var tag := stream.readInt32BE()
	// if tag < 0 || tag >= len(stream.aseen) {
		// stream.err = fmt.Errorf('Object //%s out of range 0..%d', tag, len(stream.aseen)-1))
	// }
	// return stream.aseen[tag]
// }

// func retrieve_integer(stream *StorableReader) {
	// stream.pos += SIZE_OF_LONG
	// return stream.seen(struct.unpack_from("<q", stream.storable, stream.pos - SIZE_OF_LONG)[0])
// }

// func retrieve_double(stream *StorableReader) {
	// stream.pos += SIZE_OF_LONG
	// return stream.seen(struct.unpack_from("<d", stream.storable, stream.pos - SIZE_OF_LONG)[0])
// }

// func retrieve_scalar(stream *StorableReader) {
	// size = readUInt8(stream)
	// return stream.seen(get_lstring(stream, size))
// }

// func retrieve_lscalar(stream *StorableReader) {
	// size = stream.readInt32LE()
	// return stream.seen(get_lstring(stream, size))
// }

// func retrieve_utf8str(stream *StorableReader) {
	// size = readUInt8(stream)
	// return stream.seen(stream.read(size).decode('utf8'))
// }

// func retrieve_lutf8str(stream *StorableReader) {
	// size = stream.readInt32LE()
	// return stream.seen(stream.read(size).decode('utf8'))
// }

// func retrieve_array(stream *StorableReader) {
	// size = stream.readInt32LE()
	// array = stream.seen([])
	// for i in range(0, size):
		// array.append(stream.retrieve())

	// return array
// }

// func retrieve_ref(stream *StorableReader) {
	// """ Аналога ссылки perl-а в Python нет, в Python всё ссылки, поэтому возвращаем значение так """
	// stream.seen(nil)
	// sv = stream.retrieve()
	// stream.aseen[-1] = sv
	// return sv
// }

// func retrieve_hash(stream *StorableReader) {
	// length = stream.readInt32LE()
	// hash = stream.seen({})
	// for i in range(0, length):
		// value = stream.retrieve()
		// size = stream.readInt32LE()
		// key = get_lstring(stream, size)
		// hash[key] = value

	// return hash
// }

// func retrieve_flag_hash(stream *StorableReader) {
	// hash_flags = readUInt8(stream)
	// length = stream.readInt32LE()
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
			// size = stream.readInt32LE()
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

	// // делаем класс F одинаковым с классом stream.bless[classname]
	// // объекты класса F будут "instanceof A"
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
		// length = stream.readInt32LE()

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
		// raise PerlStorableException("Повреждена структура Storable: битый индекс в aclass: " + idx)
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

// func readInt32LE(stream *StorableReader) int32 {
	//err := binary.Read(buf, binary.BigEndian, &myInt)
	// stream.pos += SIZE_OF_INT
	// return struct.unpack_from("<i", stream.storable, stream.pos - SIZE_OF_INT)[0]
// }

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
		stream.err = &StorableError{fmt.Sprintf("Неожиданный конец данных: прочитано %d байт, когда требовалось %d.", n, length), stream.path}
		return []byte{}
	}
	return result
}

// func get_lstring(stream *StorableReader, length uint32, in_utf8=False) {
	// if length == 0:
		// return ''
	// s = stream.read(length)

	// if in_utf8:
		// return s.decode('utf8')
	// if is_ascii(s):
		// return s.decode('ascii')
	// if stream.iconv:
		// return stream.iconv(s)
	// return s
// }

func end(stream *StorableReader) {
	buf := make([]byte, 1)
	_, err := stream.storable.Read(buf)
	if err == nil {
		stream.err = &StorableError{"Структура не разобрана до конца", stream.path}
	} else if len(stream.path) == 0 {
		stream.err = &StorableError{"Нет результата", stream.path}
	} else if len(stream.path) > 1 {
		stream.err = &StorableError{"В пути осталось несколько результатов", stream.path}
	}
}

var RETRIVE_METHOD = []func(*StorableReader)() {
retrieve_other, // retrieve_object,  // /* SX_OBJECT -- entry unused dynamically */
retrieve_other, // retrieve_lscalar,  // /* SX_LSCALAR */
retrieve_other, // retrieve_array,  // /* SX_ARRAY */
retrieve_other, // retrieve_hash,  // /* SX_HASH */
retrieve_other, // retrieve_ref,  // /* SX_REF */
retrieve_other, // retrieve_undef,  // /* SX_UNDEF */
retrieve_other, // retrieve_integer,  // /* SX_INTEGER */
retrieve_other, // retrieve_double,  // /* SX_DOUBLE */
    retrieve_byte,  // /* SX_BYTE */
retrieve_other, // retrieve_netint,  // /* SX_NETINT */
retrieve_other, // retrieve_scalar,  // /* SX_SCALAR */
retrieve_other, // retrieve_tied_array,  // /* SX_TIED_ARRAY */
retrieve_other, // retrieve_tied_hash,  // /* SX_TIED_HASH */
retrieve_other, // retrieve_tied_scalar,  // /* SX_TIED_SCALAR */
retrieve_other, // retrieve_sv_undef,  // /* SX_SV_UNDEF */
retrieve_other, // retrieve_sv_yes,  // /* SX_SV_YES */
retrieve_other, // retrieve_sv_no,  // /* SX_SV_NO */
retrieve_other, // retrieve_blessed,  // /* SX_BLESS */
retrieve_other, // retrieve_idx_blessed,  // /* SX_IX_BLESS */
retrieve_other, // retrieve_hook,  // /* SX_HOOK */
retrieve_other, // retrieve_overloaded,  // /* SX_OVERLOAD */
retrieve_other, // retrieve_tied_key,  // /* SX_TIED_KEY */
retrieve_other, // retrieve_tied_idx,  // /* SX_TIED_IDX */
retrieve_other, // retrieve_utf8str,  // /* SX_UTF8STR  */
retrieve_other, // retrieve_lutf8str,  // /* SX_LUTF8STR */
retrieve_other, // retrieve_flag_hash,  // /* SX_FLAG_HASH */
retrieve_other, // retrieve_code,  // /* SX_CODE */
retrieve_other, // retrieve_weakref,  // /* SX_WEAKREF */
retrieve_other, // retrieve_weakoverloaded,  // /* SX_WEAKOVERLOAD */
retrieve_other, // retrieve_vstring,  // /* SX_VSTRING */
retrieve_other, // retrieve_lvstring,  // /* SX_LVSTRING */
retrieve_other, // retrieve_svundef_elem,  // /* SX_SVUNDEF_ELEM */
retrieve_other, // retrieve_regexp,  // /* SX_REGEXP */
retrieve_other, // retrieve_lobject,  // /* SX_LOBJECT */
    retrieve_other,  // /* SX_LAST */
}


func Unmarshal(storable []byte, result any) error {
	if result == nil {
		return &StorableError{
			message: "Требуется указатель на структуру или интерфейс, а получен nil",
		}
	}
	
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &StorableError{
			fmt.Sprintf("Требуется указатель на структуру или интерфейс, а получен %v", rv.Type()),
			[]reflect.Value{rv},
		}
	}

	reader := bytes.NewReader(storable)
    stream := &StorableReader{
		storable: reader,
		path: []reflect.Value{rv},
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
