import { DatastarEvent, DatastarEventOptions } from "./types.ts";

export const fragmentMergeModes = [
    "morph"
    , "inner"
    , "outer"
    , "prepend"
    , "append"
    , "before"
    , "after"
] as const;

export type FragmentMergeMode = typeof fragmentMergeModes[number];

export interface FragmentOptions extends DatastarEventOptions {
    selector: string;
    useViewTransitions: boolean;
}

export interface RenderFragmentOptions extends FragmentOptions {
    settleDuration: number;
    mergeMode: FragmentMergeMode;
};

export interface RenderFragmentEvent extends DatastarEvent, RenderFragmentOptions {
    type: "datastar-fragment";
    fragment: string;
};

export interface RemoveFragmentsEvent extends DatastarEvent, RenderFragmentOptions {
    type: "datastar-remove";
};
