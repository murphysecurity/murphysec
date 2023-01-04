import { Sha256 } from "https://deno.land/std@0.159.0/hash/sha256.ts"

const gitTag = Deno.env.get("CI_BUILD_REF_NAME")
const larkPushKey = Deno.env.get("LARK_PUSH_KEY")
const qCloudUrl = (fn) => `https://${Deno.env.get('QCLOUD_COS_DOMAIN')}/client/${gitTag}/${fn}`

const contentText = [
    ['GitTag', gitTag],
    ['pro.zip', qCloudUrl('pro.zip')],
    ['SHA-256', new Sha256().update(await Deno.readFileSync('out/zip/pro.zip')).hex()],
].filter(it => it[1]).map(it => `**${it[0]}: **${it[1]}`).join('\n')

const cardContent = {
    "config": {
        "wide_screen_mode": true
    },
    "header": {
        "template": "orange",
        "title": {
            "content": `上传推送: Client(${gitTag})`,
            "tag": "plain_text"
        }
    },
    "elements": [
        {
            "tag": "div",
            "text": {
                "tag": "lark_md",
                "content": contentText,
            }
        }
    ]
}
const r = await fetch(
    "https://open.feishu.cn/open-apis/bot/v2/hook/" + larkPushKey,
    {
        method: "POST",
        headers: {"content-type": "application/json"},
        body: JSON.stringify({msg_type: "interactive", card: cardContent}),
    }
);
console.log('pushed!', await r.text())
