const gitTag = Deno.env.get("CI_BUILD_REF_NAME")
const fileName = Deno.env.get("UPLOAD_FILENAME")
const larkPushKey = Deno.env.get("LARK_PUSH_KEY")

const contentText = [
    ['File', fileName],
    ['GitTag', gitTag],
    ['URL', `https://${Deno.env.get('QCLOUD_COS_DOMAIN')}/client/${gitTag}/${fileName}`],
].filter(it => it[1]).map(it => `**${it[0]}: **${it[1]}`).join('\n')

const cardContent = {
    "config": {
        "wide_screen_mode": true
    },
    "header": {
        "template": "orange",
        "title": {
            "content": `上传推送: ${fileName}(${gitTag})`,
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
