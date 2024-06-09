import { useLocalStorage, type RemovableRef } from "@vueuse/core";
export default function useProxiesStore() {
    const proxies = useLocalStorage<string[]>("proxies", []);
    const add = (proxy: string) => {
        proxy = proxy.trim();
        if (proxy && !proxies.value.includes(proxy)) {
            proxies.value.unshift(proxy);
            if (proxies.value.length > 6) {
                proxies.value.splice(6);
            }
        }
    };
    const remove = (proxy: string) => {
        proxies.value = proxies.value.filter((p) => p !== proxy);
    };
    return {
        proxies,
        add,
        remove
    };
}
