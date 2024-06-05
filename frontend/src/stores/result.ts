import { ref } from "vue";
import { type Conf } from "./conf";
import http from "@/helpers/http";

type TestingResult = {
    conf: Conf;
    start: number; // unix milliseconds
    end: number; // unix milliseconds
    groups: {
        category: string;
        items: {
            name: string;
            url: string;
            status: "pending" | "done";
            result: {
                total_duration: number; // milliseconds
                records: {
                    url: string;
                    remote_addr: {
                        type: "tcp" | "udp";
                        ip: string;
                        port: number;
                    } | null;
                    duration: number; // milliseconds
                    resp: {
                        status_code: number;
                        status: string;
                        content_length: number; // -1 means unknown
                        content_type: string;
                    } | null;
                }[];
                size: number;
                title: string;
                err?: string;
                last_resp_redirect?: string; // url
            } | null;
            err?: string;
        }[];
    }[];
    err?: string;
};
export type { TestingResult };

const result = ref<{
    loading: boolean;
    testingResult: TestingResult | null;
}>({
    loading: false,
    testingResult: null
});
export default function useResultStore() {
    return {
        result,
        loadResult() {
            if (result.value.loading) {
                return Promise.resolve();
            }
            result.value.loading = true;
            return http<TestingResult>("/api/result", {
                method: "GET"
            })
                .then((data) => {
                    result.value.testingResult = data.data;
                })
                .finally(() => {
                    result.value.loading = false;
                });
        }
    };
}
