package log

import (
	"time"

	"go.uber.org/zap/zapcore"
)

type sentryObjectEncoder struct {
	err     error
	strings map[string]string
	objects map[string]interface{}
}

func (s *sentryObjectEncoder) AddArray(key string, marshaler zapcore.ArrayMarshaler) error {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddObject(key string, marshaler zapcore.ObjectMarshaler) error {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddBinary(key string, value []byte) {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddByteString(key string, value []byte) {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddBool(key string, value bool) {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddComplex128(key string, value complex128) {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddComplex64(key string, value complex64) {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddDuration(key string, value time.Duration) {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddFloat64(key string, value float64) {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddFloat32(key string, value float32) {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddInt(key string, value int) {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddInt64(key string, value int64) {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddInt32(key string, value int32) {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddInt16(key string, value int16) {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddInt8(key string, value int8) {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddString(key, value string) {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddTime(key string, value time.Time) {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddUint(key string, value uint) {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddUint64(key string, value uint64) {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddUint32(key string, value uint32) {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddUint16(key string, value uint16) {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddUint8(key string, value uint8) {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddUintptr(key string, value uintptr) {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) AddReflected(key string, value interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (s *sentryObjectEncoder) OpenNamespace(key string) {
	//TODO implement me
	panic("implement me")
}
