package regex

const (
	Import    = `[A-Za-z](\/?[A-Z-a-z0-9._-])*`
	GoToken   = `[A-Za-z][A-Za-z0-9_]*`
	YamlToken = `[A-Za-z](\.?[A-Za-z0-9_])*`

	MetaPkg           = GoToken
	MetaContainerType = GoToken
	MetaImport        = Import
	MetaImportAlias   = YamlToken
	MetaFn            = GoToken
	MetaGoFn          = `((?P<import>` + Import + `)\.)?(?P<fn>` + GoToken + `)`

	MetaParamName = YamlToken

	MetaServiceName   = YamlToken
	MetaServiceGetter = MetaPkg
	MetaServiceType   = `\*?` + `(` + Import + `\.)?` + GoToken
	MetaServiceValue  = `(` + Import + `\.)?(` + GoToken + `\{\}\.)?` + GoToken
)
