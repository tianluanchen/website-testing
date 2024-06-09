import { ref } from "vue";
import http from "@/helpers/http";
type LabelValue = {
    label: string;
    value: string;
};
export type Options = {
    dns_servers: LabelValue[];
    user_agents: LabelValue[];
    origins: string[];
};
const options = ref<Options>({
    dns_servers: [],
    user_agents: [],
    origins: []
});
export default function useOptionsStore() {
    return {
        options,
        loadOptions() {
            return http<Options>("/api/options", {
                method: "GET"
            }).then((data) => {
                options.value = data.data;
            });
        }
    };
}
