<script setup lang="ts">
import { ref, h, type VNodeChild } from "vue";
import type { TestingResult } from "@/stores/result";
import {
    type DataTableColumns,
    NSpin,
    NIcon,
    NTooltip,
    NPopover,
    NA,
    NText,
    NButton,
    NPerformantEllipsis,
    useThemeVars,
    NQrCode
} from "naive-ui";
import { useLocalStorage, useWindowSize } from "@vueuse/core";
import { formatDuration, formatSize } from "@/helpers/format";
import {
    CheckmarkCircle24Filled,
    Warning24Filled,
    ErrorCircle24Filled,
    LockOpen24Regular,
    Info24Regular
} from "@vicons/fluent";
import http from "@/helpers/http";
import { LoadingOne } from "@icon-park/vue-next";
type Item = TestingResult["groups"][0]["items"][0];
const props = defineProps<{
    items: Item[];
    category: string;
}>();
const { width: windowWidth, height: windowHeight } = useWindowSize();
const existErr = (v: Item) => {
    return !!v?.err || !!v.result?.err;
};
const existRedirect = (v: Item) => {
    return (v.result?.records.length || 0) > 1;
};
const getLastResp = (v: Item) => {
    if (v?.result?.records) {
        return v.result.records[v.result.records.length - 1].resp;
    }
    return null;
};

const themeVars = useThemeVars();

