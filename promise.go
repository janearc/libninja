package libninja

import (
	"go/types"
)

// gatherPromises finds every method that keeps a promise. Take each type
// the module declares and each interface in sight that the type or its
// pointer satisfies. The method that answers each interface method keeps
// a promise wherever it is declared.
//
// So a method promoted through an embedded field is found where it
// lives. Such a method stays exported though no one outside names it.
func (read *reading) gatherPromises() {
	faces := read.interfacesInSight()
	read.promises = map[string]bool{}
	for _, pkg := range read.plain() {
		for _, typeName := range typeNames(pkg.Types) {
			// typeNames leaves aliases out, so each is a named
			// type. Interface types are asked too. They keep
			// promises only for their own methods, which no
			// scroll counts.
			declared := typeName.Type().(*types.Named)
			for _, face := range faces {
				read.keep(declared, face)
			}
		}
	}
}

// keep records the methods by which a type satisfies an interface, if
// it does. Generic types cannot be asked whether they implement
// anything, so their own methods are matched by name alone.
func (read *reading) keep(declared *types.Named, face *types.Interface) {
	if positive(declared.TypeParams().Len()) {
		for method := range declared.Methods() {
			if hasMethod(face, method.Name()) {
				read.promises[read.key(method)] = true
			}
		}
		return
	}
	var held types.Type = declared
	if !types.Implements(held, face) {
		held = types.NewPointer(declared)
		if !types.Implements(held, face) {
			return
		}
	}
	for method := range face.Methods() {
		found, _, _ := types.LookupFieldOrMethod(held, true,
			method.Pkg(), method.Name())
		if answer, ok := found.(*types.Func); ok {
			read.promises[read.key(answer)] = true
		}
	}
}

// interfacesInSight is every interface with methods that the module
// can see. That is its own, and those of every package it imports
// however deep. The error interface is in sight too. Generic interfaces
// and interfaces written inline are not in sight.
func (read reading) interfacesInSight() []*types.Interface {
	sight := inSight{visited: map[*types.Package]bool{}}
	for _, pkg := range read.pkgs {
		sight.visit(pkg.Types)
	}
	failure := types.Universe.Lookup("error").Type().Underlying()
	return append(sight.found, failure.(*types.Interface))
}

// inSight is a walk through packages and everything they import: the
// packages seen, and the interfaces found in them.
type inSight struct {
	visited map[*types.Package]bool
	found   []*types.Interface
}

// visit takes in a package's interfaces, then those of everything it
// imports. It visits each package once. The guard holds the whole walk,
// so a walk without it finds nothing.
func (sight *inSight) visit(pkg *types.Package) {
	if !sight.visited[pkg] {
		sight.visited[pkg] = true
		sight.found = append(sight.found, interfacesOf(pkg)...)
		for _, imported := range pkg.Imports() {
			sight.visit(imported)
		}
	}
}

// interfacesOf is every named, non-generic interface with methods that
// a package declares.
func interfacesOf(pkg *types.Package) []*types.Interface {
	var found []*types.Interface
	for _, typeName := range typeNames(pkg) {
		declared, ok := typeName.Type().(*types.Named)
		if !ok || positive(declared.TypeParams().Len()) {
			continue
		}
		face, ok := declared.Underlying().(*types.Interface)
		if ok && positive(face.NumMethods()) && face.IsMethodSet() {
			found = append(found, face)
		}
	}
	return found
}

// hasMethod is whether an interface asks for a method by this name.
func hasMethod(face *types.Interface, name string) bool {
	for method := range face.Methods() {
		if method.Name() == name {
			return true
		}
	}
	return false
}
