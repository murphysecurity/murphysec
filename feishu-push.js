const gitTag = Deno.env.get("CI_BUILD_REF_NAME")
const larkPushKey = Deno.env.get("LARK_PUSH_KEY")
const qCloudUrl = (fn) => `https://${Deno.env.get('QCLOUD_COS_DOMAIN')}/client/${gitTag}/${fn}`

const contentText = [
    ['GitTag', gitTag],
    ['Windows', qCloudUrl('murphysec-windows-amd64.exe')],
    ['Linux', qCloudUrl('murphysec-linux-amd64')],
    ['Apple', qCloudUrl('murphysec-darwin-amd64')],
    ['SaaS-Windows', qCloudUrl('murphysec-saas-windows-amd64.exe')],
    ['SaaS-Linux', qCloudUrl('murphysec-saas-linux-amd64')],
    ['SaaS-Apple', qCloudUrl('murphysec-saas-darwin-amd64')],
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
