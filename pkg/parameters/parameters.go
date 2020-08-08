package parameters

type RawParameter interface{}

type RawParameters map[string]RawParameter

type ResolvedParam struct {
	Code string
	Raw  RawParameter
}

type ResolvedParams map[string]ResolvedParam
