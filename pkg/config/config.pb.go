// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: config.proto

package config

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import wrappers "github.com/golang/protobuf/ptypes/wrappers"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type Basic struct {
	Username             string                `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password             *wrappers.StringValue `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Basic) Reset()         { *m = Basic{} }
func (m *Basic) String() string { return proto.CompactTextString(m) }
func (*Basic) ProtoMessage()    {}
func (*Basic) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_75e1cfcba55f6f26, []int{0}
}
func (m *Basic) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Basic.Unmarshal(m, b)
}
func (m *Basic) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Basic.Marshal(b, m, deterministic)
}
func (dst *Basic) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Basic.Merge(dst, src)
}
func (m *Basic) XXX_Size() int {
	return xxx_messageInfo_Basic.Size(m)
}
func (m *Basic) XXX_DiscardUnknown() {
	xxx_messageInfo_Basic.DiscardUnknown(m)
}

var xxx_messageInfo_Basic proto.InternalMessageInfo

func (m *Basic) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *Basic) GetPassword() *wrappers.StringValue {
	if m != nil {
		return m.Password
	}
	return nil
}

type OAuthToken struct {
	Token                string                `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	ApplicationId        *wrappers.StringValue `protobuf:"bytes,2,opt,name=application_id,json=applicationId" json:"application_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *OAuthToken) Reset()         { *m = OAuthToken{} }
func (m *OAuthToken) String() string { return proto.CompactTextString(m) }
func (*OAuthToken) ProtoMessage()    {}
func (*OAuthToken) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_75e1cfcba55f6f26, []int{1}
}
func (m *OAuthToken) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OAuthToken.Unmarshal(m, b)
}
func (m *OAuthToken) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OAuthToken.Marshal(b, m, deterministic)
}
func (dst *OAuthToken) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OAuthToken.Merge(dst, src)
}
func (m *OAuthToken) XXX_Size() int {
	return xxx_messageInfo_OAuthToken.Size(m)
}
func (m *OAuthToken) XXX_DiscardUnknown() {
	xxx_messageInfo_OAuthToken.DiscardUnknown(m)
}

var xxx_messageInfo_OAuthToken proto.InternalMessageInfo

func (m *OAuthToken) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *OAuthToken) GetApplicationId() *wrappers.StringValue {
	if m != nil {
		return m.ApplicationId
	}
	return nil
}

type OAuth2Token struct {
	Token                string                `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	TokenType            *wrappers.StringValue `protobuf:"bytes,2,opt,name=token_type,json=tokenType" json:"token_type,omitempty"`
	RefreshToken         *wrappers.StringValue `protobuf:"bytes,3,opt,name=refresh_token,json=refreshToken" json:"refresh_token,omitempty"`
	Expiry               *wrappers.StringValue `protobuf:"bytes,4,opt,name=expiry" json:"expiry,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *OAuth2Token) Reset()         { *m = OAuth2Token{} }
func (m *OAuth2Token) String() string { return proto.CompactTextString(m) }
func (*OAuth2Token) ProtoMessage()    {}
func (*OAuth2Token) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_75e1cfcba55f6f26, []int{2}
}
func (m *OAuth2Token) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OAuth2Token.Unmarshal(m, b)
}
func (m *OAuth2Token) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OAuth2Token.Marshal(b, m, deterministic)
}
func (dst *OAuth2Token) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OAuth2Token.Merge(dst, src)
}
func (m *OAuth2Token) XXX_Size() int {
	return xxx_messageInfo_OAuth2Token.Size(m)
}
func (m *OAuth2Token) XXX_DiscardUnknown() {
	xxx_messageInfo_OAuth2Token.DiscardUnknown(m)
}

var xxx_messageInfo_OAuth2Token proto.InternalMessageInfo

func (m *OAuth2Token) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *OAuth2Token) GetTokenType() *wrappers.StringValue {
	if m != nil {
		return m.TokenType
	}
	return nil
}

