package main

type Option struct {
	Name           string
	Proto3Filename string
	Output         string
	GenClientCode  bool
	GenServerCode  bool
	Prefix         string
}
