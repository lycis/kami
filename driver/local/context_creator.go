package local

import (
	"github.com/lycis/kami/entity"
	"github.com/lycis/kami/script"
)

// implementation of the context creator interface

func (d LocalDriver) GetScriptPrivilegeLevel() script.PrivilegeLevel {
	return script.PrivilegeRoot
}

func (d LocalDriver) GetScriptReferenceEntity() *entity.Entity {
	return nil
}
