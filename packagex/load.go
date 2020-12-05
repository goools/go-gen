package packagex

import "golang.org/x/tools/go/packages"

func Load(pattern string) (*Package, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.LoadAllSyntax,
	}, pattern)

	if err != nil {
		return nil, err
	}

	return NewPackage(pkgs[0]), nil

}
