// This file is automatically generated using github.com/mjibson/esc.

package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/admin.html": {
		local:   "../assets.min/admin.html",
		size:    0,
		modtime: 1505064091,
		compressed: `
H4sIAAAAAAAA/wEAAP//AAAAAAAAAAA=
`,
	},

	"/css/flat.css": {
		local:   "../assets.min/css/flat.css",
		size:    692,
		modtime: 1505064091,
		compressed: `
H4sIAAAAAAAA/3yQ3Y6jMAyFX6WrvdmVGkQL2+6kEu9icKDRBDsyoT9CvPuobToDnXbuwmfrOwcnBAcV
oOwKZwsYShY0ogTQ9p1Ox8SS74NqhHtfdB5oChQgMi1nO9d3UrO0qmIKwu670wOZJ7gMpLzYFuS8vH4g
UGPk9j6CkKUmDkwNvQsD98FZMjr9ZVvPEoDCTKOhCvZgljNWc9V3QwnV+6UxoarYsejfWbaFcruLtSJc
mw1CPk7qzJwRvVLi278srx+UmGd1BuP0r2bOO3slrVPA3DxIjYF1tomnTW614oX1yp8WHTuLi1gorqmS
8ZyUfCoS4aMmDn90baULqtpbh3+HFqSxpAJ7vU79aecB0VLzBWJAei9zGUzSzMrk5v/TtIRY7eXnzHQW
mH6mjR8BAAD//8bgaAm0AgAA
`,
	},

	"/css/pipeline.css": {
		local:   "../assets.min/css/pipeline.css",
		size:    1259,
		modtime: 1505064091,
		compressed: `
H4sIAAAAAAAA/6RU4W6jPBB8lUifqv75jCANrc48jcFrsqrjteylIUX37iewe8k1bnXSiR+A7dnZmR3Q
+FZ59GDRwa6KrEaIC8PMQlkcnRzAMYTufEQGEb0aQDo6B+V/lpD5vmiM3qqLRLfuit7S8NoZciwivoNs
Ws/p9Qw4Hlm+1HV3UmFEJ3pippNsaj93Z9R8lM1z6+c7OvBlFk8RGcnJAFYxvkGuclgLHhPd9pz5wrbS
wql7F+g0zLLuegoagghK4xRlWz90A1kK8j9jTDdMIVKQnnBz5qpq3/q52zq5pbnzsqCkwtFRAL3k2mpi
KhyTPRgK8H9hRxmGsAzkGBzLx8csepWVxWRfWz/vIlnUV6NUH8lODB2Tl81LtcpIesegLqU+EpsFw7Kp
64evW12Su18dMhgiCzKCLx6+UWfVzbFE/jF7Rw5KpW8hyx+jrksDAD1CtQnKjWfXtqWrZ3mwayA3t8Rh
/zum+6vXOTj1br2acnoTZWopSfoAp0D+A+fKuCvK9CpG0Euvhtcx0OS0yMF++qFU+1yCGIW2CNH9oT20
JUiYnEM3FjCmf/ou/neANYCfvrBPeKsuEP7uh/UrAAD//8EK/7DrBAAA
`,
	},

	"/css/style.css": {
		local:   "../assets.min/css/style.css",
		size:    2768,
		modtime: 1505064091,
		compressed: `
H4sIAAAAAAAA/5xWi27rKBD9FUtVpEYKkfNweoOl6P7AfgQ2g8MWA4JJE6/lf19hm7yc9q62kZrJGeYB
zJzhiLVaFIY37RFkdUS6StNZF9DWfIETypxJQ33pjFJdv84yzqWuCBpLs9ReuqWQSrW11OTORX4vF+ZC
vPxH6ooWxnFwpDCX7ncNXLL3YHiWHI90v1/by7zl8mtpnfkbSlTS48RzzIseJeeg85q5SmpSGERTU7IP
f/aSxzxHPMJDAnRlL4k3SvLkjXM+omMeaRK0aZJ23VMuB2FcHfaL4B4OYr22l8nio2sfc0ufk0pj4ODj
MaVnb8kxW3JAJpU/SKoNvlMhnUdSHqXi8xhJgUD6y166t9HWxzxpesMOSl7hVRo2m4cYxGOjgGBjgWqj
IWbXw7TPLUKlUcbRV2fXn95DqMPSMu+BL16BCWtHX5s9Y9nuyVIwqSaWA3iz5MU222ZPlu6ktdTV4iV6
sxXF5snwzCRODUf0alg51kxufKkMC4cab6O/1sxecoQLEqZkpWkJGsHdm5IvCedDAJAVpDQaQeODi9Bm
x+wFdPCW6Ye7z0JHIisUJBj69YDugLz9AoeyZGpMopacKxiSYBoUCUvHG2kLVn5Wzpw0p0pqYI5UjnEJ
Gt/RJC504mK8rSS9SiGjxZsQ4kFKZ/NEG+LAAsPneMM9/rd4wx2HeKP0P+KNl/9TwG3GoVqgY9pb5kDj
vZyss9kiVMxNyNLZw4rn3x9xYRTm+S16IEUYGjD8u9cMG6DDF2lypmXNUBpNx02QK5LsfDJsI5FaSC0R
ut+f0AjHavDJZH2bzu5OgFjjZe84EF44wG+U6zSU3NORBlYqTV1LbIXRSM4DT3+k6evauq1f/EE/4YWX
tfO9v2f9hC1e18b3DicLHjlkaU9Kkb5YD8UJ0ejl0JaxO3tVT83BcXGSipOR0JOAKHNfmJFg1+t1HqOs
wifvz1mwWqqG/mU0K82iNtp4y0rIH0n9+0jJMlQMPYZZ+iLqdrv9k612cfe73e6e3/p9XgedG+f23UDu
aapHnkgJjf0pag3eswra81EikH6/1DogZ8dsfjaO9xItHLBPEn7nw0jqXzTf+w3fUFts2ilJ93weTmh8
NOxuI1pswudhft9pN9v9L15EbSCXk6dZOhtT+hV2P75penmoEcpOaO5a3Vupk/WP3R2WhIbuOSe8Tqgz
yBDe0/nQzBPFZpdyqOZd9zYOK1IwR4IfDa7l0lvFmn78d0umvVwWRvFJdw8qB/zaA2JX7iJeOQAdNcVK
8I/9qGlAKXO+GQmxX8U46gRRsd+VhYgJ1KwCjbduEx8bwUdd2bAYKIjdvwEAAP//h990XNAKAAA=
`,
	},

	"/home.html": {
		local:   "../assets.min/home.html",
		size:    236,
		modtime: 1505064091,
		compressed: `
H4sIAAAAAAAA/ySOMWoDMRBFrzK5wAqTKmEtMHhJ5cYuUoaRPNIoljRmNJvg24eNiw+/eA/ezDv/STVK
IzCBgfEWy8vseOcvVCkaINxVvrcnHYwJtGS2Df4p9AthLfUKXIaJPkAUImPPNGCQWel5THNQv+1CBLIq
zAislPZsdh/vzuVivIYpSnMJ9SH97dU9O1yoElzDYaTuvByOp2VqVzDUTLb/ChX7zX/863BeDqfjMjv0
kEQhrWpMCqUn0YZWpE9/AQAA///Irf0y7AAAAA==
`,
	},

	"/img/failed.svg": {
		local:   "../assets.min/img/failed.svg",
		size:    726,
		modtime: 1505064091,
		compressed: `
H4sIAAAAAAAA/4yS34vcIBDH/xWZpdDCJmpiuGuJeSrcU58K926iRotrgrr5cX99sb0e2+1uOQRhxvl8
5zuDbVxGtJ2cjxxMSvMXjNd1Lde6nMKIK0IIjssIaLUyGQ6PFJBRdjSJQ0Wga531SoSnIKRVPiErOQhA
W8WBANorDpSQD9C1MU0zmrSOKv16ynExTG4KHA59379mplkMNu0cSgr4iqK3a/DfFro2qCGhsHGo79lG
2jrH4dA0TVa4BPL98MYxdpNTpJGMZXQWyVxlkeTwrX5AxLClIqZgL9c9/mvqHNzHg/iUmfGPstYaUFJb
KoQfTF7ZyUrpFCA9+VRocbJu5/BV/RDPZ/Rd+Hh8VkEKL45PyqtFHKPwsYgqWP3KRPuiOFAKXZuF89z0
c9kA2jnQ5m1SQvP5HV4svoauP1snW5zhWxLsTkXz+P4WWlhn/XhPgv1Tgceuzd+1+xkAAP//R7Dg19YC
AAA=
`,
	},

	"/img/passed.svg": {
		local:   "../assets.min/img/passed.svg",
		size:    724,
		modtime: 1505064091,
		compressed: `
H4sIAAAAAAAA/5SSXasbIRCG/4pMKLSQXXU/ejhl3avCuepV4dy7q64W4y5q9iO/vpimIU2TQhGEGed5
553BJswDWg/WBQY6xukLxsuy5EuZj37ABSEEh3kAtBgRNYNXAkhLM+jIoCDQNtY4yf2b58JIF5ERDDig
tWBAAG0FA0rIB2ibEMcJjUoFGc9PKc760Y6ewa7ruktmnHhv4sYgp4DvKPq4Bv9poW287CPyK4PymW2k
jLUMdnVdJ4VbIN0vV64uH3JVf+488ahvU0gw+Fa+IKKruSA6q0736v+0c/T2445/SszwW1YpBSjKNWbc
9Tot62CEsBKQGl3MFD8YuzH4Kn/w9yP6zl3Yv0svuOP7N+nkzPeBu5AF6Y26MMGcJANKoW2ScJqYvuY1
oI0Bra8zEprOr/Bm5SW03dFY0eAEP5KonlR8Lv6nycRDMG54LlL9VYOHtkmftf0ZAAD//yNy3HPUAgAA
`,
	},

	"/img/running.svg": {
		local:   "../assets.min/img/running.svg",
		size:    724,
		modtime: 1505064091,
		compressed: `
H4sIAAAAAAAA/5SS3arcIBCAX0VmKbSwiZqfHk6JuSqcq14Vzr2JGi2uCWr+9umL7bJst7uFEgjMON/n
zGATlgFtJ+sCAx3j9AXjdV3ztcxHP+CCEILDMgBajYiawSsBpKUZdGRQEGgba5zk/s1zYaSLyAgGHNBW
MCCA9oIBJeQDtE2I44RGpYKMv45SnPWjHT2DQ9d1l8w48d7EnUFOAd9R9HEN/rOFtvGyj8hvDMpnbSNl
rGVwqOs6GW6B9H+5cnX5kFNdmbiJR32bQoLBt/IFEV0tBdFZdb63/7Od2duPB/4pMcNVqxSgKLeYcdfr
tKyTEcJKQGp0MVP8ZOzO4Kv8wd9n9J27cHyXXnDHj2/SyYUfA3chC9IbdWGCOUsGlELbJHGamL7mNaCd
Aa2vMxKavt/hzcpLaLvZWNHgBD9SVE8qPhf/c4mfnTNueC6p/qrBQ9ukx9r+DAAA//8uS+XC1AIAAA==
`,
	},

	"/img/waiting.svg": {
		local:   "../assets.min/img/waiting.svg",
		size:    730,
		modtime: 1505064091,
		compressed: `
H4sIAAAAAAAA/5SSXavbMAyG/4pRGWzQxHY+djgjztXgXO1q0HsntmMP1wm2m4/++uGtK13XDg6CgBQ9
r14JN2Ee0Hq0LjDQMU5fMF6WJV/KfPQDLgghOMwDoMWIqBm8EkBamkFHBgWBtrHGSe7fPBdGuoiMYMAB
rQUDAmgrGFBCPkDbhDhOaFQqyPjrV8qzfrSjZ7Druu5SGSfem7gxyCngO4o+7sF/W2gbL/uI/MqgfGYb
KWMtg11d10nhFkjflytXlw85UaZI6MSjvqsiweBb+YKIruaC6Kw638/4r6mTtx93/FNihj/KSilAUa4x
467X6WRHI4SVgNToYqb40diNwVf5gx9O6Dt3YX+QXnDH92/SyZnvA3chC9IbdWGCOUsGlELbJOG0N33N
a0AbA1pfNyU0xe/05vAltN3JWNHgBD+SqJ50fC7eM2ThJho3PBep/unBQ9ukJ9v+DAAA//9SPclr2gIA
AA==
`,
	},

	"/index.html": {
		local:   "../assets.min/index.html",
		size:    3408,
		modtime: 1505064091,
		compressed: `
H4sIAAAAAAAA/7RXS2/duA7e31+hOsBNC1zZadO0t6ltoJsBBuhigMGsCx6JtpnIkqqHk6LNfx/Ir+Oc
NI9iMhtb4uMjP5pkcsoX0ojwzSLrQq/qMj2ZAt1WqJluOVhbgbV12WMAJjpwHkMVQ8P/P8u6ECzHr5GG
6ppH4ML0FgLtFDJhdEAdqoywQtliNrto6LEaCK+scWFvdUUydJXEgQTy8fI/0hQIFPcCFFavbwFI9MKR
DWT0grFVQwydcXtNoKCw9iAuBZXFdCt34JF1DpuqqEtF+pI5VJUP3xT6DjFMukTRnxdFD9dC6nxnTPDB
gU0XYfpiFRSn+Wn+vhDe72V5TzoX3v8D/MbowOEKvemxeJu/z0/GEFvxL0QRUl/4XCgTZaPA4RgCLuC6
ULTzBeg2KnBcGZCkW74DV5zkH/KTYiN5QrgxQwXhCWaj7Al2liwq0rPp9PWZd+IXqeUXvnidv8vPVkFi
c+HrspgwnxObOxMD/rsR0JvoxLME0e3vuklTh38KZ5QqXuen+UmhW06znPtR8YyMeCS+H6E3+Vl+UmxF
PFjlnzPe4619X6CLBOLpS7QPmtgDtYaBCQXeV5mGYQeOTS9OekDncbk2dI2SB2OztHvT6kqlRldpGKiF
tOrqUtIClgyANDreqEhyq5oBOwSJri7htnjnQMtpro7WlQh1WUgaJhSSs+2StzBKgd2nutyzuoxqQ27R
axiyNM4jkUn7HUSgAc8Z+U/j6eVxcfzqJkvpzTv4D2cuUAQ/ZvOoM7Sog78NMcnqT+PriTCyJ32AkkT1
p/Sk1IOp8lOForqf8HJ01HZhZe87c1UpGvAvKyGgX6McsXHJVcIo486PTj8AnL37KKLzxp1LbCCqUJe0
xGqANcA9tRpUVpcF1f/VO28/fqYB2YQ9891HqEtvQS8IrfpmOxJGs/XElUmjnfCS5QJpWhPDnvDUF/NT
w7DttOygC1lDKuW3sXDm6idSYRTvJX/D7PTRFflw2PcbVV02xvWzc0MqpL7ew5G2MfDWmWhvkd7IOUiZ
BuiwoghOdFNF5yqUo9cSy7h+yYlZBQI7oyS6KpuyWPL3Y/K9kemv14iZp39E1rolnLrsbmU97519m0z3
Q5NEdULp3tafZ595VthviDLP87Lo3i5WUaUBXvJa0V+s8FNnOrQIocpmQ0Z65cJ+sInd+Vqe7mw/R7NZ
7gOE6O+UVHQoLsdyULPCz8asqtixBe9RHs9FP/B22JsBH3JvgNS97ngtOtDtgwAuak26vQfha0SfJv4h
hCugsEGY2mZZH7N98f374kny5qbeX1Nj3NxMAza6/ucR32IXScmNcLx/0bG/uTnMclWxmp1k9cujn3q9
mqJ3Z+OXnfkHvA68jwElkxiAlM/uflxQqCW4DfOUwY603OeA1yhiquGXQD2yHyy9oDXZfsQOUJURl9zc
KbmM0/aduDwQcTX8wZbjo2BTJmMNNnvu7o76sPyM4OlXywibDgeLMT3/DgAA//9w66k/UA0AAA==
`,
	},

	"/js/ansi_up.js": {
		local:   "../assets.min/js/ansi_up.js",
		size:    8345,
		modtime: 1505064091,
		compressed: `
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
7DThxJkgAAA=
`,
	},

	"/js/app.js": {
		local:   "../assets.min/js/app.js",
		size:    12976,
		modtime: 1505064091,
		compressed: `
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
`,
	},

	"/project.html": {
		local:   "../assets.min/project.html",
		size:    2159,
		modtime: 1505064091,
		compressed: `
H4sIAAAAAAAA/5yVUW+jOBDH3/spfD6pbh8outce9kmVetp9W20/wMrAhMyuM0bG0EYp331lMCRNadLN
S2Jm/vzMjMczWYkdK4xuGlk7+xMKn3QIz4yqpLDknTUG3OQK1mZtn6WxuoRSZet/VJar3S7670hvoO+z
NFdXGUYsX2m20ommykBS2jY3kDis1p6rLEWVNbWmAMaV5F7nTEom8hZNKbi6pryp/30IT1kahOpqUb/G
xlu3FVwNWvZlfD75To01GCQQXH2Ly5N6oE5w9UgdOksbIH9S3YD3SFUjuHqKy0X9ECn7S7KWSlghQcmu
r9lsFUZ7aLzg6ubv3W4w9/1tJGVpyH9rpkST7hjpLvE6b7jKDA6nOPp2uvDYwT17m+Keq0wH2drBSqbx
HNP9iWLZ9+mgTef9xyTfXxjIPo4TCZGHnNdXNlv3nHExp0KfC3iqkTHkk/FG6XExnd1iLqlP7DFpD4rv
LD+U4CfQQN3bOj0Lnqv1E/RJe1DXIz92kdaY8X6rLG+9t8T8tgYZ1/ELck8s95TUDjfabfnBx3mHVQXu
yWvfNmPYzlkn2H9MhFdKTRU4we7Hx0gQEYHFLxkJN7fBlJvWTZYH0wZrbWvbgUs8egOSfyUPjrRhT+A6
cOwxbMdZi3kSlXK3i4TB1/czoja6gJBkaWDl9+RRLbkgGzrM7MAmsTXQx1Hy435yrHtzLd5TBFfHnbc2
ejv22rGdfm9puXW9gz1rDCe8wCxsFf6aGuky9BTtMRheCqM32qOlkMVhahzu8L9GA+V868eqUlnaGnUw
ykIbHuYXkGeFAe1W+MLfKWpNwMa7wD8YQYOZCtOWILmIF+HH4Lxb+40JIaQldn/EnsfVMj26L+cPo2qZ
DdRdzt0PtWX45D/a4eD3dwAAAP//JqIZ7m8IAAA=
`,
	},

	"/project_build.html": {
		local:   "../assets.min/project_build.html",
		size:    4270,
		modtime: 1505064091,
		compressed: `
H4sIAAAAAAAA/6RX23LbNhO+z1Oslf+37Gko2bWTmSokO2maaabjJB0n9xmQWJKoQYAFlrYVWe/eAXgQ
JcqHJjcUBGCxx+9bIOTiGlLJrI2SWkgecCQmpAWVB6lWZLSUaKLK6L8xpd/clnggM2ESDYH/BpypHM3E
iYosmqAx2sBBBLXimAmFHA4PoZ+dTvu/bzVHN3V+cj6JQ0tGqzx+d3n56XIRztu/sFr5vet1OOfiemiE
0TdgaSkxKpnJhQpIV4vTk5PqdssUryXqtAx8SLUMSh6cQzvQWWaRgvNJeyrhLQVMilwtUlSEJg5FJ5sx
yFiQ6ixD7PZnWlFgxTf0RlAczkUcFmfxRw0+xhZIg0U8gC/sCoFBWlegM2hOAW2gZMsEgYzIczTAGjEo
maqZlMtfw3lxFrdh2ARjK+zRMOx3d9DPTqfbzldMoQT/DThmrJa0taFZSTRf+pLYVMrMEqPatj5XjHOh
ch/6k9fdv0QT6XJxelLd9nMSM1qcnVS3Oyncl5JXPiW2DF5N4rB42a6muiwFjZNQYHrVV9+Wjc7vilmL
3Hnv0rEja7DU13i/cMaEvFcYb9PClf794qZWSqi8k7cVU25vIhQPCiplJ9J4NivRWpYj3EGtrpS+UU7O
CcWHKrHV66Pnq1UjoOpyvT4O58XLQXx8vZY1IR+ZWls0gR6YEYdJ/KamQptFOE9ieNQ25jfvMa3/eaot
qTCpxEAHSlNa7Nj01mvzNj0LGRQGs2iyWrU81BrztTZyZrCSLMWj6cpgtp6+gC1zDWbH63WfmbF8Q0UT
IGZypOhrIpm6ivsAb46BO7CFNuQIiMXP+kDtP3YmUeVUdIh75LxxAHc5rkPE2QYR+5D3xMgziYozsxPz
L6LEPVXQFQCJ0pUkZ4SL6XK5XAYfPgScw/v3i7Kcfn8dFLo2uZsICiazHZt+rw0jodUDdvF2i7OtHX6/
MRbNNe4G5qPmbWB63GmOg8RtMjb+DnLoGlUlKpRC4RbJtsn95TuYuUDm0h//1Z4Lnz3tjJQPeLypG32N
JpP6ZlEIzlFt9ZCewHK0o1YyXByU+YmnaNfkOhehFXoBqUFGCFohCAVUIHSYfn4wb8Ez36Bb8PV63p3S
e+ZgB8SS2X6EtNY2HuyxcbTXbTVYIaNo0kwIteXcJF6t/GimWIn9xWNXsWRLNA/p9UXbKcZq00ebw9s+
8ag1jZxIryKLElP67KaP/icUx9tjqEUSVLpyWY0mPWqg82APSBwtthJBe9WIJtNS1xb9LWe6WfYMW6Ki
qOnoI9j8U6P1wOsqaMs1R4IiV9o8tf+OpB/sv9uNfyT7YONPde5hX4kHbN9p3x3qHwV8x9k/AOs3qQvr
Q3COw6Qm0gpoWWHUjltdCSlISHV6/DiROr2a7MlAZtC2bRgu0dQK/HU/nDdn+j6s8sDDds4qMb8+vQ+7
vm7nW5eUOTMkMpaSnREzs/zb4zb2PdmizIDrGyU1G1N2t9CY3tyQ/kCCN51CRxxPyNbOHbRN3unPP5C9
C51vINiSyGBvLWVgRF5QHLLtC8iPBFrqfPhw4Zjqlg2UVjgKXyYkBn6nHkbwU4UKLtmNv+zcDXP/fENE
Uud/1mV1dPxfFKYFXhutfN6GGt1J7mX0TvFRxnYz0DQyqfMJCJUJJQgDm7rHalRqgxc6PzoeSmFZ0bLD
90GTD6nzrn0dHsKBKyGh8gudT1wTc1HUNVU1Qcoqqg1yyLSBhnuRg6Py2ci8oaLNieNrRiWUQtNTzyAM
F40UeHqHC53PZp0aYonENtQ3glPhHpj/j0NqeIBM15Fc8x30E+eLUNC7DXcgRSnoi144Py/cGMiw9AqS
JTQdpWfDQZhi3+SJD9QoE69WjQT8BKfr9Wa5e8VsPyYa9UxZkWqp/U3Lu3U/QgfffwMAAP//JVzJxK4Q
AAA=
`,
	},

	"/project_env.html": {
		local:   "../assets.min/project_env.html",
		size:    1122,
		modtime: 1505064091,
		compressed: `
H4sIAAAAAAAA/4xSzWrcSBC+z1PUzmGxwZqBva4kMLsx5JJLfiDHUnfNqDOtLlFVLaMYv3toaeIMyZD4
Iqqbr7+/Uu3DBOlYOU4mHCNJMwp/IWdv0tR+5gwOE2QloDQF4TRQMphQAnaRFIxhRMGBTMJXgpmzQJdD
9Dv40JPSBRSFAOHI7GGM6Ki8VWMhGFH1kcXr3aZntYQDKbCA5GRhIFAyC+mo0JHD4sV6mhe+xLZyeAgJ
MM0w5i4GBzhhiEUXhEbWYCzzYml9h+MYA/lNSEDo+tUylBIwJBJABZ3VaLgeewcPWawnGYr9+aIlhzFe
f7N5DNaHVLyDoju5AB/f7mpbTLqIqs12PSzfqueJZNvW1hP6tjYpY/vpTAfvcKD1IuZluHcWOGltHfu5
4MtahUZCa7Y3J5rvYCrY29LUecc7SlOR8O3T04nm5+d1XHDL4WzMcawGX/3T1l024wQ2j9Sc57P3zhJ0
liqP6UiyLWUeggyNp0hG7aYO35EHhANWJqh9xdu23of279Tp+O//C7Ter8TnzL6tQxqzlTja82NTPg8s
Q7kY2FNs0PtSx6vBS2cr+noGzc6R6vaF5q9LHheDOxWem9v251RjzHoZ6d77lzyb16r9KqY40RU115M7
Xcq9x4n+qOfpgDna7/QcJkfx/mpEoYEnulT9b0H/2Nt++YPbeu/D1H4LAAD//6ea5IpiBAAA
`,
	},

	"/project_history.html": {
		local:   "../assets.min/project_history.html",
		size:    2003,
		modtime: 1505064091,
		compressed: `
H4sIAAAAAAAA/5yV23LcNgyGX4VWZrzJBVfrOsnFRlKnh8z4okk7nrwAJUISaxLQkJC96nrfvaOT7ZXX
beIbDQ8A8eEHKCba3AqsZEHInqwFnzae/oaCr0xg8l02WBRWhZBGyoJnMXylVliBj3rnUNNdGoH35MVZ
KlrUUBoELc7PxcPqahVlSWBPWGWfr6//vN4m8TQV+/1gdjgksTa3T0N6uhOBOwupU74yKJma7cVm0+we
A5+Fokf/tQ1dH/Esb43VYW0BK66jowQKstJp+V5MAyrLACzfR1MQhh1LZU2F2wKQwWeJmX1LJUolCypL
gNm+JGQZzD8wMHGWxCZL6svsK4mRQjCJAHAmvqkbEEoUbSOoFOMpgrxwqstBsDdVBV6o0U04ha2ytvs5
ievLbFLlURuDpUHDIMfMU4Qd/6UqePtuuSW1CaywgHRzakvlFnQaPRHw/l6Ms8+oj7VrFIIVw1dqKFVr
eai+hwYUp9FIbnDK/Mh59MpJd0OzDWuD2Tqw4jZMcjZKa4PVUOTNp3mWEzO57cWm2T2sWSh5e7lpdotm
OVXtj0O1g5MfoyypP0y7BTln+Hl9ayhuhrxMOeU0M6apWDUqBNB9K/eVXvh6cHQLLzuXytgXnWFX1P2V
etndt4gGq9k/NAqzRInaQ5nG062N9/tptDb6cIiHQ+L9fjwMW3c4jJ59lNygljU7OwcbNVk7CEFVIO5F
izdId9hHHMKdYx6aT2/fHB34LonVbJDE9YcnIg/3ybUM+lm+bQAv6SiXPPul5Zr8NonzTPwvphqMT1D+
MEthfGFBkkTiol4w/TZEG5neHEEteDyU4l6EmjyfRln+3OYGvXxs0FMX4TtzUBZQK7+g/2YcnNBzRmfj
+jprxbBddV3XyS9fpNbi6mrr3Or1itbU+qpfkLWy5YLp99YrNoT/waUnk55tGr4eJoC/haUwX0lPwjw0
M2non6BnFXv+PX6g5pfo8Tf61KIgq/sqX/z03Y9MaAwi+Hk4kv9Bqu8J4ciDWK/XL8H9GwAA//8cuyPd
0wcAAA==
`,
	},

	"/project_settings.html": {
		local:   "../assets.min/project_settings.html",
		size:    5512,
		modtime: 1505064091,
		compressed: `
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
AAA=
`,
	},

	"/": {
		isDir: true,
		local: "../assets.min",
	},

	"/css": {
		isDir: true,
		local: "../assets.min/css",
	},

	"/img": {
		isDir: true,
		local: "../assets.min/img",
	},

	"/js": {
		isDir: true,
		local: "../assets.min/js",
	},
}
