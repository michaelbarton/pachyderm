// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: internal/storage/fileset/index/index.proto

package index

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	chunk "github.com/pachyderm/pachyderm/v2/src/internal/storage/chunk"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Index stores an index to and metadata about a file.
type Index struct {
	Path                 string   `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	Range                *Range   `protobuf:"bytes,2,opt,name=range,proto3" json:"range,omitempty"`
	File                 *File    `protobuf:"bytes,3,opt,name=file,proto3" json:"file,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Index) Reset()         { *m = Index{} }
func (m *Index) String() string { return proto.CompactTextString(m) }
func (*Index) ProtoMessage()    {}
func (*Index) Descriptor() ([]byte, []int) {
	return fileDescriptor_dfa1b84c403551af, []int{0}
}
func (m *Index) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Index) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Index.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Index) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Index.Merge(m, src)
}
func (m *Index) XXX_Size() int {
	return m.Size()
}
func (m *Index) XXX_DiscardUnknown() {
	xxx_messageInfo_Index.DiscardUnknown(m)
}

var xxx_messageInfo_Index proto.InternalMessageInfo

func (m *Index) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *Index) GetRange() *Range {
	if m != nil {
		return m.Range
	}
	return nil
}

func (m *Index) GetFile() *File {
	if m != nil {
		return m.File
	}
	return nil
}

type Range struct {
	Offset               int64          `protobuf:"varint,1,opt,name=offset,proto3" json:"offset,omitempty"`
	LastPath             string         `protobuf:"bytes,2,opt,name=last_path,json=lastPath,proto3" json:"last_path,omitempty"`
	ChunkRef             *chunk.DataRef `protobuf:"bytes,3,opt,name=chunk_ref,json=chunkRef,proto3" json:"chunk_ref,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Range) Reset()         { *m = Range{} }
func (m *Range) String() string { return proto.CompactTextString(m) }
func (*Range) ProtoMessage()    {}
func (*Range) Descriptor() ([]byte, []int) {
	return fileDescriptor_dfa1b84c403551af, []int{1}
}
func (m *Range) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Range) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Range.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Range) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Range.Merge(m, src)
}
func (m *Range) XXX_Size() int {
	return m.Size()
}
func (m *Range) XXX_DiscardUnknown() {
	xxx_messageInfo_Range.DiscardUnknown(m)
}

var xxx_messageInfo_Range proto.InternalMessageInfo

func (m *Range) GetOffset() int64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *Range) GetLastPath() string {
	if m != nil {
		return m.LastPath
	}
	return ""
}

func (m *Range) GetChunkRef() *chunk.DataRef {
	if m != nil {
		return m.ChunkRef
	}
	return nil
}

type File struct {
	Datum                string           `protobuf:"bytes,1,opt,name=datum,proto3" json:"datum,omitempty"`
	DataRefs             []*chunk.DataRef `protobuf:"bytes,2,rep,name=data_refs,json=dataRefs,proto3" json:"data_refs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *File) Reset()         { *m = File{} }
func (m *File) String() string { return proto.CompactTextString(m) }
func (*File) ProtoMessage()    {}
func (*File) Descriptor() ([]byte, []int) {
	return fileDescriptor_dfa1b84c403551af, []int{2}
}
func (m *File) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *File) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_File.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *File) XXX_Merge(src proto.Message) {
	xxx_messageInfo_File.Merge(m, src)
}
func (m *File) XXX_Size() int {
	return m.Size()
}
func (m *File) XXX_DiscardUnknown() {
	xxx_messageInfo_File.DiscardUnknown(m)
}

var xxx_messageInfo_File proto.InternalMessageInfo

func (m *File) GetDatum() string {
	if m != nil {
		return m.Datum
	}
	return ""
}

func (m *File) GetDataRefs() []*chunk.DataRef {
	if m != nil {
		return m.DataRefs
	}
	return nil
}

