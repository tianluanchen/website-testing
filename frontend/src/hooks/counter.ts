import { onUnmounted, ref } from "vue";
/** 返回下一个时间和添加值绝对值，如果只返回一个值则表示下一个等待时间且默认添加值绝对值为1 */
type NextTimeout = (
    start: number,
    end: number,
    originalStart: number,
    randomInt: (a: number, b: number) => number
) => number | [number, number];
const regular: NextTimeout = (start, end, originalStart, r) => {
    const ratio = Math.abs((end - start) / (end - originalStart));
    const seg = Math.min(Math.abs((end - start) / 5), 10);
    return [r(10, 50) + parseInt(Math.max(Math.max(seg, 1), 5) + "") * 1000 * (1 - ratio), seg];
};
const randomInt = (a: number, b: number) => {
    return Math.floor(Math.random() * (b - a + 1)) + a;
};
export default function useCounter(option?: { init: number; max: number }) {
    option = Object.assign({}, option, {
        init: 0,
        max: 100
    });
    const opt = option!;
    opt.init = Math.min(opt.init, opt.max);
    opt.max = Math.max(1, opt.max);
    const count = ref(opt.init);

    let id: number;
    let lastResolve: () => void;
    let originalStart = 0;
    let lastTarget = 0;
    const clear = () => {
        clearTimeout(id);
        lastResolve?.();
    };

    onUnmounted(() => {
        clear();
    });
    const wrap = (v: number | [number, number]): [number, number] => {
        if (v instanceof Array) {
            v[0] = Math.max(v[0], 0);
            v[1] = parseInt(Math.max(Math.abs(v[1]), 1) + "");
            return v;
        }
        return [Math.max(v, 0), 1];
    };
    const go = (target: number, fn: NextTimeout, cb?: () => void) => {
        const diff = target - count.value;
        const [timeout, seg] = wrap(fn(count.value, target, originalStart, randomInt));
        // console.log(target, count.value, diff, timeout, seg)
        if (Math.abs(diff) <= seg) {
            count.value = target;
            cb?.();
            return;
        }
        count.value += (diff > 0 ? 1 : -1) * seg;
        id = setTimeout(go, Math.min(timeout, 0), target, fn, cb);
    };
    const increaseTo = (to: number, fn?: NextTimeout) => {
        clear();
        return new Promise<void>((resolve) => {
            const target = Math.max(Math.min(to, opt.max), 0);
            lastResolve = resolve;
            lastTarget = target;
            originalStart = count.value;
            go(target, fn || regular, resolve);
        });
    };
    const increase = (n: number, fn?: NextTimeout) => {
        return increaseTo(count.value + n, fn);
    };
    const decreaseTo = (to: number, fn?: NextTimeout) => {
        return increaseTo(to, fn);
    };
    const decrease = (n: number, fn?: NextTimeout) => {
        return increaseTo(count.value - n, fn);
    };
    return {
        count,
        set(v: number) {
            clear();
            count.value = parseInt(Math.max(Math.min(v, opt.max), 0) + "");
            lastTarget = count.value;
        },
        done() {
            clear();
            if (lastTarget >= 0) {
                count.value = lastTarget;
            }
            lastTarget = -1;
        },
        clear,
        increaseTo,
        increase,
        decreaseTo,
        decrease
    };
}
