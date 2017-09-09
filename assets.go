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
		modtime: 1504952589,
		compressed: `
H4sIAAAAAAAA/wEAAP//AAAAAAAAAAA=
`,
	},

	"/css/flat.css": {
		local:   "../assets.min/css/flat.css",
		size:    581,
		modtime: 1504952589,
		compressed: `
H4sIAAAAAAAA/3yP0U7rMAyGX2VH5wakJerWwiCT+i5OnXYRrR056dhU7d3RRoAVBnfJlz//Z2uCvUpg
Y937GibLgk6UAPoxmuKkPYUxqU54DHUMQNdAASLTcpa5nHXLMqiGKQn3PzsDkLuBbSIVxA8gx+XlgkCd
k4nH1Htypvjnh8CSgNIsbKBJfu+WM9ZyM8bJQvNynotQNdyzmP9luQG72WZ5hmv3iFCdrqSzzox+q8Tn
h7Jqv1ViVbYl5F113uQ9YlbhsIjce1zkvzmmLONRWz7UWvjVEKc703qJSTU73+P9NIB0nlTiYNZFOGwD
IHrqvkAWFB/DnB+ubG7lKvd006aJ1U7+dhYzYfFpO70FAAD//91SfZ5FAgAA
`,
	},

	"/css/pipeline.css": {
		local:   "../assets.min/css/pipeline.css",
		size:    1259,
		modtime: 1504952589,
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
		modtime: 1504952589,
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

	"/img/failed.svg": {
		local:   "../assets.min/img/failed.svg",
		size:    726,
		modtime: 1504952589,
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
		modtime: 1504952589,
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
		modtime: 1504952589,
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
		modtime: 1504952589,
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
		modtime: 1504952589,
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
		modtime: 1504952589,
		compressed: `
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
7DThxJkgAAA=
`,
	},

	"/js/app.js": {
		local:   "../assets.min/js/app.js",
		size:    11321,
		modtime: 1504952589,
		compressed: `
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
/wsAAP//KYkcKDksAAA=
`,
	},

	"/project.html": {
		local:   "../assets.min/project.html",
		size:    2159,
		modtime: 1504952589,
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
		modtime: 1504952589,
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
		size:    682,
		modtime: 1504952589,
		compressed: `
H4sIAAAAAAAA/2yRQYvcPgzF7/MpxBz+/AubGei1GUOhLfTSU1voUbG1E3UcychKlnTZ716cTGEPewnP
4en556c+8QJy7aKKm+ZMdimmvyn6Z1nCL50hosBcCUgWNpWJxGFBYxwyVXCFgoYTufEfglVng2HmnE7w
faRKr6xoBAhX1QQlY6Q2W12NoGCtT2qpPhxGrS44UQU1sFmcJ4JK7izXCgNFbCw+0rrlifqekYAFUFYo
85A5Ai7Iud0LRkUru9q6Ie1zWEpmSgcWIIzjjgytBGQhA6xQ1+o0vf3sE3yZzUeyqeGvr1qKmPPbM4cn
9pGlsUPFeIsMP76eet8gY8ZaL8f9sH27UReyY+h9JEyhd2sy/LzHwTecaP+R5018jM4qtfdB09r8ba1G
hdAvx/9vtD7A0rzvWlP3HZ9IlnZFCs/PN1pfXna5+bbDHSxq7qbUvQ/9MLurgK+FLnd9Zx9cYHDpEsqV
7NjKfGSbLokyOYVDz/+cjwiP2LlhHTs9hv7M4T8ZavnwabP25z049OetidCfEy/hbwAAAP//H9oMFaoC
AAA=
`,
	},

	"/project_history.html": {
		local:   "../assets.min/project_history.html",
		size:    2003,
		modtime: 1504952589,
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
		size:    5513,
		modtime: 1504952589,
		compressed: `
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
AA==
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
