import { useLocalStorage, usePreferredDark, type RemovableRef } from "@vueuse/core";
const theme = useLocalStorage<"dark" | "light">(
    "theme",
    usePreferredDark().value ? "dark" : "light"
);
const isDark = () => {
    return theme.value === "dark";
};
const toggleTheme = () => {
    theme.value = isDark() ? "light" : "dark";
};
export default function useThemeStore() {
    return {
        theme,
        isDark,
        toggleTheme
    };
}
