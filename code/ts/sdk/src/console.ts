import { requireThatString, assertThatString } from "@cowwoc/requirements";
import { DatastarEvent, DatastarEventOptions } from "./types.ts";

export const consoleModes = [
    "assert"
    , "clear"
    , "count"
    , "countReset"
    , "debug"
    , "dir"
    , "dirxml"
    , "error"
    , "group"
    , "groupCollapsed"
    , "groupEnd"
    , "info"
    , "log"
    , "table"
    , "time"
    , "timeEnd"
    , "timeLog"
    , "trace"
    , "warn"
] as const;

export type ConsoleMode = typeof consoleModes[number];

interface ConsoleGroupEndEvent extends DatastarEvent {
    type: "datastar-console";
    mode: "groupEnd";
}

interface ConsoleMessageEvent extends DatastarEvent {
    type: "datastar-console";
    mode: Omit<ConsoleMode, "groupEnd">;
    message: String;
};

export type ConsoleEvent = ConsoleGroupEndEvent | ConsoleMessageEvent;