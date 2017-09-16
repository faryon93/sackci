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

	"/config.html": {
		local:   "../assets.min/config.html",
		size:    0,
		modtime: 1505568521,
		compressed: `
H4sIAAAAAAAA/wEAAP//AAAAAAAAAAA=
`,
	},

	"/css/flat.css": {
		local:   "../assets.min/css/flat.css",
		size:    692,
		modtime: 1505568521,
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
		modtime: 1505568521,
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
		size:    2760,
		modtime: 1505568521,
		compressed: `
H4sIAAAAAAAA/5xW3Y7qNhB+lUgrpKXCKMsSzq6joiP1uk9Q9cKJJ8E9ju3awxIa5d0rOwnkh91TVUhg
Ps+fPTOf54SV3GSaX5sTiPKE9CWOV61HG/0BtpD6Qq7U5VZL2QY5wzgXqiSoDU1iU7fbQkjZVEKRkYk0
0zVx4h+hSpppy8GSTNft9wq4YM9e+CI4nuj7+87U64aLj62x+i/IUQqHD6x5E/TF1JHTUvDoiXPeo72l
OPK7cRS37czasdC28lEi2En4u52pF8In21TMlkKRTCPqisbpoHMDesfexjSkubXolGw5IBPSHQVVGp9p
IaxDkp+E5OvBk4QC6Zup26de1w1x0viOHaW4wS+xP2zqfRCHVwkErwao0gqG6AJMQ2wDlGupLX10d+H2
Jq6OW8OcA755BEas6W29vjOWHGaaBRNyodmBd02e7ZN9MtO0Z6WEKjcP0btukb3OFC9M4FKxR2+KpWXX
Rca3UjN/qUM2QloTU6cINRImRaloDgrBjlXJh4DL0QPIMpJrhaBwYsI3xyl5AB2dYWqS+8T3EbJMQoS+
y45oj8ibD7Aocib7ICrBuYQuCKZAEi/aZ6TJWP6jtPqsOJVCAbOktIwLUPiMOrK+lzZ9tqL4tvIRbZ6K
opis4tU6UppYMMBw7q/L43/z1+XY++tX/8Nfn/yvHO4TDuUGLVPOMAsKx+tol6w2vmLuiyReTSTm/78N
gsNind69e1qDrgH913inOwDtfsg1ZUpUDIVWtD8EuSHRwUXdMSKhCqEEQvv9B1wLyypw0UK+iVejGyBG
OxEMe8LzF/jJ5i72JTe7Us9Kua4qgU2hFZJLx7Tf4vhxbd3lNz/ZX/DCw9r53N58f8EWj2vjc4MLgSmH
bM1ZShKK9ZidEbXadm05dGfYCtTsDWdnITnpCT3yiNTjwhwIdrfbpYOXF/9Jwz0XrBLySn/XiuV6U2ml
nWE5pFNS/9xTtPUVQ0/+ZX7gdb/f/0xX2eH0h8NhzG/hnLeHzvYvr6lvUKCpgMxICbX5ymsFzrESmstJ
IJBwXmoskItlJr1oy8OKZhbYD+L/p92TFOaQz+36X6gMXpslSQc+9zfUDw2H+xNdvPrP5P0e7b7u3994
Nux6cjk7msSrPqQ3f/p+KgnrrkYoO6MetbozQkW7L7vbi/iGDpzjpxNqNTKE53jdNfNi4/UQcyjXbfvU
P1YkY5Z4Owpsw4Uzkl3D899uc82hmRQcKKk3fdn9ppXTkrlNrs9WgI0UXO61uJzQvh3e/IS2tfqyhb/P
TN68FRLq1H91KQwZnY5OuzAf3jSpZA6JLsKsMp+yRnLHP3LJnPvl11xL8uex6+fFVDYeWNt/AwAA//91
Z4hFyAoAAA==
`,
	},

	"/home.html": {
		local:   "../assets.min/home.html",
		size:    236,
		modtime: 1505568521,
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
		modtime: 1505568521,
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
		modtime: 1505568521,
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
		modtime: 1505568521,
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
		modtime: 1505568521,
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
		size:    3443,
		modtime: 1505568521,
		compressed: `
H4sIAAAAAAAA/7RX227cNhO+/5+CkfHXCVBK6zhOGkcSEBQoUCAXBYpeB7PkSBqbIhkeZAeJ372gTrte
14eg7o1EznzzzUEzY2/5QhoRvlpkXehVXaYnU6DbCjXTLQdrK7C2LnsMwEQHzmOoYmj4L7OsC8Fy/BJp
qK55BC5MbyHQViETRgfUocoIK5QtZrOJhh6rgfDKGhd2qCuSoaskDiSQj5efSVMgUNwLUFid3CKQ6IUj
G8johWNfDTF0xu00gYLC2oO4FFQW063cgkfWOWyqoi4V6UvmUFU+fFXoO8Qw6VKK/rwoergWUudbY4IP
Dmy6CNMXq6A4zU/zd4XwfifLe9K58P5f8DdGBw5X6E2PxZv8Xb4ZXeyLf8CLkPrC50KZKBsFDkcXcAHX
haKtL0C3UYHjyoAk3fItuGKTv883xZ7kCe7GCBWEJ8BG2RNwliwq0jN0+vrMO/GDqeUXvjjJ3+ZnqyBl
c+Hrspg4n5ObOxMD/rce0JvoxLM40e3vuklTh38KZ5QqTvLTfFPoltMs535UPGNGPBLfjdDr/CzfFPsi
Hqzyz+nv8da+z9FFIvH0OdoHIfZArWFgQoH3VaZh2IJj04uTHtB5XK4NXaPkwdgs7d60ulKp0VUaBmoh
rbq6lLSQJQCQRscbFUnuq2bCDkGiq0u4Ld460HKaq6N1JUJdFpKGiYXkjF3iFkYpsLtQl3tWl1HtJbfo
NQxZGucxkUn7DUSgAc8Z+Y/j6eVxcfzqJkvhzTv4D2cuUAQ/RvOoMbSog79NMcnqj+PraTTC6Iba2zST
rP51fEU3ln4qUVT3Z7wcHbVdWNP3nbmqFA34l5UQ0K9ejti45SphlHHnR6fvAc7efhDReePOJTYQVahL
Wnw1wBrgnloNKqvLguqf9NbbD59oQDZxzwnvPNSlt6AXhlZ9tR0Jo9l64sqk2U58CblQmtbEsEt4aoz5
qWF4oAvnnDpMNTg/2Wz+vw/OnLli+CWCymZgT5rfBxZG8V7y18xOfaHIh8PR2FPVZWNcPxs3pEJq/R0d
aRsDb52J9lZZ9uQcpEwzdlhzBCe6qeZzncrRavFlXL/ExKwCgZ1REl2VTVEs8fsx+N7I9Adu5MzT/ypr
ZRNPXXa3op5X066RpvshJKU6sXRv6k+zzTxO7DdEmed5WXRvFlRUacaXuFb2Fyv91LsOLUKoshnISK+5
sO9syu58LU93thu1GZb7ACH6OyUVHYrLsRzUrPQzmFUVO7bgPcrjuegH1g57M+BD5g2Qutccr0UHun2Q
wEWtSbf3MHyJ6NNOeIjhCijsMUxtsyyYGV98+7ZYkry5qXfX1Bg3N9MIjqb/e8S22EZSck843j/r2N/c
HEa5qljNNln98ugfrV5N3ruz8cvO+Qe8DryPASWTGICUz+5+XFCoJbi9zFMEW9JyFwNeo4iphp8D9ci+
s/SC1mS7ETtgVUZccnOn5HLez1MuD3hcgd/ZcnyUbIpkrMHeJry7o94vvzR4+mEz0qbDwepMz78DAAD/
/xEGomJzDQAA
`,
	},

	"/js/ansi_up.js": {
		local:   "../assets.min/js/ansi_up.js",
		size:    8345,
		modtime: 1505568521,
		compressed: `
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
7DThxJkgAAA=
`,
	},

	"/js/app.js": {
		local:   "../assets.min/js/app.js",
		size:    13071,
		modtime: 1505568521,
		compressed: `
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
x4v6lpIiZSsk09b/BQAA//9GCrevDzMAAA==
`,
	},

	"/project.html": {
		local:   "../assets.min/project.html",
		size:    2159,
		modtime: 1505568521,
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
		size:    4276,
		modtime: 1505568521,
		compressed: `
H4sIAAAAAAAA/6RXW3PbthJ+z69YK+dY9pxQso+dzFQh2UnTTDMdJ+k4ec+AxJJEDQIMsLStyPrvHYAX
UaJ8afJCQQAW++3tWyDk4hpSyayNklpIHnAkJqQFlQepVmS0lGiiyui/MaXf3JZ4IDNhEg2B/wacqRzN
xImKLJqgMdrAQQS14pgJhRwOD6GfnU77v281Rzd1fnI+iUNLRqs8fnd5+elyEc7bv7Ba+b3rdTjn4noI
wugbsLSUGJXM5EIFpKvF6clJdbsFxWuJOi0DG1Itg5IH59AOdJZZpOB80p5KeEsBkyJXixQVoYlD0clm
DDIWpDrLELv9mVYUWPEdPQiKw7mIw+Is/qjB+9gCabCIB/CFXSEwSOsKdAbNKaANlGyZIJAReY4GWCMG
JVM1k3L5azgvzuLWDRtnbLk9Grr97g762el02/iKKZTgvwHHjNWStjY0K4nmS58Sm0yZWWJU29bminEu
VO5df/K6+5doIl0uTk+q235OYkaLs5PqdieE+0LyyofElsGrSRwWL9vVVJeloHEQCkyv+uzbwujsrpi1
yJ31Lhw7sgZLfY33C2dMyHuF8TYtXOrfL25qpYTKO3lbMeX2JkLxoKBSdiKNZbMSrWU5wh3U6krpG+Xk
nFB8qBJbvT56vlo1Aqou1+vjcF68HPjH52tZE/IR1NqiCfQARhwm8ZuaCm0W4TyJ4VFszG/eA63/eSqW
VJhUYqADpSktdjC99do8pmchg8JgFk1Wq5aHWjBfayNnBivJUjyargxm6+kL2IJrMDter/vIjOUbKpoA
MZMjRV8TydRV3Dt4cwzcgS20IUdALH7WO2r/sTOJKqeiq7hHzhs7cJfjuoo421TEvsp7oueZRMWZ2fH5
F1HinizoEoBE6VKSM8LFdLlcLoMPHwLO4f37RVlOfzwPCl2b3E0EBZPZDqbfa8NIaPUALt5ucdja4Y+D
sWiucdcxHzVvHdPXneY4CNwmYuPvIIauUeG3mkmoRIVSKNyi2jbEv/wAPxfIXBLEf7XnwmdPPiMIAzZv
skdfo8mkvlkUgnNUW52kp7Ec7aihDBcHyX7iidq1us5EaIVeQGqQEYJWCEIBFQhdZT8/mLclNN/UuODr
9bw7pbfMFR8QS2b766RF21iwB+Nor9tqsEJG0aSZEGrLuEm8WvnRTLES++vHrmLJlmge0utTt1OM1aab
Noe33eJRNI2cSK8iixJT+uymj/4jFMfbY6hFElS6clGNJn3tQGfBnlJx5NhKBO2FI5pMS11b9Hed6WbZ
82yJiqKmr4+K51uN1pdfl0FbpjkqFLnS5qldeCT9YBfebv8j2Qfbf6pzX/yVeAD7ThPvav/Rsu+Y+yfK
+k3q3PpQOcdhUhNpBbSsMGrHra6EFCSkOj1+nEidXk32RCAzaNtmDJdoagX+0h/OmzN9N1Z54Mt2ziox
vz69r3Z93s63ripzZkhkLCU7I2Zm+ffHMfad2aLMgOsbJTUbE3e30EBv7kl/IMGbTqEjjidEa+cm2gbv
9P8/Eb0LnW9KsCWRwd5aysCIvKA4ZNvXkJ9xtNT58PnCMdUtGyitcOS+TEgM/E499OCnChVcsht/5bkb
xv75hoikzv+sy+ro+N8oTAu8Nlr5uA01upPc++id4qOI7UagaWRS5xMQKhNKEAY2dU/WqNQGL3R+dDyU
wrKiZVffB008pM679nV4CAcuhYTKL3Q+cU3MeVHXVNUEKauoNsgh0wYa7kUOjspnI3hDRZsTx5eNSiiF
pqeegRsuGinw9A4XOp/NOjXEEomtq28Ep8I9M/8bh9TwAJmuI7nmO+gnzhahoDcb7kCKUtAXvXB2Xrgx
kGHpFSRLaDpKz4YDN8W+yRMfqFEmXq0aCfgfnK7Xm+XuLbP9pGjUM2VFqqX29y1v1v0VOvj+EwAA//+f
Mt/UtBAAAA==
`,
	},

	"/project_env.html": {
		local:   "../assets.min/project_env.html",
		size:    1174,
		modtime: 1505568521,
		compressed: `
H4sIAAAAAAAA/4xSzWrcMBC+71NM91ASiHeh0FNtQWgb6KWX/kCPY2l2ra4smZmRgxvy7kXebWqyS5uL
kMcz35+mdn6EuK9sisopBOJm4PSTrH6Mo/mRMliMkIWA4ug5xZ6iwojssQ0koAkGZOxJ2f8imFJmaLMP
bgNfOxJatCITIOxTcjAEtFRmRRMTDChyn9jJzapLohF7EkgMnKP6nkBI1ce9QEsWixbtaJrxYtIjhgMf
AeMEQ26Dt4Aj+lB4gWlI4jXxNEs6zuEwBE9u5SMQ2u4oGUoI6CMxoIBMotRftr2Bu8zaEfdF/rRIyWII
l2dW9147H4t2ELQH6+Hbp02ts0gbUKRZHz/ms+rSSLw2tXaEztTK5XpqtClUvavemu8ndPiMPV38H/J5
/Y25tepTlFrb5KYCXnaAaSDUZn11oOkGxjJ6XWI9LcSG4lj0uCcwR+bh4UDT4+NZdZ6e66Zus2qKoNNA
zel+8ttqhFZj5TDuidflAXae+8ZRICWzqv2fzh3CDitllK5Ka1NvvXkdWxnefZhb6+0R+JSTM7WPQ9bi
Srp035TjLnFfCn1yFBp0rmT24uY5SbP088yDZGtJZP0E82qJY4O3h4JzdW2euxpClqWlW+ee/KxeynZO
JjjSBTbbkT0s6b7gSP/lc7TDHPRffBajpXB70SJTn0Zasr6fu/++23beelNvnR/N7wAAAP//t74M55YE
AAA=
`,
	},

	"/project_history.html": {
		local:   "../assets.min/project_history.html",
		size:    2003,
		modtime: 1505568521,
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
		modtime: 1505568521,
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