func init() {
	proto.RegisterType((*Index)(nil), "index.Index")
	proto.RegisterType((*Range)(nil), "index.Range")
	proto.RegisterType((*File)(nil), "index.File")
}

func init() {
	proto.RegisterFile("internal/storage/fileset/index/index.proto", fileDescriptor_dfa1b84c403551af)
}

var fileDescriptor_dfa1b84c403551af = []byte{
	// 300 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x51, 0xc1, 0x4a, 0xf4, 0x30,
	0x10, 0x26, 0xdb, 0x76, 0xd9, 0x66, 0x7f, 0xfe, 0x43, 0x10, 0x29, 0x0a, 0xb5, 0xf4, 0xb4, 0x28,
	0x34, 0xb0, 0xbe, 0x81, 0x2c, 0xc2, 0xde, 0x24, 0x47, 0x2f, 0x6b, 0xb6, 0x9d, 0xb4, 0xc1, 0x6e,
	0x5b, 0x92, 0x54, 0xf4, 0x0d, 0x3d, 0xfa, 0x08, 0xd2, 0x27, 0x91, 0x24, 0x3d, 0x08, 0x8a, 0x97,
	0x61, 0xbe, 0x99, 0x6f, 0xe6, 0xfb, 0x26, 0xc1, 0xd7, 0xb2, 0x33, 0xa0, 0x3a, 0xde, 0x52, 0x6d,
	0x7a, 0xc5, 0x6b, 0xa0, 0x42, 0xb6, 0xa0, 0xc1, 0x50, 0xd9, 0x55, 0xf0, 0xea, 0x63, 0x31, 0xa8,
	0xde, 0xf4, 0x24, 0x72, 0xe0, 0x22, 0xff, 0x31, 0x52, 0x36, 0x63, 0xf7, 0xec, 0xa3, 0xa7, 0xe6,
	0x4f, 0x38, 0xda, 0x5b, 0x32, 0x21, 0x38, 0x1c, 0xb8, 0x69, 0x12, 0x94, 0xa1, 0x4d, 0xcc, 0x5c,
	0x4e, 0x72, 0x1c, 0x29, 0xde, 0xd5, 0x90, 0x2c, 0x32, 0xb4, 0x59, 0x6f, 0xff, 0x15, 0x5e, 0x84,
	0xd9, 0x1a, 0xf3, 0x2d, 0x72, 0x85, 0x43, 0x6b, 0x24, 0x09, 0x1c, 0x65, 0x3d, 0x53, 0xee, 0x65,
	0x0b, 0xcc, 0x35, 0x72, 0x89, 0x23, 0x37, 0x40, 0xce, 0xf1, 0xb2, 0x17, 0x42, 0x83, 0x71, 0x1a,
	0x01, 0x9b, 0x11, 0xb9, 0xc4, 0x71, 0xcb, 0xb5, 0x39, 0x38, 0xf9, 0x85, 0x93, 0x5f, 0xd9, 0xc2,
	0x83, 0xb5, 0x70, 0x83, 0x63, 0x67, 0xf7, 0xa0, 0x40, 0xcc, 0x1a, 0xff, 0x0b, 0x7f, 0xc0, 0x8e,
	0x1b, 0xce, 0x40, 0xb0, 0x95, 0x83, 0x0c, 0x44, 0xbe, 0xc7, 0xa1, 0x15, 0x26, 0x67, 0x38, 0xaa,
	0xb8, 0x19, 0x4f, 0xf3, 0x31, 0x1e, 0xd8, 0x55, 0x15, 0x37, 0xdc, 0x6e, 0xd2, 0xc9, 0x22, 0x0b,
	0x7e, 0x5b, 0x55, 0xf9, 0x44, 0xdf, 0xb1, 0xf7, 0x29, 0x45, 0x1f, 0x53, 0x8a, 0x3e, 0xa7, 0x14,
	0x3d, 0xee, 0x6a, 0x69, 0x9a, 0xf1, 0x58, 0x94, 0xfd, 0x89, 0x0e, 0xbc, 0x6c, 0xde, 0x2a, 0x50,
	0xdf, 0xb3, 0x97, 0x2d, 0xd5, 0xaa, 0xa4, 0x7f, 0x7f, 0xd1, 0x71, 0xe9, 0x9e, 0xfc, 0xf6, 0x2b,
	0x00, 0x00, 0xff, 0xff, 0x33, 0x84, 0xa8, 0x77, 0xcb, 0x01, 0x00, 0x00,
}