const columns: DataTableColumns<Item> = [
    {
        title: "序列",
        key: "index",
        fixed: "left",
        width: 60,
        render(row) {
            return props.items.indexOf(row) + 1;
        }
    },
    {
        title: "名称",
        key: "name",
        fixed: windowWidth.value >= 720 ? "left" : undefined,
        width: 128,
        render(row) {
            let icon: any = CheckmarkCircle24Filled,
                type = "success";
            if (row.status === "pending") {
                icon = LoadingOne;
                type = "info";
            } else if (existErr(row)) {
                icon = ErrorCircle24Filled;
                type = "error";
            } else if (existRedirect(row) || getLastResp(row)?.status_code !== 200) {
                icon = Warning24Filled;
                type = "warning";
            }

            return h(NText, { type, class: "icon-with-text", style: "--size:4px" }, () => [
                h(NIcon, { size: type === "info" ? "large" : 22.5 }, () => h(icon)),
                h(
                    NPerformantEllipsis,
                    { style: "max-width: 100px;color:var(--n-td-text-color)" },
                    () => row.name
                )
            ]);
        }
    },
    {
        title: "网址",
        key: "url",
        width: 168,
        render(row) {
            const tlsEnabled = row.url.startsWith("https://");
            return h(
                NText,
                {
                    class: tlsEnabled ? "" : "icon-with-text",
                    style: "--size:3px;max-width: 120px",
                    type:
                        row.status === "pending"
                            ? "info"
                            : existErr(row)
                              ? "error"
                              : existRedirect(row) || getLastResp(row)?.status_code !== 200
                                ? "warning"
                                : "info"
                },
                () => [
                    tlsEnabled
                        ? ""
                        : h(
                              NTooltip,
                              {
                                  style: "max-width: 190px",
                                  arrowPointToCenter: true
                              },
                              {
                                  default: () => "网站未使用TLS加密，连接存在安全风险",
                                  trigger: () =>
                                      h(
                                          NIcon,
                                          { color: themeVars.value.warningColor, size: "large" },
                                          () => h(LockOpen24Regular)
                                      )
                              }
                          ),
                    h(
                        NPopover,
                        { trigger: "hover", placement: "top", style: "padding:3px" },
                        {
                            default: () =>
                                h(NQrCode, {
                                    padding: 6,
                                    value: row.url,
                                    style: "box-sizing: content-box;"
                                }),
                            trigger: () =>
                                h(
                                    NA,
                                    {
                                        style: "transition:0s;text-decoration-color:inherit;color:inherit",
                                        href: row.url,
                                        target: "_blank"
                                    },
                                    () => new URL(row.url).hostname
                                )
                        }
                    )
                ]
            );
        }
    },
    {
        title: "用时",
        key: "result.total_duration",
        width: 90,
        render(row) {
            if (row.status === "pending") {
                return h(NSpin, { size: 19, strokeWidth: 20 });
            }
            if (!row.result) {
                return;
            }
            const formatted = formatDuration(row.result.total_duration);
            return h(
                NText,
                {
                    type:
                        row.result.total_duration <= 1000
                            ? "success"
                            : row.result.total_duration <= 2000
                              ? "info"
                              : "warning"
                },
                () => formatted
            );
        }
    },
    {
        title: "网页标题",
        key: "result.title",
        width: 128,
        render(row) {
            if (!row.result) {
                return;
            }
            const title = row.result.title;
            return h(
                NPerformantEllipsis,
                {
                    style: "max-width:100px",
                    lineClamp: 2,
                    tooltip: {
                        // @ts-ignore
                        style: "max-width: 200px"
                    }
                },
                () => title
            );
        }
    },
    {
        title: "响应链",
        key: "result.records",
        width: 158,
        render(row) {
            const records = row.result?.records;
            if (!records) {
                return;
            }
            const vnodes: VNodeChild[] = [];

            for (let i = 0; i < records.length; i++) {
                const r = records[i];
                // h(NText, { type: "info", class: "icon-with-text", style: "--size:2px" }, () => [h(NIcon, { size: 18 }, () => h(Time)), `${formatDuration(r.duration)}`]),
                vnodes.push(
                    h(NText, { type: "info" }, () =>
                        h(
                            NA,
                            {
                                style: "transition:0s;text-decoration-color:inherit;color:inherit",
                                target: "_blank",
                                href: r.url
                            },
                            () => new URL(r.url).host
                        )
                    )
                );
                const code = r.resp?.status_code;
                vnodes.push(
                    "[",
                    h(
                        NText,
                        {
                            type: code
                                ? code === 200
                                    ? "success"
                                    : code >= 500
                                      ? "error"
                                      : "warning"
                                : "error"
                        },
                        () => (code ? code : "失败")
                    ),
                    "]"
                );
                if (i !== records.length - 1) {
                    vnodes.push(" => ");
                }
            }
            return h(
                NPopover,
                {},
                {
                    default: () =>
                        h(
                            NButton,
                            {
                                onClick() {
                                    chainModal.value.name = row.name;
                                    chainModal.value.records = row.result?.records || [];
                                    chainModal.value.show = true;
                                },
                                size: "small",
                                strong: true,
                                secondary: true,
                                type: "info"
                            },
                            {
                                default: () => "查看详情",
                                icon: () => h(NIcon, { size: "large" }, () => h(Info24Regular))
                            }
                        ),
                    trigger: () => h("div", {}, vnodes)
                }
            );
        }
    },
    {
        title: "响应内容",
        key: "result.size",
        width: 120,
        render(row) {
            if (!row.result || row.result.size === 0) {
                return;
            }
            const formatted = formatSize(row.result.size, " ");
            return h(
                NButton,
                {
                    onClick() {
                        loadContent(row);
                    },
                    size: "small",
                    tertiary: true,
                    strong: true
                },
                () => formatted
            );
        }
    },
    {
        title: "错误描述",
        key: "error",
        render(row) {
            let reason: VNodeChild = [];
            let advice = "";
            let type = "error";
            if (row.err) {
                reason.push("构建请求时出错：" + row.err);
            } else if (row.result?.err) {
                if (row.result.err.indexOf("closed by the remote host") > -1) {
                    advice = "尝试使用网络代理或VPN进行访问";
                } else if (row.result.err.indexOf("Client.Timeout") > -1) {
                    advice = "尝试配置更久的访问超时时长 ";
                }
                reason.push(
                    (row.result.records.length > 1
                        ? "在请求重定向后的网址时出错："
                        : "请求失败：") + row.result.err
                );
            } else if (row.result?.last_resp_redirect) {
                type = "warning";
                const url = row.result.last_resp_redirect;
                reason.push(
                    "超出测试允许的最大重定向次数，最后一次重定向是",
                    h(NText, { type: "info" }, () =>
                        h(
                            NA,
                            {
                                target: "_blank",
                                href: url
                            },
                            () => new URL(url).host
                        )
                    )
                );
            } else if (existRedirect(row)) {
                type = "warning";
                reason.push(`存在${row.result!.records.length - 1}次重定向`);
            } else {
                const resp = getLastResp(row);
                if (resp && resp.status_code !== 200) {
                    type = "warning";
                    reason.push("最后一次响应状态码异常：" + resp.status);
                }
            }
            if (reason.length > 0) {
                return [h(NText, { type }, () => reason)].concat(
                    advice.length > 0
                        ? [
                              h("br"),
                              h(
                                  NText,
                                  { italic: true, depth: 2 },
                                  () => "解决方案（仅供参考）：" + advice
                              )
                          ]
                        : []
                );
            }
        }
    }
];
const pageSize = useLocalStorage("pageSize", 10);

const chainModal = ref({
    show: false,
    name: "",
    records: [] as NonNullable<Item["result"]>["records"]
});

const contentModal = ref({
    timestamp: 0,
    ctrl: new AbortController(),
    loading: false,
    show: false,
    size: 0,
    content: "",
    name: "",
    err: ""
});
function loadContent(row: Item) {
    const ts = Date.now();
    contentModal.value.ctrl?.abort();
    const ctrl = new AbortController();
    contentModal.value = {
        timestamp: ts,
        ctrl,
        loading: true,
        show: true,
        size: row.result?.size || 0,
        content: "",
        err: "",
        name: row.name
    };

    return http<string>(`/api/content?category=${props.category}&name=${row.name}`, {
        signal: ctrl.signal,
        method: "GET"
    })
        .then((data) => {
            if (ts !== contentModal.value.timestamp) {
                return;
            }
            contentModal.value.content = data.data;
        })
        .catch((err) => {
            if (http.isAborted(err) || ts !== contentModal.value.timestamp) {
                return;
            }
            const cause = http.checkError(err);
            let reason = err?.message || String(err);
            if (cause && cause.parseJSONError) {
                reason =
                    `响应状态: ${cause.originalResponse.statusText} 内容格式: ${cause.originalResponse.headers.get("Content-Type")} 解析JSON错误:` +
                    cause.parseJSONError.message;
            }
            contentModal.value.err = reason;
        })
        .finally(() => {
            if (ts !== contentModal.value.timestamp) {
                return;
            }
            contentModal.value.loading = false;
        });
}
</script>

