package db

type PriviligeStatus int

const (
	PriviligeStatus_FreeUser PriviligeStatus = iota
	PriviligeStatus_Operator
	PriviligeStatus_NetAdmin
	PriviligeStatus_Developer
)

type AccessLevel int

const (
	AccessLevel_ReleaseGA AccessLevel = iota
	AccessLevel_ConsumerPreview
	AccessLevel_TechnologyPreview
	AccessLevel_DeveloperPreview
)

type PasswordType string

// Password Types are set by "[protocol]/[type]"
const (
// PasswordType_MySpaceIM         = "gsp-msim/rc4"
// PasswordType_MySpaceIM_Legacy  = "gsp-msim/legacy"
// PasswordType_OSCAR_Legacy      = "oscar/roasted"
// PasswordType_OSCAR_BUCP_Legacy = "oscar/bucp-1"
// PasswordType_OSCAR_BUCP        = "oscar/bucp-2"
// PasswordType_MSNP_Legacy       = "msnp/md5"
)
