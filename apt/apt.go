package apt

import (
	"fmt"
	"strings"

	"github.com/krilor/gossh/machine"
	"github.com/pkg/errors"
)

//go:generate stringer -type=PackageStatus

// PackageStatus is the state of a package
type PackageStatus byte

// Each state represents package status
// This enum matches the statuses for http://man7.org/linux/man-pages/man1/dpkg-query.1.html
const (
	StatusInstalled    PackageStatus = 'i'
	StatusNotInstalled PackageStatus = 'n'
)

// Package is a apt/dpkg package
type Package struct {
	Name   string
	Status PackageStatus
}

// Check checks if package is in the desired state
func (p Package) Check(m *machine.Machine, sudo bool) (bool, error) {
	cmd := fmt.Sprintf("dpkg-query -f '${Package}\t${db:Status-Abbrev}\t${Version}\t${Name}' -W %s", p.Name)

	r, err := m.Run(cmd, false)

	if err != nil {
		return false, errors.Wrapf(err, "could not check package status for %s", p.Name)
	}

	if r.ExitStatus != 0 && p.Status == StatusInstalled {
		return false, nil
	}

	// at this point, the package info has been returned, so we need to do some string-fiddling to get the status byte
	status := strings.Split(r.Stdout.String(), "\t")[1][1]

	if status != byte(p.Status) {
		return false, nil
	}

	return true, nil
}

// Ensure ensures that the package is in the desiStatusInstalledred state
func (p Package) Ensure(m *machine.Machine, sudo bool) error {

	actions := map[PackageStatus]string{
		StatusInstalled:    "install",
		StatusNotInstalled: "remove",
	}

	cmd := fmt.Sprintf("apt %s -y %s", actions[p.Status], p.Name)

	r, err := m.Run(cmd, true)

	if err != nil || r.ExitStatus != 0 {
		return errors.Wrapf(err, "could not %s package %s", actions[p.Status], p.Name)
	}

	return nil
}
