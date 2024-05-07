// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package common

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson2220f231DecodeGithubComAlekseyt9YpmetricsInternalCommon(in *jlexer.Lexer, out *MetricsBatch) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "counters":
			if in.IsNull() {
				in.Skip()
				out.Counters = nil
			} else {
				in.Delim('[')
				if out.Counters == nil {
					if !in.IsDelim(']') {
						out.Counters = make([]CounterItem, 0, 2)
					} else {
						out.Counters = []CounterItem{}
					}
				} else {
					out.Counters = (out.Counters)[:0]
				}
				for !in.IsDelim(']') {
					var v1 CounterItem
					(v1).UnmarshalEasyJSON(in)
					out.Counters = append(out.Counters, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "gauges":
			if in.IsNull() {
				in.Skip()
				out.Gauges = nil
			} else {
				in.Delim('[')
				if out.Gauges == nil {
					if !in.IsDelim(']') {
						out.Gauges = make([]GaugeItem, 0, 2)
					} else {
						out.Gauges = []GaugeItem{}
					}
				} else {
					out.Gauges = (out.Gauges)[:0]
				}
				for !in.IsDelim(']') {
					var v2 GaugeItem
					(v2).UnmarshalEasyJSON(in)
					out.Gauges = append(out.Gauges, v2)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2220f231EncodeGithubComAlekseyt9YpmetricsInternalCommon(out *jwriter.Writer, in MetricsBatch) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"counters\":"
		out.RawString(prefix[1:])
		if in.Counters == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v3, v4 := range in.Counters {
				if v3 > 0 {
					out.RawByte(',')
				}
				(v4).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"gauges\":"
		out.RawString(prefix)
		if in.Gauges == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.Gauges {
				if v5 > 0 {
					out.RawByte(',')
				}
				(v6).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MetricsBatch) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2220f231EncodeGithubComAlekseyt9YpmetricsInternalCommon(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MetricsBatch) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2220f231EncodeGithubComAlekseyt9YpmetricsInternalCommon(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MetricsBatch) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2220f231DecodeGithubComAlekseyt9YpmetricsInternalCommon(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MetricsBatch) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2220f231DecodeGithubComAlekseyt9YpmetricsInternalCommon(l, v)
}
func easyjson2220f231DecodeGithubComAlekseyt9YpmetricsInternalCommon1(in *jlexer.Lexer, out *Metrics) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = string(in.String())
		case "type":
			out.MType = string(in.String())
		case "delta":
			if in.IsNull() {
				in.Skip()
				out.Delta = nil
			} else {
				if out.Delta == nil {
					out.Delta = new(int64)
				}
				*out.Delta = int64(in.Int64())
			}
		case "value":
			if in.IsNull() {
				in.Skip()
				out.Value = nil
			} else {
				if out.Value == nil {
					out.Value = new(float64)
				}
				*out.Value = float64(in.Float64())
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2220f231EncodeGithubComAlekseyt9YpmetricsInternalCommon1(out *jwriter.Writer, in Metrics) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.MType))
	}
	if in.Delta != nil {
		const prefix string = ",\"delta\":"
		out.RawString(prefix)
		out.Int64(int64(*in.Delta))
	}
	if in.Value != nil {
		const prefix string = ",\"value\":"
		out.RawString(prefix)
		out.Float64(float64(*in.Value))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Metrics) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2220f231EncodeGithubComAlekseyt9YpmetricsInternalCommon1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Metrics) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2220f231EncodeGithubComAlekseyt9YpmetricsInternalCommon1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Metrics) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2220f231DecodeGithubComAlekseyt9YpmetricsInternalCommon1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Metrics) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2220f231DecodeGithubComAlekseyt9YpmetricsInternalCommon1(l, v)
}
func easyjson2220f231DecodeGithubComAlekseyt9YpmetricsInternalCommon2(in *jlexer.Lexer, out *GaugeItem) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			out.Name = string(in.String())
		case "value":
			out.Value = float64(in.Float64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2220f231EncodeGithubComAlekseyt9YpmetricsInternalCommon2(out *jwriter.Writer, in GaugeItem) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"value\":"
		out.RawString(prefix)
		out.Float64(float64(in.Value))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GaugeItem) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2220f231EncodeGithubComAlekseyt9YpmetricsInternalCommon2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GaugeItem) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2220f231EncodeGithubComAlekseyt9YpmetricsInternalCommon2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GaugeItem) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2220f231DecodeGithubComAlekseyt9YpmetricsInternalCommon2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GaugeItem) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2220f231DecodeGithubComAlekseyt9YpmetricsInternalCommon2(l, v)
}
func easyjson2220f231DecodeGithubComAlekseyt9YpmetricsInternalCommon3(in *jlexer.Lexer, out *CounterItem) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			out.Name = string(in.String())
		case "value":
			out.Value = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2220f231EncodeGithubComAlekseyt9YpmetricsInternalCommon3(out *jwriter.Writer, in CounterItem) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"value\":"
		out.RawString(prefix)
		out.Int64(int64(in.Value))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CounterItem) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2220f231EncodeGithubComAlekseyt9YpmetricsInternalCommon3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CounterItem) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2220f231EncodeGithubComAlekseyt9YpmetricsInternalCommon3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CounterItem) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2220f231DecodeGithubComAlekseyt9YpmetricsInternalCommon3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CounterItem) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2220f231DecodeGithubComAlekseyt9YpmetricsInternalCommon3(l, v)
}