func (m *OAuth2Token) GetRefreshToken() *wrappers.StringValue {
	if m != nil {
		return m.RefreshToken
	}
	return nil
}

func (m *OAuth2Token) GetExpiry() *wrappers.StringValue {
	if m != nil {
		return m.Expiry
	}
	return nil
}

type Github struct {
	BaseUrl              *wrappers.StringValue `protobuf:"bytes,1,opt,name=base_url,json=baseUrl" json:"base_url,omitempty"`
	UploadUrl            *wrappers.StringValue `protobuf:"bytes,2,opt,name=upload_url,json=uploadUrl" json:"upload_url,omitempty"`
	Oauth2               *OAuth2Token          `protobuf:"bytes,5,opt,name=oauth2" json:"oauth2,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Github) Reset()         { *m = Github{} }
func (m *Github) String() string { return proto.CompactTextString(m) }
func (*Github) ProtoMessage()    {}
func (*Github) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_75e1cfcba55f6f26, []int{3}
}
func (m *Github) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Github.Unmarshal(m, b)
}
func (m *Github) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Github.Marshal(b, m, deterministic)
}
func (dst *Github) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Github.Merge(dst, src)
}
func (m *Github) XXX_Size() int {
	return xxx_messageInfo_Github.Size(m)
}
func (m *Github) XXX_DiscardUnknown() {
	xxx_messageInfo_Github.DiscardUnknown(m)
}

var xxx_messageInfo_Github proto.InternalMessageInfo

func (m *Github) GetBaseUrl() *wrappers.StringValue {
	if m != nil {
		return m.BaseUrl
	}
	return nil
}

func (m *Github) GetUploadUrl() *wrappers.StringValue {
	if m != nil {
		return m.UploadUrl
	}
	return nil
}

func (m *Github) GetOauth2() *OAuth2Token {
	if m != nil {
		return m.Oauth2
	}
	return nil
}

type Gitlab struct {
	BaseUrl              *wrappers.StringValue `protobuf:"bytes,1,opt,name=base_url,json=baseUrl" json:"base_url,omitempty"`
	Private              *OAuthToken           `protobuf:"bytes,5,opt,name=private" json:"private,omitempty"`
	Oauth                *OAuthToken           `protobuf:"bytes,6,opt,name=oauth" json:"oauth,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Gitlab) Reset()         { *m = Gitlab{} }
func (m *Gitlab) String() string { return proto.CompactTextString(m) }
func (*Gitlab) ProtoMessage()    {}
func (*Gitlab) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_75e1cfcba55f6f26, []int{4}
}
func (m *Gitlab) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Gitlab.Unmarshal(m, b)
}
func (m *Gitlab) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Gitlab.Marshal(b, m, deterministic)
}
func (dst *Gitlab) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Gitlab.Merge(dst, src)
}
func (m *Gitlab) XXX_Size() int {
	return xxx_messageInfo_Gitlab.Size(m)
}
func (m *Gitlab) XXX_DiscardUnknown() {
	xxx_messageInfo_Gitlab.DiscardUnknown(m)
}

var xxx_messageInfo_Gitlab proto.InternalMessageInfo

func (m *Gitlab) GetBaseUrl() *wrappers.StringValue {
	if m != nil {
		return m.BaseUrl
	}
	return nil
}

func (m *Gitlab) GetPrivate() *OAuthToken {
	if m != nil {
		return m.Private
	}
	return nil
}

func (m *Gitlab) GetOauth() *OAuthToken {
	if m != nil {
		return m.Oauth
	}
	return nil
}

