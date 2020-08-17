package regex

const (
	MetaPkg           = "[A-Za-z][A-Za-z0-9_]*"
	MetaContainerType = MetaPkg
	MetaImport        = "[a-zA-Z0-9_./]+"
	MetaImportAlias   = "[a-zA-Z0-9_]+"
	MetaFn            = `[A-Za-z][A-Za-z0-9]*`
	MetaGoFn          = `((?P<import>[A-Za-z][A-Z-a-z0-9._\/-]*)\.)?(?P<fn>[A-Za-z][A-Za-z0-9]*)` // todo split to smaller consts

	MetaParamName = "[A-Za-z]([_.]?[A-Za-z0-9])*"
)
