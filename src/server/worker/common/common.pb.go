// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: server/worker/common/common.proto

package common

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	pfs "github.com/pachyderm/pachyderm/src/pfs"
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

type Input struct {
	FileInfo             *pfs.FileInfo `protobuf:"bytes,1,opt,name=file_info,json=fileInfo,proto3" json:"file_info,omitempty"`
	ParentCommit         *pfs.Commit   `protobuf:"bytes,5,opt,name=parent_commit,json=parentCommit,proto3" json:"parent_commit,omitempty"`
	Name                 string        `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	JoinOn               string        `protobuf:"bytes,8,opt,name=join_on,json=joinOn,proto3" json:"join_on,omitempty"`
	GroupBy              string        `protobuf:"bytes,10,opt,name=group_by,json=groupBy,proto3" json:"group_by,omitempty"`
	Lazy                 bool          `protobuf:"varint,3,opt,name=lazy,proto3" json:"lazy,omitempty"`
	Branch               string        `protobuf:"bytes,4,opt,name=branch,proto3" json:"branch,omitempty"`
	GitURL               string        `protobuf:"bytes,6,opt,name=git_url,json=gitUrl,proto3" json:"git_url,omitempty"`
	EmptyFiles           bool          `protobuf:"varint,7,opt,name=empty_files,json=emptyFiles,proto3" json:"empty_files,omitempty"`
	S3                   bool          `protobuf:"varint,9,opt,name=s3,proto3" json:"s3,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Input) Reset()         { *m = Input{} }
func (m *Input) String() string { return proto.CompactTextString(m) }
func (*Input) ProtoMessage()    {}
func (*Input) Descriptor() ([]byte, []int) {
	return fileDescriptor_91fb6c79ddd9db74, []int{0}
}
func (m *Input) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Input) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Input.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Input) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Input.Merge(m, src)
}
func (m *Input) XXX_Size() int {
	return m.Size()
}
func (m *Input) XXX_DiscardUnknown() {
	xxx_messageInfo_Input.DiscardUnknown(m)
}

var xxx_messageInfo_Input proto.InternalMessageInfo

func (m *Input) GetFileInfo() *pfs.FileInfo {
	if m != nil {
		return m.FileInfo
	}
	return nil
}

func (m *Input) GetParentCommit() *pfs.Commit {
	if m != nil {
		return m.ParentCommit
	}
	return nil
}

func (m *Input) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Input) GetJoinOn() string {
	if m != nil {
		return m.JoinOn
	}
	return ""
}

func (m *Input) GetGroupBy() string {
	if m != nil {
		return m.GroupBy
	}
	return ""
}

func (m *Input) GetLazy() bool {
	if m != nil {
		return m.Lazy
	}
	return false
}

func (m *Input) GetBranch() string {
	if m != nil {
		return m.Branch
	}
	return ""
}

func (m *Input) GetGitURL() string {
	if m != nil {
		return m.GitURL
	}
	return ""
}

func (m *Input) GetEmptyFiles() bool {
	if m != nil {
		return m.EmptyFiles
	}
	return false
}

func (m *Input) GetS3() bool {
	if m != nil {
		return m.S3
	}
	return false
}

func init() {
	proto.RegisterType((*Input)(nil), "common.Input")
}

func init() { proto.RegisterFile("server/worker/common/common.proto", fileDescriptor_91fb6c79ddd9db74) }

