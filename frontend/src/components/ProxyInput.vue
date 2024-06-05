<script setup lang="ts">
import { ref } from "vue";
import { DeleteOne } from "@icon-park/vue-next";
import { NPopover, NInput } from "naive-ui";
import useProxiesStore from "@/stores/proxies";
import { onClickOutside } from "@vueuse/core";
const show = ref(false);
const proxyURL = defineModel<string>("value");
const proxiesStore = useProxiesStore();
const popoverRef = ref<InstanceType<typeof NPopover> | null>(null);
const wrapperRef = ref<HTMLDivElement | null>(null);
const inputRef = ref<InstanceType<typeof NInput> | null>(null);
const setShow = (v: boolean) => {
    show.value = v;
};
function setProxy(proxy: string) {
    inputRef.value?.focus();
    proxyURL.value = proxy;
    setShow(false);
}
onClickOutside(inputRef, (e) => {
    const p = e.composedPath();
    if (p.includes(wrapperRef.value as any)) {
        return;
    }
    setShow(false);
});
</script>

<template>
    <n-popover
        :show="show"
        ref="popoverRef"
        width="trigger"
        trigger="manual"
        raw
        placement="bottom"
        :show-arrow="false"
    >
        <template #trigger>
            <n-input
                ref="inputRef"
                @input="setShow(true)"
                @focus="setShow(true)"
                placeholder="留空则跟随系统环境"
                v-model:value="proxyURL"
                type="text"
                clearable
            />
        </template>
        <div ref="wrapperRef" class="wrapper">
            <n-list hoverable :bordered="false">
                <n-list-item v-for="proxy in proxiesStore.proxies.value" :key="proxy">
                    <span
                        style="cursor: pointer; width: calc(100% - 25px)"
                        @click="setProxy(proxy)"
                    >
                        <n-ellipsis>
                            {{ proxy }}
                        </n-ellipsis>
                    </span>
                    <n-icon @click="proxiesStore.remove(proxy)" :size="16">
                        <DeleteOne />
                    </n-icon>
                </n-list-item>
            </n-list>
        </div>
    </n-popover>
</template>

<style scoped lang="less">
.wrapper {
    border-radius: 3px;
    overflow: hidden;

    :deep(.n-list-item) {
        padding: 10px 12px;

        .n-icon {
            cursor: pointer;
            opacity: 0.75;

            &:hover {
                opacity: 1;
            }
        }

        & > .n-list-item__main {
            width: 100%;
            display: flex;
            align-items: center;
            justify-content: space-between;
        }
    }
}
</style>
