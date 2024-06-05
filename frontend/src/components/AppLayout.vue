<script setup lang="ts">
import AppHeader from "./AppHeader.vue";
import AppForm from "./AppForm.vue";
import AppStatus from "./AppStatus.vue";
import AppResult from "./AppResult.vue";
import IntroPopover from "./IntroPopover.vue";
</script>

<template>
    <n-layout style="height: 100vh">
        <app-header></app-header>
        <n-layout
            position="absolute"
            style="top: 64px"
            embedded
            content-style="padding: 24px;"
            :native-scrollbar="false"
        >
            <n-card style="max-width: 1280px; margin: auto; min-height: calc(100vh - 64px - 48px)">
                <n-space vertical :size="28">
                    <app-form></app-form>
                    <app-status />
                    <app-result />
                </n-space>
            </n-card>
            <div style="text-align: center; padding-top: 24px">
                <n-text depth="2">
                    Copyright © {{ new Date().getFullYear() }} by
                    <IntroPopover>
                        <n-text type="primary" depth="2" style="cursor: pointer">Ayouth</n-text>
                    </IntroPopover>
                    , All Rights Reserved.
                </n-text>
            </div>
            <n-back-top :right="40" />
        </n-layout>
    </n-layout>
</template>

<style lang="less"></style>
