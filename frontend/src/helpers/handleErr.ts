import { type DialogApi, type DialogOptions, NText } from "naive-ui";
import { type VNodeChild, h } from "vue";
import http from "./http";
export function handleErrWithDialog(d: DialogApi, err: Error, options?: Partial<DialogOptions>) {
    if (http.isAborted(err)) {
        return;
    }
    d.error(
        Object.assign(
            {
                closable: false,
                closeOnEsc: false,
                maskClosable: false,
                positiveText: "刷新页面",
                negativeText: "忽略",
                content: generateContent(err),
                onPositiveClick: () => {
                    location.reload();
                }
            } as DialogOptions,
            options || {}
        )
    );
}

function generateContent(err: Error) {
    const cause = http.checkError(err);
    const vnodes: VNodeChild[] = [];
    if (cause && cause.parseJSONError) {
        const message = cause.parseJSONError.message;
        vnodes.push(
            "响应状态：",
            h(
                NText,
                { type: "warning" },
                () => cause.originalResponse.status + " " + cause.originalResponse.statusText
            ),
            h("br"),
            "内容格式：",
            h(NText, { type: "info" }, () => cause.originalResponse.headers.get("Content-Type")),
            h("br"),
            "解析JSON出错：",
            h(NText, { type: "error" }, () => message)
        );
    } else {
        vnodes.push(h(NText, { type: "error" }, () => err?.message || String(err)));
    }
    return () => vnodes;
}