type Bitbucket struct {
	Basic                *Basic      `protobuf:"bytes,5,opt,name=basic" json:"basic,omitempty"`
	Oauth                *OAuthToken `protobuf:"bytes,6,opt,name=oauth" json:"oauth,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Bitbucket) Reset()         { *m = Bitbucket{} }
func (m *Bitbucket) String() string { return proto.CompactTextString(m) }
func (*Bitbucket) ProtoMessage()    {}
func (*Bitbucket) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_75e1cfcba55f6f26, []int{5}
}
func (m *Bitbucket) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Bitbucket.Unmarshal(m, b)
}
func (m *Bitbucket) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Bitbucket.Marshal(b, m, deterministic)
}
func (dst *Bitbucket) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Bitbucket.Merge(dst, src)
}
func (m *Bitbucket) XXX_Size() int {
	return xxx_messageInfo_Bitbucket.Size(m)
}
func (m *Bitbucket) XXX_DiscardUnknown() {
	xxx_messageInfo_Bitbucket.DiscardUnknown(m)
}

var xxx_messageInfo_Bitbucket proto.InternalMessageInfo

func (m *Bitbucket) GetBasic() *Basic {
	if m != nil {
		return m.Basic
	}
	return nil
}

func (m *Bitbucket) GetOauth() *OAuthToken {
	if m != nil {
		return m.Oauth
	}
	return nil
}

type Generic struct {
	BaseUrl              string   `protobuf:"bytes,1,opt,name=base_url,json=baseUrl,proto3" json:"base_url,omitempty"`
	Path                 string   `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	PerPageParameter     string   `protobuf:"bytes,3,opt,name=per_page_parameter,json=perPageParameter,proto3" json:"per_page_parameter,omitempty"`
	PageParameter        string   `protobuf:"bytes,4,opt,name=page_parameter,json=pageParameter,proto3" json:"page_parameter,omitempty"`
	PageSize             int32    `protobuf:"varint,5,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	Selector             string   `protobuf:"bytes,6,opt,name=selector,proto3" json:"selector,omitempty"`
	Basic                *Basic   `protobuf:"bytes,10,opt,name=basic" json:"basic,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Generic) Reset()         { *m = Generic{} }
func (m *Generic) String() string { return proto.CompactTextString(m) }
func (*Generic) ProtoMessage()    {}
func (*Generic) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_75e1cfcba55f6f26, []int{6}
}
func (m *Generic) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Generic.Unmarshal(m, b)
}
func (m *Generic) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Generic.Marshal(b, m, deterministic)
}
func (dst *Generic) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Generic.Merge(dst, src)
}
func (m *Generic) XXX_Size() int {
	return xxx_messageInfo_Generic.Size(m)
}
func (m *Generic) XXX_DiscardUnknown() {
	xxx_messageInfo_Generic.DiscardUnknown(m)
}

var xxx_messageInfo_Generic proto.InternalMessageInfo

func (m *Generic) GetBaseUrl() string {
	if m != nil {
		return m.BaseUrl
	}
	return ""
}

func (m *Generic) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *Generic) GetPerPageParameter() string {
	if m != nil {
		return m.PerPageParameter
	}
	return ""
}

func (m *Generic) GetPageParameter() string {
	if m != nil {
		return m.PageParameter
	}
	return ""
}

func (m *Generic) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *Generic) GetSelector() string {
	if m != nil {
		return m.Selector
	}
	return ""
}

func (m *Generic) GetBasic() *Basic {
	if m != nil {
		return m.Basic
	}
	return nil
}

