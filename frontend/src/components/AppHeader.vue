<script setup lang="ts">
import ThemeBtn from "./ThemeBtn.vue";
import IntroPopover from "./IntroPopover.vue";
const reload = () => (location.href = "/");
</script>

<template>
    <n-layout-header class="nav" style="height: 64px" bordered>
        <n-space align="center" justify="space-between" style="width: 100%">
            <IntroPopover>
                <n-avatar
                    aria-label="网站标识"
                    style="cursor: pointer; display: block"
                    round
                    :size="32"
                    src="./favicon.ico"
                />
            </IntroPopover>
            <n-button
                @click="reload"
                text
                title="刷新页面"
                style="font-size: 18px; font-weight: bold"
            >
                网站本地测试
            </n-button>
            <n-space align="center">
                <theme-btn></theme-btn>
            </n-space>
        </n-space>
    </n-layout-header>
</template>

<style lang="less">
.nav {
    padding: 0 1.25rem;
    display: flex;
    justify-content: space-between;
}
</style>
