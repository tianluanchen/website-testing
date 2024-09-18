<script setup lang="ts">
import { watchEffect, ref, h, type FunctionalComponent } from "vue";
import { NButton, NText, type ProgressProps } from "naive-ui";
import { formatDuration } from "@/helpers/format";
import useCounter from "@/hooks/counter";
import { useTitle, useWebSocket, useThrottleFn } from "@vueuse/core";
import useResultStore from "@/stores/result";
import { ExperimentOne } from "@icon-park/vue-next";
// import http from "@/helpers/http";
// import { useIntervalFn } from "@vueuse/core";
import { useDialog } from "naive-ui";
import { handleErrWithDialog } from "@/helpers/handleErr";
const resultStore = useResultStore();
const loadResult = useThrottleFn(() => resultStore.loadResult(), 1250, true);
const counter = useCounter();
const processStatus = ref<NonNullable<ProgressProps["status"]>>("info");

const content = ref<FunctionalComponent>(() => h(NText, { type: "info" }, () => "构建连接中..."));
counter.increaseTo(10);
const dialog = useDialog();
const testing = ref(false);
const title = useTitle();
watchEffect(() => {
    title.value = testing.value ? "正在测试中..." : "网站本地测试";
});

// const { pause: pauseRefreshState, resume: resumeRefreshState } = useIntervalFn(() => {
//     pauseRefreshState();
//     http<boolean>("/api/state", {
//         method: "GET"
//     })
//         .then((data) => {
//             testing.value = data.data;
//             resumeRefreshState();
//         })
//         .catch((error) =>
//             handleErrWithDialog(dialog, error, {
//                 title: "查询状态出错",
//                 onNegativeClick: () => resumeRefreshState()
//             })
//         );
// }, 5000);

useWebSocket(`ws${location.protocol === "https:" ? "s" : ""}://${location.host}/api/watch`, {
    immediate: true,
    heartbeat: {
        message: "ping",
        interval: 3000,
        pongTimeout: 3000
    },
    autoReconnect: {
        delay: 3000
    },
    onDisconnected() {
        counter.set(0);
        counter.increaseTo(10);
        processStatus.value = "warning";
        content.value = () => h(NText, { type: "warning" }, () => "连接已断开，尝试重连...");
    },
    onConnected() {
        counter.set(100);
        processStatus.value = "info";
        content.value = () => h(NText, { type: "info" }, () => "已成功连接，接收消息中...");
    },
    onError() {
        counter.set(0);
        counter.increaseTo(10);
        processStatus.value = "error";
        content.value = () => h(NText, { type: "warning" }, () => "连接出错，尝试重连...");
    },
    onMessage(_, evt) {
        if (evt.data === "pong") {
            return;
        }
        const { event, data }: { event: string; data: any } = JSON.parse(evt.data);
        loadResult().catch((err) => handleErrWithDialog(dialog, err, { title: "查询结果出错" }));
        switch (event) {
            case "State":
                if (!data) {
                    testing.value = false;
                    content.value = () => h(NText, { type: "info" }, () => "暂无进行中的测试");
                    return;
                }
                content.value = () => h(NText, { type: "info" }, () => "正在测试中...");
                break;
            case "Finish": {
                counter.set(100);
                const { err, duration } = data;
                if (err) {
                    processStatus.value = "error";
                    content.value = () =>
                        h(
                            NText,
                            { type: "error" },
                            () => `测试出错: ${err}，用时${formatDuration(duration)}`
                        );
                } else {
                    processStatus.value = "success";
                    content.value = () =>
                        h(
                            NText,
                            { type: "success" },
                            () => `已完成测试，总用时${formatDuration(duration)}`
                        );
                }
                testing.value = false;
                return;
            }
            case "Start":
                counter.set(0);
                counter.increaseTo(10);
                content.value = () => "开始测试中...";
                processStatus.value = "info";
                break;
            case "PickFastestAPI":
                counter.set(10);
                counter.increaseTo(20);
                content.value = () => [
                    "已优选API ",
                    h(
                        NButton,
                        {
                            style: "font-size:inherit",
                            tag: "a",
                            text: true,
                            href: data.api,
                            target: "_blank",
                            type: "primary"
                        },
                        () => data.api
                    ),
                    "，用时",
                    h(NText, { type: "info" }, () => formatDuration(data.duration))
                ];
                break;
            case "FetchWebsites":
                counter.set(20);
                counter.increaseTo(30);
                content.value = () => [
                    "已获取",
                    h(NText, { type: "success" }, () => data),
                    "个站点"
                ];
                break;
            case "Test":
                counter.set(30 + (70 * data.finished) / data.count);
                content.value = () => [
                    h(NText, { type: "success" }, () => data.finished + ""),
                    "/",
                    h(NText, { type: "info" }, () => data.count + ""),
                    " 已对",
                    `${data.category}`,
                    "分类下",
                    h(
                        NButton,
                        {
                            style: "font-size:inherit",
                            tag: "a",
                            text: true,
                            href: data.url,
                            target: "_blank",
                            type: "info"
                        },
                        () => `${data.name}`
                    ),
                    `完成测试`
                ];
                break;
        }
        processStatus.value = "info";
        testing.value = true;
    }
});
</script>

<template>
    <n-h6>
        <transition name="fade-in" mode="out-in">
            <span v-if="!testing"> 实时状态 </span>
            <n-tag v-else type="warning" strong :bordered="false">
                测试中
                <template #icon>
                    <n-icon :component="ExperimentOne" />
                </template>
            </n-tag>
        </transition>
    </n-h6>
    <n-space
        vertical
        :size="12"
        style="margin: auto; padding: 0.5rem 0; max-width: 660px; width: 100%"
    >
        <div style="text-align: center">
            <n-ellipsis :line-clamp="2" :tooltip="false" style="font-size: 15px">
                <n-text depth="2">
                    <component :is="content"></component>
                </n-text>
            </n-ellipsis>
        </div>
        <n-progress
            type="line"
            :percentage="counter.count.value"
            :status="processStatus"
            :height="28"
            :indicator-placement="'inside'"
            processing
        />
    </n-space>
</template>

<style scoped lang="less"></style>
