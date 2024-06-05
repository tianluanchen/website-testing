<script setup lang="ts">
import { NAutoComplete, type SelectOption } from "naive-ui";
const model = defineModel<string>("value", {
    required: true
});
const props = defineProps<{
    placeholder?: string;
    options: {
        label: string;
        value: any;
    }[];
}>();

const renderLabel = (option: SelectOption) => {
    return props.options.find((o) => o.value == option.value)?.label;
};
</script>

<template>
    <n-auto-complete
        ref="autoCompleteRef"
        :render-label="renderLabel"
        :placeholder="props.placeholder"
        :get-show="() => true"
        v-model:value="model"
        :input-props="{
            autocomplete: 'disabled'
        }"
        :options="props.options.map((e) => ({ label: e.value, value: e.value }))"
        clearable
    />
</template>

<style scoped lang="less"></style>
