// github issue: https://github.com/AMuzykus/risor/issues/6
// expected value: 11
// expected type: int

s := "\ntest\t\"str\\"

raw := `
test	"str\`

assert(s == raw)

len(s)