<template>
    <n-data-table
        :scroll-x="1080"
        :columns="columns"
        :data="props.items"
        @update-page-size="(v) => (pageSize = v)"
        :pagination="{
            pageSizes: [10, 15, 20, 30],
            pageSize: pageSize,
            simple: windowWidth < 600,
            pageSlot: 6,
            showSizePicker: true,
            showQuickJumper: true
        }"
        :bordered="false"
    />
    <n-modal
        v-model:show="chainModal.show"
        display-directive="if"
        style="width: 80vw; max-width: 800px; min-width: 300px"
        preset="card"
        :title="`${chainModal.name}响应链（${chainModal.records.length}）`"
        size="huge"
        :bordered="false"
    >
        <n-scrollbar :style="`max-height: ${Math.min(windowHeight * 0.6, 700)}px`" trigger="none">
            <n-timeline :icon-size="20">
                <n-timeline-item
                    v-for="(r, i) in chainModal.records"
                    :key="i + r.url"
                    :type="r.resp ? 'success' : 'error'"
                >
                    <template #icon>
                        <n-icon>
                            <CheckmarkCircle24Filled v-if="r.resp" />
                            <ErrorCircle24Filled v-else />
                        </n-icon>
                    </template>
                    <template #header>
                        <n-button
                            icon-placement="right"
                            text
                            tag="a"
                            :type="r.resp ? 'success' : 'error'"
                            :href="r.url"
                            target="_blank"
                            >{{ r.url }}
                        </n-button>
                    </template>
                    <n-space vertical>
                        <n-space>
                            <template v-if="r.remote_addr">
                                <span>
                                    协议: <n-text type="info">{{ r.remote_addr.type }}</n-text>
                                </span>
                                <span>
                                    IP: <n-text type="info">{{ r.remote_addr.ip }}</n-text>
                                </span>
                                <span>
                                    端口: <n-text type="info">{{ r.remote_addr.port }}</n-text>
                                </span>
                            </template>
                            <n-text v-else type="warning"> 未连接成功，远程地址未知 </n-text>
                        </n-space>
                        <n-space>
                            <template v-if="r.resp">
                                <span>
                                    响应状态:
                                    <n-text
                                        :type="
                                            r.resp.status_code === 200
                                                ? 'success'
                                                : r.resp.status_code >= 500
                                                  ? 'error'
                                                  : 'warning'
                                        "
                                        >{{ r.resp.status }}</n-text
                                    >
                                </span>
                                <span>
                                    内容格式:
                                    <n-text :type="r.resp.content_type ? 'info' : 'warning'">{{
                                        r.resp.content_type || "未知"
                                    }}</n-text>
                                </span>
                                <span>
                                    响应大小:
                                    <n-text type="info">{{
                                        r.resp.content_length !== -1
                                            ? formatSize(r.resp.content_length)
                                            : "未知"
                                    }}</n-text>
                                </span>
                            </template>
                            <n-text v-else type="error"> 未获取响应 </n-text>
                        </n-space>
                        <span>
                            用时: <n-text type="info">{{ formatDuration(r.duration) }}</n-text>
                        </span>
                    </n-space>
                </n-timeline-item>
            </n-timeline>
        </n-scrollbar>
    </n-modal>

    <n-modal
        v-model:show="contentModal.show"
        display-directive="if"
        :on-after-leave="() => contentModal.ctrl?.abort()"
        style="width: 80vw; max-width: 800px; min-width: 300px"
        preset="card"
        :title="`${contentModal.name}响应内容${contentModal.size > 0 ? '（' + formatSize(contentModal.size) + '）' : ''}`"
        size="huge"
        :bordered="false"
    >
        <n-skeleton v-if="contentModal.loading" text :repeat="6" />

        <n-scrollbar
            v-else-if="contentModal.err"
            :style="`max-height: ${Math.min(windowHeight * 0.6, 700)}px`"
            trigger="none"
        >
            <n-result status="error">
                <n-text strong style="text-align: center; font-size: 16px" tag="div">
                    {{ contentModal.err }}
                </n-text>
            </n-result>
        </n-scrollbar>
        <n-input
            v-else
            :input-props="{ spellcheck: 'false' }"
            type="textarea"
            size="small"
            :style="`height:${Math.min(windowHeight * 0.6, 700)}px`"
            v-model:value="contentModal.content"
            round
        />
    </n-modal>
</template>

<style scoped lang="less">
:deep(.n-data-table__pagination) {
    justify-content: flex-start;
}
</style>
