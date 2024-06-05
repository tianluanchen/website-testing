<script setup lang="ts">
import useOptionsStore from "@/stores/options";
import { LinkSquare24Filled } from "@vicons/fluent";
import { Github } from "@icon-park/vue-next";
const { options } = useOptionsStore();
</script>

<template>
    <n-popover trigger="hover">
        <template #trigger>
            <slot></slot>
        </template>
        <n-button
            type="info"
            text
            tag="a"
            href="https://github.com/tianluanchen/website-testing"
            target="_blank"
        >
            本项目仓库
            <template #icon>
                <n-icon>
                    <Github />
                </n-icon>
            </template>
        </n-button>
        <template v-for="url in options.origins" :key="url">
            <n-divider style="margin: 3px 0" />
            <n-button
                title="Ayouth的个人网站"
                type="primary"
                text
                tag="a"
                :href="url"
                target="_blank"
                >{{ url }}
                <template #icon>
                    <n-icon>
                        <LinkSquare24Filled />
                    </n-icon>
                </template>
            </n-button>
        </template>
    </n-popover>
</template>

<style scoped lang="less"></style>