var fileDescriptor_91fb6c79ddd9db74 = []byte{
	// 347 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0xb1, 0x4e, 0xf3, 0x30,
	0x10, 0xc7, 0x95, 0x7c, 0x6d, 0x92, 0x5e, 0xbf, 0x32, 0x58, 0x08, 0x4c, 0x87, 0xb6, 0xc0, 0x52,
	0x31, 0x34, 0x88, 0x0e, 0xec, 0x45, 0x80, 0x2a, 0x21, 0x21, 0x45, 0xea, 0xc2, 0x12, 0x25, 0xc1,
	0x49, 0x0d, 0x89, 0x6d, 0xd9, 0x0e, 0x28, 0xbc, 0x14, 0xaf, 0xc1, 0xc8, 0x13, 0x20, 0x94, 0x27,
	0x41, 0x76, 0x3a, 0x30, 0x30, 0xe5, 0xf7, 0xff, 0xdd, 0xe5, 0xac, 0xb3, 0xe1, 0x58, 0x11, 0xf9,
	0x42, 0x64, 0xf8, 0xca, 0xe5, 0x33, 0x91, 0x61, 0xc6, 0xab, 0x8a, 0xb3, 0xdd, 0x67, 0x21, 0x24,
	0xd7, 0x1c, 0x79, 0x5d, 0x1a, 0x8f, 0x44, 0xae, 0x42, 0x91, 0xab, 0x4e, 0x8f, 0xf7, 0x0b, 0x5e,
	0x70, 0x8b, 0xa1, 0xa1, 0xce, 0x9e, 0xbc, 0xbb, 0xd0, 0x5f, 0x33, 0x51, 0x6b, 0x74, 0x06, 0x83,
	0x9c, 0x96, 0x24, 0xa6, 0x2c, 0xe7, 0xd8, 0x99, 0x39, 0xf3, 0xe1, 0xc5, 0x68, 0x61, 0x7e, 0xbf,
	0xa1, 0x25, 0x59, 0xb3, 0x9c, 0x47, 0x41, 0xbe, 0x23, 0x74, 0x0e, 0x23, 0x91, 0x48, 0xc2, 0x74,
	0x6c, 0xce, 0xa2, 0x1a, 0xf7, 0x6d, 0xff, 0xd0, 0xf6, 0x5f, 0x59, 0x15, 0xfd, 0xef, 0x3a, 0xba,
	0x84, 0x10, 0xf4, 0x58, 0x52, 0x11, 0xec, 0xce, 0x9c, 0xf9, 0x20, 0xb2, 0x8c, 0x0e, 0xc1, 0x7f,
	0xe2, 0x94, 0xc5, 0x9c, 0xe1, 0xc0, 0x6a, 0xcf, 0xc4, 0x7b, 0x86, 0x8e, 0x20, 0x28, 0x24, 0xaf,
	0x45, 0x9c, 0x36, 0x18, 0x6c, 0xc5, 0xb7, 0x79, 0xd5, 0x98, 0x39, 0x65, 0xf2, 0xd6, 0xe0, 0x7f,
	0x33, 0x67, 0x1e, 0x44, 0x96, 0xd1, 0x01, 0x78, 0xa9, 0x4c, 0x58, 0xb6, 0xc5, 0xbd, 0x6e, 0x4c,
	0x97, 0xd0, 0x29, 0xf8, 0x05, 0xd5, 0x71, 0x2d, 0x4b, 0xec, 0x99, 0xc2, 0x0a, 0xda, 0xaf, 0xa9,
	0x77, 0x4b, 0xf5, 0x26, 0xba, 0x8b, 0xbc, 0x82, 0xea, 0x8d, 0x2c, 0xd1, 0x14, 0x86, 0xa4, 0x12,
	0xba, 0x89, 0xcd, 0x72, 0x0a, 0xfb, 0x76, 0x2e, 0x58, 0x65, 0x16, 0x57, 0x68, 0x0f, 0x5c, 0xb5,
	0xc4, 0x03, 0xeb, 0x5d, 0xb5, 0x5c, 0x5d, 0x7f, 0xb4, 0x13, 0xe7, 0xb3, 0x9d, 0x38, 0xdf, 0xed,
	0xc4, 0x79, 0xb8, 0x2c, 0xa8, 0xde, 0xd6, 0xe9, 0x22, 0xe3, 0x55, 0x28, 0x92, 0x6c, 0xdb, 0x3c,
	0x12, 0xf9, 0x9b, 0x94, 0xcc, 0xc2, 0xbf, 0x9e, 0x2c, 0xf5, 0xec, 0xfd, 0x2f, 0x7f, 0x02, 0x00,
	0x00, 0xff, 0xff, 0x1c, 0x06, 0xad, 0xf5, 0xd1, 0x01, 0x00, 0x00,
}