func (m *Index) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Index) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Index) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.File != nil {
		{
			size, err := m.File.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintIndex(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Range != nil {
		{
			size, err := m.Range.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintIndex(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Path) > 0 {
		i -= len(m.Path)
		copy(dAtA[i:], m.Path)
		i = encodeVarintIndex(dAtA, i, uint64(len(m.Path)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Range) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Range) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Range) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.ChunkRef != nil {
		{
			size, err := m.ChunkRef.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintIndex(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.LastPath) > 0 {
		i -= len(m.LastPath)
		copy(dAtA[i:], m.LastPath)
		i = encodeVarintIndex(dAtA, i, uint64(len(m.LastPath)))
		i--
		dAtA[i] = 0x12
	}
	if m.Offset != 0 {
		i = encodeVarintIndex(dAtA, i, uint64(m.Offset))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *File) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *File) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *File) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.DataRefs) > 0 {
		for iNdEx := len(m.DataRefs) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.DataRefs[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintIndex(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Datum) > 0 {
		i -= len(m.Datum)
		copy(dAtA[i:], m.Datum)
		i = encodeVarintIndex(dAtA, i, uint64(len(m.Datum)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintIndex(dAtA []byte, offset int, v uint64) int {
	offset -= sovIndex(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Index) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Path)
	if l > 0 {
		n += 1 + l + sovIndex(uint64(l))
	}
	if m.Range != nil {
		l = m.Range.Size()
		n += 1 + l + sovIndex(uint64(l))
	}
	if m.File != nil {
		l = m.File.Size()
		n += 1 + l + sovIndex(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Range) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Offset != 0 {
		n += 1 + sovIndex(uint64(m.Offset))
	}
	l = len(m.LastPath)
	if l > 0 {
		n += 1 + l + sovIndex(uint64(l))
	}
	if m.ChunkRef != nil {
		l = m.ChunkRef.Size()
		n += 1 + l + sovIndex(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *File) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Datum)
	if l > 0 {
		n += 1 + l + sovIndex(uint64(l))
	}
	if len(m.DataRefs) > 0 {
		for _, e := range m.DataRefs {
			l = e.Size()
			n += 1 + l + sovIndex(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovIndex(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozIndex(x uint64) (n int) {
	return sovIndex(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Index) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIndex
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Index: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Index: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Path", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIndex
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIndex
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthIndex
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Path = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Range", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIndex
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIndex
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthIndex
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Range == nil {
				m.Range = &Range{}
			}
			if err := m.Range.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field File", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIndex
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIndex
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthIndex
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.File == nil {
				m.File = &File{}
			}
			if err := m.File.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIndex(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthIndex
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Range) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIndex
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Range: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Range: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Offset", wireType)
			}
			m.Offset = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIndex
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Offset |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastPath", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIndex
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIndex
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthIndex
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LastPath = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChunkRef", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIndex
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIndex
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthIndex
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ChunkRef == nil {
				m.ChunkRef = &chunk.DataRef{}
			}
			if err := m.ChunkRef.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIndex(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthIndex
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *File) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIndex
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: File: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: File: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Datum", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIndex
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIndex
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthIndex
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Datum = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DataRefs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIndex
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIndex
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthIndex
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DataRefs = append(m.DataRefs, &chunk.DataRef{})
			if err := m.DataRefs[len(m.DataRefs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIndex(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthIndex
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipIndex(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowIndex
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowIndex
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowIndex
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthIndex
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupIndex
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthIndex
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthIndex        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowIndex          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupIndex = fmt.Errorf("proto: unexpected end of group")
)
