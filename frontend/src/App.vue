<script setup lang="ts">
import { zhCN, darkTheme } from "naive-ui";
import { watchEffect } from "vue";
import AppLayout from "./components/AppLayout.vue";
import useThemeStore from "@/stores/theme";
const { isDark } = useThemeStore();

watchEffect(() => {
    isDark()
        ? document.documentElement.classList.add("dark")
        : document.documentElement.classList.remove("dark");
});
</script>
<template>
    <n-config-provider :locale="zhCN" :theme="isDark() ? darkTheme : null">
        <n-loading-bar-provider>
            <n-message-provider>
                <n-dialog-provider>
                    <app-layout></app-layout>
                </n-dialog-provider>
            </n-message-provider>
        </n-loading-bar-provider>
    </n-config-provider>
</template>

<style scoped lang="less"></style>
