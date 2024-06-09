<script setup lang="ts">
import useThemeStore from "@/stores/theme";
const { isDark, toggleTheme } = useThemeStore();
</script>

<template>
    <n-switch :value="isDark()" @update:value="toggleTheme()">
        <template #checked>
            <span style="font-size: 13px">🌙</span>
        </template>
        <template #unchecked>
            <span style="font-size: 13px">☀️</span>
        </template>
    </n-switch>
</template>

<style scoped lang="less"></style>
