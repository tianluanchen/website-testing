type AppResponse<T = any> = {
    code: number;
    message: string;
    data: T;
};
enum Status {
    OK = 0,
    BAD = 1
}

type Cause = {
    originalResponse: Response;
    data?: AppResponse;
    parseJSONError?: Error;
};

export default async function http<T = any>(
    input: RequestInfo | URL,
    init?: RequestInit | undefined
) {
    const resp = await fetch(input, init);
    let data: AppResponse<T>;
    try {
        data = await resp.json();
    } catch (err) {
        const error = new Error((err as Error)?.message || String(err));
        (error as any).__cause = {
            originalResponse: resp,
            parseJSONError: err
        } as Cause;
        throw error;
    }
    if (data.code !== Status.OK) {
        const error = new Error(data.message);
        (error as any).__cause = {
            originalResponse: resp,
            data
        } as Cause;
        throw error;
    }
    return data;
}
http.isAborted = (err: Error) => {
    return err.name === "AbortError";
};
http.checkError = (err: Error): Cause | null => {
    const cause = (err as any)?.__cause;
    if (!cause) {
        return null;
    }
    return cause;
};

http.isOK = <T>(data: AppResponse<T>) => {
    return data.code === Status.OK;
};
http.jsonify = (data: any, headers?: Record<string, string>) => {
    return {
        body: JSON.stringify(data),
        headers: {
            "Content-Type": "application/json; charset=utf-8",
            ...headers
        }
    };
};
