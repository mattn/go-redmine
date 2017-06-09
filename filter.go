package redmine

import "strings"

type Filter struct {
	filters map[string]string
}

func NewFilter(args ...string) *Filter {
	f := &Filter{}
	if len(args)%2 == 0 {
		for i := 0; i < len(args); i += 2 {
			f.AddPair(args[i], args[i+1])
		}
	}
	return f
}

func (f *Filter) AddPair(key, value string) {
	if f.filters == nil {
		f.filters = make(map[string]string)
	}
	f.filters[key] = encode4Redmine(value)
}

func (f *Filter) ToURLParams() string {
	params := ""
	for k, v := range f.filters {
		params += "&" + k + "=" + v
	}
	return params
}

func encode4Redmine(s string) string {
	a := strings.Replace(s, ">", "%3E", -1)
	a = strings.Replace(a, "<", "%3C", -1)
	a = strings.Replace(a, "=", "%3D", -1)
	return a
}
