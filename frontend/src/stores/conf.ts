import { useLocalStorage, type RemovableRef } from "@vueuse/core";
export type Conf = {
    dns_server: string;
    proxy_url: string;
    timeout_seconds: number;
    user_agent: string;
};
const defaultConf = () =>
    ({
        dns_server: "",
        proxy_url: "",
        timeout_seconds: 60,
        user_agent:
            "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0"
    }) as Conf;
export default function useConfStore() {
    const conf = useLocalStorage("conf", defaultConf());
    return {
        conf,
        resetConf() {
            conf.value = defaultConf();
        }
    };
}
