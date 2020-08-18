package regex

const (
	BaseImport = `[A-Za-z](\/?[A-Z-a-z0-9._-])*`
	Import     = `(` + BaseImport + `)|("` + BaseImport + `")`
	GoToken    = `[A-Za-z][A-Za-z0-9_]*`
	YamlToken  = `[A-Za-z](\.?[A-Za-z0-9_])*`
	GoFunc     = `((?P<import>` + Import + `)\.)?(?P<fn>` + GoToken + `)`

	MetaPkg           = GoToken
	MetaContainerType = GoToken
	MetaImport        = Import
	MetaImportAlias   = YamlToken
	MetaFn            = GoToken
	MetaGoFn          = GoFunc

	MetaParamName = YamlToken

	MetaServiceName   = YamlToken
	MetaServiceGetter = GoToken
	MetaServiceType   = `\*?` + `(` + Import + `\.)?` + GoToken
	// MetaServiceValue supports
	// Value
	// my/import.Value
	// my/import.StructName{}.Method
	MetaServiceValue       = `((?P<import>` + Import + `)\.)?((?P<struct>` + GoToken + `)\{\}\.)?(?P<value>` + GoToken + `)`
	MetaServiceConstructor = GoFunc
	MetaServiceCallName    = GoToken
	MetaServiceFieldName   = GoToken
	MetaServiceTag         = YamlToken
)
