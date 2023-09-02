package vix

import "strconv"

type StringValue struct {
	Value string
	Err   error
}

func (s StringValue) AsInt() (int, error) {
	if s.Err != nil {
		return 0, s.Err
	}
	return strconv.Atoi(s.Value)
}

func (s StringValue) AsInt32() (int32, error) {
	if s.Err != nil {
		return int32(0), s.Err
	}
	num, err := strconv.ParseInt(s.Value, 10, 32)
	return int32(num), err
}

func (s StringValue) AsInt64() (int64, error) {
	if s.Err != nil {
		return int64(0), s.Err
	}
	return strconv.ParseInt(s.Value, 10, 64)
}

func (s StringValue) AsFloat32() (float32, error) {
	if s.Err != nil {
		return float32(0), s.Err
	}
	num, err := strconv.ParseFloat(s.Value, 32)
	return float32(num), err
}

func (s StringValue) AsFloat64() (float64, error) {
	if s.Err != nil {
		return float64(0), s.Err
	}
	return strconv.ParseFloat(s.Value, 64)
}

func (s StringValue) AsByte() ([]byte, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	return []byte(s.Value), nil
}

func (s StringValue) IsError() bool {
	return s.Err != nil
}
