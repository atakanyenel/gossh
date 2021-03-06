package file

import (
	"fmt"

	"github.com/krilor/gossh/machine"
	"github.com/pkg/errors"
)

// Exists is a struct that implements the rule to check if a file exists or not
type Exists string

// Check if file exists
func (e Exists) Check(m *machine.Machine, sudo bool) (bool, error) {
	cmd := fmt.Sprintf("stat %s", string(e))

	r, err := m.Run(cmd, sudo)

	if err != nil {
		return false, errors.Wrap(err, "stat errored")
	}

	if r.ExitStatus != 0 {
		return false, nil
	}

	return true, nil
}

// Ensure that file exists
func (e Exists) Ensure(m *machine.Machine, sudo bool) error {
	cmd := fmt.Sprintf("touch %s", string(e))
	r, err := m.Run(cmd, sudo)
	if err != nil {
		return errors.Wrap(err, "could not ensure file")
	}
	if r.ExitStatus != 0 {
		return errors.New("something went wrong with touch")
	}
	return nil
}
