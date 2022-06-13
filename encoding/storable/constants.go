package storable

var (
	MAGICSTR_BYTES = []byte("pst0")
	BYTEORDER_BYTES = []byte("12345678")
)

const (

	STORABLE_BIN_MAJOR = 2  // Binary major "version"
	STORABLE_BIN_MINOR = 11  // Binary minor "version"

	SIZE_OF_INT = 4
	SIZE_OF_LONG = 8
	SIZE_OF_CHAR_PTR = 8
	SIZE_OF_NV = 8

	LG_BLESS = 127	// Large classname bless limit

	//define svis_REF		0
	// svis_SCALAR	= 1
	//define svis_ARRAY		2
	//define svis_HASH		3
	//define svis_TIED		4
	//define svis_TIED_ITEM		5
	//define svis_CODE		6
	//define svis_REGEXP		7
	//define svis_OTHER		8

	SHV_RESTRICTED = 0x01
	SHV_K_UTF8 = 0x01
	SHV_K_WASUTF8 = 0x02
	SHV_K_LOCKED = 0x04
	SHV_K_ISSV = 0x08
	// SHV_K_PLACEHOLDER = 0x10

	SX_OBJECT = 0 // Already stored object
	SX_LSCALAR = 1 // Scalar (large binary) follows (length, data)
	SX_ARRAY = 2 // Array forthcoming (size, item list)
	SX_HASH = 3 // Hash forthcoming (size, key/value pair list)
	SX_REF = 4 // Reference to object forthcoming
	SX_UNDEF = 5 // Undefined scalar
	SX_INTEGER = 6 // Integer forthcoming
	SX_DOUBLE = 7 // Double forthcoming
	SX_BYTE = 8 // (signed) byte forthcoming
	SX_NETINT = 9 // Integer in network order forthcoming
	SX_SCALAR = 10 // Scalar (binary, small) follows (length, data)
	SX_TIED_ARRAY = 11 // Tied array forthcoming
	SX_TIED_HASH = 12 // Tied hash forthcoming
	SX_TIED_SCALAR = 13 // Tied scalar forthcoming
	SX_SV_UNDEF = 14 // Perl's immortal PL_sv_undef
	SX_SV_YES = 15 // Perl's immortal PL_sv_yes
	SX_SV_NO = 16 // Perl's immortal PL_sv_no
	SX_BLESS = 17 // Object is blessed
	SX_IX_BLESS = 18 // Object is blessed, classname given by index
	SX_HOOK = 19 // Stored via hook, user-defined
	SX_OVERLOAD = 20 // Overloaded reference
	SX_TIED_KEY = 21 // Tied magic key forthcoming
	SX_TIED_IDX = 22 // Tied magic index forthcoming
	SX_UTF8STR = 23 // UTF-8 string forthcoming (small)
	SX_LUTF8STR = 24 // UTF-8 string forthcoming (large)
	SX_FLAG_HASH = 25 // Hash with flags forthcoming (size, flags, key/flags/value triplist)
	SX_CODE = 26 // Code references as perl source code
	SX_WEAKREF = 27 // Weak reference to object forthcoming
	SX_WEAKOVERLOAD = 28 // Overloaded weak reference
	SX_VSTRING = 29 // vstring forthcoming (small)
	SX_LVSTRING = 30 // vstring forthcoming (large)
	SX_SVUNDEF_ELEM = 31 // array element set to &PL_sv_undef
	SX_REGEXP = 32 // Regexp
	SX_LOBJECT = 33 // Large object: string, array or hash (size >2G)
	SX_LAST = 34  // invalid. marker only
)

var SX = map[uint8]string {
	0: "SX_OBJECT", // Already stored object
	1: "SX_LSCALAR", // Scalar (large binary) follows (length, data)
	2: "SX_ARRAY", // Array forthcoming (size, item list)
	3: "SX_HASH", // Hash forthcoming (size, key/value pair list)
	4: "SX_REF", // Reference to object forthcoming
	5: "SX_UNDEF", // Undefined scalar
	6: "SX_INTEGER", // Integer forthcoming
	7: "SX_DOUBLE", // Double forthcoming
	8: "SX_BYTE", // (signed) byte forthcoming
	9: "SX_NETINT", // Integer in network order forthcoming
	10:" SX_SCALAR", // Scalar (binary, small) follows (length, data)
	11: "SX_TIED_ARRAY", // Tied array forthcoming
	12: "SX_TIED_HASH", // Tied hash forthcoming
	13: "SX_TIED_SCALAR", // Tied scalar forthcoming
	14: "SX_SV_UNDEF", // Perl"s immortal PL_sv_undef
	15: "SX_SV_YES", // Perl"s immortal PL_sv_yes
	16: "SX_SV_NO", // Perl"s immortal PL_sv_no
	17: "SX_BLESS", // Object is blessed
	18: "SX_IX_BLESS", // Object is blessed, classname given by index
	19: "SX_HOOK", // Stored via hook, user-defined
	20: "SX_OVERLOAD", // Overloaded reference
	21: "SX_TIED_KEY", // Tied magic key forthcoming
	22: "SX_TIED_IDX", // Tied magic index forthcoming
	23: "SX_UTF8STR", // UTF-8 string forthcoming (small)
	24: "SX_LUTF8STR", // UTF-8 string forthcoming (large)
	25: "SX_FLAG_HASH", // Hash with flags forthcoming (size, flags, key/flags/value triplist)
	26: "SX_CODE", // Code references as perl source code
	27: "SX_WEAKREF", // Weak reference to object forthcoming
	28: "SX_WEAKOVERLOAD", // Overloaded weak reference
	29: "SX_VSTRING", // vstring forthcoming (small)
	30: "SX_LVSTRING", // vstring forthcoming (large)
	31: "SX_SVUNDEF_ELEM", // array element set to &PL_sv_undef
	32: "SX_REGEXP", // Regexp
	33: "SX_LOBJECT", // Large object: string, array or hash (size >2G)
	34: "SX_LAST",  // invalid. marker onl"y"
} 

/*
def is_ascii(buffer):
    for c in buffer:
        if c > 127:
            return False
    return True
*/