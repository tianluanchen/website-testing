<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useDialog, useMessage, type FormInst, type FormRules } from "naive-ui";
import isIP from "validator/es/lib/isIP";
import isURL from "validator/es/lib/isURL";
import http from "@/helpers/http";
import useConfStore from "@/stores/conf";
import {
    Server,
    NetworkTree,
    StopwatchStart,
    Browser,
    Clear,
    Experiment
} from "@icon-park/vue-next";
import AutoComplete from "./AutoComplete.vue";
import { handleErrWithDialog } from "@/helpers/handleErr";
import ProxyInput from "./ProxyInput.vue";
import useOptionsStore from "@/stores/options";
import useProxiesStore from "@/stores/proxies";
const proxiesStore = useProxiesStore();
const dialog = useDialog();
const message = useMessage();
const { options, loadOptions } = useOptionsStore();
const { conf, resetConf } = useConfStore();
const formRef = ref<FormInst | null>(null);
const rules = ref<FormRules>({
    dns_server: {
        required: false,
        trigger: ["input", "blur"],
        message: "请输入有效IP",
        validator: (_, value) => {
            if (typeof value === "string" && value.length > 0 && !isIP(value)) {
                return new Error("invalid IP");
            }
            return true;
        }
    },
    proxy_url: {
        required: false,
        trigger: ["input", "blur"],
        message: "请输入http(s)或socks5 URL",
        validator: (_, value) => {
            if (
                typeof value === "string" &&
                value.length > 0 &&
                !isURL(value, {
                    protocols: ["https", "http", "socks5"]
                })
            ) {
                return new Error("invalid URL");
            }
            return true;
        }
    }
});

onMounted(() => {
    loadOptions().catch((err) => handleErrWithDialog(dialog, err, { title: "获取选项出错" }));
});

const submitting = ref(false);

function submit() {
    if (!formRef.value) {
        return;
    }
    formRef.value.validate((errors) => {
        if (errors) {
            message.error(errors[0][0].message || "表单验证错误");
            return;
        }
        proxiesStore.add(conf.value.proxy_url);
        submitting.value = true;
        http<{
            already: boolean;
        }>("/api/start", {
            method: "POST",
            ...http.jsonify(conf.value)
        })
            .then((data) => {
                if (data.data.already) {
                    message.warning("提交失败，已存在进行中的测试");
                } else {
                    message.success("已启动测试");
                }
            })
            .catch((err) => handleErrWithDialog(dialog, err, { title: "提交测试出错" }))
            .finally(() => {
                submitting.value = false;
            });
    });
}

function reset() {
    resetConf();
    formRef.value?.restoreValidation();
    message.success("重置成功");
}
</script>

<template>
    <n-form ref="formRef" :rules="rules" :model="conf" label-placement="top" :disabled="submitting">
        <n-grid cols="1 600:2" x-gap="0 600:32">
            <n-form-item-gi path="dns_server">
                <template #label>
                    <div class="icon-with-text">
                        <n-icon :size="15">
                            <server />
                        </n-icon>
                        DNS服务器
                    </div>
                </template>
                <AutoComplete
                    placeholder="留空则跟随系统"
                    v-model:value="conf.dns_server"
                    :options="options.dns_servers"
                />
                <!-- <n-auto-complete placeholder="留空则跟随系统" :get-show="() => true" v-model:value="conf.dns_server"
                    :input-props="{
        autocomplete: 'disabled'
    }" :options="options?.dns_servers || []" clearable /> -->
            </n-form-item-gi>
            <n-form-item-gi path="proxy_url">
                <template #label>
                    <div class="icon-with-text">
                        <n-icon :size="15">
                            <network-tree />
                        </n-icon>
                        网络代理
                    </div>
                </template>
                <ProxyInput v-model:value="conf.proxy_url"></ProxyInput>
            </n-form-item-gi>
            <n-form-item-gi path="timeout_seconds">
                <template #label>
                    <div class="icon-with-text">
                        <n-icon :size="16">
                            <stopwatch-start />
                        </n-icon>
                        访问超时时长
                    </div>
                </template>
                <n-input-number :precision="0" v-model:value="conf.timeout_seconds" :min="1">
                    <template #suffix> 秒 </template>
                </n-input-number>
            </n-form-item-gi>
            <n-form-item-gi path="user_agent">
                <template #label>
                    <div class="icon-with-text">
                        <n-icon :size="16">
                            <browser />
                        </n-icon>
                        用户代理
                    </div>
                </template>
                <AutoComplete
                    placeholder="留空可能会被待测网站屏蔽"
                    v-model:value="conf.user_agent"
                    :options="options.user_agents"
                >
                </AutoComplete>
                <!-- <n-auto-complete v-model:value="conf.user_agent" clearable :get-show="() => true"
                    :options="options?.user_agents || []"></n-auto-complete> -->
            </n-form-item-gi>
        </n-grid>
        <n-space>
            <n-button type="primary" @click="submit" :loading="submitting">
                <template #icon>
                    <n-icon>
                        <experiment />
                    </n-icon>
                </template>
                提交测试
            </n-button>
            <n-button type="info" @click="reset" :disabled="submitting">
                <template #icon>
                    <n-icon>
                        <clear />
                    </n-icon>
                </template>
                重置
            </n-button>
        </n-space>
    </n-form>
</template>

<style scoped lang="less"></style>
