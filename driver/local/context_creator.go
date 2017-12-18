package local

import (
	"gitlab.com/lycis/kami/entity"
	"gitlab.com/lycis/kami/privilege"
)

// implementation of the context creator interface

func (d LocalDriver) GetScriptPrivilegeLevel() privilege.Level {
	return privilege.PrivilegeRoot
}

func (d LocalDriver) GetScriptReferenceEntity() *entity.Entity {
	return nil
}
