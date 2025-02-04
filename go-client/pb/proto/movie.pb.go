// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v3.6.1
// source: proto/movie.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Request to generate an embedding.
type MovieRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MovieId       int32                  `protobuf:"varint,1,opt,name=movie_id,json=movieId,proto3" json:"movie_id,omitempty"`
	Title         string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Overview      string                 `protobuf:"bytes,3,opt,name=overview,proto3" json:"overview,omitempty"`
	Keywords      string                 `protobuf:"bytes,4,opt,name=keywords,proto3" json:"keywords,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MovieRequest) Reset() {
	*x = MovieRequest{}
	mi := &file_proto_movie_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MovieRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MovieRequest) ProtoMessage() {}

func (x *MovieRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_movie_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MovieRequest.ProtoReflect.Descriptor instead.
func (*MovieRequest) Descriptor() ([]byte, []int) {
	return file_proto_movie_proto_rawDescGZIP(), []int{0}
}

func (x *MovieRequest) GetMovieId() int32 {
	if x != nil {
		return x.MovieId
	}
	return 0
}

func (x *MovieRequest) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *MovieRequest) GetOverview() string {
	if x != nil {
		return x.Overview
	}
	return ""
}

func (x *MovieRequest) GetKeywords() string {
	if x != nil {
		return x.Keywords
	}
	return ""
}

// Response containing the embedding.
type EmbeddingResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Embedding     []float32              `protobuf:"fixed32,1,rep,packed,name=embedding,proto3" json:"embedding,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EmbeddingResponse) Reset() {
	*x = EmbeddingResponse{}
	mi := &file_proto_movie_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EmbeddingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmbeddingResponse) ProtoMessage() {}

func (x *EmbeddingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_movie_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmbeddingResponse.ProtoReflect.Descriptor instead.
func (*EmbeddingResponse) Descriptor() ([]byte, []int) {
	return file_proto_movie_proto_rawDescGZIP(), []int{1}
}

func (x *EmbeddingResponse) GetEmbedding() []float32 {
	if x != nil {
		return x.Embedding
	}
	return nil
}

// Request to add an embedding (if you want to separate this functionality).
type AddMovieRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MovieId       int32                  `protobuf:"varint,1,opt,name=movie_id,json=movieId,proto3" json:"movie_id,omitempty"`
	Embedding     []float32              `protobuf:"fixed32,2,rep,packed,name=embedding,proto3" json:"embedding,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddMovieRequest) Reset() {
	*x = AddMovieRequest{}
	mi := &file_proto_movie_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddMovieRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddMovieRequest) ProtoMessage() {}

func (x *AddMovieRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_movie_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddMovieRequest.ProtoReflect.Descriptor instead.
func (*AddMovieRequest) Descriptor() ([]byte, []int) {
	return file_proto_movie_proto_rawDescGZIP(), []int{2}
}

func (x *AddMovieRequest) GetMovieId() int32 {
	if x != nil {
		return x.MovieId
	}
	return 0
}

func (x *AddMovieRequest) GetEmbedding() []float32 {
	if x != nil {
		return x.Embedding
	}
	return nil
}

// A simple response for the add operation.
type AddMovieResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddMovieResponse) Reset() {
	*x = AddMovieResponse{}
	mi := &file_proto_movie_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddMovieResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddMovieResponse) ProtoMessage() {}

func (x *AddMovieResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_movie_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddMovieResponse.ProtoReflect.Descriptor instead.
func (*AddMovieResponse) Descriptor() ([]byte, []int) {
	return file_proto_movie_proto_rawDescGZIP(), []int{3}
}

func (x *AddMovieResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_proto_movie_proto protoreflect.FileDescriptor

var file_proto_movie_proto_rawDesc = string([]byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x22, 0x77, 0x0a, 0x0c, 0x4d, 0x6f,
	0x76, 0x69, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x6f,
	0x76, 0x69, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6d, 0x6f,
	0x76, 0x69, 0x65, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6f,
	0x76, 0x65, 0x72, 0x76, 0x69, 0x65, 0x77, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6f,
	0x76, 0x65, 0x72, 0x76, 0x69, 0x65, 0x77, 0x12, 0x1a, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x77, 0x6f,
	0x72, 0x64, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6b, 0x65, 0x79, 0x77, 0x6f,
	0x72, 0x64, 0x73, 0x22, 0x31, 0x0a, 0x11, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x6d, 0x62, 0x65,
	0x64, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x03, 0x28, 0x02, 0x52, 0x09, 0x65, 0x6d, 0x62,
	0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x22, 0x4a, 0x0a, 0x0f, 0x41, 0x64, 0x64, 0x4d, 0x6f, 0x76,
	0x69, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x6f, 0x76,
	0x69, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6d, 0x6f, 0x76,
	0x69, 0x65, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e,
	0x67, 0x18, 0x02, 0x20, 0x03, 0x28, 0x02, 0x52, 0x09, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69,
	0x6e, 0x67, 0x22, 0x2c, 0x0a, 0x10, 0x41, 0x64, 0x64, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x32, 0x9c, 0x01, 0x0a, 0x10, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x42, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x4d, 0x6f, 0x76, 0x69,
	0x65, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x13, 0x2e, 0x6d, 0x6f, 0x76,
	0x69, 0x65, 0x2e, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x18, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x2e, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e,
	0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a, 0x11, 0x41, 0x64, 0x64,
	0x4d, 0x6f, 0x76, 0x69, 0x65, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x16,
	0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x2e, 0x41, 0x64, 0x64, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x2e, 0x41,
	0x64, 0x64, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42,
	0x23, 0x5a, 0x21, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x2d, 0x72, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x2f, 0x67, 0x6f, 0x2d, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x70,
	0x62, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_proto_movie_proto_rawDescOnce sync.Once
	file_proto_movie_proto_rawDescData []byte
)

func file_proto_movie_proto_rawDescGZIP() []byte {
	file_proto_movie_proto_rawDescOnce.Do(func() {
		file_proto_movie_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_movie_proto_rawDesc), len(file_proto_movie_proto_rawDesc)))
	})
	return file_proto_movie_proto_rawDescData
}

var file_proto_movie_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_movie_proto_goTypes = []any{
	(*MovieRequest)(nil),      // 0: movie.MovieRequest
	(*EmbeddingResponse)(nil), // 1: movie.EmbeddingResponse
	(*AddMovieRequest)(nil),   // 2: movie.AddMovieRequest
	(*AddMovieResponse)(nil),  // 3: movie.AddMovieResponse
}
var file_proto_movie_proto_depIdxs = []int32{
	0, // 0: movie.EmbeddingService.GetMovieEmbedding:input_type -> movie.MovieRequest
	2, // 1: movie.EmbeddingService.AddMovieEmbedding:input_type -> movie.AddMovieRequest
	1, // 2: movie.EmbeddingService.GetMovieEmbedding:output_type -> movie.EmbeddingResponse
	3, // 3: movie.EmbeddingService.AddMovieEmbedding:output_type -> movie.AddMovieResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_movie_proto_init() }
func file_proto_movie_proto_init() {
	if File_proto_movie_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_movie_proto_rawDesc), len(file_proto_movie_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_movie_proto_goTypes,
		DependencyIndexes: file_proto_movie_proto_depIdxs,
		MessageInfos:      file_proto_movie_proto_msgTypes,
	}.Build()
	File_proto_movie_proto = out.File
	file_proto_movie_proto_goTypes = nil
	file_proto_movie_proto_depIdxs = nil
}