type Account struct {
	Github               *Github    `protobuf:"bytes,1,opt,name=github" json:"github,omitempty"`
	Gitlab               *Gitlab    `protobuf:"bytes,2,opt,name=gitlab" json:"gitlab,omitempty"`
	Bitbucket            *Bitbucket `protobuf:"bytes,3,opt,name=bitbucket" json:"bitbucket,omitempty"`
	Generic              *Generic   `protobuf:"bytes,4,opt,name=generic" json:"generic,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Account) Reset()         { *m = Account{} }
func (m *Account) String() string { return proto.CompactTextString(m) }
func (*Account) ProtoMessage()    {}
func (*Account) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_75e1cfcba55f6f26, []int{7}
}
func (m *Account) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Account.Unmarshal(m, b)
}
func (m *Account) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Account.Marshal(b, m, deterministic)
}
func (dst *Account) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Account.Merge(dst, src)
}
func (m *Account) XXX_Size() int {
	return xxx_messageInfo_Account.Size(m)
}
func (m *Account) XXX_DiscardUnknown() {
	xxx_messageInfo_Account.DiscardUnknown(m)
}

var xxx_messageInfo_Account proto.InternalMessageInfo

func (m *Account) GetGithub() *Github {
	if m != nil {
		return m.Github
	}
	return nil
}

func (m *Account) GetGitlab() *Gitlab {
	if m != nil {
		return m.Gitlab
	}
	return nil
}

func (m *Account) GetBitbucket() *Bitbucket {
	if m != nil {
		return m.Bitbucket
	}
	return nil
}

func (m *Account) GetGeneric() *Generic {
	if m != nil {
		return m.Generic
	}
	return nil
}

type Configuration struct {
	Accounts             []*Account `protobuf:"bytes,1,rep,name=accounts" json:"accounts,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Configuration) Reset()         { *m = Configuration{} }
func (m *Configuration) String() string { return proto.CompactTextString(m) }
func (*Configuration) ProtoMessage()    {}
func (*Configuration) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_75e1cfcba55f6f26, []int{8}
}
func (m *Configuration) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Configuration.Unmarshal(m, b)
}
func (m *Configuration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Configuration.Marshal(b, m, deterministic)
}
func (dst *Configuration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Configuration.Merge(dst, src)
}
func (m *Configuration) XXX_Size() int {
	return xxx_messageInfo_Configuration.Size(m)
}
func (m *Configuration) XXX_DiscardUnknown() {
	xxx_messageInfo_Configuration.DiscardUnknown(m)
}

var xxx_messageInfo_Configuration proto.InternalMessageInfo

func (m *Configuration) GetAccounts() []*Account {
	if m != nil {
		return m.Accounts
	}
	return nil
}

func init() {
	proto.RegisterType((*Basic)(nil), "mjpitz.gitfs.Basic")
	proto.RegisterType((*OAuthToken)(nil), "mjpitz.gitfs.OAuthToken")
	proto.RegisterType((*OAuth2Token)(nil), "mjpitz.gitfs.OAuth2Token")
	proto.RegisterType((*Github)(nil), "mjpitz.gitfs.Github")
	proto.RegisterType((*Gitlab)(nil), "mjpitz.gitfs.Gitlab")
	proto.RegisterType((*Bitbucket)(nil), "mjpitz.gitfs.Bitbucket")
	proto.RegisterType((*Generic)(nil), "mjpitz.gitfs.Generic")
	proto.RegisterType((*Account)(nil), "mjpitz.gitfs.Account")
	proto.RegisterType((*Configuration)(nil), "mjpitz.gitfs.Configuration")
}

func init() { proto.RegisterFile("config.proto", fileDescriptor_config_75e1cfcba55f6f26) }

