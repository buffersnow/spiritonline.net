package db

type AccountFlags int

const (
	StatusFlag_Banned     = AccountFlags(1 << 0)
	StatusFlag_Unverified = AccountFlags(1 << 1)
	StatusFlag_FreeUser   = AccountFlags(1 << 2)
	StatusFlag_Operator   = AccountFlags(1 << 3)
	StatusFlag_NetAdmin   = AccountFlags(1 << 4)
	StatusFlag_Developer  = AccountFlags(1 << 5)
)

const (
	AccessFlag_ReleasePreview    = AccountFlags(1 << 6)
	AccessFlag_ConsumerPreview   = AccountFlags(1 << 7)
	AccessFlag_TechnologyPreview = AccountFlags(1 << 8)
	AccessFlag_DeveloperPreview  = AccountFlags(1 << 9)
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
