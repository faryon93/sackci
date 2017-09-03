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
		modtime: 1504452996,
		compressed: `
H4sIAAAAAAAA/wEAAP//AAAAAAAAAAA=
`,
	},

	"/css/flat.css": {
		local:   "../assets.min/css/flat.css",
		size:    581,
		modtime: 1504452996,
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
		modtime: 1504452996,
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
		size:    2787,
		modtime: 1504452996,
		compressed: `
H4sIAAAAAAAA/5xWj27rKg9/lUhTpVUqVfon3SmVqvMC30OQYFK+EUDgrs2N8u5XJKFtmu6c6S7S5tnY
/mHjH5ywUovc8Lo5gSxPSFdpOmuDtjFf4IQyF1JTXzijVNuts4xzqUuCxtIstdd2KaRSTSU1eQhxeJRz
cyVe/iN1SXPjODiSm2v7uwIu2XtwvEiOJ7rfr+113nD5tbTO/B8KVNLjJHLERU+Sc9CHirlSapIbRFNR
sg8/9nqIOAd9VPcA6MpeE2+U5Mkb53zQDjjSJFjTJG3bJyxHYVwV9ovgRoVYr+11svjkmjG29BlUGhOH
GGNIz9GSU7bkgEwqf5RUG3ynQjqPpDhJxecxkwKB9Je9tm+Dr484aXrXHZW8qVdp2Owh5CAeawUEawtU
Gw0RXaemHbaoKowyjr6qXVe9Uarj0jLvgS9eKRPWDLE2e8ay3ZOnYFJNPHvl3ZPn22ybPXm6s9ZSl4uX
2ruvyDdPjhcmceo4aG+OpWP1pONLZVgoauxG19bMXg8IVyRMyVLTAjSCe3QlXxIux6BAlpPCaASNoxBh
zE7ZC9XRW6ZHvc/CRCLLFSQY5vWI7oi8+QKHsmBqAFFJzhX0IJgGRcLSoSNNzorP0pmz5lRJDcyR0jEu
QeM7msSFSVwM3UrSmxQQLd6EECMpnc0TbYgDCwyf8/V9/Fm+vsch3yD9h3xD8/+UcJtxKBfomPaWOdD4
KCfrbLYIJ+YuZOlstOL5/4+4MArzwz17IEXoBzD8erT0G6D9H1IfmJYVQ2k0HTZBbppk55N+G4nUQmqJ
0P7+hFo4VoFPJuubdPZQAWKNl13gQHihgN8Y12k4ck8lDaxUmKqS2AijkVx6nv5I09dn675+8Rf7hBde
np3v4z3bJ2zx+mx8H3CyYMwhS3tWinSH9ZifEY1e9mMZp7MzddQcAudnqTgZCD0JGmUeD2Yk2PV6fYhZ
VuE7dHUWrJKqpv8zmhVmURltvGUF3O6Xjgru102fO/0+c8c+4RDRU7heXwDZbrc/cO/4qJO0i9XZ7XZj
XB0jjqryRJGd8tBfKJswGFx6q1hNpQ6hSa5M8fljNBV4z0poLieJQLo6UevgL/5QWaybKW93FB8qNLwj
dvdbW2zCN7rSH6yb7f4Xz6M18M3Z0yydDfv8FfY5PHM6uS8QZWc0D9PvrdTJ+o8DH5aEGe9oKDxYqDPI
EN7TeT/fE8Nml3Io5237NtxfJGeOhDgaXBOLH14E7ZJpL5e5UXwy8L3JAb+NhdgVu6gvHYCOlnwl+Md+
sNSglLncnYTYr2IedYZo2O+KXEQAFStB430AxcdG8MFW1CwmCmL7bwAAAP//WYaZa+MKAAA=
`,
	},

	"/img/failed.svg": {
		local:   "../assets.min/img/failed.svg",
		size:    726,
		modtime: 1504452996,
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
		modtime: 1504452996,
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
		modtime: 1504452996,
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
		modtime: 1504452996,
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
		size:    3402,
		modtime: 1504452996,
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
xGqANcA9tRpUVpcF1f/VO28/fqYB2YQ9891HqEtvQS8IrfpmOxJGs/XElUmjnfCSZf3ZtCaGPdWpI+an
hmHbY9lB/7GGVMpsY+HM1U+kwijeS/6G2elzK/LhsOM3qrpsjOtn54ZUSB29hyNtY+CtM9HeoruRc5Ay
jc5hLRGc6KZazvzL0WuJZVy/5MSsAoGdURJdlU1ZLPn7MfneyPR3a8TM078ga90STl12t7KeN86+Qab7
oUmiOqF0b+vPs888Jew3RJnneVl0bxerqNLoLnmt6C9W+KknHVqEUGWzISO9cmE/2MTufC1Pd7afoNks
9wFC9HdKKjoUl2M5qFnhZ2NWVezYgvcoj+eiH3g77M2AD7k3QOped7wWHej2QQAXtSbd3oPwNaJPs/4Q
whVQ2CBMbbMsjtm++P598SR5c1Pvr6kxbm6mARtd//OIb7GLpORGON6/6Njf3BxmuapYzU6y+uXRT71e
TdG7s/HLzvwDXgfex4CSSQxAymd3Py4o1BLchnnKYEda7nPAaxQx1fBLoB7ZD5Ze0JpsP2IHqMqIS27u
lFzGae9OXB6IuBr+YMvxUbApk7EGmz13d0d9WH5A8PR7ZYRNh4PFmJ5/BwAA//9jHTu8Sg0AAA==
`,
	},

	"/js/ansi_up.js": {
		local:   "../assets.min/js/ansi_up.js",
		size:    8345,
		modtime: 1504452996,
		compressed: `
H4sIAAAAAAAA/6xZ23LbONK+11PQnCqJHFEUJdsTxxTsf/6t2aq52GRrsrs3JIcFUiCFmKcFoEguRfvs
WzjwKNpxtpJyJBLo/vrrBtBoQEayL2KGy8IgZcmsBMasJM/mCScGe65QmWhblOACAQBmtehsOpWNNsy3
5kk+G94MHauSMDoLGhj3jDKKtBZMiXC0MvqMYjabTvtddlFu0QeYoysAZpQRXKQz86QADSWkgE9fIOF6
4HR2OxKmy32xfy0o/mcF0LGytyiB+4y557PBdphajdM13knfU6RxczHT3bpbI+nRYHmVfYw+m8IY3UeU
AS9wk5IYvCHEYOWGeANJus9RwaidoSJlOzfE87l5EvJeiBerADQiXogD9zzh6gSl6PgPdGRAmbEJPHhO
4PLOAyXpESz/9On8q0/nfsE/f/J86n8K5o9+sUxzIceObA0aJJugKoMxMoS6NZuZLkFsTwqtQAftD5T+
dqwMrmLN8pmpeKhYtbPBPDVBkH2GeeKhs//12x+ffv/4Aehr27HXuisaYUFxGJdZSSjwvBNJo3vPsRzL
Caw4g5SGBczRvc7FFlEG4yf9bEmp1d27F+QI2jZSjiXkRqRSglDRQ3tJ8hllWXnoQArQUYJ7NOA3LpjD
FBUMDliOy8bPsKW5vr211P8R0cMOM6SfA0sF8u7WEn9jVAlOd2wQUo78DYV+bLnC6EhJ4X6Ia+6vwQ9C
LR0Yd7bxoBNz5cCrCsPY370WUaXzHWOgNOqhCOQ037E8+2tJcsgYIuDECCxoUpL8vlk2CYEpX+MWLiiD
RYxk2mBHBuoum6Ejc3lKPDKVLAAAjjlRq5TJ3hrADhGNYYXCpCQhJ2BOOFrTXWbboQRHNjnGVWNT+jOd
Ng1JCgAo9lnWaYvqth4XkfbYc4Yoz3v8TUSrfU1S0IEVTVGnKUo5lVF7gpQ5SdLWnU4m8VaB9y4QfrTB
2FMUKvtik0pScyLZ2dWe7gxd6N6TNDL0eZLaJI3szyUujJk1M+e6qYvAREOtCMZPKSn3xXbRBYgEgFQ7
T8SWI23KL7udN3y3YmSPhPbMPCmOEr4nOdcXSdrifSf583ki+Z/E11sZRAMG0UsMXg3EJRm5e0houVmD
2aydMp02nBg1IznpzUlPTZcowNf1eS0ozWncnK/rHEJxrRF6RnRpUyAoOQng1gBqVusbWsFCmVHa8y7U
XH/Q5+zI5vpmyUUfdPdsxWVelRS1a50iMY1pZ60rA3WPtK+LSMkMwtf+d2cQhdpPID+IEEVsX4Xr21/C
CmaIMWSYsqO70kACM4pU+yDZAD7nZJdcz13hJAWyJxVrXyFE+yRBhE+K8+SjqAJtWUP+nZQVIuzZkCWH
XZGSlbw8tPQOG906pYjdd0qVOlsNebtni3YlIUlVFdPzDpLUPVuo2OeIwChD99wlKy6LBKf7tuVsum9l
O4jRNxgPpF9hPYz925kPOV4OPOjQE5UtNymGTw6bEuMqPPMP6z47KclvMN61JaRSME/q4VKCoNg8hUN0
mYl4n3sWf5xNhr6gjALPsd7fWqvrW2v17tZar+T+3VTjBDgu2fzizufEPNWNKXDcVDSmbWMEHDcSjZF0
Ny4zICsDacojgaWe0uYpCnq1Qiffni+CJN2Iy6zNkylBz6GAAncNZwwcF2/WN+58jq1WYg5WjmSWkmfF
rO21Rh+/l1wqjmjny8kxUlS004NXF6e2QGjOG0tvunkIlmneHq8oI2KjoowAAPSp3lQW+hTmlatyuujc
dDsz1ut76PalvI/PjHHeGS6ecPL8Br7GjrGKPt77S3/p/enTYG5esm/2DKjtCErU7kL4hvKgnjZL+PAS
oy1iKGYhLni2zhBDIV8zo+SujKX986PnH2+chX98h4KlzRBlQuSN2Nz3Abao0Ph+rvJyO+3a+nOxcvGD
4+LFQozW0qdf/ePq/1v7Hg5MvnY4jMj3EUHwSRUjV6Jdao6FtONFPYiOyysP9bZYqUNoRVCCj4KXODcT
A89X1o2o2GTfSL3MuSeGzu3qNi626Pixlja7chxrNIoih7GyFzYrqfdoGcDqiYHu3sWLg4vNTCQgeAj5
9hxWT4yC6onZtMowM5Y8nr63FL70hFqXVuak3yPWKD+/y8MHLLYZ6o42Rf/eoyJGtI8oE2aCCWVhvNsX
T5L6AbNdSBlkaECA7nDCDFOqRVkZP1FQoIP2KyHweZSs2U1eVoYKMCbl4k2GCp7WzJOE9XAAjCa0dlP9
GDI9kTJGlIoF0jfLp5/FZfg64FOu45yohpTFB8ecSEv2vpBujVnraNeodYnVSKviypBwUuyV2fNaimw2
SlaKuXV5lHwVmfv3Pci9EnMUuZ0JHVzUAp9kIXffKeqsJL1XFZ0VqccotbjWPf84j9l5Zca2hsUwULnM
MthMWdnsya8mTYnToNGKteM+nQpO4+m2o2HWpZRatzpfmJ4+byVcZbMqK8NsXuql2J4COyxohmNkLFY8
4QhAfcyK7sZlQcsM2VmZGjqBB91S3r9ss2tF5kTHuvB/sTLNHngD22nrxEAuol5sAQDr6VRFfRVwR/Tm
3Qne7OLFENbppfaqXphvd+w8VqF0k0U7mcRqFTvRlSRHUxKKm9GabdMAjBACT/cLrfPvT+3lfz9pEUpx
UeAi1cpEy3CB+sqGd7VZPAaP5kANahXBXyBDi7zcIi3eQaIZV5a2sTRgaQ+W9mgOgbZu8LM5TgIWz9oW
p5hRrSQaRTnmNV5BhxDaYhk8jmhruGCI5GiLIUNaXm5xghHpa3v/t/hPcGH/J43tkBaXeQ6L7dAc/TTO
WBLmKUJLyiwrDzx8fCi0v3z6XavzQRdND6wQ2gQefvjo+P7xOl7wz+SHjJLv/4Bx8v3j2uGk1sn/Pl6+
LBpF1fh94+b71Pd/5NiR9GiE0Kx/WMghi3eIArEybfFmDJemuOJTgv3NrVOzyKWtUEuC05AdGVBa3o3Y
Geq3VXAFwGz29WvdcC0a8tmL6DVgbYCTi/MtbQysA1XJzdyZ6R52OENGLdTuQqeObsjPLo1InQW5QLHP
QQUJRb8XzOhIWytHxALTD/CDUexz8+tXLsur2NPL9ym9axe5P2k4MaTmSmkqIVG8D2XW677QONL1+5aE
MD4UuKkFogsBLvEArvk2zR8313emjFWEt0fQsf24unf6V0fdC2EuHXgcYXHtmMGFgfeNgfd35iBkg4tl
AfJ+DOSmAblpQKIREEeB3IyBrJwGZeW8ClNz4RrBZdDv6ilwcydP0uNzDlN+UEfy5hY0ynK+8YTGlS5n
I18zdScAs9vZdPqCgfruQJywLmZvg1hP4Z74A3Cm017LBqxvb4U/PeLmpDdknfsKr6ceyONjL6avCMtj
atfR9Yija+ko+YZz8kbp2zLRt4NkEBEZIqMxnRqpeE+b90i8R/JdXVOpqyBipVb08n3Py5GN+7GLeXl1
Pk++mRfFlSB0z/VRSVZk7pkfHMevRdVv65Yehoj+rdzuM6Rbpy8w2ze3kvVv/ur3eVCDmqb73wAAAP//
7DThxJkgAAA=
`,
	},

	"/js/app.js": {
		local:   "../assets.min/js/app.js",
		size:    11021,
		modtime: 1504452996,
		compressed: `
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
wK18WeGWGN/D0qGZ+/n/BQAA///ftdh4DSsAAA==
`,
	},

	"/project.html": {
		local:   "../assets.min/project.html",
		size:    2153,
		modtime: 1504452996,
		compressed: `
H4sIAAAAAAAA/5yVUW+jOBDH3/spfD6ptA8U3WsP+6RKPe2+rbYfYGVgQmbXGVvG0EYp331lMCRNadLN
S2Jm/vzMjMczeYUdK7VqGmGd+QmlTzuEZ0Z1WhryzmgNbnIFa7M2z0IbVUEl8/U/Mi/kbhf9d6Q20Pd5
VsirHCOWrxRbqVRRrSGtTFtoSB3Wa89lnqHMG6sogHEluFcFE4IlRYu6Srh8CP95FiTyalG5xsYbt520
7Mv4fPIdixY0EiRcfovLk3qgLuHykTp0hjZA/qS6Ae+R6ibh8ikuF/VDjOwvwVqqYIUEFbu+ZrM10cpD
4xMub/7e7QZz399GUp6FzLd6SjGpjpHqUq+Khstc43B+o2+nSo8d3LO3ye25zFWQrR2sRBZPMNufJVZ9
nw3abN5/TPL9hYHs4ziREHHIeX1ls3XPGRdzKtS5gKcaGUM+GW+UHhfT2S3mkvrEHpP2oPjO8kMJfgIN
1L2t07PguVo/QZ+0B3U98mP/aLUeb7bMi9Z7Q8xvLYi4jl9QeGKFp9Q63Ci35Qcf5x3WNbgnr3zbjGE7
Z1zC/mNJeKVSVINL2P34GAlJRGD5S0TCzW0wFbp1k+VBt8FqjTUduNSj1yD4V/LgSGn2BK4Dxx7Ddpy1
WKRRKXa7SBh8fT8jrFYlhCQLDSu/J49qwRMyocPMDmxSY4E+jpIf95Nj3Ztr8Z6ScHncc61W27HLXlPR
2H+/t7Tcut7BnhWGE15glqYOf41Fugw9RXsMhpdSq43yaChkcZgXhzv8r1BDNd/6sapknrVaHgyx0IaH
yQXkWalBuRW+8HcKqwjYeBf4B8NnMFOp2woET+JF+DE479Z+o0MIWYXdH7HncbVMj+7L+cOoWmYDdZdz
90NtGT75j3Y4+P0dAAD//wgBR7ZpCAAA
`,
	},

	"/project_build.html": {
		local:   "../assets.min/project_build.html",
		size:    3805,
		modtime: 1504452996,
		compressed: `
H4sIAAAAAAAA/6xW3XLbNhO991OsmW+sZL5QkmunM1VIdtK001w4acfJfQYkliQqEGCBpW1V1rt3wD+R
omJ7nN5QJISze3YXZxcBFzeQSGZtGFdCcp8jMSEtqMxPtCKjpUQTlkb/hQn94rZEA4zHJBqC+ulzpjI0
noOKNPTQGG3gNIRKcUyFQg5nZ9Cvzmb953vN0S1dLi+9KLBktMqi366v/7heBYv2E7bbeu9uFyy4uBmS
MPoWLG0khgUzmVA+6XJ1vlyWdyMqtZew8zKIIdHSL7h/Ce2LTlOL5F96rVXCO/KZFJlaJagITRSIDpsy
SJmf6DRF7PanWpFvxT9Yk6AoWIgoyC+iTxrqHFsgDRbxFL6wNQKDpCpBp9BYAW2gYJsYgYzIMjTAGhgU
TFVMys3PwSK/iNo07JMxSns4TPv9PfSrs9k4+JIplFA/fY4pqySNNjT/xJpv6iOxPylzS4wq28ZcMs6F
yurUL992X7Em0sXqfFne9WsSU1pdLMu7gxIOP9tC/BgF+Zt+qSgETTOfY7Luj9yImAu2ZNYidyG7Ghxg
DRb6Br8NTpmQ3wTjXZK78/5tuKmUEirr8LZkyu2NheJ+ToXsIE1k8wKtZRnCPVRqrfStcjgHis5UbMu3
L19stw1AVcVu9ypY5G8G+akPaVER8gnVyqLx9YBGFMTRu4pybVbBIo7gUW6s3nyEWv/zVC6JMIlEX/tK
U5IfcHpfe6s5nQQMcoNp6G23bfNpyXytjJwbLCVL8OVsazDdzV7DiK7B9NVu11dmim/6jwfETIYUfo0l
U+uoT/DeDNyDzbUh13VYdNIn6rjZuUSVUd7J7BF70wQeNrZWBhfHNPbEdDOJijNzkOgvosAjpe+qTqJw
55AzwtVss9ls/I8ffc7hw4dVUcyeX/xcVyZzC37OZHrA6dfKMBJaPcCLt1sct/b1+WQsmhs8TMwnzdvE
9GLTHAfV2pdp+hx0VTeSSlGiFAq9IyX96Rk9OEfmyh/92dqFz3WvmTgfdOzm3OgbNKnUt6tccI5qNC36
rpWhnQyN4Z+Ds730XIbdOOtChBb0GhKDjBC0QhAKKEfohPzidNEqZrGXtOC73aKz0kfmtAbE4vlxWbRs
mwiOcJzsdVsNlsgo9JoFoUbBedF2W7/NFSuwv2IcOpZsg+Yhv/Wh7RxjuZ+YjfF2ODzKpsGJZB1alJjQ
Z7f88n9Ccbx7BZWI/VKXrqqh16sGugiOiMT1whbht5eK0JsVurJY32dm+7/rtlqgorCZ3RPZ/F2hrYXX
naBRaK7ziUxp89ShO0E/OHTH036CfXDaJzqrZV+KB7gfzOxO9Y8KvuvU3yHrd4lL60NyjoK4ItIKaFNi
2L63vmJSEJPq/NTvsdTJ2jtSgdSgbWcvXKOpFNQX+2DR2IxO/htHXN8qqRlvPf2OBO8MiZQlZHtfT0ju
8cvh+Q/fkewrne0VM9L8ocmmkUqdjaxhUdKmO0WnjRmps65Jnp3BqYtcqOzKIT9pZwF0RWVFkLCSKoMc
Um2gUThycA1jPiExdLS3OB1mpVAKTX/A64w398arBgV1E4Ernc3nUzd1Ex/0JcdWKOgDAzIsWUO8gaYL
9QoaBB3Vg2HQAp1RZaLttsHA/+G8n6Qnw43d1Xd8A3Vu74EpKxIttfEelePk+W8AAAD//4YQoCLdDgAA
`,
	},

	"/project_env.html": {
		local:   "../assets.min/project_env.html",
		size:    682,
		modtime: 1504452996,
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
		size:    1981,
		modtime: 1504452996,
		compressed: `
H4sIAAAAAAAA/5yVzXLcNgzHXwXezNjJgat1nfSwkdTpR2Z8aNKOJy8AiZDEmgI1JGSvut537+jL9srr
Ns2FQ1IA8cMfoBhrcwdcqtyxeGct+aTx7i/K5doEcb5LB4vcYgjJCi15gWFUGrkkv+qdQ+XukxV57zyc
JdCypsIwaTg/h8fdi4tVGgfxjsv0083NHzfbOJqWsN8PZodDHGlz9zykd/cQpLOU1OhLw0pcs73cbJrd
U+CzkPfov7Sh6yOeZa2xOqwtcSnV6iiB3FlVa/UepokrikCi3q+mIEI7UWhNyducWMinsZl9C4QCVe6K
gmi2LxyLCuZvGpgkjSOTxtVV+sXBSAHiIBCdwVe8JUDI2wZcAeMp4DzU2GUE4k1Zkgcc3aBGbtHa7qc4
qq7SSZUnbQwXho2QGjNPmHbyJ5b09t3yk9ImCHJOyebUJ8ws6WT1TMCHBxhXn1gfa9cgk4VhVJoKbK0M
1ffUEEqyGskNT5kfOY9emdPd0GzD3mC2DoLShknOBrU2XA5F3nycV5kTcfX2ctPsHvcsFbK92jS7RbM8
X041/jGNqw+PW3Vt5GVRK8pvh2RMMSUygyUJXDQYAum+f/vyLnw91e6OXncu0NhXnWmXV/09et3dt8yG
y9k/NMhpjFB5KpJouqrRfj/N1kYfDtFwSLTfj4dxWx8Oo2cfJTOsVSW1nYONmqxrCgFLggdo+ZbdPfcR
h3DnnIXm49s3Rwe+iyOcDeKo+vBM5OES1a2QfpFvG8grd5RLlv7cSuX8No6yFP4TEwfjE5T/myU3Prek
nGInebVg+nWINjK9OYJa8Hgq4AFC5bycRln+0aauvDrV8t8IjpZYo18gfzU1nRBx5hVT98XVKLS96Lqu
U58/K63h+npb1xffL2PlWl/2G6pCWyyYfms9inH8L1x6MunZpun3wwTyd7QU5ovTkzCPHew09Y/NizK9
HI+fovnNefphLkqr+9pe/vDNz0loDDP5eTqS/+6w7wmonSdYr9evwf0TAAD//yCRZ/q9BwAA
`,
	},

	"/project_settings.html": {
		local:   "../assets.min/project_settings.html",
		size:    5417,
		modtime: 1504452996,
		compressed: `
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
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
