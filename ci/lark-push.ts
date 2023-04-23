// import {Sha256} from "https://deno.land/std@0.159.0/hash/sha256.ts"

const larkPushKey = Deno.env.get("LARK_PUSH_KEY")
const qCloudUrl = (fn) => `https://${Deno.env.get('QCLOUD_COS_DOMAIN')}/client/${Deno.env.get('CI_COMMIT_REF_NAME')}/${fn}`
const main = async () => {

    const r = await fetch(
        "https://open.feishu.cn/open-apis/bot/v2/hook/" + larkPushKey,
        {
            method: "POST",
            headers: {"content-type": "application/json"},
            body: JSON.stringify({msg_type: "interactive", card: data}),
        }
    );
}

const sign = Deno.env.get('CI_COMMIT_TAG') !== undefined ? '🔖' : '✔'

const data = {
    "config": {
        "wide_screen_mode": true
    },
    "elements": [
        {
            "tag": "markdown",
            "content": `**📦Bundle：** ${qCloudUrl('pro.zip')}`
        }
    ],
    "header": {
        "template": "green",
        "title": {
            "content": `${sign}Client - ${Deno.env.get('CI_COMMIT_REF_NAME')}`,
            "tag": "plain_text"
        }
    }
}

await main()