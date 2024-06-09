<script setup lang="ts">
import { NEllipsis, NText, useDialog, useMessage, useThemeVars } from "naive-ui";
import { ref, onMounted, h, type FunctionalComponent } from "vue";
import useResultStore, { type TestingResult } from "@/stores/result";
import { Refresh } from "@icon-park/vue-next";
import { useTimestamp, useDateFormat, useLocalStorage } from "@vueuse/core";
import { computed } from "vue";
import { formatDuration } from "@/helpers/format";
import useConfStore from "@/stores/conf";
import { handleErrWithDialog } from "@/helpers/handleErr";

const { result, loadResult } = useResultStore();

const load = () =>
    loadResult().catch((error) => handleErrWithDialog(dialog, error, { title: "获取结果出错" }));
onMounted(() => load());

const message = useMessage();
const dialog = useDialog();
const confStore = useConfStore();

const nowTimestamp = useTimestamp();
const now = useDateFormat(nowTimestamp, "YYYY-MM-DD HH:mm:ss");
const testingTimestamp = computed(() => {
    return result.value.testingResult?.start || 0;
});
const testDate = useDateFormat(testingTimestamp, "YYYY-MM-DD HH:mm:ss");
const formatDurationOption = formatDuration.optionZh({
    intSeconds: true
});

const themeVars = useThemeVars();

function applyConf() {
    result.value.testingResult?.conf && (confStore.conf.value = result.value.testingResult?.conf);
    message.success("应用成功");
}
const showConf = useLocalStorage("showConf", false);
const CodeEllipsis: FunctionalComponent = (_, { slots }) => {
    return h(
        NText,
        {
            code: true,
            style: "font-size:inherit"
        },
        () =>
            h(
                NEllipsis,
                {
                    style: "max-width:210px",
                    tooltip: {
                        // @ts-ignore
                        style: "max-width: 280px",
                        arrowPointToCenter: true
                    }
                },
                () => slots.default?.()
            )
    );
};

const category = useLocalStorage<"animation" | "video">("category", "animation");
const keyword = ref("");

enum CheckOption {
    AccessibleWebsite = 0b1,
    SortByAccessSpeed = 0b10,
    RedirectedWebsite = 0b100
}
const checked = useLocalStorage<CheckOption[]>("checked", [CheckOption.SortByAccessSpeed]);

const getItemsByCategory = (category: "animation" | "video") => {
    return result.value?.testingResult?.groups?.find((r) => r.category === category)?.items || [];
};

const filteredItems = computed(() => {
    let items = getItemsByCategory(category.value);
    if (items.length === 0) {
        return [];
    }
    items = [...items];
    const kw = keyword.value.trim().toLowerCase();
    if (kw !== "") {
        items = items.filter((v) => v.name.toLowerCase().indexOf(kw) > -1);
    }
    const flag = checked.value.reduce((a, b) => a | b, 0);
    if (flag & CheckOption.RedirectedWebsite) {
        items = items.filter(
            (v) => (v.result?.records.length || 0) > 1 || !!v.result?.last_resp_redirect
        );
    }

    if (flag & CheckOption.AccessibleWebsite) {
        items = items.filter((v) => !v.err && !v.result?.err);
    }
    if (flag & CheckOption.SortByAccessSpeed) {
        const existErr = (v: TestingResult["groups"][0]["items"][0]) => {
            return !v.result || !!v.result.err;
        };
        items.sort((a, b) => {
            if (existErr(a) && existErr(b)) {
                return 0;
            }
            if (existErr(a) && !existErr(b)) {
                return 1;
            }
            if (existErr(b) && !existErr(a)) {
                return -1;
            }
            return a.result!.total_duration - b.result!.total_duration;
        });
    }
    return items;
});

const options = computed(() => {
    const arr = [
        {
            label: `动漫网站`,
            value: "animation"
        },
        { label: "影视网站", value: "video" }
    ];
    for (const v of arr) {
        if (v.value === category.value) {
            v.label += `（${filteredItems.value.length}/${getItemsByCategory(v.value).length}）`;
        } else {
            v.label += `（${getItemsByCategory(v.value as any).length}）`;
        }
    }
    return arr;
});
</script>

