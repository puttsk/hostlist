# go-hostlist

Go utility for hostlist expression.

## Installation

```bash
go get github.com/puttsk/hostlist
```

## Hostlist Expression

Hostlist expression is an expression for specifying a group of hostnames with in various applications, including utilities like [pdsh](https://code.google.com/archive/p/pdsh/wikis/HostListExpressions.wiki) and resource managers like [Slurm](https://slurm.schedmd.com/slurm.conf.html#OPT_NodeAddr). This expression allows the concise representation of multiple hostnames. For instance, the expression `host-[001-003]` is a shorthand for `host-001`, `host-002`,and `host-003`.

In the given example, `[001-003]` is referred to as a *range expression*. This expression defines a range of values in the format `[i-j,n-m,..]`, where `i`,`j`,`n`,and `m` are integer with the contraint `i < j`, and `n < m`. Although a range expression can accommodate string values, it does not support character ranges like `[a-c]`.

Example of valid range expression: `[1-10,11,100-101]`, `[a,b,c]`, `[a,22-25]`

## Usage

Hostlist package contains two main functions `Expand` and `Compress`.

`Expand` recieves a hostlist expression and returns a list of hostname contained in the expression.

**Example:**

```go
hosts, err := hostlist.Expand("host-[001-003]")
if err != nil {
    fmt.Print("Error: " + err.Error())
}

// Print host-001 host-002 host-003
fmt.Printf("%s\n", strings.Join(hosts, " ")) 
```

`Compress` recieves a list of hostnames and return a hostlist expression representing the list.

**Example:**

```go
hosts := []string{"host-001","host-oo2","host-003"}
expr, _ := hostlist.Compress(hosts)

// Print host-[001-003]
fmt.Println(expr)
```

## Command Line Interface

```bash
> hostlist -h
Usage of ./hostlist:
  -c    See. -compress
  -compress
        Compress list of hostnames to hostlist expression
  -e    See. -expand
  -expand
        Expand hostlist expression
```

### Expand hostlist expression

```bash
> hostlist -e "host[001-004]"
host001 host002 host003 host004

> hostlist -e "192.168.[0-1].[211-213]"
192.168.0.211 192.168.0.212 192.168.0.213 192.168.1.211 192.168.1.212 192.168.1.213
```

### Compress hostlist expression

```bash
> hostlist -c host001 host002 host003 host004 
host[001-004]

> hostlist -c 192.168.0.211 192.168.0.212 192.168.0.213 192.168.1.211 192.168.1.212 192.168.1.213
192.168.[0-1].[211-213]
```

## Contributing

TBD

## Authors

* **Putt Sakdhnagool** - *Initial work*

See also the list of [contributors](https://github.com/puttsk/hostlist/graphs/contributors) who participated in this project.

## Issues / Feature request

You can submit bug / issues / feature request using [Tracker](https://github.com/puttsk/hostlist/issues).

## Known Issues

* Current implementation of `Expand` does not accept a nested range expression, e.g., `[01-02,a[03-04]]`
* The result of `Expand` an hostlist expression following by `Compress` might not results in the same input hostlist expression.
* In some scenario, `Compress` generates nested range expression, which is not supported by `Expand` function

## License

MIT License.

