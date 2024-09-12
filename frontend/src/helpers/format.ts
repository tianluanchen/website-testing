/** unit */
type FormatDurationOption = {
    intSeconds?: boolean;
    second: string;
    minute: string;
    hour: string;
    day: string;
};
export type { FormatDurationOption };
/**
 *
 * @param duration milliseconds
 * @param option default formatDuration.Zh
 * @returns
 */
export function formatDuration(duration: number, option?: FormatDurationOption) {
    const opt = option || formatDuration.optionZh();
    const durationSeconds = Math.floor(duration / 1000);
    // allow float
    const seconds = (duration % 60000) / 1000;
    const minutes = Math.floor((durationSeconds % (60 * 60)) / 60);
    const hours = Math.floor(durationSeconds / (60 * 60));
    const days = Math.floor(durationSeconds / (60 * 60 * 24));
    let str = "";
    if (days) {
        str += days + opt.day;
    }
    if (hours) {
        str += hours + opt.hour;
    }
    if (minutes) {
        str += minutes + opt.minute;
    }
    if (seconds) {
        if (opt.intSeconds) {
            str += Math.floor(seconds) + opt.second;
        } else {
            const fixed = seconds.toFixed(1);
            const secondsStr =
                Math.floor(seconds) === seconds
                    ? seconds
                    : fixed.endsWith(".0")
                      ? fixed.slice(0, fixed.length - 2)
                      : fixed;
            if (secondsStr !== "0" || str === "") {
                str += secondsStr + opt.second;
            }
        }
    }
    return str;
}

formatDuration.optionZh = (opt?: Partial<FormatDurationOption>) =>
    Object.assign(
        {},
        {
            second: "秒",
            minute: "分钟",
            hour: "小时",
            day: "天"
        } as Readonly<FormatDurationOption>,
        opt || {}
    );

formatDuration.optionEn = (opt?: Partial<FormatDurationOption>) =>
    Object.assign(
        {},
        {
            second: "s",
            minute: "m",
            hour: "h",
            day: "d"
        } as Readonly<FormatDurationOption>,
        opt || {}
    );

export function formatSize(size: number, separate = "") {
    const clean = (v: string) => (v.endsWith(".0") ? v.slice(0, -2) : v);
    if (size < 1024) {
        return size + separate + "B";
    } else if (size < 1024 * 1024) {
        return clean((size / 1024).toFixed(1)) + separate + "KB";
    } else {
        return clean((size / (1024 * 1024)).toFixed(1)) + separate + "MB";
    }
}