<template>
    <n-h6>
        <span class="icon-with-text">
            测试结果
            <n-tooltip trigger="hover">
                <template #trigger>
                    <n-icon
                        id="refreshResult"
                        :class="{ rotate: result.loading }"
                        :style="{ '--color': themeVars.primaryColor }"
                        @click="!result.loading && load()"
                    >
                        <Refresh />
                    </n-icon>
                </template>
                {{ result.loading ? "正在刷新中..." : "点击刷新" }}
            </n-tooltip>
        </span>
    </n-h6>
    <Transition name="fade-in" mode="out-in">
        <n-space vertical size="large" v-if="result.testingResult">
            <transition name="fade-in" mode="out-in">
                <n-space v-if="result.testingResult?.conf" vertical>
                    <!-- <n-switch v-model:value="showConf">
                        <template #checked> 显示测试配置 </template>
                        <template #unchecked> 隐藏测试配置 </template>
                    </n-switch> -->
                    <n-collapse-transition :show="showConf || true">
                        <n-space>
                            <n-tag>
                                DNS服务器
                                <CodeEllipsis>
                                    {{ result.testingResult.conf.dns_server || "跟随系统" }}
                                </CodeEllipsis>
                            </n-tag>
                            <n-tag>
                                网络代理
                                <CodeEllipsis>
                                    {{ result.testingResult.conf.proxy_url || "无" }}
                                </CodeEllipsis>
                            </n-tag>
                            <n-tag>
                                访问超时
                                <CodeEllipsis>
                                    {{
                                        result.testingResult.conf.timeout_seconds === 0
                                            ? "无"
                                            : formatDuration(
                                                  result.testingResult.conf.timeout_seconds * 1000,
                                                  formatDurationOption
                                              )
                                    }}
                                </CodeEllipsis>
                            </n-tag>
                            <n-tag>
                                用户代理
                                <CodeEllipsis>
                                    {{ result.testingResult.conf.user_agent || "无" }}
                                </CodeEllipsis>
                            </n-tag>
                            <n-button @click="applyConf" size="small" type="info" strong secondary
                                >应用此配置到表单</n-button
                            >
                        </n-space>
                    </n-collapse-transition>
                </n-space>
            </transition>
            <n-space>
                <n-tag round :bordered="false" type="info" strong> 当前时间: {{ now }} </n-tag>
                <template v-if="result.testingResult?.start">
                    <n-tag round :bordered="false" type="success" strong>
                        测试时间: {{ testDate }}
                    </n-tag>
                    <n-tag round :bordered="false" type="warning" strong>
                        距离上次测试已过去{{
                            formatDuration(nowTimestamp - testingTimestamp, formatDurationOption)
                        }}
                    </n-tag>
                </template>
            </n-space>
            <n-result
                v-if="result.testingResult.err"
                status="error"
                style="max-width: 500px; margin: auto"
            >
                <n-text strong style="text-align: center; font-size: 16px" tag="div">
                    测试出现致命错误：<n-text type="error">
                        {{ result.testingResult.err }}
                    </n-text>
                </n-text>
            </n-result>
            <template v-else>
                <n-space>
                    <n-select v-model:value="category" style="width: 210px" :options="options" />
                    <n-input
                        v-model:value="keyword"
                        clearable
                        type="text"
                        placeholder="输入关键字筛选"
                    />
                </n-space>
                <n-checkbox-group v-model:value="checked" size="large">
                    <n-space item-style="display: flex;">
                        <n-checkbox
                            :value="CheckOption.SortByAccessSpeed"
                            label="依据访问速度排序"
                        />
                        <n-checkbox
                            :value="CheckOption.AccessibleWebsite"
                            label="仅查看可访问的网站"
                        />

                        <n-checkbox
                            :value="CheckOption.RedirectedWebsite"
                            label="仅查看发生重定向的网站"
                        />
                    </n-space>
                </n-checkbox-group>
                <AppTable :category="category" :items="filteredItems" />
            </template>
        </n-space>
        <n-empty description="您还未进行过测试" v-else> </n-empty>
    </Transition>
</template>

<style lang="less">
#refreshResult {
    cursor: pointer;

    &:hover {
        color: var(--color);
    }

    &.rotate {
        color: var(--color);
        animation: rotate 1s linear infinite;

        @keyframes rotate {
            0% {
                transform: rotate(0deg);
            }

            100% {
                transform: rotate(180deg);
            }
        }
    }
}
</style>
