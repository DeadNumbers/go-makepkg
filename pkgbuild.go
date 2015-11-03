package main

import "text/template"

var pkgbuildTemplate = template.Must(
	template.New("pkgbuild").Parse(
		`{{if ne .Maintainer ""}}# Maintainer: {{.Maintainer}}
{{end}}pkgname={{.PkgName}}
pkgver="autogenerated"
pkgrel={{.PkgRel}}
pkgdesc="{{.PkgDesc}}"
arch=('i686' 'x86_64')
license=('{{.License}}')
makedepends=('go' 'git')

source=(
	"{{.PkgName}}::{{.RepoURL}}"{{range .Files}}
	"{{.Name}}"{{end}}
)

md5sums=(
	'SKIP'{{range .Files}}
	'{{.Hash}}'{{end}}
)

backup=({{range .Backup}}
	"{{.}}"{{end}}
)

pkgver() {
	cd "$srcdir/$pkgname"
	local date=$(git log -1 --format="%cd" --date=short | sed s/-//g)
	local count=$(git rev-list --count HEAD)
	local commit=$(git rev-parse --short HEAD)
	echo "$date.${count}_$commit"
}

build() {
	cd "$srcdir/$pkgname"

	if [ -L "$srcdir/$pkgname" ]; then
		rm "$srcdir/$pkgname" -rf
		mv "$srcdir/.go/src/$pkgname/" "$srcdir/$pkgname"
	fi

	rm -rf "$srcdir/.go/src"

	mkdir -p "$srcdir/.go/src"

	export GOPATH="$srcdir/.go"

	mv "$srcdir/$pkgname" "$srcdir/.go/src/"

	cd "$srcdir/.go/src/$pkgname/"
	ln -sf "$srcdir/.go/src/$pkgname/" "$srcdir/$pkgname"

	git submodule init
	git submodule update

	echo "Running 'go get'..."
	GO15VENDOREXPERIMENT=1 go get{{if .IsWildcardBuild}} ./...{{end}}
}

package() {
	find "$srcdir/.go/bin/" -type f -executable | while read filename; do
		install -DT "$filename" "$pkgdir/usr/bin/$(basename $filename)"
	done{{range .Files}}
	install -DT -m0755 "$srcdir/{{.Name}}" "$pkgdir/{{.Path}}"{{end}}
}
`))
