import { useLocalStorage, usePreferredDark, type RemovableRef } from "@vueuse/core";
export default function useThemeStore() {
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
    return {
        theme,
        isDark,
        toggleTheme
    };
}
