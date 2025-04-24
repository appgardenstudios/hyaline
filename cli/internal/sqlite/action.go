package sqlite

import (
	"github.com/go-git/go-git/v5/utils/merkletrie"
)

type Action string

const (
	ActionNone   Action = ""
	ActionInsert Action = "Insert"
	ActionModify Action = "Modify"
	ActionRename Action = "Rename"
	ActionDelete Action = "Delete"
)

func MapAction(changeAction merkletrie.Action, fromName string, toName string) Action {
	switch changeAction {
	case merkletrie.Insert:
		return ActionInsert
	case merkletrie.Modify:
		if fromName != toName {
			return ActionRename
		} else {
			return ActionModify
		}
	case merkletrie.Delete:
		return ActionDelete
	default:
		return ActionNone
	}
}
