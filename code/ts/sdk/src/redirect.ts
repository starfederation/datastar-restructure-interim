import { DatastarEvent, DatastarEventOptions } from "./types.ts";

export interface RedirectEventOptions extends DatastarEventOptions {
    url: string;
}

export interface RedirectEvent extends DatastarEvent, RedirectEventOptions {
    type: "datastar-redirect";
}
