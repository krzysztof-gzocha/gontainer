package regex

const (
	Import    = `[A-Za-z](\/?[A-Z-a-z0-9._-])*`
	GoToken   = `[A-Za-z][A-Za-z0-9_]*`
	YamlToken = `[A-Za-z](\.?[A-Za-z0-9_])*`
	GoFunc    = `((?P<import>` + Import + `)\.)?(?P<fn>` + GoToken + `)`

	MetaPkg           = GoToken
	MetaContainerType = GoToken
	MetaImport        = Import
	MetaImportAlias   = YamlToken
	MetaFn            = GoToken
	MetaGoFn          = GoFunc

	MetaParamName = YamlToken

	MetaServiceName        = YamlToken
	MetaServiceGetter      = MetaPkg
	MetaServiceType        = `\*?` + `(` + Import + `\.)?` + GoToken
	MetaServiceValue       = `(` + Import + `\.)?(` + GoToken + `\{\}\.)?` + GoToken
	MetaServiceConstructor = GoFunc
)
