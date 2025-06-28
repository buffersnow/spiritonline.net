package version

import (
	"fmt"
)

func (bi BuildTag) GetFullTag() string {
	return fmt.Sprintf("%s.%s.%s.%s",
		bi.GetVersion(),
		bi.GetConfig(),
		bi.GetLab(),
		bi.GetIdentifier(),
	)
}

func (bi BuildTag) GetPartialTag() string {
	return fmt.Sprintf("%s.%s.%s",
		bi.GetConfig(),
		bi.GetLab(),
		bi.GetIdentifier(),
	)
}

func (bi BuildTag) GetService() string {
	return bi.service
}

func (bi BuildTag) GetVersion() string {
	return fmt.Sprintf("%s.%s.%s.%s",
		bi.version.major,
		bi.version.minor,
		bi.version.hotfix,
		bi.version.commit,
	)
}

func (bi BuildTag) GetConfig() string {
	return fmt.Sprintf("%s%s.%s",
		bi.cfg.os_arch,
		bi.cfg.build_mode,
		bi.cfg.os_name,
	)
}

func (bi BuildTag) GetLab() string {
	return bi.lab.name
}

func (bi BuildTag) GetCISlave() string {
	return bi.lab.host + "\\" + bi.lab.username
}

func (bi BuildTag) GetIdentifier() string {
	return fmt.Sprintf("%s.%s",
		bi.iden.timestamp,
		bi.iden.dev_fingerprint,
	)
}