var fileDescriptor_config_75e1cfcba55f6f26 = []byte{
	// 619 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x94, 0xd1, 0x6e, 0xd3, 0x3c,
	0x14, 0xc7, 0x95, 0x6f, 0x6b, 0xda, 0x9c, 0xad, 0xd3, 0x27, 0x33, 0x44, 0x36, 0xb8, 0x98, 0x22,
	0x21, 0x0d, 0x69, 0x4a, 0xb5, 0x01, 0x02, 0xc4, 0xd5, 0xba, 0x8b, 0x89, 0x2b, 0xa6, 0x6c, 0x70,
	0x81, 0x84, 0x22, 0x27, 0x3b, 0x75, 0xcd, 0xd2, 0xd8, 0x72, 0x1c, 0xc6, 0xf6, 0x3e, 0x3c, 0x00,
	0x0f, 0xc2, 0x03, 0xf0, 0x18, 0xbc, 0x01, 0x8a, 0xed, 0xb4, 0x5b, 0x87, 0x68, 0x11, 0x77, 0x76,
	0xf2, 0xfb, 0xfb, 0xef, 0x73, 0x7c, 0xce, 0x81, 0xf5, 0x5c, 0x94, 0x23, 0xce, 0x62, 0xa9, 0x84,
	0x16, 0x64, 0x7d, 0xf2, 0x49, 0x72, 0x7d, 0x1d, 0x33, 0xae, 0x47, 0xd5, 0xf6, 0x2b, 0xc6, 0xf5,
	0xb8, 0xce, 0xe2, 0x5c, 0x4c, 0x06, 0x4c, 0x14, 0xb4, 0x64, 0x03, 0x83, 0x65, 0xf5, 0x68, 0x20,
	0xf5, 0x95, 0xc4, 0x6a, 0x70, 0xa9, 0xa8, 0x94, 0xa8, 0x66, 0x0b, 0x7b, 0x50, 0xf4, 0x11, 0x3a,
	0x43, 0x5a, 0xf1, 0x9c, 0x6c, 0x43, 0xaf, 0xae, 0x50, 0x95, 0x74, 0x82, 0xa1, 0xb7, 0xe3, 0xed,
	0x06, 0xc9, 0x74, 0x4f, 0x5e, 0x42, 0x4f, 0xd2, 0xaa, 0xba, 0x14, 0xea, 0x3c, 0xfc, 0x6f, 0xc7,
	0xdb, 0x5d, 0x3b, 0x78, 0x14, 0x33, 0x21, 0x58, 0x81, 0x71, 0xeb, 0x13, 0x9f, 0x6a, 0xc5, 0x4b,
	0xf6, 0x9e, 0x16, 0x35, 0x26, 0x53, 0x3a, 0x62, 0x00, 0x6f, 0x0f, 0x6b, 0x3d, 0x3e, 0x13, 0x17,
	0x58, 0x92, 0x4d, 0xe8, 0xe8, 0x66, 0xe1, 0x0c, 0xec, 0x86, 0x1c, 0xc1, 0x06, 0x95, 0xb2, 0xe0,
	0x39, 0xd5, 0x5c, 0x94, 0x29, 0x5f, 0xce, 0xa3, 0x7f, 0x43, 0xf3, 0xe6, 0x3c, 0xfa, 0xe1, 0xc1,
	0x9a, 0x71, 0x3a, 0xf8, 0x93, 0xd5, 0x6b, 0x00, 0xb3, 0x48, 0x9b, 0xa4, 0x2c, 0x65, 0x13, 0x18,
	0xfe, 0xec, 0x4a, 0x22, 0x39, 0x84, 0xbe, 0xc2, 0x91, 0xc2, 0x6a, 0x9c, 0xda, 0xa3, 0x57, 0x96,
	0xd0, 0xaf, 0x3b, 0x89, 0xbd, 0xd5, 0x33, 0xf0, 0xf1, 0x8b, 0xe4, 0xea, 0x2a, 0x5c, 0x5d, 0x42,
	0xeb, 0xd8, 0xe8, 0x9b, 0x07, 0xfe, 0xb1, 0x79, 0x61, 0xf2, 0x02, 0x7a, 0x19, 0xad, 0x30, 0xad,
	0x55, 0x61, 0x22, 0x5b, 0x74, 0x44, 0xb7, 0xa1, 0xdf, 0xa9, 0xa2, 0x89, 0xbc, 0x96, 0x85, 0xa0,
	0xe7, 0x46, 0xba, 0x54, 0xe4, 0x96, 0x6f, 0xc4, 0xfb, 0xe0, 0x0b, 0xda, 0xe4, 0x36, 0xec, 0x18,
	0xe1, 0x56, 0x7c, 0xb3, 0xfc, 0xe2, 0x1b, 0x79, 0x4f, 0x1c, 0x18, 0x7d, 0xb5, 0x77, 0x2e, 0xe8,
	0x3f, 0xdc, 0xf9, 0x00, 0xba, 0x52, 0xf1, 0xcf, 0x54, 0xa3, 0xf3, 0x0d, 0x7f, 0xe3, 0x6b, 0x6d,
	0x5b, 0x90, 0xc4, 0xd0, 0x31, 0x37, 0x08, 0xfd, 0x05, 0x0a, 0x8b, 0x45, 0x23, 0x08, 0x86, 0x5c,
	0x67, 0x75, 0x7e, 0x81, 0x9a, 0x3c, 0x81, 0x4e, 0xd6, 0x34, 0x83, 0xb3, 0xbb, 0x77, 0x5b, 0x6c,
	0xfa, 0x24, 0xb1, 0xc4, 0x5f, 0xfb, 0xfc, 0xf4, 0xa0, 0x7b, 0x8c, 0x25, 0x2a, 0x9e, 0x93, 0xad,
	0xb9, 0x84, 0x04, 0xb3, 0x90, 0x09, 0xac, 0x4a, 0xaa, 0xc7, 0xe6, 0x81, 0x82, 0xc4, 0xac, 0xc9,
	0x1e, 0x10, 0x89, 0x2a, 0x95, 0x94, 0x61, 0x2a, 0xa9, 0xa2, 0x13, 0xd4, 0xa8, 0x4c, 0xf1, 0x05,
	0xc9, 0xff, 0x12, 0xd5, 0x09, 0x65, 0x78, 0xd2, 0x7e, 0x27, 0x8f, 0x61, 0x63, 0x8e, 0x5c, 0x35,
	0x64, 0x5f, 0xde, 0xc2, 0x1e, 0x42, 0x60, 0xb0, 0x8a, 0x5f, 0xdb, 0xec, 0x76, 0x9a, 0xae, 0x65,
	0x78, 0xca, 0xaf, 0xb1, 0x99, 0x05, 0x15, 0x16, 0x98, 0x6b, 0xa1, 0x4c, 0x7c, 0x41, 0x32, 0xdd,
	0xcf, 0x72, 0x04, 0x8b, 0x72, 0x14, 0x7d, 0xf7, 0xa0, 0x7b, 0x98, 0xe7, 0xa2, 0x2e, 0x35, 0xd9,
	0x03, 0xdf, 0x0e, 0x29, 0x57, 0x02, 0x9b, 0xb7, 0x75, 0xb6, 0xbc, 0x13, 0xc7, 0x38, 0xba, 0xa0,
	0x99, 0xab, 0xd4, 0xbb, 0x74, 0x41, 0x2d, 0xdd, 0x14, 0xd8, 0x73, 0x08, 0xb2, 0xf6, 0x0d, 0x5d,
	0x53, 0x3e, 0x98, 0xbb, 0x56, 0xfb, 0x3b, 0x99, 0x91, 0x64, 0x00, 0x5d, 0x66, 0x5f, 0xc4, 0x75,
	0xe3, 0xfd, 0x39, 0x17, 0xfb, 0x33, 0x69, 0xa9, 0x68, 0x08, 0xfd, 0x23, 0x33, 0x84, 0x6b, 0x65,
	0xc6, 0x0e, 0xd9, 0x87, 0x1e, 0xb5, 0xf1, 0x55, 0xa1, 0xb7, 0xb3, 0x72, 0xf7, 0x08, 0x17, 0x7d,
	0x32, 0xc5, 0x86, 0xbd, 0x0f, 0xbe, 0x1d, 0xe4, 0x99, 0x6f, 0x8a, 0xff, 0xe9, 0xaf, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xcd, 0x6f, 0xb6, 0x25, 0xd9, 0x05, 0x00, 0x00,
}
