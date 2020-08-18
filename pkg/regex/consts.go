package regex

const (
	Import = `[A-Za-z][A-Z-a-z0-9._\/-]*`
	Token  = `[A-Za-z][A-Za-z0-9]*`

	MetaPkg           = "[A-Za-z][A-Za-z0-9_]*"
	MetaContainerType = MetaPkg
	MetaImport        = Import
	MetaImportAlias   = "[a-zA-Z0-9_]+"
	MetaFn            = Token
	MetaGoFn          = `((?P<import>` + Import + `)\.)?(?P<fn>` + Token + `)`

	MetaParamName = "[A-Za-z]([_.]?[A-Za-z0-9])*"

	MetaServiceName   = "(?P<service>[A-Za-z]([._]?[A-Za-z0-9])*)"
	MetaServiceGetter = MetaPkg
	MetaServiceType   = `\*?` + `(` + Import + `\.)?` + Token + ``
)