func (m *Input) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Input) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Input) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.GroupBy) > 0 {
		i -= len(m.GroupBy)
		copy(dAtA[i:], m.GroupBy)
		i = encodeVarintCommon(dAtA, i, uint64(len(m.GroupBy)))
		i--
		dAtA[i] = 0x52
	}
	if m.S3 {
		i--
		if m.S3 {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x48
	}
	if len(m.JoinOn) > 0 {
		i -= len(m.JoinOn)
		copy(dAtA[i:], m.JoinOn)
		i = encodeVarintCommon(dAtA, i, uint64(len(m.JoinOn)))
		i--
		dAtA[i] = 0x42
	}
	if m.EmptyFiles {
		i--
		if m.EmptyFiles {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x38
	}
	if len(m.GitURL) > 0 {
		i -= len(m.GitURL)
		copy(dAtA[i:], m.GitURL)
		i = encodeVarintCommon(dAtA, i, uint64(len(m.GitURL)))
		i--
		dAtA[i] = 0x32
	}
	if m.ParentCommit != nil {
		{
			size, err := m.ParentCommit.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintCommon(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Branch) > 0 {
		i -= len(m.Branch)
		copy(dAtA[i:], m.Branch)
		i = encodeVarintCommon(dAtA, i, uint64(len(m.Branch)))
		i--
		dAtA[i] = 0x22
	}
	if m.Lazy {
		i--
		if m.Lazy {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintCommon(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x12
	}
	if m.FileInfo != nil {
		{
			size, err := m.FileInfo.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintCommon(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintCommon(dAtA []byte, offset int, v uint64) int {
	offset -= sovCommon(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Input) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.FileInfo != nil {
		l = m.FileInfo.Size()
		n += 1 + l + sovCommon(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovCommon(uint64(l))
	}
	if m.Lazy {
		n += 2
	}
	l = len(m.Branch)
	if l > 0 {
		n += 1 + l + sovCommon(uint64(l))
	}
	if m.ParentCommit != nil {
		l = m.ParentCommit.Size()
		n += 1 + l + sovCommon(uint64(l))
	}
	l = len(m.GitURL)
	if l > 0 {
		n += 1 + l + sovCommon(uint64(l))
	}
	if m.EmptyFiles {
		n += 2
	}
	l = len(m.JoinOn)
	if l > 0 {
		n += 1 + l + sovCommon(uint64(l))
	}
	if m.S3 {
		n += 2
	}
	l = len(m.GroupBy)
	if l > 0 {
		n += 1 + l + sovCommon(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovCommon(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozCommon(x uint64) (n int) {
	return sovCommon(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Input) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommon
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
			return fmt.Errorf("proto: Input: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Input: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FileInfo", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
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
				return ErrInvalidLengthCommon
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCommon
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.FileInfo == nil {
				m.FileInfo = &pfs.FileInfo{}
			}
			if err := m.FileInfo.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
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
				return ErrInvalidLengthCommon
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommon
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Lazy", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Lazy = bool(v != 0)
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Branch", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
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
				return ErrInvalidLengthCommon
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommon
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Branch = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ParentCommit", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
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
				return ErrInvalidLengthCommon
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCommon
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ParentCommit == nil {
				m.ParentCommit = &pfs.Commit{}
			}
			if err := m.ParentCommit.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GitURL", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
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
				return ErrInvalidLengthCommon
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommon
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GitURL = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EmptyFiles", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.EmptyFiles = bool(v != 0)
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field JoinOn", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
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
				return ErrInvalidLengthCommon
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommon
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.JoinOn = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field S3", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.S3 = bool(v != 0)
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GroupBy", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
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
				return ErrInvalidLengthCommon
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommon
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GroupBy = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommon(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommon
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
func skipCommon(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCommon
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
					return 0, ErrIntOverflowCommon
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
					return 0, ErrIntOverflowCommon
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
				return 0, ErrInvalidLengthCommon
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupCommon
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthCommon
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthCommon        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCommon          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupCommon = fmt.Errorf("proto: unexpected end of group")
)
