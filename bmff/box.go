package bmff

import (
	"errors"
	"fmt"
)

// ErrUnknownParser is returned by Box.Parse for unrecognized box parser.
var ErrUnknownParser = errors.New("error no parser for box")

// BoxType is an ISOBMFF box
type BoxType uint8

// Common box types.
const (
	TypeUnknown BoxType = iota
	TypeAuxC            // 'auxC'
	TypeAuxl            // 'auxl'
	TypeAv01            // 'av01'
	TypeAv1C            // 'av1C'
	TypeCdsc            // 'cdsc'
	TypeClap            // 'clap'
	TypeColr            // 'colr'
	TypeDimg            // 'dimg'
	TypeDinf            // 'dinf'
	TypeFtyp            // 'ftyp'
	TypeHdlr            // 'hdlr'
	TypeHvcC            // 'hvcC'
	TypeIdat            // 'idat'
	TypeIinf            // 'iinf'
	TypeIloc            // 'iloc'
	TypeImir            // 'imir'
	TypeInfe            // 'infe'
	TypeIovl            // 'iovl'
	TypeIpco            // 'ipco
	TypeIpma            // 'ipma'
	TypeIprp            // 'iprp'
	TypeIref            // 'iref'
	TypeIrot            // 'irot'
	TypeIspe            // 'ispe'
	TypeMeta            // 'meta'
	TypeMoov            // 'moov'
	TypePasp            // 'pasp'
	TypePitm            // 'pitm'
	TypePixi            // 'pixi'
	TypeThmb            // 'thmb'
	TypeTrak            // 'trak'
	TypeUUID            // 'uuid'
)

var mapStringBoxType = map[string]BoxType{
	"auxC": TypeAuxC,
	"auxl": TypeAuxl,
	"av01": TypeAv01,
	"av1C": TypeAv1C,
	"cdsc": TypeCdsc,
	"clap": TypeClap,
	"colr": TypeColr,
	"dimg": TypeDimg,
	"dinf": TypeDinf,
	"ftyp": TypeFtyp,
	"hdlr": TypeHdlr,
	"hvcC": TypeHvcC,
	"idat": TypeIdat,
	"iinf": TypeIinf,
	"iloc": TypeIloc,
	"imir": TypeImir,
	"infe": TypeInfe,
	"iovl": TypeIovl,
	"ipco": TypeIpco,
	"ipma": TypeIpma,
	"iprp": TypeIprp,
	"iref": TypeIref,
	"irot": TypeIrot,
	"ispe": TypeIspe,
	"meta": TypeMeta,
	"moov": TypeMoov,
	"pasp": TypePasp,
	"pitm": TypePitm,
	"pixi": TypePixi,
	"thmb": TypeThmb,
	"trak": TypeTrak,
	"uuid": TypeUUID,
}

var mapBoxTypeString = map[BoxType]string{
	TypeAuxC: "auxC",
	TypeAuxl: "auxl",
	TypeAv01: "av01",
	TypeAv1C: "av1C",
	TypeCdsc: "cdsc",
	TypeClap: "clap",
	TypeColr: "colr",
	TypeDimg: "dimg",
	TypeDinf: "dinf",
	TypeFtyp: "ftyp",
	TypeHdlr: "hdlr",
	TypeHvcC: "hvcC",
	TypeIdat: "idat",
	TypeIinf: "iinf",
	TypeIloc: "iloc",
	TypeImir: "imir",
	TypeInfe: "infe",
	TypeIovl: "iovl",
	TypeIpco: "ipco",
	TypeIpma: "ipma",
	TypeIprp: "iprp",
	TypeIref: "iref",
	TypeIrot: "irot",
	TypeIspe: "ispe",
	TypeMeta: "meta",
	TypeMoov: "moov",
	TypePasp: "pasp",
	TypePitm: "pitm",
	TypePixi: "pixi",
	TypeThmb: "thmb",
	TypeTrak: "trak",
	TypeUUID: "uuid",
}

func (t BoxType) String() string {
	str, ok := mapBoxTypeString[t]
	if ok {
		return str
	}
	return "nnnn"
}

func isBoxinfe(buf []byte) bool {
	return buf[0] == 'i' && buf[1] == 'n' && buf[2] == 'f' && buf[3] == 'e'
}

func boxType(buf []byte) BoxType {
	if isBoxinfe(buf) { // inital check for performance reasons, TODO: confirm with benchmarks
		return TypeInfe
	}
	b, ok := mapStringBoxType[string(buf)]
	if ok {
		return b
	}
	if Debug {
		fmt.Println(string(buf))
	}
	return TypeUnknown
}

type box struct {
	bufReader
	size    int64 // 0 means unknown, will read to end of file (box container)
	boxType BoxType
}

func (b box) String() string {
	return fmt.Sprintf("(Box) type: \"%s\", size: %d", b.boxType, b.size)
}

func (b box) Size() int64   { return b.size }
func (b box) Type() BoxType { return b.boxType }

func (b *box) Parse() (Box, error) {
	if !b.anyRemain() {
		return nil, ErrBufReaderLength
	}
	parser, ok := parsers[b.Type()]
	if !ok {
		return UnknownBox{t: b.Type(), s: b.size}, ErrUnknownParser
	}

	v, err := parser(b)
	if err != nil {
		// Write error with parser
		return nil, err
	}
	return v, nil
}

// Parsers
type parserFunc func(b *box) (Box, error)

var parsers map[BoxType]parserFunc

func init() {
	parsers = map[BoxType]parserFunc{
		TypeDinf: parseDataInformationBox,
		//boxType("dref"): parseDataReferenceBox,
		TypeFtyp: parseFtyp,
		TypeHdlr: parseHdlr,
		TypeIinf: parseIinf,
		TypeInfe: parseInfe,
		TypeIloc: parseIloc,
		TypeIpco: parseIpco,
		TypeIpma: parseIpma,
		TypeIprp: parseIprp,
		TypeIref: parseIref,
		TypeIrot: parseImageRotation,
		TypeIspe: parseImageSpatialExtentsProperty,
		TypeMeta: parseMeta,
		TypePitm: parsePitm,
	}
}

// DataInformationBox is a "dinf" box
type DataInformationBox struct {
	Children []Box
}

// Type returns TypeDinf
func (dinf DataInformationBox) Type() BoxType {
	return TypeDinf
}

func parseDataInformationBox(outer *box) (Box, error) {
	dib := DataInformationBox{}
	err := outer.discard(outer.remain)
	return dib, err //br.parseAppendBoxes(&dib.Children)
}

// UnknownBox is a box that was unable to be parsed.
type UnknownBox struct {
	t BoxType
	s int64
}

// Type returns the BoxType of the UnknownBox
func (ub UnknownBox) Type() BoxType {
	return ub.t
}

// Size returns the Size of the UnknownBox
func (ub UnknownBox) Size() int64 {
	return ub.s
}

func (ub UnknownBox) String() string {
	return fmt.Sprintf(" Type: %s, Size: %d", ub.t, ub.s)
}