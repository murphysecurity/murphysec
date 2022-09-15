const t = (n: string) => `
package module

import "github.com/murphysecurity/murphysec/module/${n}"

func init() {
	Inspectors = append(Inspectors, ${n}.Instance)
}
`

Array.from(Deno.readDirSync('.'))
    .filter((i: any) => i.isDirectory && ['base'].indexOf(i.name) === -1)
    .map((i: any) => i.name)
    .forEach((i: string) => {
        Deno.writeTextFileSync(`${i}_init.go`, t(i))
    })
