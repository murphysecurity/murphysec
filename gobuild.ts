const targets: [string, string][] = [
    ["linux", "amd64"],
    ["darwin", "amd64"],
    ["windows", "amd64"],
];

const isSaaS = Deno.args.indexOf("saas") > -1;
const bn = isSaaS ? "murphysec-saas" : "murphysec";

const tags = [!isSaaS ? "pro" : ""].filter((it) => it);
const opts = targets.map((it) => ({
    cmd: [
        "go",
        "build",
        "-trimpath",
        "-ldflags",
        "-s -w",
        tags.length > 0 ? ["-tags", ...tags] : [],
        "-o",
        `out/${bn}-${it[0]}-${it[1]}${it[0] === "windows" ? ".exe" : ""}`,
        ".",
    ].flat().filter((it) => it),
    env: {GOOS: it[0], GOARCH: it[1]},
    stdin: "null" as "null",
}));
console.log(opts);
const process = opts.map((it) => Deno.run(it).status());
const status = await Promise.all(process);
if (status.every((it) => it.success && it.code === 0)) {
    console.log("编译成功");
} else {
    console.error("编译失败", status);
    Deno.exit(-1);
}
