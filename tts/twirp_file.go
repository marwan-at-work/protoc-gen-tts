package tts

const twirpFile = `/* tslint:disable */

// This file has been generated by marwan.io/protoc-gen-tts.
// Do not edit.

export interface TwirpErrorJSON {
    code: string;
    msg: string;
    meta: {
        [index: string]: string;
    };
}

export class TwirpError extends Error {
    code: string;
    meta: {
        [index: string]: string;
    };

    constructor(te: TwirpErrorJSON) {
        super(te.msg);

        this.code = te.code;
        this.meta = te.meta;
    }
}

export const throwTwirpError = (resp: Response) => {
    return resp.json().then((err: TwirpErrorJSON) => {
        throw new TwirpError(err);
    });
};

export const createTwirpRequest = (body: object = {}, headers: object = {}, opts: RequestInit = {}): RequestInit => {
    return {
        method: 'POST',
        headers: { ...headers, 'Content-Type': 'application/json' },
        body: JSON.stringify(body || {}),
        ...opts,
    };
};

export type Fetch = (input: RequestInfo, init?: RequestInit) => Promise<Response>;
`
