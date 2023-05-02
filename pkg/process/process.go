//////////////////////////////////////////////////////////////////////////
//
//  (c) TMList 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)
//  Copyright under MIT Lincese. Please see LICENSE file for details
//
//  process.go - include C1WS exclusion lists one into another
//
//////////////////////////////////////////////////////////////////////////

package process

import (
	"errors"
	"fmt"

	"github.com/mpkondrashin/tmlist/pkg/c1ews"
	"golang.org/x/exp/maps"
)

const (
	DependencePrefix = "Do not delete this list! It is used to populate the following lists:"
)

var (
	ErrCycleDependence = errors.New("cycle dependence")
	ErrListNotFound    = errors.New("not found")
)

type Process struct {
	in  []c1ews.ListResponse
	out []c1ews.ListResponse
}

func NewProcess(in []c1ews.ListResponse) *Process {
	p := &Process{
		in: in,
	}
	p.populateOut()
	return p
}

func (p *Process) populateOut() {
	p.out = make([]c1ews.ListResponse, len(p.in))
	copy(p.out, p.in)
	for i := range p.in {
		Cleanup(&p.out[i])
	}
}

func (p *Process) Process() error {
	for n := range p.out {
		if err := p.GetAllItems(&p.out[n]); err != nil {
			return err
		}
	}
	return nil
}

func (p *Process) GetAllItemsWithMap(l *c1ews.ListResponse, seen map[string]struct{}) error {
	AddDependences(l, maps.Keys(seen)...)
	/*	if len(l.Items) > 0 {
		return nil
	}*/
	includes := Includes(l)
	if len(includes) == 0 {
		return nil
	}
	seen[l.Name] = struct{}{}
	for _, name := range includes {
		list := p.FindList(name)
		if list == nil {
			return fmt.Errorf("%w: \"%s\" included in \"%s\" list", ErrListNotFound, name, l.Name)
		}
		_, found := seen[list.Name]
		if found {
			return fmt.Errorf("list %s refers to %s: %w", l.Name, list.Name, ErrCycleDependence)
		}
		if err := p.GetAllItemsWithMap(list, seen); err != nil {
			return err
		}
		AddToTheList(l, list.Items)
	}
	delete(seen, l.Name)
	return nil
}

func (p *Process) GetAllItems(l *c1ews.ListResponse) error {
	return p.GetAllItemsWithMap(l, make(map[string]struct{}))
}

func (p *Process) FindList(name string) *c1ews.ListResponse {
	for i := range p.out {
		if p.out[i].Name == name {
			return &p.out[i]
		}
	}
	return nil
}

func (p *Process) IterateChanged(callback func(*c1ews.ListResponse) error) error {
	for i := range p.in {
		if !Equal(&p.in[i], &p.out[i]) {
			if err := callback(&p.out[i]); err != nil {
				return err
			}
		}
	}
	return nil
}
