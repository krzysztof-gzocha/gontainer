package compiled

// getAllServiceArgs returns arguments passed to constructor, calls and fields to fetch information about all dependencies.
func getAllServiceArgs(s Service) []Arg {
	var res []Arg
	res = append(res, s.Args...)
	for _, c := range s.Calls {
		res = append(res, c.Args...)
	}
	for _, f := range s.Fields {
		res = append(res, f.Value)
	}
	return res
}
