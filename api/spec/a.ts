import { gzip } from "https://deno.land/x/compress@v0.4.5/mod.ts"

const last = <T>(arr: T[], n: number) => {
    let i = arr.length - n
    if (i < 0 || n >= arr.length) {
        return undefined
    }
    return arr[i]
}
const f = (node: any, path: string[]): any => {
    if (node  === undefined|| node === null) return node
    if (Array.isArray(node)) {
        return node.map(i => f(i, path))
    }
    if (typeof node === 'object') {
        let t: any = {}
        Object.keys(node).forEach(i => t[i] = f(node[i], [...path, i]))
        node = t;
    }
    if (typeof node === "object" && last(path, 1) !== 'properties' && path.length>1) {
        if(node.description!==undefined) node.description = ''
        delete node.title
        delete node.summary
        delete node.example
        delete node.examples
        delete node.tag
        delete node.tags
        Object.keys(node).filter(i=>i.startsWith('x-')).forEach(i=>delete node[i])
        return node
    }
    return node
}

const root = JSON.parse(Deno.readTextFileSync('SAAS3.0.openapi.json'));
const target = f(root, []);
const targetData = gzip(new TextEncoder().encode(JSON.stringify(target)))
Deno.writeFileSync("api.json.gz", targetData);