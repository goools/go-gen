package packagex

import (
	"go/types"
	"golang.org/x/tools/go/packages"
)

type PackageSet map[string]*packages.Package

func (s PackageSet) add(pkg *packages.Package) {
	s[pkg.ID] = pkg

	for k := range pkg.Imports {
		if _, ok := s[k]; !ok {
			s.add(pkg.Imports[k])
		}
	}
}

func (s PackageSet) allPackages() []*packages.Package {
	list := make([]*packages.Package, 0)
	for id := range s {
		list = append(list, s[id])
	}
	return list
}

type Package struct {
	*packages.Package
	AllPackages []*packages.Package
}

func NewPackage(pkg *packages.Package) *Package {
	p := &Package{
		Package: pkg,
	}

	s := PackageSet{}
	s.add(pkg)

	p.AllPackages = s.allPackages()

	return p
}

func (p *Package) Const(name string) *types.Const {
	for ident, def := range p.TypesInfo.Defs {
		if typeConst, ok := def.(*types.Const); ok {
			if ident.Name == name {
				return typeConst
			}
		}
	}
	return nil
}

func (p *Package) TypeName(name string) *types.TypeName {
	for ident, def := range p.TypesInfo.Defs {
		if typeName, ok := def.(*types.TypeName); ok {
			if ident.Name == name {
				return typeName
			}
		}
	}
	return nil
}

func (p *Package) Var(name string) *types.Var {
	for ident, def := range p.TypesInfo.Defs {
		if typeVar, ok := def.(*types.Var); ok {
			if ident.Name == name {
				return typeVar
			}
		}
	}
	return nil
}

func (p *Package) Func(name string) *types.Func {
	for ident, def := range p.TypesInfo.Defs {
		if typeFunc, ok := def.(*types.Func); ok {
			if ident.Name == name {
				return typeFunc
			}
		}
	}
	return nil
}

func (p *Package) Pkg(pkgPath string) *packages.Package {
	for _, pkg := range p.AllPackages {
		if pkg.PkgPath == pkgPath {
			return pkg
		}
	}
	return nil
}
