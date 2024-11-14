import { DatastarEvent } from "./types.ts";

export interface PatchStoreOptions extends DatastarEventOptions {
    onlyIfMissing: boolean;
};

export interface PatchStoreEvent extends DatastarEvent, PatchStoreOptions {
    type: "datastar-signal";
    data: string;
};

export interface RemoveFromStoreEvent extends DatastarEvent  {
    type: "datastar-remove";
}
