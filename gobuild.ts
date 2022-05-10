const p = [
    Deno.run({
        cmd: ['go', 'build', '-tags', 'pro', '-o', 'out/murphysec-linux-amd64', '.'],
        env: {GOOS: 'linux', GOARCH: 'amd64'},
        stdin: 'null'
    }),
    Deno.run({
        cmd: ['go', 'build', '-tags', 'pro', '-o', 'out/murphysec-windows-amd64.exe', '.'],
        env: {GOOS: 'windows', GOARCH: 'amd64'},
        stdin: 'null'
    }),
    Deno.run({
        cmd: ['go', 'build', '-tags', 'pro', '-o', 'out/murphysec-darwin-amd64', '.'],
        env: {GOOS: 'darwin', GOARCH: 'amd64'},
        stdin: 'null'
    }),
]

console.log(await Promise.all(p.map(it => it.status())))
