const p = [
    Deno.run({
        cmd: ['go', 'build', '-o', 'out/murphysec-saas-linux-amd64', '.'],
        env: {GOOS: 'linux', GOARCH: 'amd64'},
        stdin: 'null'
    }),
    Deno.run({
        cmd: ['go', 'build', '-o', 'out/murphysec-saas-windows-amd64.exe', '.'],
        env: {GOOS: 'windows', GOARCH: 'amd64'},
        stdin: 'null'
    }),
    Deno.run({
        cmd: ['go', 'build', '-o', 'out/murphysec-saas-darwin-amd64', '.'],
        env: {GOOS: 'darwin', GOARCH: 'amd64'},
        stdin: 'null'
    }),
]

console.log(await Promise.all(p.map(it=>it.status())))
